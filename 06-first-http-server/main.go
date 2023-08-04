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

	//ctx := extism.NewContext()
	ctx := context.Background()

	//defer ctx.Free() // this will free the context and all associated plugins

	config := extism.PluginConfig{
		ModuleConfig: wazero.NewModuleConfig().WithSysWalltime(),
		EnableWasi:   true,
	}

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

	//plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
	pluginInst, err := extism.NewPlugin(ctx, manifest, config, nil) // new

	if err != nil {
		panic(err)
	}

	/*
		app := fiber.New(fiber.Config{
			DisableStartupMessage: true,
			DisableKeepalive:      true,
			Concurrency:           100000,
		})
	*/
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Post("/", func(c *fiber.Ctx) error {

		params := c.Body()

		_, out, err := pluginInst.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
			//os.Exit(1)
		} else {
			c.Status(http.StatusOK)

			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
