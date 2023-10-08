#!/bin/bash
snippets generate \
  --input vanilla-wasm.yml \
  --output ../.vscode/vanilla-wasm.code-snippets 

snippets generate \
  --input wazero-wasm.yml \
  --output ../.vscode/wazero-wasm.code-snippets 
