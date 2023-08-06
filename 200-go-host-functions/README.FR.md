# Ecrire des Host Functions en Go avec Extism

Dans cet article :

- Sur le même principe que l'article précédent : [Write a host function with the Extism Host SDK](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk), nous allons modifier l'application hôte développée en **Go** dans l'article [Run Extism WebAssembly plugins from a Go application](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application) pour lui ajouter une **host function** développée par nos soins.
- Nous utiliserons exactement le même plugin que celui modifié lors de l'article précédent à la section [Modify the Wasm plugin](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk#heading-modify-the-wasm-plugin) pour pouvoir appeler la host function.

## Pré-requis

Vous aurez besoin de 

- Go (v1.20) et TinyGo (v0.28.1) pour compiler les plugins
- Extism 0.4.0 : [Install Extism](https://extism.org/docs/install)
- Et avoir au moins lu l'article précédent : [Write a host function with the Extism Host SDK](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk) (mais probablement aussi tous les articles de la série).

## Modification de l'application hôte écrite en Go

L'objectif est le même que pour l'article précédent : développer une host function qui permet de récupérer des messages stockés dans la mémoire de l'application hôte à partir d'une clé. Nous allons pour cela utiliser une [Map Go](https://go.dev/blog/maps). Et cette fonction sera utilisée (appelée) par le plugin Wasm.

**Important** : Pour implémenter des host functions, le Host SDK Go d'Extism utilise le package "Golang CGO" (qui permet d'invoquer du code C à partir de Go et inversement)
> cf. la documentation : [go-host-sdk/#host-functions](https://extism.org/docs/integrate-into-your-codebase/go-host-sdk/#host-functions)

Voici donc le code modifié de l'application :

```golang
package main

import (
	"fmt"
	"unsafe"
	"github.com/extism/extism"
)

// 1️⃣
/* 
#include <extism.h>
EXTISM_GO_FUNCTION(memory_get);
*/
import "C" // 2️⃣

// 3️⃣ define a map with some records
var memoryMap = map[string]string{
	"hello": "👋 Hello World 🌍",
	"message": "I 💜 Extism 😍",
}

// 4️⃣ host function definition (callable by the Wasm plugin)
//export memory_get
func memory_get(plugin unsafe.Pointer, inputs *C.ExtismVal, nInputs C.ExtismSize, outputs *C.ExtismVal, nOutputs C.ExtismSize, userData uintptr) {

    // input parameters
	inputSlice := unsafe.Slice(inputs, nInputs)
    // output value
	outputSlice := unsafe.Slice(outputs, nOutputs)

	currentPlugin := extism.GetCurrentPlugin(plugin)

    // 5️⃣ Read the value of inputs from the memory
	keyStr := currentPlugin.InputString(unsafe.Pointer(&inputSlice[0]))

    // 6️⃣ get the associated string value
	returnValue := memoryMap[keyStr]

    // 7️⃣ copy the return value to the memory
	currentPlugin.ReturnString(unsafe.Pointer(&outputSlice[0]), returnValue)

}

func main() {

	// Function is used to define host functions
    // 8️⃣ define a slice of host functions
	hostFunctions := []extism.Function{
		extism.NewFunction(
			"hostMemoryGet",
			[]extism.ValType{extism.I64},
			[]extism.ValType{extism.I64},
			C.memory_get,
			"",
		),
	}

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

    // 9️⃣ use the updated plugin
	path := "../12-simple-go-mem-plugin/simple.wasm"

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

    // 1️⃣0️⃣ 
	plugin, err := ctx.PluginFromManifest(
		manifest,
		hostFunctions,
		true,
	)

	if err != nil {
		panic(err)
	}

	res, err := plugin.Call(
		"say_hello",
		[]byte("👋 Hello from the Go Host app 🤗"),
	)

	if err != nil {
		fmt.Println("😡", err)
		//os.Exit(1)
	} else {
		//fmt.Println("🙂", res)
		fmt.Println("🙂", string(res))
	}
}
```

- 1: la première étape est de déclarer l'utilisation de `EXTISM_GO_FUNCTION` avec le nom de la fonction qui sera utilisée. 
- 2: ne pas oublier d'importer le package `"C"`.
- 3: créer une `map` avec quelques éléments. Cette `map` sera utilisée par la host function.
- 4: définition de la host function `memory_get`. Ne pas oublier d'exporter la fonction avec `//export memory_get` (une host function aura toujours la même signature).
- 5: lorsque la host function est appelée par le plugin Wasm, le passage de paramètres se fait à l'aide de la mémoire partagée entre le plugin et l'hôte. `currentPlugin.InputString(unsafe.Pointer(&inputSlice[0]))` sert à aller chercher cette information dans la mémoire partagée. `keyStr` est une string qui contient la clé pour retrouver une valeur dans la `map`. 
- 6: aller lire la valeur associée à la clé dans la `map`.
- 7: copier la valeur obtenue en mémoire pour permettre au plugin Wasm de la lire
- 8: on définit un tableau de host functions. Dans notre cas nous en créons une seule, où `"hostMemoryGet"` sera l'alias de la fonction "vue" par le plugin Wasm, `[]extism.ValType{extism.I64}` représente le type du paramètre d'entrée et le type du paramètre de retour (on se souvient que les fonctions Wasm n'acceptent que des nombres - et dans notre cas ces nombres contiennent les positions et tailles des valeurs dans la mémoire partagée) et enfin `C.memory_get` qui est la définition de notre host function.
- 9: utiliser le plugin Wasm modifié
- 10: créer une instance du plugin wasm en lui passant en paramètre le tableau de host functions.


**Rappel** : le code du plugin wasm modifié (écrit en Go) est ici : [Plugin Wasm Go](https://k33g.hashnode.dev/write-a-host-function-with-the-extism-host-sdk#heading-modify-the-wasm-plugin)

## Exécuter l'application

Pour tester votre nouvelle application hôte il vous suffit de lancer la commande suivante :

```bash
go run main.go 
```

Et vous obtiendrez ceci, avec les messages de chacune des clés de la `map` Go:

```bash
🙂 👋 Hello 👋 Hello from the Go Host app 🤗
key: hello, value: 👋 Hello World 🌍
key: message, value: I 💜 Extism 😍
```

🎉 voilà, nous avons écrit une host function en Go, utilisable avec le même plugin wasm (sans le modifier). Donc ce plugin pourra indifférement appeler une host function écrite an JavaScript, Go, Rust ... si l'application qui l'utilise a implémenter cette host function avec la même signature.

