package main

import (
	"encoding/xml"
	"strconv"
	"strings"
	"syscall/js"
)

// Use Struct so we cn use define functions
type jsGlobal struct {
	Doc        js.Value
	DivGridMap map[string]DivGrid // Dictionary of ivGrids so multiple Grids can be handled in one htmlpage
}

// Store DivGrid
func (j jsGlobal) GetGrid(DivID string) DivGrid {
	return (j.DivGridMap[DivID])
}

// Store DivGrid
func (j *jsGlobal) SetGrid(DivID string, value DivGrid) {
	j.DivGridMap[DivID] = value
}

// Exist the Grid
func (j jsGlobal) ExistGrid(DivID string) bool {
	if _, ok := j.DivGridMap[DivID]; ok {
		return true
	}
	return false
}

// Create Dictionary of Settings from Javascript
// Example (divTest, ps=10, fi=filter, sc=string, etc)
// Two step approach to split syscall function and go functions
// build constraints exclude all Go files in /usr/local/go/src/syscall/js
func getArgumentsJS(args []js.Value) map[string]string {
	arrSetting := make(map[string]string)

	for i := range args {
		value := args[i].String()
		s2("struct.go", value)

		// Skip the first argument which is Grid ID and if no = is in string
		if i == 0 || !strings.Contains(value, "=") {
			continue
		}

		lValues := strings.Split(args[i].String(), "=")
		if len(lValues) == 2 {
			arrSetting[lValues[0]] = lValues[1]

		}
	}

	return arrSetting
}

func (dg *DivGrid) setProperties(mapProp map[string]string) error {

	for key, val := range mapProp {
		switch key {
		case "ps":
			i, err := strconv.Atoi(val)
			if err == nil {
				dg.PageSize = i
			}
		case "sc":
			i, err := strconv.Atoi(val)
			if err == nil {
				dg.SortCol = i
			}
		case "st":
			dg.SortColType = val
		case "fi":
			dg.FilterID = val
		case "ft":
			dg.FilterText = val

		}
	}

	return nil
}

// TODO Documentation
type DivGrid struct {
	ID          string       // Div ID which is handled
	Header      DivGridRow   // First Row stored as XML Struct needs to be on top
	Rows        []DivGridRow // All the rows except header
	FilterID    string       // HTML Dom id of Filter Input Text Box
	FilterText  string       // Search Filter to define FilterRows
	FilterRows  []DivGridRow // Only the rows matching the search criteria
	SortCol     int          // Index of Column which is sorted
	SortColType string       // How sort will be handled option int and string
	PageSize    int          // Number of rows showed per page
	Debug       bool         // Property to use for debug
	DebugFile   string       // File used instead of html parts
	/*
			var _ID = '#toGrid';
		    var _sortcol = 0; // Index Column to besorted
		    var _sortcoltype = 'string';
		    var _direction = 'asc'; // TODO
		    // var _direction = 'desc';
		    var _search = '';
		    var _page = 1; //Paged to be showed
		    var _size = 0;
		    var _collen = 0;
		    var _rowlen = 0;
		    var _values = []; // ['March', 'Jan'];
		    //  var _numbers = []; // [67, 7];
		    var _maxSize = 3000;
		    var _defaultSize = 10;
		    var _minsearchchars = 1;
		    var _debug = true;
		    const _oddcss = "odd";
	*/

}

type DivGridParent struct {
	XMLName xml.Name     `xml:"div"`
	ID      string       `xml:"id,attr"`
	Class   string       `xml:"class,attr"`
	Rows    []DivGridRow `xml:"div"`
}

type DivGridRow struct {
	XMLName xml.Name      `xml:"div"`
	ID      string        `xml:"id,attr,omitempty"`
	Class   string        `xml:"class,attr"`
	Cells   []DivGridCell `xml:"div"`
}

type DivGridCell struct {
	XMLName     xml.Name `xml:"div"`
	Text        string   `xml:",chardata"`
	Type        string   `xml:"type,attr,omitempty"`
	Onclick     string   `xml:"onclick,attr,omitempty"`
	Onmouseover string   `xml:"onmouseover,attr,omitempty"`
}
