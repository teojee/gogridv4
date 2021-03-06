package main

import (
	"strconv"
	"strings"
	"syscall/js"
)

//#region DivGrid Handler
func dgFilterFunc(pGridID string, pProperties map[string]string) error {
	var lGrid DivGrid

	lGrid = mGlobal.GetGrid(pGridID)
	lGrid.setProperties(pProperties)

	errFilter := lGrid.Filter()
	if errFilter != nil {
		return errFilter
	}

	/*
		errSort := lGrid.Sort()
		if errSort != nil {
			return errSort
		}
	*/

	lGrid.OutputTotal()

	mGlobal.SetGrid(pGridID, lGrid)

	return nil
}

func (dg *DivGrid) Filter() error {
	s3("filter.go", "Filter ID", dg.FilterID)

	if len(dg.FilterID) == 0 {
		return nil
	}

	objFilter, err := mGlobal.Get(dg.FilterID, "value")

	if err != nil {
		return err
	}

	dg.FilterText = objFilter.String()
	s3("filter.go", "Filter Content", dg.FilterText)

	// s3("filter.go", "Len Filter", strconv.Itoa(len(dg.FilterText)))

	// Init an array of filter rows which will be used in output
	dg.FilterRows = make([]DivGridRow, 0)

	// Fill the Filterrows by traversing rows and  cells,  if match is found continue loop
	if len(dg.FilterText) > 0 {
		for _, row := range dg.Rows {
			for _, cell := range row.Cells {
				if strings.Contains(strings.ToLower(cell.Text), strings.ToLower(dg.FilterText)) {
					dg.FilterRows = append(dg.FilterRows, row)
					continue
				}
			}
		}
	}

	return nil
}

func dgFilterHandler(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		// handleMessage("Please int Div")
		s2("filter.go", "Please use references GridDiv ID")
		return nil
	}

	s3("filter.go", "Length args", strconv.Itoa(len(args)))

	argGridID := args[0].String()
	s3("filter.go", "Grid ID", argGridID)

	arrArguments := getArgumentsJS(args)

	hFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		argResolve := args[0]
		argReject := args[1]

		err := dgFilterFunc(argGridID, arrArguments)
		if err != nil {
			argReject.Invoke("Filter Error")
		}

		argResolve.Invoke("Filter OK")

		return nil
	})

	// Create and return the Promise object
	objProm := js.Global().Get("Promise")
	return objProm.New(hFunc)
}
