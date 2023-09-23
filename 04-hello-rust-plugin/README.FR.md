# Créer un plugin Webassembly avec Extism et Rust

Avant de passer à un article plus compliqué (créer un serveur HTTP en Go pour servir des services WASM), nous allons voir comment développer un plug-in Extism en Rust. Vous allez voir, c'est très simple et c'est une bonne opportunité pour faire vos premiers pas en Rust.

## Pré-requis

- Extism 0.4.0 : [Installer Extism](https://extism.org/docs/install)
- Rust : 
  - [Installer Rust](https://www.rust-lang.org/tools/install)
    ```bash
    curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh -s -- -y

    echo 'export CARGO_HOME="~/.cargo"' >> ~/.bashrc
    echo 'export PATH=\$CARGO_HOME/bin:\$PATH' >> ~/.bashrc

    source ~/.cargo/env
    source ~/.bashrc
    ```
    > Je suis sous Linux
  - [Installer Wasm Pack](https://rustwasm.github.io/wasm-pack/installer/)
    ```bash
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
    ```
  - Installer les targets wasm nécessaires :
    ```bash
    rustup target add wasm32-wasi
    rustup target add wasm32-unknown-unknown
    ```
  - Installer Wasm Bindgen :
    ```bash
    cargo install -f wasm-bindgen-cli
    ```
  - Installer les composants de la toolchain de build en fonction de votre architecture :
    ```bash
    rustup component add rust-analysis --toolchain stable-aarch64-unknown-linux-gnu 
    rustup component add rust-src --toolchain stable-aarch64-unknown-linux-gnu 
    rustup component add rls --toolchain stable-aarch64-unknown-linux-gnu
    ```
    > Dans mon cas, j'utilise une architecture ARM

## Génération du projet de plugin

Pour générer le projet Rust, utilisez la commande ci-dessous :

```bash
cargo new --lib hello-rust-plugin --name hello
```

Cela va créer un dossier `hello-rust-plugin`. Dans ce dossier, ajouter ceci au fichier `Cargo.toml` :

```toml
[lib]
crate_type = ["cdylib"]
```

Donc, le fichier `Cargo.toml` devrait présenter le contenu suivant :

```toml
[package]
name = "hello"
version = "0.1.0"
edition = "2021"

[lib]
crate_type = ["cdylib"]

[dependencies]
```

Ensuite, ajoutez les dépendances :

```bash
cd hello-rust-plugin
cargo add extism-pdk
cargo add serde
cargo add serde_json
```

La section `[dependencies]` de `Cargo.toml` devrait ressembler à ceci :

```toml
[dependencies]
extism-pdk = "0.3.3"
serde = "1.0.178"
serde_json = "1.0.104"
```

## Modification du code source

Allez ensuite modifier le code source de `src/lib.rs` de la façon ci-dessous :

```rust
#![no_main]

use extism_pdk::*;
use serde::Serialize;

#[derive(Serialize)]
struct Output {
    pub message: String,
}

#[plugin_fn]
pub fn hello(input: String) -> FnResult<Json<Output>> {

    let msg: String = "🦀 Hello ".to_string() + &input;

    let output = Output { message: msg };
    
    Ok(Json(output))
}
```

👋 Vous pouvez remarquer que le code Rust d'un plugin Extism est en fait plus simple que celui d'un plugin [Extism en Go](https://k33g.hashnode.dev/extism-webassembly-plugins) en ce qui concerne le passage de paramètres et le retour de valeur de fonction.

## Compilation du plugin

Pour compiler le plugin, tapez les commandes ci-dessous :

```bash
cargo clean
cargo build --release --target wasm32-wasi
```

Le plugin `hello.wasm` a été généré dans le répertoire suivant : `./target/wasm32-wasi/release`

Pour le tester, il vous suffit d'utiliser la commande suivante :

```bash
extism call ./target/wasm32-wasi/release/hello.wasm \
  hello --input "Bob Morane"  \
  --wasi
```

Et vous obtiendrez :

```bash
{"message":"🦀 Hello Bob Morane"}
```

Voilà, c'est terminé. Vous avez pu voir que commencer le développement d'un plugin Extism avec Rust est relativement aisé. Dans le prochain article, qui sera beaucoup plus long, nous utiliserons ce plugin pour créer un MicroService à l'aide d'une application hôte développée en Go.
