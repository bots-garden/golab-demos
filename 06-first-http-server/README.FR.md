# MicroServices Wasm avec Extism et Fiber

Aujourd'hui, je vais vous montrer rapidement comment servir des plugins Extism (donc des plugins Webassembly) avec l'excellent framework [Fiber](https://docs.gofiber.io/). **Fiber** est un framework web pour faire des serveurs HTTP avec un esprit similaire aux frameworkx Node.js comme [Express](https://expressjs.com/) (que j'ai maintes fois utilisé dans le passé) ou [Fastify](https://fastify.dev/).

Cer article sera "légèrement" plus long que les précédents, car je voudrais aussi vous parler de mes erreurs lors de mon apprentissage avec Wasi.

## Pré-requis

- Au mieux : avoir lu tous les articles de blog de cette série ["Discovery of Extism (The Universal Plug-in System)"](https://k33g.hashnode.dev/series/extism-discovery)
- À minima : 
  - [Extism & WebAssembly Plugins](https://k33g.hashnode.dev/extism-webassembly-plugins)
  - [Run Extism WebAssembly plugins from a Go application](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application)
  - [Create a Webassembly plugin with Extism and Rust](https://k33g.hashnode.dev/create-a-webassembly-plugin-with-extism-and-rust)

## Création d'un serveur HTTP comme application hôte

Commencez par créer un fichier `go.mod` avec la commande `go mod init first-http-server`, puis un fichier `main.go` avec le contenu suivant :

```go
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

func main() {
    // Parameters of the program 0️⃣
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
    httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free()

    // Define the path to the wasm file 1️⃣
	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

    // Load the wasm plugin 2️⃣
	plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
	if err != nil {
		panic(err)
	}

    // Create an instance of Fiber application 3️⃣
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

    // Create a route "/" and a handler to call the wasm function 4️⃣
	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Call the wasm function 5️⃣
        // with a string parameter
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
            // Send the HTTP response to the client 6️⃣
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

    // Start the HTTP server 7️⃣
	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

Si vous avez lu les articles précédents, une partie du code est déjà familière pour vous.

- 0: utiliser les paramètres du programme pour lui passer les informations suivantes : le chemin du plugin wasm, le nom de la fonction à appeler et le port HTTP.
- 1: définir un manifest avec des propriétés dont le chemin pour accéder au fichier Wasm.
- 2: charger le plugin Wasm.
- 3: créer une application Fiber.
- 4: créer une route "/" qui sera déclenchée par une requête HTTP de type `POST`.
- 5: appeler la fonction du plugin.
- 6: renvoyer le résultat (la réponse HTTP).
- 7: démarrer le serveur.

### Démarrer le serveur et servir le plugin WASM

Nous allons utiliser le plugin wasm développé en Rust de notre précédent article [Create a Webassembly plugin with Extism and Rust](https://k33g.hashnode.dev/create-a-webassembly-plugin-with-extism-and-rust).

Démarrez l'application de la façon suivante :

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

Vous devriez obtenir ceci :

```bash
🌍 http server is listening on: 8080
```

Et maintenant, faites une requête HTTP :

```bash
curl -X POST \
http://localhost:8080 \
-H 'content-type: text/plain; charset=utf-8' \
-d '😄 Bob Morane'
echo ""
```

Et vous obtiendrez ceci :

```bash
{"message":"🦀 Hello 😄 Bob Morane"}
```

### Stresser l'application ... C'est le drame !

Je vérifie toujours le comportement de mes services web en les "stressant" grâce à l'utilitaire [Hey](https://github.com/rakyll/hey) qui est extrêmement facile à utiliser (notamment avec des jobs de CI pour par exemple, vérifier les performances avant et après modifications).

Je vais donc stresser une première fois mon service avec la commande suivante :

```bash
hey -n 300 -c 1 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

Je vais donc faire 300 requêtes HTTP vers mon service avec 1 seule connexion.

Et je vais obtenir un rapport de ce type (c'est juste un extrait) :

```bash
Summary:
  Total:        0.0973 secs
  Slowest:      0.0125 secs
  Fastest:      0.0001 secs
  Average:      0.0003 secs
  Requests/sec: 3082.8745
  
  Total data:   9900 bytes
  Size/request: 33 bytes

Status code distribution:
  [200] 300 responses
```
> Je travaille avec un Mac M1 Max


Maintenant, vérifions le comportement du service avec plusieurs connexions en même temps :

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

Je vais donc faire 300 requêtes HTTP vers mon service avec 100 connexions simultannées.

Et là cette fois çi mon serveur HTTP va crasher ! Et dans le rapport de test de charge vous allez voir que la plupart des requêtes sont en erreur (c'est juste un extrait) :

```bash
Status code distribution:
  [200] 3 responses

Error distribution:
  [3]   Post "http://localhost:8080": EOF
  [196] Post "http://localhost:8080": dial tcp 127.0.0.1:8080: connect: connection refused
  [1]   Post "http://localhost:8080": read tcp 127.0.0.1:38552->127.0.0.1:8080: read: connection reset by peer
  [1]   Post "http://localhost:8080": read tcp 127.0.0.1:38568->127.0.0.1:8080: read: connection reset by peer
```

### Mais que c'est-il passé ?

Dans le paragraphe ["WASI" du premier article de la série](https://k33g.hashnode.dev/extism-webassembly-plugins#heading-wasi), j'expliquais que le moyen d'échanger des valeurs autres que des chiffres entre l'application hôte et lt le plugin wasm "invité" est d'utiliser la mémoire webassemble partagée. 

> Je vous engage à lire cet excellent article sur le sujet [A practical guide to WebAssembly memory](https://radu-matei.com/blog/practical-guide-to-wasm-memory/) par [Radu Matei](https://twitter.com/matei_radu) (CTO chez [FermyonTech](https://twitter.com/fermyontech)).


> Vous pouvez lire aussi celui-ci, écrit par votre serviteur : [WASI, Communication between Node.js and WASM modules with the WASM buffer memory](https://k33g.hashnode.dev/wasi-communication-between-nodejs-and-wasm-modules-with-the-wasm-buffer-memory)

Mais revenons à notre problème. En fait il est très simple : il y a eu 100 connexions essayant simultanément d'accéder à cette mémoire partagée et donc il y a eu "collision", car cette mémoire est faite pour être partagé entre l'application hôte et un seul "invité" à la fois.

Nous devons donc résoudre ce problème pour rendre notre application réellement exploitable.

## Création d'un deuxième serveur HTTP, la solution "naïve"

Ma première approche fut de déplacer le chargement du plugin à partir du manifest et son instanciation à l'intérieur du handler HTTP pour me garantir que pour une requête donnée il n'y aura qu'un seul accès à la mémoire partagée :

```golang
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

func main() {
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Load the wasm plugin 1️⃣
		plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		}

        // Call the wasm function 2️⃣
        // with a string parameter
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

- 1: charger le plugin Wasm.
- 2: appeler la fonction du plugin.

J'ai donc lancé mon nouveau serveur HTTP :

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

Et j'ai refais des tests de charge :

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

Et j'ai obtenu ce rapport :

```bash
Summary:
  Total:        7.6182 secs
  Slowest:      4.6650 secs
  Fastest:      0.0857 secs
  Average:      2.0480 secs
  Requests/sec: 39.3794

Status code distribution:
  [200] 300 responses
```

Donc, c'est magnifique, tout fonctionne ! 🎉 Mais, cependant, le nombre de requêtes par seconde semble vraiment petit. Moins de 40 requêtes par seconde, comparées au 3000 requêtes par seconde du premier test, c'est ridicule 😞. Mais, au moins mon application fonctionne.

Mais n'hésitez jamais à demander de l'aide (c'est pour ça que l'Open Source est un fabuleux modèle).

## Création d'un troisième serveur HTTP, la solution "intelligente"

J'étais quand même ennuyé par les piètres performances de mon MicroService. J'avais fait une application similaire avec Node.js (rappelez-vous de l'article suivant : [Writing Wasm MicroServices with Node.js and Extism](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism#heading-developing-the-application)) et les tests de charge me donnait du 1800 requêtes par seconde.

Et avec la version Node.js de l'application, le plugin wasm n'était instancié [qu'une seule fois](https://k33g.hashnode.dev/writing-wasm-microservices-with-nodejs-and-extism#heading-developing-the-application) et je n'avais aucun problème de collision mémoire 🤔.

Cela aurait du me mettre sur la piste, car en effet, les applications Node.js utilisent un "Single Threaded Event Loop Model" contrairement à Fiber qui utilise une architecture de type "Multi-Threaded Request-Response" pour gérer les accès concurrents. Donc voilà pourquoi mon application Node.js ne "plante" pas.

Ce fut [Steve Manuel](https://twitter.com/nilslice) (CEO de [Dylibso](https://twitter.com/dylibso), mais aussi le créateur d'Extism) à qui j'expliquais mon problème lors d'une discussion sur Discord qui me donna la solution : 

***"So if you want thread-safety in Go HTTP handlers re-using plugins, you need to protect them with a mutex"***

En fait, mais oui, c'était tellement évident (at aussi l'occasion de commencer à étudier ce qu'était un mutex).

J'ai donc suivi les conseils de Steve, et j'ai modifié mon code de la façon suivante :

```golang
package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

// Store all your plugins in a normal Go hash map, 
// protected by a Mutex 1️⃣
var m sync.Mutex
var plugins = make(map[string]extism.Plugin)

// Store the plugin 2️⃣
func StorePlugin(plugin extism.Plugin) {
	plugins["code"] = plugin
}

// Retrieve the plugin 3️⃣
func GetPlugin() (extism.Plugin, error) {
	if plugin, ok := plugins["code"]; ok {
		return plugin, nil
	} else {
		return extism.Plugin{}, errors.New("🔴 no plugin")
	}
}

func main() {
	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	ctx := extism.NewContext()

	defer ctx.Free()

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

    // Create an instance of the plugin 4️⃣
	plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
	if err != nil {
		log.Println("🔴 !!! Error when loading the plugin", err)
		os.Exit(1)
	}
    // Sauvegarder le plugin dans la map 5️⃣
	StorePlugin(plugin)
	
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

        // Lock the mutex 6️⃣
		m.Lock()
		defer m.Unlock()
		
        // Get the plugin 7️⃣
		plugin, err := GetPlugin()

		if err != nil {
			log.Println("🔴 !!! Error when getting the plugin", err)
			c.Status(http.StatusInternalServerError)
			return c.SendString(err.Error())
		}
		
		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		} else {
			c.Status(http.StatusOK)
			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
```

- 1: créer une map protégée par un mutex. Cette map servira à "protéger" le plugin wasm.
- 2: créer une fonction pour sauvegarder le plugin dans la map.
- 3: créer une fonction pour récupérer le plugin à partir de la map.
- 4: créer une instance du plugin wasm.
- 5: sauvegarder cette instance dans la map.
- 6: vérouiller le mutex et utiliser `defer` pour le dévérouiller à la fin de l'exécution.
- 7: obtenir le plugin à partir de map protégée.

Une fois cette modification effectuée, j'ai à nouveau lancé mon serveur HTTP :

```bash
LD_LIBRARY_PATH=/usr/local/lib go run main.go \
path_to_the_plugin/hello.wasm \
hello \
8080
```

Et j'ai lancé à nouveau des tests de charge :

```bash
hey -n 300 -c 100 -m POST \
-d 'John Doe' \
"http://localhost:8080" 
```

Et j'ai obtenu ceci :

```bash
Summary:
  Total:        0.0365 secs
  Slowest:      0.0280 secs
  Fastest:      0.0001 secs
  Average:      0.0092 secs
  Requests/sec: 8207.9604
  
  Total data:   9900 bytes
  Size/request: 33 bytes

Status code distribution:
  [200] 300 responses
```

La nouvelle version du serveur HTTP pouvait aller jusqu'à plus de 8000 requêtes secondes ! 🚀

Pas mal, non ? Ce sera tout pour aujourd'hui. Encore un énorme merci à Steve Manuel pour son aide. J'ai appris énormément, car j'ai osé demander de l'aide. Donc, lorsque que vous bataillez avec quelque chose et que pous ne parvenait pas à trouver une solution, n'hésitez pas à demander autour de vous.

À bientôt pour le prochain article. 👋
