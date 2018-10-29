package gocsvj

import (
	"strings"
	"testing"
)

func TestSimpleWriter(t *testing.T) {
	var sw strings.Builder

	w := NewWriter(&sw)
	w.WriteHeader([]string{"h1", "h2", "h3"})
	err := w.Write([]int{2, 4, 5})

	if err != nil {
		t.Error(err)
	}

	w.Flush()

	if sw.String() != `"h1","h2","h3"`+"\r\n2,4,5\r\n" {
		t.Error("unexpected CSVJ")
	}

	err = w.Error()
	if err != nil {
		t.Error(err)
	}
}

func TestInterface(t *testing.T) {
	var sw strings.Builder

	w := NewWriter(&sw)
	w.WriteHeader([]string{"h1", "h2", "h3"})
	err := w.Write([]interface{}{"test", nil, 42})

	if err != nil {
		t.Error(err)
	}

	w.Flush()

	if sw.String() != `"h1","h2","h3"`+"\r\n"+`"test",null,42`+"\r\n" {
		t.Error("unexpected CSVJ: ", sw.String())
	}

	err = w.Error()
	if err != nil {
		t.Error(err)
	}
}

func TestWriterNonSlice(t *testing.T) {
	var sw strings.Builder
	w := NewWriter(&sw)
	w.WriteHeader([]string{"h1"})
	err := w.Write(42)

	if err == nil {
		t.Error("Expected error, but none returned")
	}
}

func TestWriterBadHeader(t *testing.T) {
	var sw strings.Builder
	w := NewWriter(&sw)
	w.WriteHeader([]string{""})
	err := w.Write([]string{"item", "item2"})

	if err == nil {
		t.Error("Expected header error")
	}
}

func TestWriteNonCSVJ(t *testing.T) {
	var sw strings.Builder
	w := NewWriter(&sw)

	w.WriteHeader([]string{"h1", "h2", "h3"})

	mp := make(map[string]string)
	mp["test"] = "test"
	err := w.Write([]interface{}{2, 3, mp})

	if err == nil {
		t.Error("Expected error, but none returned")
	}
}
