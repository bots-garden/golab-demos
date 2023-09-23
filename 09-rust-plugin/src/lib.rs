#![no_main]

use extism_pdk::*;

extern "C" {
    fn hostRobotMessage(ptr: u64) -> u64;
}

pub fn robot_message(text: String) {
    let mut memory_text: Memory = extism_pdk::Memory::new(text.len());
    memory_text.store(text);
    unsafe { hostRobotMessage(memory_text.offset) };
}


#[plugin_fn]
pub fn hello(input: String) -> FnResult<u64> {

    let msg: String = "🦀 Hello ".to_string() + &input;

    robot_message(msg);
    
    Ok(0)
}
