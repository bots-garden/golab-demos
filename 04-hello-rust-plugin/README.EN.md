# Create a Webassembly plugin with Extism and Rust

Before moving on to a more complicated article (creating an HTTP server in Go to serve WASM services), we will see how to develop an Extism plugin in Rust. You will see, it is very simple and it is a good opportunity to take your first steps in Rust.

## Prerequisites

- Extism 0.4.0: [Install Extism](https://extism.org/docs/install)
- Rust:
  - [Install Rust](https://www.rust-lang.org/tools/install)
    ```bash
    curl --proto '=https' --tlsv1.2 https://sh.rustup.rs -sSf | sh -s -- -y

    echo 'export CARGO_HOME="~/.cargo"' >> ~/.bashrc
    echo 'export PATH=\$CARGO_HOME/bin:\$PATH' >> ~/.bashrc

    source ~/.cargo/env
    source ~/.bashrc
    ```
    > I'm on Linux
  - [Install Wasm Pack](https://rustwasm.github.io/wasm-pack/installer/)
    ```bash
    curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh
    ```
  - Install the necessary wasm targets:
    ```bash
    rustup target add wasm32-wasi
    rustup target add wasm32-unknown-unknown
    ```
  - Install Wasm Bindgen:
    ```bash
    cargo install -f wasm-bindgen-cli
    ```
  - Install the components of the build toolchain according to your architecture:
    ```bash
    rustup component add rust-analysis --toolchain stable-aarch64-unknown-linux-gnu 
    rustup component add rust-src --toolchain stable-aarch64-unknown-linux-gnu 
    rustup component add rls --toolchain stable-aarch64-unknown-linux-gnu
    ```
    > In my case, I use an ARM architecture

## Generating the plugin project

```bash
cargo new --lib hello-rust-plugin --name hello
```

This will create a folder `hello-rust-plugin`. In this folder, add this to the file `Cargo.toml`:

```toml
[lib]
crate_type = ["cdylib"]
```

So, the file `Cargo.toml` should have the following content:

```toml
[package]
name = "hello"
version = "0.1.0"
edition = "2021"

[lib]
crate_type = ["cdylib"]

[dependencies]
```

Then, add the dependencies:

```bash
cd hello-rust-plugin
cargo add extism-pdk
cargo add serde
cargo add serde_json
```

The `[dependencies]` section of `Cargo.toml` should look like this:

```toml
[dependencies]
extism-pdk = "0.3.3"
serde = "1.0.178"
serde_json = "1.0.104"
```

## Modifying the source code

Then modify the source code of `src/lib.rs` as follows:

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

👋 You can notice that the Rust code of an Extism plugin is actually simpler than that of an [Extism plugin in Go](https://k33g.hashnode.dev/extism-webassembly-plugins) regarding the passing of parameters and the return value of function.

## Compiling the plugin

To compile the plugin, type the commands below:

```bash
cargo clean
cargo build --release --target wasm32-wasi
```

The plugin `hello.wasm` has been generated in the following directory: `./target/wasm32-wasi/release`

To test it, you just need to use the following command:

```bash
extism call ./target/wasm32-wasi/release/hello.wasm \
  hello --input "Bob Morane"  \
  --wasi
```

And you will get:

```bash
{"message":"🦀 Hello Bob Morane"}
```

That's it, it's over. You have seen that starting the development of an Extism plugin with Rust is relatively easy. In the next article, which will be much longer, we will use this plugin to create a MicroService using a host application developed in Go.
