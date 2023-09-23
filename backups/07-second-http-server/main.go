package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
	"github.com/tetratelabs/wazero"
)

func main() {

	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	ctx := context.Background()

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
		
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
		AllowedHosts:  []string{"*"}, 
	}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

		pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil) // new
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		}

		_, out, err := pluginInst.Call(wasmFunctionName, params)

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
