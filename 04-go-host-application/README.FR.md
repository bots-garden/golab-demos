# Exécuter des plugins WebAssembly Extism à partir d'une application Go

Depuis quelques jours, nous avons vu qu'il était possible de développer des plugins WebAssembly avec le Plugin Development Kit d'Extism et les exécuter avec la CLI Extism. Aujourd'hui, il est temps de passer à un niveau supérieur : nous allons créer une application en Go qui pourra charger ces plugins et les exécuter comme le fait la CLI. 

Pour cela nous allons utiliser le **Host SDK** d'Extism pour le langage Go. Pour rappel Extism fournit des Host SDK pour de nombreux langages (https://extism.org/docs/category/integrate-into-your-codebase).

Pour rappel, une application hôte est une application qui grâce à un SDK de runtime Wasm, est capable d'éxécuter des programmes WebAssembly. les **Host SDK** d'Extism sont des "sur-couches" au SDK de runtime Wasm pour vous simplifier la vie (éviter la plomberie compliquée).

À l'heure actuelle, Extism utilise le runtime **[WasmTime](https://wasmtime.dev/)**. 

> Si je me réfère à cette [issue (WASI threads support)](https://github.com/extism/extism/issues/357), il n'est pas impossible que le support d'autres runtime Wasm soient pris en compte, et notamment [Wazero](https://wazero.io/).

Mais assez parlé, passons à la pratique.

## Pré-requis

Vous aurez besoin de 

- Go (v1.20)
- Extism 0.4.0 : [Install Extism](https://extism.org/docs/install)

## Création de l'application

Commencez par créer un fichier `go.mod` avec la commande `go mod init go-host-application`, puis un fichier `main.go` avec le contenu suivant :

```golang
package main

import (
	"fmt"

	"github.com/extism/extism"
)

func main() {

	ctx := extism.NewContext()

    // This will free the context and all associated plugins
	defer ctx.Free() 

    // Path to the wasm file 0️⃣
	path := "../03-even-with-javascript/hello-js.wasm"
    
    // Define the path to the wasm file 1️⃣
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: path},
		}}

    // Load the wasm plugin 2️⃣
	plugin, err := ctx.PluginFromManifest(
        manifest, 
        []extism.Function{}, 
        true,
    )

	if err != nil {
		panic(err)
	}

    // Call the `say_hello` function 3️⃣
    // with a string parameter
	res, err := plugin.Call(
		"say_hello",
		[]byte("👋 Hello from the Go Host app 🤗"),
	)

	if err != nil {
		fmt.Println("😡", err)
	} else {
        // Display the return value 4️⃣
		fmt.Println("🙂", string(res))
	}

}
```

Vous voyez, le code est très très simple :

- 0: utilisons le plugin Wasm JavaScript que nous avons développé dans le précédent article.
- 1: définir un manifest avec des propriétés dont le chemin pour accéder au fichier Wasm.
- 2: charger le plugin Wasm.
- 3: appeler la fonction `say_hello` du plugin.
- 4: afficher le résultat (le type de `res` est `[]byte`).

### Lancer le programme

Utiliser tout simplement cette commande:

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go
```
> You need to set the linker lookup path env var explicitly.


Et vous obtiendrez ceci:

```bash
🙂 param: 👋 Hello from the Go Host app 🤗
```

Vous pouvez bien sûr faire le test avec le premier plugin développé avec TinyGo. Changez la valeur de la variable `	path := "../01-simple-go-plugin/simple.wasm"` et relancez :

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go
```

Et vous devrier obtenir ceci:

```bash
🙂 👋 Hello 👋 Hello from the Go Host app 🤗
```

🎉 vous voyez, il est facile de créer des applications en Go qui soient capables d'exécuter des plugins Wasm écrits dans différents langages.

Si j'arrive à conserver le rythme, demain je vous explique comment faire la même chose mais cette fois ci avec Node.js.

