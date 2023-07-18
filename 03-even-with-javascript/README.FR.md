# Plugin WebAssembly en JavaScript avec Extism

Ces derniers jours nous avons utilisé le PDK (Plugin Development Kit) Go pour développer des applications WebAssembly et les exécuter avec la CLI Extism.

Comme je le disais dans le premier article, il existe des PDKs pour plusieurs langages, et notamment pour [JavaScript](https://github.com/extism/js-pdk) (et je suis un gros fan de JavaScript). Mais comment est-ce possible ? En effet on ne peut pas compiler du JavaScript en code natif Wasm.

En fait, ce PDK utilise (entre autre chose) le projet [QuickJS](https://bellard.org/quickjs/) pour exécuter du code JavaScript dans un programme Wasm.

Certes, cela ne s'exécutera pas aussi rapidement qu'un programme Wasm compilé avec TinyGo ou Rust, mais cela permet d'exécuter des fonctions JavaScript dans un environnement complètement "sandboxé". Comme le fait Shopify (ce PDK est un fork du projet [Javy](https://github.com/bytecodealliance/javy) initialisé par [Shopify](https://www.shopify.com/).

> Shopify a développé le projet Javy pour apporter le support de JavaScript aux Shopify Functions. Les Shopify Functions permettent aux développeurs de créer des extensions et des fonctionnalités sur mesure pour les besoins spécifiques des marchands en JavaScript qui est un langage populaire et familier pour beaucoup de développeurs web.

Donc, imaginez que vous decidiez de créer une plateforme FaaS orientée JavaScript et que vous vouliez donner à vos utilisateurs la possibilité de créer et publier leurs propres fonctions, passer par un mécanisme similaire aura au minimum deux avantages :

- Favoriser l'adoption (JavaScritp est bien connu)
- Garantir l'intégrité de votre plateforme (les fonctions sont exécutées dans un environnement sandboxé)

> Un peu de lecture : [Bringing Javascript to WebAssembly for Shopify Functions](https://shopify.engineering/javascript-in-webassembly-for-shopify-functions)

Mais revenons à nos moutons et laissez moi vous expliquer comment créer un plugin Extism en JavaScript.

## Pré-requis

Il vous faudra Extism 0.4.0 et Extims-js PDK pour builder des modules Wasm avec du JavaScript
  - [Install Extism](https://extism.org/docs/install)
  - [Install Extism-js PDK](https://extism.org/docs/write-a-plugin/js-pdk#how-to-install-and-use-the-extism-js-pdk)

> **Remarque**, si vous rencontriez un problème pour l'installation du PDK, vous pouvez le faire manuellement de cette façon (modifiez selon votre environnement) :
> ```bash
> export TAG="v0.5.0"
> export ARCH="aarch64"
> export  OS="linux"
> curl -L -O "https://github.com/extism/js-pdk/releases/download/$TAG/extism-js-$ARCH-$OS-$TAG.gz"
> gunzip extism-js*.gz
> sudo mv extism-js-* /usr/local/bin/extism-js
> chmod +x /usr/local/bin/extism-js
> ```

## Le plus simple des plugins Extism

Dans un répertoire, créez un fichier `index.js` avec le contenu ci-dessous :

```javascript
function say_hello() {

	// read function argument from the memory
	let input = Host.inputString()

	let output = "param: " + input

	console.log("👋 Hey, I'm a JS function into a wasm module 💜")

	// copy output to host memory
	Host.outputString(output)

	return 0
}

module.exports = {say_hello}
```

Vous pouvez voir que le code est très simple et possède une logique complètement similaire à ce que nous avons vu dans les précédents articles.


### Compiler le plugin Wasm

Pour compiler le programme, utiliser la commande ci-dessous, qui produira un fichier `hello-js.wasm` :

```bash
extism-js index.js -o hello-js.wasm
```

Et maintenant nous allons exécuter notre plugin Wasm comme nous l'avaons fait dans les exemples précédents.

### Exécuter la fonction `say_hello` du plugin Wasm

Pour cela nous allons utiliser la CLI Extism.

Pour exécuter la fonction `say_hello` avec comme paramètre la chaine de caractères `"😀 Hello World 🌍! (from JavaScript)"`, utilisez la commande suivante :

```bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello World 🌍! (from JavaScript)" \
  --wasi \
  --log-level info
```

Vous obtiendrez :

```bash
extism_runtime::pdk INFO 2023-07-18T07:08:34.347325607+02:00 - 👋 Hey, I'm a JS function into a wasm module 💜
param: 😀 Hello World 🌍! (from JavaScript)
```

Si vous essayez de lancer le plugin sans préciser le niveau de log, comme cela :

```bash
extism call ./hello-js.wasm \
  say_hello --input "😀 Hello World 🌍! (from JavaScript)" \
  --wasi
```

vous n'aurez que ceci :

```bash
param: 😀 Hello World 🌍! (from JavaScript)
```

> Dans le cas du PDK JavaScript `console.log()` est une sorte d'alias pour appeler la host function de log de la CLI Extism.


Voilà, c'est tout pour aujourd'hui. Je vous laisse vous familiariser avec le développement des plugins Extism. Si tout va bien, demain nous attaquons l'écriture d'une application hôte en Go.


