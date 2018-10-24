package lint

import (
	"fmt"
	"github.com/csvj-org/gocsvj"
	"reflect"
	"strings"
	"testing"
)

func TestSimple(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `"Row1", "Row2", "Row3"` + "\n"
	csvj += `"Row1", "Row2", "Row3"` + "\n"

	result := Do(gocsvj.NewReader(strings.NewReader(csvj)))

	er := []Message{
		{Info, "header contains 3 columns"},
		{Info, "read 2 rows"},
	}

	if !reflect.DeepEqual(result, er) {
		t.Error("Bad Basic Info Messages")
	}
}

func TestRowType(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `"Row1", "Row2", "Row3"` + "\n"
	csvj += `"Row1", 2, "Row3"` + "\n"

	result := Do(gocsvj.NewReader(strings.NewReader(csvj)))

	er := []Message{
		{Info, "header contains 3 columns"},
		{Warning, "row 2 column 2 type float64 differs from first row type string"},
		{Info, "read 2 rows"},
	}

	if !reflect.DeepEqual(result, er) {
		t.Error("Bad Basic Info Messages")
	}
}

func TestRowCount(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `"Row1", "Row2", "Row3"` + "\n"
	csvj += `"Row1", "Row2"` + "\n"

	result := Do(gocsvj.NewReader(strings.NewReader(csvj)))

	er := []Message{
		{Info, "header contains 3 columns"},
		{Warning, "row 2 contains different number of items 2 then header 3"},
		{Warning, "row 2 contains less number of elements 2 than first row 3"},
		{Info, "read 2 rows"},
	}

	if !reflect.DeepEqual(result, er) {
		fmt.Println(result)
		t.Error("Bad Basic Info Messages")
	}
}
