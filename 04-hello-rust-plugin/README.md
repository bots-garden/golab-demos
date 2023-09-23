

ref: https://extism.org/docs/write-a-plugin/rust-pdk

```bash
cargo new --lib 05-hello-rust-plugin --name hello-rust-plugin

rustup target add wasm32-wasi # if needed
```

In the generated Cargo.toml, be sure to include:

```toml
[lib]
crate_type = ["cdylib"]
```
> ref: https://doc.rust-lang.org/reference/linkage.html


```bash
cd 05-hello-rust-plugin
cargo add extism-pdk
cargo add serde
cargo add serde_json
```

## Build 

```bash
# build
cargo clean
cargo build --release --target wasm32-wasi #--offline
# ls -lh *.wasm
ls -lh ./target/wasm32-wasi/release/*.wasm
```

