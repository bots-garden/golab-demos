# Ecrire des MicroServices Wasm avec Node.js et Extism

> Où comment écrire une application hôte avec Node.js

Avec l'aide d'Extism, écrire une application hôte (donc une application capable d'exécuter des plugins WebAssembly) est plutôt facile. Nous avons vu dans un [précédent article](https://k33g.hashnode.dev/run-extism-webassembly-plugins-from-a-go-application) comment le faire en Go. Aujourd'hui nous allons le faire avec Node.js. Vous allez voir que c'est très simple, mais cet exemple nous permettra ensuite d'aller plus loin pour découvrir comment écrire des hosts functions.

**Cette application est un serveur HTTP qui servira le plugin WebAssembly comme un MicroService**. Et nous utiliserons [le plugin WebAssembly développé avec TinyGo](https://k33g.hashnode.dev/extism-webassembly-plugins) que nous avons fait dans un article précédent.

## Pré-requis

Vous aurez besoin de 

- Go (v1.20) et TinyGo (v0.28.1) pour compiler les plugins
- Extism 0.4.0 : [Install Extism](https://extism.org/docs/install)
- Node.js (v19.9.0) (c'est la version que j'utilise)

## Création de l'application

### Installation des dépendances

Dans un répertoire, créez un fichier `package.json` avec le contenu suivant :

```json
{
  "dependencies": {
    "@extism/extism": "^0.4.0",
    "fastify": "^4.20.0"
  },
  "type": "module"
}
```

> **[Fastify](https://fastify.dev/)** est un projet Node.js qui permet de développer des serveurs applicatifs web (côté serveur), comme [Express.js](https://expressjs.com/). Mais vous pouvez utiliser ce que vous voulez.

Ensuite, tapez la commande ci-dessous pour installer les dépendances nécessaires :

```bash
npm install
```

### Développer l'application

Créer un fichier `server.js` avec le contenu suivant :

```javascript
import Fastify from 'fastify'
import process from "node:process"

import { Context } from '@extism/extism'
import { readFileSync } from 'fs'

let wasmFile = "../01-simple-go-plugin/simple.wasm"
let functionName = "say_hello"
let httpPort = 7070

let wasm = readFileSync(wasmFile) // 1️⃣

const fastify = Fastify({
  logger: true
})

const opts = {}


// 2️⃣
let ctx = new Context()
let plugin = ctx.plugin(wasm, true, [])

// Create and start the HTTP server
const start = async () => {

  fastify.post('/', opts, async (request, reply) => { // 3️⃣

    // 4️⃣
    let buf = await plugin.call(functionName, request.body); 
    let result = buf.toString()

    return result
  })

  try { // 5️⃣
    await fastify.listen({ port: httpPort, host: '0.0.0.0'})
  } catch (err) {
    fastify.log.error(err)
    process.exit(1)
  }
}
start().then(r => console.log("😄 started"))
```

- 1: charger le fichier du plugin WebAssembly.
- 2: créer un contexte Extism et l'utiliser pour initialiser le plugin Wasm.
- 3: définir une route (endpoint) du serveur HTTP. Le code sera exécuté à chaque appel HTTP de type POST de `http://localhost:7070`.
- 4: appeler la fonction du module avec comme paramètres le nom de la fonction (`say_hello`) et les données postées par la requête HTTP, et retourner le résultat.
- 5: démarrer le serveur HTTP.

### Démarrer le serveur HTTP

Utiliser tout simplement cette commande:

```bash
node server.js
```

### Appeler le MicroService

Pour appeler le MicroService, utilisez cette simple commande `curl` :

```bash
curl -X POST http://localhost:7070 \
-H 'Content-Type: text/plain; charset=utf-8' \
-d 'Jane Doe'
```

Et vous obtiendrez :

```bash
👋 Hello Jane Doe
```

😍 Vous voyez qu'avec le système de plugins proposé par Extism, il devient très facile d'écrire des MicroServices polyglottes et de les proposer à l'aide de Node.js. Vous n'êtes pas loin d'avoir les bases pour écrire un FaaS (mais ce sera une autre histoire, un peu plus tard probablement 😉). Je vous laisse déjà expérimenter avec ce que nous venons de voir aujourd'hui.

L'article suivant ré-utilisera cet exemple, et je rappelerais le concept de **host function** et comment en écrire pour apporter des fonctionnalités supplémentaires aux plugins WebAssembly.
