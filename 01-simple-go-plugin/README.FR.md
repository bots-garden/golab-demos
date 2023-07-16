# Extism & Plugins WebAssembly

Extism est un ensemble de projets de SDK permettant de développer des applications qui exécutent des plugins WebAssembly, mais aussi de développer les plugins WebAssembly.

Extism fournit plusieurs SDK pour réaliser les applications hôtes (donc celles qui vont charger et exécuter les plugins WebAssembly) et ce, pour différents langages (Rust, Go, Ruby, PHP, JavaScript, Java, Erlang, Haskell, Zig, .Net, C, Swift, OCaml).

Et comme je l'écrivait plus haut, Extism fourni aussi des PDK (Plug-in Development Kit) pour développer des plugins WebAssembly en Go, Rust, Haskell, C, Zig, AssemblyScript et JavaScript.

## Wasi

Pour exécuter des plugins WebAssembly à partir d'applications hôtes qui ne sont pas un navigateur, vous de vez utiliser la spécification Wasi. Les runtimes WebAssembly (WasmEdge, WasmTimen, Wazero, Wasmer, ...) implémentent cette spécification et fournissent des SDK pour créer les applications hôtes. Néanmoins, la spécification Wasi est loin d'être terminée et viens avec quelques limitations. 

Par exemple, les fonctions d'un programme Wasm n'acceptent que des nombres en paramètres et ne peuvent retourner qu'un seul nombre, ce qui signifie qu'utiliser des strings comme paramètres et valeurs de retour n'est pas trivial. Sachez qu'en jouant avec la mémoire partagée entre l'hôte et le programme Wasm il est possible de contourner ça (et c'est très formateur d'apprendre à le faire).

Autre exemple de limitation : une fonction Wasm ne pourra pas faire de requêtes HTTP, ou écrire dans la console (afficher un résultat). Le workaround est de créer des fonctions dans l'applications hôte pour exécuter ces traitement, et de les exposer au module wasm pour lui permettre de les utiliser. Là non plus ce n'est pas trivial.

## Heureusement nous avons Extism !

Toute la "plomberie" nécessaire pour contourner les limitations de la spécification Wasi est offerte par Extism, et finalement il devient simple de développer par exemple une application Go qui pourra indifférement exécuter des plugins Wasm développés en Go, Rust et même JavaScript! (et d'autres langages bien sûr).

Extism vient aussi avec une CLI qui vous permet de tester vos plugins. Donc aujourd'hui, nous ne traiterons que du développement des plugins WebAssembly.

## Pré-requis

Pour reproduire les exemples de cet articles vous aurez besoin de:

- Go (v1.20) & TinyGo (v0.28.1)
- Node.js (v19.9.0)
- Extism 0.4.0 & Extims-js PDK pour builder des modules Wasm avec du JavaScript
  - [Install Extism](https://extism.org/docs/install)
  - [Install Extism-js PDK](https://extism.org/docs/write-a-plugin/js-pdk#how-to-install-and-use-the-extism-js-pdk)


Mais voyons comment créer notre premier plugin Wasm.

## Premier plugin en Go

Commencez par créer un fichier `go.mod` avec la commande `go mod init 01-simple-go-plugin/README.md`, puis un fichier `main.go` avec le contenu suivant :

```golang
package main

import (
	"github.com/extism/go-pdk"
)

//export say_hello 1️⃣
func say_hello() int32 {

	// read function argument from the memory
	input := pdk.Input() //2️⃣

	output := "👋 Hello " + string(input)

	mem := pdk.AllocateString(output) //3️⃣
	// copy output to host memory
	pdk.OutputMemory(mem) //4️⃣

	return 0
}

func main() {}
```

**Remarques**:
- 1️⃣: l'annotation `//export say_hello` est obligatoire, pour que la fonction `say_hello` soit "visible" par l'application hôte (qui sera la CLI Extism).
- 2️⃣: `pdk.Input()` permet de lire la mémoire partagée entre le module Wasm et l'application hôte pour en extraire un buffer (`[]byte`) contenant le paramètre envoyé par la fonction hôte.
- 3️⃣: allouer une place mémoire pour lea valeur de retour
- 4️⃣: copier la valeur en mémoire (elle sera utilisable par l'application hôte)

### Compiler le plugin Wasm

Pour compiler le programme, utiliser TinyGo et la commande ci-dessous, qui produira un fichier `simple.wasm` :

```bash
tinygo build -scheduler=none --no-debug \
  -o simple.wasm \
  -target wasi main.go
```

### Exécuter la fonction `say_hello` du plugin Wasm

Pour cela nous allons utiliser la CLI Extism (nous verrons comment développer notre propre application hôte dans un prochain article).

Pour exécuter la fonction `say_hello` avec comme paramètre la chaine de caractères `"Lisa"`, utilisez la commande suivante:

```bash
extism call ./simple.wasm \
  say_hello --input "Lisa" \
  --wasi
```

Et vous obtiendrez en sortie :

```bash
👋 Hello Lisa
```

Voilà, c'est tout pour aujourd'hui. Dans les prochains articles nous verrons :
- Comment utiliser les host functions "ready to use" apportées par Extism
- Comment faire un plugin Wasm avec du JavaScript
- Comment développer une application hôte en Go
- Et probablement plus 😉

