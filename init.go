package main

import (
	"encoding/xml"
	"fmt"
	"syscall/js"
)

func divgridInitTest() {
	lDebugGridDivID := "divDebug"

	lDebugGridDiv := DivGrid{
		ID:          lDebugGridDivID,
		SortCol:     0,
		SortColType: "string",
		PageSize:    10,
		FilterText:  "",
		Debug:       true,
		DebugFile:   "divgrid.xml",
	}

	mGlobal.SetGrid(lDebugGridDivID, lDebugGridDiv)
	dgInitFunc(lDebugGridDivID, nil)
}

// divgridInit create a new DivGrid struct in case it not already exist
// Function is the main part of divgridInitHandler which is called via Javascript
//
func dgInitFunc(pDivID string, pProperties map[string]string) error {
	var lGridDiv DivGrid
	if mGlobal.ExistGrid(pDivID) {
		lGridDiv = mGlobal.GetGrid(pDivID)
		s2("init.go", "Grid Exist")
	} else {
		s2("init.go", "Grid Doesn't Exist")

		lGridDiv = DivGrid{
			ID:          pDivID,
			SortCol:     0,
			SortColType: "string",
			PageSize:    10,
			FilterID:    "",
			FilterText:  "",
			Debug:       false,
			DebugFile:   "divgrid.xml",
		}

		err := lGridDiv.Init()
		if err != nil {
			return err
		}
	}

	lGridDiv.setProperties(pProperties)
	lGridDiv.Filter()
	lGridDiv.Sort()
	lGridDiv.OutputTotal()

	// Store the Grid
	mGlobal.SetGrid(pDivID, lGridDiv)

	return nil
}

// Init TODO
func (dg *DivGrid) Init() error {
	var lGridHTML string
	var lGridHTMLBytes []byte
	var lGridTotal DivGridParent

	// Read the Rows html code of actual grid it has to be in xml format
	if !dg.Debug {
		objOuterHTML, err := mGlobal.Get(dg.ID, "outerHTML")
		if err != nil {
			return err
		}

		lGridHTML = objOuterHTML.String()
	} else {
		lGridHTMLBytes = readxml(dg.DebugFile)
	}

	// Convert to Rows Struct so it can be handled
	lGridTotal = DivGridParent{}
	lGridHTMLBytes = []byte(lGridHTML)

	// xml.Unmarshal([]byte(divSortHTML), &d.Rows)
	err := xml.Unmarshal(lGridHTMLBytes, &lGridTotal)
	if err != nil {
		return err
	}

	/*
		mesHelp("Size:")
		mesHelp(unsafe.Sizeof(lDivTotal))
		mesHelp("Length:")
		mesHelp(len(lDivTotal.Div))
	*/

	// Get the first header row  and store the header so it can be appended after sorting
	// NICE mabye there is no headerrow make setting
	dg.Header = lGridTotal.Rows[0]
	// mesHelp(d.Header)

	// Loop through the header cells and fill the onclick event
	// Result Gv4.dgSort('divSort', 1, 'int')
	// Par 1 'divSort' is ID of grid to be sorted
	// Par 2 1 is Column to be sorted
	// Par 3 is Sorttype so way of sorting

	for index, hCell := range dg.Header.Cells {
		dg.Header.Cells[index].Onclick = fmt.Sprintf("%s.%s('%s', 'sc=%d', 'st=%s')", Gv4, dgSort, dg.ID, index, hCell.Type)
	}

	// d.Header.Div[0].Text = "Changed"
	// d.Header.Div[0].Onclick = fmt.Sprintf("Gv4.dgSort('%s', 0, '%s')", d.ID, d.Header.Div[0].Type)
	// d.Header.Div[1].Onclick = fmt.Sprintf("Gv4.dgSort('%s', 1, '%s')", d.ID, d.Header.Div[1].Type)
	// d.Header.Div[2].Onclick = fmt.Sprintf("Gv4.dgdSort('%s', 2, '%s' )", d.ID, d.Header.Div[2].Type)
	// mesHelp(d.Header)

	// d.Header.Div[2].Onclick = fmt.Sprintf("Gv4.dgSort('%s', %d)", d.ID, d.SortCol)

	// Set all the rows except the first header row  and store can be sorted
	// [1:]	 all Divs after the first (header) one
	// d.Rows = lDivTotal.Div[1:]
	// WARNING MAX is 4000

	// WARNING wasm results in errors when rows exceed 4000-5000 3 test col
	// TODO make  warning max of 2000 rows
	lRowLength := len(lGridTotal.Rows)
	if lRowLength > 2000 {
		lRowLength = 2000
	}

	// Store all rows (except header) in struct so it can be reused
	dg.Rows = lGridTotal.Rows[1:lRowLength]

	// mesHelp(d.Rows)

	return nil
}

// Init receives Div ID from javascript and initialize the DIV
func dgInitHandler(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		s2("init.go", "Please use reference GridDiv ID")
		return nil
	}

	parDivID := args[0].String()
	s2("init.go", parDivID)

	arrSettings := getArgumentsJS(args)

	// Promise handler
	hFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		argResolve := args[0]
		argReject := args[1]

		err := dgInitFunc(parDivID, arrSettings)

		if err != nil {
			argReject.Invoke("Init Error")
		}

		argResolve.Invoke("Init OK")

		return nil
	})

	// Handle as a new promise
	objProm := js.Global().Get("Promise")
	return objProm.New(hFunc)
}
