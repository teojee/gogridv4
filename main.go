package main

import (
	"syscall/js"
)

var mGlobal jsGlobal

const Gv4 = "Gv4"
const dgHello = "dgHello"
const dgInit = "dgInit"
const dgSort = "dgSort"
const dgFilter = "dgFilter"

func registerMain() {
	// Register Callback. Make sure Gv4 is Set first
	// TODO how to make mutiple libraries ? Only one object can be efined why ? The O 2021.02.24
	// TeeJeeObject := make(map[string]interface{})
	// js.Global().Set("TeeJeeGo", js.ValueOf(TeeJeeObject))

	objModuleGv4 := make(map[string]interface{})
	js.Global().Set(Gv4, js.ValueOf(objModuleGv4))
	js.Global().Get(Gv4).Set(dgHello, js.FuncOf(testHelloFn))
}

func registerDivGrid() {
	// Make sure Gv4 is set  in main.go -> js.Global().Set("Gv4", js.ValueOf(objModule))
	// To keep code organized use seperate files  divgrid.go etc.

	// js.Global().Get(Gv4).Set(dgHello, js.FuncOf(testHelloFn))
	js.Global().Get(Gv4).Set(dgInit, js.FuncOf(dgInitHandler))
	js.Global().Get(Gv4).Set(dgSort, js.FuncOf(dgSortHandler))
	js.Global().Get(Gv4).Set(dgFilter, js.FuncOf(dgFilterHandler))
}

func main() {
	// Create a channel -> long-running pogram
	c := make(chan struct{}, 0)

	// Create reference to webpage via sysall js wasm model
	// Will be different in case multiple webpages are used
	// Reused in different get en set struct functions
	// Create empty hashGrid f type XmlDivGrid // The Grid will be ued in divgrid.go
	mGlobal = jsGlobal{
		Doc:        js.Global().Get("document"),
		DivGridMap: make(map[string]DivGrid),
	}

	// Info loading WASM file
	s2("main.go", "Init DivGridv4 is OK")

	// Register functions
	registerMain()
	registerDivGrid()

	// divgridInitTest()

	<-c
}
