package gocsvj

import (
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"
)

func TestSimple(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `"Row1", "Row2", "Row3"` + "\n"
	csvj += " " // empty last line, just in case

	r := NewReader(strings.NewReader(csvj))

	// test initial parse and cache
	for l := 0; l < 2; l++ {
		hdr, err := r.Headers()
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(hdr, []string{"Header1", "Header2", "Header3"}) {
			t.Error("Unexpected Header")
		}
	}

	row, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	erow := []CSVJValue{"Row1", "Row2", "Row3"}
	if !reflect.DeepEqual(row, erow) {
		t.Error("Bad Row", row, "expected", erow)
	}

	_, eofErr := r.Read()
	if eofErr != io.EOF {
		t.Error("EOF is expected on empty line")
	}
}

func TestSimpleNoNewline(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `42, 42, false`

	r := NewReader(strings.NewReader(csvj))

	hdr, err := r.Headers()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(hdr, []string{"Header1", "Header2", "Header3"}) {
		t.Error("Unexpected Header")
	}

	row, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	erow := []CSVJValue{42.0, 42.0, false}

	if !reflect.DeepEqual(row, erow) {
		t.Error("Bad Row", row, "expected", erow, "reason")
	}

	_, eofErr := r.Read()
	if eofErr != io.EOF {
		t.Error("EOF is expected on empty line")
	}
}

func TestEmptyLineInMiddle(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += "\n"
	csvj += `null, null, true`

	r := NewReader(strings.NewReader(csvj))

	row, err := r.Read()
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(row, []CSVJValue{}) {
		t.Error("Bad Row", row, "expected empty array")
	}

	row, err = r.Read()
	if err != nil {
		t.Error(err)
	}

	erow := []CSVJValue{nil, nil, true}

	if !reflect.DeepEqual(row, erow) {
		t.Error("Bad Row", row, "expected", erow, "reason")
	}

	_, eofErr := r.Read()
	if eofErr != io.EOF {
		t.Error("EOF is expected on empty line")
	}
}

func TestParseError(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `42, $, false`

	r := NewReader(strings.NewReader(csvj))

	_, err := r.Headers()
	if err != nil {
		t.Error(err)
	}

	_, err = r.Read()
	if err == nil {
		t.Error("expected error, but none returned")
	}
}

func TestParseJSLikeError(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `42, [], false`

	r := NewReader(strings.NewReader(csvj))

	_, err := r.Headers()
	if err != nil {
		t.Error(err)
	}

	_, err = r.Read()
	if err == nil {
		t.Error("expected error, but none returned")
	}
}

func TestHeaderError(t *testing.T) {
	csvj := `"Header1", 1, "Header2", "Header3"` + "\n"

	r := NewReader(strings.NewReader(csvj))

	_, err := r.Headers()
	if err == nil {
		t.Error("expected error, but none returned")
	}
}

func TestReadError(t *testing.T) {
	csvj := `"Header1", "Header2", "Header3"` + "\n"
	csvj += `42, 1, false`

	r := NewReader(iotest.TimeoutReader(strings.NewReader(csvj)))

	_, err := r.Read()
	if err == nil {
		t.Error("expected error, but none returned")
	}
}

func TestEmptyError(t *testing.T) {
	csvj := ""

	r := NewReader(strings.NewReader(csvj))

	_, err := r.Headers()
	if err == nil {
		t.Error("expected error, but none returned")
	}
}
