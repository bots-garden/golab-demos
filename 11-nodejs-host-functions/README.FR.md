# Ecrire une host function avec le Host SDK d'Extism

Dans cet article, nous allons:

- Modifier l'application hôte développée avec Node.js dans [l'article précédent](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism) pour lui ajouter une **host function** développée par nos soins.
- Modifier le plugin Wasm développé en Go dans [l'article de départ de la série](https://k33g.hashnode.dev/extism-webassembly-plugins). pour qu'il utilise cette **host function**.

## Petit rappel sur les host functions

**Extrait de ["Extism, Plugins WebAssembly & Host functions"](https://k33g.hashnode.dev/extism-webassembly-plugins-host-functions)** :

*Il est possible pour l'application hôte de fournir à l'invité (le module Wasm) des pouvoirs en plus. Nous appelons ça les "host functions". C'est une fonction développée "dans le code source de l'hôte". Celui-ci l'expose (export) au module Wasm qui sera capable de l'exécuter. Par exemple vous pouvez développer une host function pour faire des affichages de message et permettre ainsi au module Wasm d'afficher des message dans un terminal pendant son exécution...*

*... le Plugin Development Kit (PDK) d'Extism apporte quelques host functions prêtes à l'emploi, notamment pour faire des logs, des requêtes HTTP ou de lire une configuration en mémoire.*

Mais avec le Host SDK d'Extism, vous pouvez développer vos propres host functions. Cela peut être utile par exemple pour de l'accès à de la base de donnée, de l'interaction avec des brokers MQTT ou Nats...

Dans cet article, nous resterons simples et allons développer une host function qui permet de récupérer des messages stockés dans la mémoire de l'application hôte à partir d'une clé. Nous allons pour cela utiliser une [Map JavaScript](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Map).

Commençons par modifier notre application Node.js.

## Développement de la host functions

Modifiez le fichier `server.js` de la façon suivante : 

```javascript
import Fastify from 'fastify'
import process from "node:process"

// 1️⃣
import { Context, HostFunction, ValType } from '@extism/extism'
import { readFileSync } from 'fs'

// 2️⃣
let memoryMap = new Map()

memoryMap.set("hello", "👋 Hello World 🌍")
memoryMap.set("message", "I 💜 Extism 😍")

// 3️⃣ Host function (callable by the WASM plugin)
function memoryGet(plugin, inputs, outputs, userData) { 

  // 4️⃣ Read the value of inputs from the memory
  let memKey = plugin.memory(inputs[0].v.i64)
  // memKey is a buffer, 
  // use toString() to get the string value
  
  // 5️⃣ This is the return value
  const returnValue = memoryMap.get(memKey.toString())
  
  // 6️⃣ Allocate memory
  let offs = plugin.memoryAlloc(Buffer.byteLength(returnValue))
  // 7️⃣ Copy the value into memory
  plugin.memory(offs).write(returnValue)
  
  // 8️⃣ return the position and the length for the wasm plugin
  outputs[0].v.i64 = offs 
}

// 9️⃣ Host functions list
let hostFunctions = [
  new HostFunction(
    "hostMemoryGet",
    [ValType.I64],
    [ValType.I64],
    memoryGet,
    "",
  )
]

// location of the new plugin
let wasmFile = "../12-simple-go-mem-plugin/simple.wasm"
let functionName = "say_hello"
let httpPort = 7070

let wasm = readFileSync(wasmFile)

const fastify = Fastify({
  logger: true
})

const opts = {}

// Create the WASM plugin
let ctx = new Context()

// 1️⃣0️⃣
let plugin = ctx.plugin(wasm, true, hostFunctions)

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => {

    // Call the WASM function, 
    // the request body is the argument of the function
    let buf = await plugin.call(functionName, request.body); 
    let result = buf.toString()

    return result
  })

  try {
    await fastify.listen({ port: httpPort, host: '0.0.0.0'})
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start().then(r => console.log("😄 started"))
```

- 1: importer `HostFunction` (qui permet à l'hôte de définir des fonctions appelables par le plugin Wasm) et `ValType` (une énumération des types possibles utilisables par la host function).
- 2: création et alimentation d'une `Map` JavaScript
- 3: définition de la host function `memoryGet`
- 4: lorsque la host function est appelée par le plugin Wasm, le passage de paramètres se fait à l'aide de la mémoire partagée entre le plugin et l'hôte. `plugin.memory(inputs[0].v.i64)` sert à aller chercher cette information dans la mémoire partagée. `memKey` est un buffer qui contient la clé pour retrouver une valeur dans la `Map` JavaScript (et on utilise `memKey.toString()` pour transformer le buffer en string).
- 5: on récupère la valeur associée à la clé.
- 6: on alloue de la mémoire pour pouvoir y copier la valeur associée à la clé. `offs` correspond à la position et la longueur de la valeur en mémoire (c'est grâce à la méthode de bit-shifting que l'on peut "faire rentrer 2 valeur dans une seule").
- 7: on copie la valeur `returnValue` dans cette mémoire à l'endroit indiqué `offs`.
- 8: on copie dans la variable de retour `outputs` (passée à la fonction par référence) la valeur de `offs` qui permettra au plugin wasm de lire en mémoire le résultat de la fonction.
- 9: on définit un tableau de host functions. Dans notre cas nous en créons une seule, où `"hostMemoryGet"` sera l'alias de la fonction "vu" par le plugin Wasm, `[ValType.I64]` représente le type du paramètre d'entrée et le type du paramètre de retout (on se souvient que les fonctions Wasm n'acceptent que des nombres - et dans notre cas ces nombres contiennent les positions et tailles des valeurs dans la mémoire partagée) et enfin `memoryGet` qui est la définition de notre host function.
- 10: En instanciant le plugin Wasm, on passe en argument le tableau de host functions.

Avant de pouvoir exécuter à nouveau notre serveur HTTP, nous allons devoir modifier notre plugin Wasm.

## Modification du plugin Wasm

```golang
package main

import (
	"strings"
	"github.com/extism/go-pdk"
)


//export hostMemoryGet // 1️⃣
func hostMemoryGet(x uint64) uint64

//export say_hello
func say_hello() int32 {

	// read function argument from the memory
	// this is the name passed to the function
	input := pdk.Input()

	// Call the host function
	// 2️⃣
	key1 := pdk.AllocateString("hello")
	// 3️⃣
	offs1 := hostMemoryGet(key1.Offset())

  // 4️⃣
	mem1 := pdk.FindMemory(offs1)
	/*
		mem1 is a struct instance
		type Memory struct {
			offset uint64
			length uint64
		}
	*/

	// 5️⃣
	buffMem1 := make([]byte, mem1.Length())
	mem1.Load(buffMem1)

	// 6️⃣ get the second message
	key2 := pdk.AllocateString("message")
	offs2 := hostMemoryGet(key2.Offset())
	mem2 := pdk.FindMemory(offs2)
	buffMem2 := make([]byte, mem2.Length())
	mem2.Load(buffMem2)

  // 7️⃣
	data := []string{
		"👋 Hello " + string(input),
		"key: hello, value: " + string(buffMem1),
		"key: message, value: " + string(buffMem2),
	}

	// Allocate space into the memory
	mem := pdk.AllocateString(strings.Join(data, "\n"))
	// copy output to host memory
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
```

- 1: la fonction `hostMemoryGet` doit être exportée pour être utilisable.
- 2: nous voulons appeler la host function pour obtenir la valeur correspondant à la clé `hello`, donc pour cela nous devons copier cette clé en mémoire.
- 3: on appelle la host function `hostMemoryGet` (`key1.Offset()` représente la position et la longueur en mémoire de la clé `key1` into only one value).
- 4: `pdk.FindMemory(offs1)` permet de récupérer une structure `mem1` contenant la position et la longueur.
- 5: on peut maintenant créer un buffer `buffMem1` avec la taille de la valeur à récupérer et le charger avec le contenue de l'emplacement mémoire (`mem1`). Il suffira ensuite de lire la chaîne de caractères avec `string(buffMem1)`.
- 6: on recommence pour lire la deuxième clé.
- 7: on construit un slice de strings que l'on transformera ensuite en une seule string pour la renvoyer à la fonction hôte.

> Si vous souhaitez approfondir le sujet de la mémoire partagée entre l'hôte et le plugin wasm, vous pouvez lire ce blog post : https://k33g.hashnode.dev/wasi-communication-between-nodejs-and-wasm-modules-with-the-wasm-buffer-memory

### Compilez le nouveau plugin

Pour compiler le programme, utiliser TinyGo et la commande ci-dessous, qui produira un fichier `simple.wasm` :

```bash
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

Il est temps de tester nos modifications.

## Lancer le serveur et appler le MicroService

Pour démarrer le serveur, utilisez tout simplement cette commande:

```bash
node server.js
```

Ensuite, pour appeler le MicroService, utilisez cette simple commande `curl` :

```bash
curl -X POST http://localhost:7070 \
-H 'Content-Type: text/plain; charset=utf-8' \
-d 'Jane Doe'
```

Et vous obtiendrez les messages de chacune des clés de la `Map` Javascript :

```bash
👋 Hello Jane Doe
key: hello, value: 👋 Hello World 🌍
key: message, value: I 💜 Extism 😍
```

Retenez bien que lorsque le plugin Wasm appelle la host function, ce n'est pas lui qui exécute le traitement, mais bien l'application hôte. Dans le cas de Node.js, cela ralentira évenuellement l'exécution du plugin, car Node.js est générallement moins rapide que du Go compilé. Néanmoins le potentiel des host functions est très intéressant.

😥 Cet article était un peu plus compliqué que les précédent, mais ce concept de host functions est incontournable. Ces deux derniers articles vous montre aussi de quelle manière vous pouvez faire évoluer vos applications Node.js avec d'autres langages. N'hésitez pas à me contacter pour plus d'explications. Mon prochain article expliquera aussi comment faire des host function, mais cette fois-çi en Go.

