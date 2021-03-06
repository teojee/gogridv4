package main

import (
	"sort"
	"strconv"
	"syscall/js"
)

func dgSortFunc(pGridID string, pProperties map[string]string) error {
	var lGrid DivGrid

	lGrid = mGlobal.GetGrid(pGridID)

	lGrid.setProperties(pProperties)

	errSort := lGrid.Sort()
	if errSort != nil {
		return errSort
	}

	lGrid.OutputTotal()

	mGlobal.SetGrid(pGridID, lGrid)

	return nil
}

func (d *DivGrid) Sort() error {
	s3("sort.go", "Amount Rows", strconv.Itoa(len(d.Rows)))
	s3("sort.go", "Amount FilterRows", strconv.Itoa(len(d.FilterRows)))
	s3("sort.go", "SortCol", strconv.Itoa(d.SortCol))

	// Sort via SliceStable
	// HOWTO this with pointers ?
	f := func(r []DivGridRow) []DivGridRow {
		if len(r) > 1 {
			//sort.SliceStable(d.Rows[:500], func(i, j int) bool {
			sort.SliceStable(r, func(i, j int) bool {
				if len((r)[i].Cells) > d.SortCol-1 {
					lValue := (r)[i].Cells[d.SortCol]

					if len((r)[j].Cells) > d.SortCol-1 {
						rValue := (r)[j].Cells[d.SortCol]

						switch d.SortColType {
						case "int":
							lInt, err1 := strconv.Atoi(lValue.Text)
							rInt, err2 := strconv.Atoi(rValue.Text)

							// s3("SORT INT", lValue.Text, rValue.Text)
							if err1 == nil && err2 == nil {
								// s1("SORT INT HIT")
								return lInt < rInt
							}

							// When type is int and value isn't a int then rank is lower drown down
							return lValue.Text < rValue.Text
						default:
							return lValue.Text < rValue.Text
						}

					}
				}

				return true

			})
		}

		return r
	}

	// When Filter is used then use FilterRows which are set in func Filter
	// Otherwise use total rows
	if len(d.FilterText) > 0 {
		d.FilterRows = f(d.FilterRows)
	} else {
		d.Rows = f(d.Rows)
	}

	return nil
}

func dgSortHandler(this js.Value, args []js.Value) interface{} {
	s3("sort.go", "Length args", strconv.Itoa(len(args)))

	if len(args) == 0 {
		s2("sort.go", "Please use reference GridDiv ID")
		return nil
	}

	argGridID := args[0].String()
	s3("sort.go", "GridID", argGridID)

	arrArguments := getArgumentsJS(args)

	// SortHandler needs to be called with two parameters
	// par1 is id of divgrid
	// par2 is index of col which has to be sorted
	// Handler for the Promise: this is a JS function
	// It receives two arguments, which are JS functions themselves: resolve and reject
	hFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		argResolve := args[0]
		argReject := args[1]

		err := dgSortFunc(argGridID, arrArguments)
		if err != nil {
			argReject.Invoke("Sort Invoke")
		}

		argResolve.Invoke("Error Sort")

		return nil
	})

	// Create and return the Promise object
	objProm := js.Global().Get("Promise")
	return objProm.New(hFunc)
}
