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
