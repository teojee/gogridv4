package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"syscall/js"
)

// Write Header back to html DIV
func (d *DivGrid) OutputHeader() {
	fr, _ := xml.MarshalIndent(d.Header, "", "    ")
	mesHelp(string(fr))
	mGlobal.Set("header", "outerHTML", string(fr))
}

// Write Header and rows back to html DIV
func (d *DivGrid) OutputTotal() {
	totalRows := make([]DivGridRow, 0)
	totalRows = append(totalRows, d.Header)
	// totalRows = append(totalRows, d.Rows...)

	f := func(r []DivGridRow) {
		lFilterLen := len(r)
		if lFilterLen > d.PageSize {
			lFilterLen = d.PageSize
		}
		totalRows = append(totalRows, r[:lFilterLen]...)
	}

	if len(d.FilterText) > 0 {
		f(d.FilterRows)
	} else {
		f(d.Rows)
	}

	t, _ := xml.MarshalIndent(totalRows, "", "    ")

	// TODO make class variable
	htmlTotal := fmt.Sprintf("<div id='%s' class='tablemaster'>%s</div>", d.ID, string(t))
	mGlobal.Set(d.ID, "outerHTML", htmlTotal)

	// WARNING innerHTML doesn't work when html div files are large
	// mGlobal.Set(d.ID, "innerHTML", string(t))

	// fr, _ := xml.MarshalIndent(d.Filter[:5], "", "    ")
	// mesHelp(string(fr))
	//	htmlTotalFilter := fmt.Sprintf("<div id='divHtml' class='tableMaster'>%s</div>", string(fr))
	//	mGlobal.Set("divHtml", "outerHTML", htmlTotalFilter)

}

// Get by ID
func (j jsGlobal) ByID(elem string) (js.Value, error) {
	jsDoc := j.Doc.Call("getElementById", elem)

	if !jsDoc.Truthy() {
		return js.Null(), errors.New(s2("Elem not found:", elem))
	}

	return jsDoc, nil
}

// Get Object
// Example Filter input  FilterText:= j.Get("filter", "value") -> document.getElementById("Filter").value
func (j jsGlobal) Get(id string, prop string) (js.Value, error) {
	jsElem, err := j.ByID(id)
	if err != nil {
		return js.Null(), err
	}

	jsProp := jsElem.Get(prop)
	// mesHelp(jsProp)

	//TODO
	/*
		if !jsProp.Truthy() {
			return js.Null(), errors.New(s3("Prop not found:", elem, prop))
		}
	*/

	return jsProp, nil
}

// Set Value
// TODO Handling Panic
func (j *jsGlobal) Set(elem string, key string, value string) error {
	jsElem, err := j.ByID(elem)
	if err != nil {
		return err
	}

	jsElem.Set(key, value)

	return nil
}

// Url: https://ian-says.com/articles/golang-in-the-browser-with-web-assembly/
/*
func AddClass(elem string, class string) {
    getElementValue(elem, "classList").Call("add", class)
}

func RemoveClass(elem string, class string) {
    classList := getElementValue(elem, "classList")
    if (classList.Call("contains", class).Bool()) {
        classList.Call("remove", class)
    }
}
*/

// func (d *DivGrid) InitArray(xml string) error {
// 	err := xml.Unmarshal([]byte(xml), d.XMLBytes)
// 	return handelError(err)
// }

/*

// Get by ID
func (d DivGrid) ByID(elem string) js.Value {
	return mDoc.Call("getElementById", elem)

}

// Get Object
func (d DivGrid) Get(elem string, value string) js.Value {
	return d.ByID(elem).Get(value)
}

func (d *DivGrid) Set(elem string, key string, value string) {
	d.ByID(elem).Set(key, value)
}

func (d DivGrid) String(elem string, value string) string {
	// return "test"
	return d.Get(elem, value).String()
}



func getFirstChild(elem string) js.Value {
	// return mDoc.Call("firstElementChild", elem)
	return getElementById(elem).Get("firstElementChild")
}

func getElementValue(elem string, value string) js.Value {
	return getElementById(elem).Get(value)
}

func Hide(elem string) {
	getElementValue(elem, "style").Call("setProperty", "display", "none")
}

func Show(elem string) {
	getElementValue(elem, "style").Call("setProperty", "display", "block")
}

func GetString(elem string, value string) string {
	return getElementValue(elem, value).String()
}

func SetValue(elem string, key string, value string) {
	getElementById(elem).Set(key, value)
}

func AddClass(elem string, class string) {
	getElementValue(elem, "classList").Call("add", class)
}

func RemoveClass(elem string, class string) {
	classList := getElementValue(elem, "classList")
	if classList.Call("contains", class).Bool() {
		classList.Call("remove", class)
	}
}
*/
