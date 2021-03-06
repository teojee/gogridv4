package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall/js"
)

func readxml(pFile string) []byte {
	// Open our xmlFile
	xmlFile, err := os.Open(pFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	s2("Successfully Opened", pFile)
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	return byteValue
}

func testHelloFn(this js.Value, args []js.Value) interface{} {
	par1 := args[0].String()

	// Handler for the Promise: this is a JS function
	// It receives two arguments, which are JS functions themselves: resolve and reject
	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// par1 := js.ValueOf(args[0])
		resolve := args[0]
		// Commented out because this Promise never fails
		//reject := args[1]
		mesHelp(fmt.Sprintf("input %s", par1))
		resolve.Invoke("Resolve Invoke Hello")

		// The handler of a Promise doesn't return any value
		return nil
	})

	// Create and return the Promise object
	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}
