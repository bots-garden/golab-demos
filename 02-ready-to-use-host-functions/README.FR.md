# Extism, Plugins WebAssembly & Host functions

Comme je l'expliquais dans le précédent article, les programmes Wasm sont "limités" par défaut (c'est aussi un principe de sécurité). La spécificaton WASI est un ensemble d'API qui permettront de développer des programmes Wasm (WebAssembly) qui pourront accéder aux ressources du système (si l'application hôte le permet). Aujourd'hui, la spécification WASI est en cours d'écriture et peu d'API sont disponibles. Donc, même s'il est possible d'utiliser l'API socket ou FileSystem, les capacités d'un programme Wasm en termes d'accès aux ressources systèmes sont réduites: pas d'affichage dans le terminal, pas d'accès HTTP, ... 

## Nous sommes sauvés, nous avons les host functions

Néanmoins, pour nous faciliter la vie, il est possible pour l'application hôte de fournir à l'invité (le module Wasm) des pouvoirs en plus. Nous appelons ça les "host functions". C'est une fonction développée "dans le code source de l'hôte". Celui-ci l'expose (export) au module Wasm qui sera capable de l'exécuter. Par exemple vous pouvez développer une host function pour faire des affichages de message et permettre ainsi au module Wasm d'afficher des message dans un terminal pendant son exécution. 

> **Attention**: vous devez noter que du moment que vous utilisez des host functions, votre module Wasm ne sera exécutable que par votre application hôte.

Je vous expliquais hier que le système de type de la spécification WASI en ce qui concerne les passages de paramètres à une fonction et sa valeur de retour, est très limité (uniquement des nombres). Cela implique une gymnastique un peu "acrobatique" pour développer une host function.

## Extism fourni des host functions prêtes à l'emploi

 Pour vous aider à développer des programmes Wasm sans vous préoccuper de la complexité, le Plugin Development Kit (PDK) d'Extism apporte quelques host functions prêtes à l'emploi, notamment pour faire des logs, des requêtes HTTP ou de lire une configuration en mémoire.

### Création d'un nouveau plugin Wasm (avec le PDK Extism)

Commencez par créer un fichier `go.mod` avec la commande `go mod init ready-to-use-host-functions`, puis un fichier `main.go` avec le contenu suivant :

```golang
package main

import (
	"github.com/extism/go-pdk"
	"github.com/valyala/fastjson"
)

//export say_hello
func say_hello() int32 {

	// read function argument from the memory
	input := pdk.Input()

    // 1️⃣ write information to the logs
	pdk.Log(pdk.LogInfo, "👋 hello this is wasm 💜") 

    // 2️⃣ get the value associated to the `route` key 
    // into the config object
	route, _ := pdk.GetConfig("route")
    // the value of `route` is
    // https://jsonplaceholder.typicode.com/todos/1

    // 3️⃣ write information to the logs
	pdk.Log(pdk.LogInfo, "🌍 calling "+route)

    // 4️⃣ make an HTTP request
	req := pdk.NewHTTPRequest("GET", route)
	res := req.Send()
	
    // Read the result of the request
	parser := fastjson.Parser{}
	jsonValue, _ := parser.Parse(string(res.Body()))
	title := string(jsonValue.GetStringBytes("title"))

    // Prepare the return value
	output := "param: " + string(input) + " title: " + title

	mem := pdk.AllocateString(output)
	// copy output to host memory
	pdk.OutputMemory(mem)

	return 0
}

func main() {}
```

- 1: `pdk.Log` est une host function présente dans la CLI Extism, elle permet d'envoyer des messages dans les logs.
- 2: la host function `pdk.GetConfig` permet d'aller lire les valeur d'une configuration passée en mémoire par l'application hôte. Dans cet exemple nous allons récupérer une URL à utiliser pour faire une requête HTTP.
- 3: nous utilisons à nouveau `pdk.Log` pour envoyer dans les logs la valeur de associée à la clé de configuration `route`
- 4: la host function `pdk.NewHTTPRequest` permet de faire des requêtes HTTP.

Lorsque vous développerez vos propres applications avec le SDK Extism, elles aussi proposeront ces même host functions (il est aussi possible de développer ses propre host functions, mais ce sera pour plus tard).

 > **Remarque**: TinyGo dispose du support de la serialisation/déserialisation JSON, mais je continue à utiliser `fastjson` qui est plus rapide pour mes cas d'usage.

Maintenant, testons notre nouveau plugin Wasm.

### Compiler le plugin Wasm

Pour compiler le programme, utiliser TinyGo et la commande ci-dessous, qui produira un fichier `simple.wasm` :

```bash
tinygo build -scheduler=none --no-debug \
  -o host-functions.wasm \
  -target wasi main.go
```

### Exécuter la fonction `say_hello` du plugin Wasm

Pour cela nous allons utiliser la CLI Extism (nous verrons comment développer notre propre application hôte dans un prochain article).

Pour exécuter la fonction `say_hello` avec comme paramètre la chaine de caractères `"😀 Hello World 🌍! (from TinyGo)"`, utilisez la commande suivante:

```bash
extism call ./host-functions.wasm \
  say_hello --input "😀 Hello World 🌍! (from TinyGo)" \
  --wasi \
  --log-level info \
  --allow-host '*' \
  --config route=https://jsonplaceholder.typicode.com/todos/1 
```

- Pour afficher les logs, vous devez préciser le niveau de log avec `--log-level info`.
- Pour permettre au module Wasm de faire une requête HTTP "vers l'expérieur", vous devez lui permettre de le faire en précisant `--allow-host '*'`.
- Et enfin avec le flag `--config`, vous pouvez "pousser" des informations de configuration à destination du programme Wasm.

Et vous obtiendrez en sortie :

```bash
extism_runtime::pdk INFO 2023-07-17T06:45:27.063609583+02:00 - 👋 hello this is wasm 💜
extism_runtime::pdk INFO 2023-07-17T06:45:27.063691542+02:00 - 🌍 calling https://jsonplaceholder.typicode.com/todos/1
param: 😀 Hello World 🌍! (from TinyGo) title: delectus aut autem
```

🎉 Et voilà, c'est terminé pour ce deuxième article de découverte d'Extism.
👋 À bientôt pour le prochain article (comment faire un plugin JavaScript).
