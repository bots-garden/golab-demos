package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/extism/extism"
	"github.com/gofiber/fiber/v2"
)

func main() {

	wasmFilePath := os.Args[1:][0]
	wasmFunctionName := os.Args[1:][1]
	httpPort := os.Args[1:][2]

	//var counter = 0

	ctx := extism.NewContext()

	defer ctx.Free() // this will free the context and all associated plugins

	manifest := extism.Manifest{
		Wasm: []extism.Wasm{
			extism.WasmFile{
				Path: wasmFilePath},
		},
	}

	/*
		plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
		if err != nil {
			panic(err)
		}
	*/

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

		plugin, err := ctx.PluginFromManifest(manifest, []extism.Function{}, true)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
		}

		out, err := plugin.Call(wasmFunctionName, params)

		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusConflict)
			return c.SendString(err.Error())
			//os.Exit(1)
		} else {
			c.Status(http.StatusOK)
			//fmt.Println(counter, string(out))
			//counter ++
			return c.SendString(string(out))
		}

	})

	fmt.Println("🌍 http server is listening on:", httpPort)
	app.Listen(":" + httpPort)
}
