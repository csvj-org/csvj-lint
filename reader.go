// Copyright CSVJ.org. All rights reserved.
// Use of this source code is governed by
// MIT license that can be found in the LICENSE file.

// Package github.com/csvj-org/gocsvj reads and writes comma-separated values files of CSVJ type.
// CSVJ is a csv-like format for storing tables that follows certain JSON encoding rules.

package gocsvj

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// A Reader reads records from a CSVJ-encoded file.
type Reader struct {
	line       int
	headerRead bool
	header     []string
	r          *bufio.Scanner
	clSet      bool
	clValues   []interface{}
	clError    error
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewScanner(r),
	}
}

// Headers reads header record from r and caches it
// so it could be returned later too
func (r *Reader) Headers() ([]string, error) {
	if r.headerRead {
		return r.header, nil
	}

	rawHeader, err := r.readLine()
	if err != nil {
		return nil, err
	}

	r.line = 0

	r.header, err = valuesAsStrings(rawHeader)
	if err != nil {
		return nil, err
	}

	r.headerRead = true
	return r.header, nil
}

// // Read reads one record (a slice of fields) from r
func (r *Reader) Read() ([]interface{}, error) {
	if r.headerRead == false {
		r.Headers()
	}
	return r.readLine()
}

func valuesAsStrings(vs []interface{}) ([]string, error) {
	strs := make([]string, len(vs))

	for i, v := range vs {
		if w, ok := v.(string); ok {
			strs[i] = w
		} else {
			return nil, errors.New("non-string item at csvj header")
		}
	}
	return strs, nil
}

func (r *Reader) scanLine() error {
	if !r.r.Scan() {
		err := r.r.Err()
		if err == nil {
			return io.EOF
		}
		return err
	}

	if err := r.r.Err(); err != nil {
		return err
	}
	return nil
}

func (r *Reader) readLine() ([]interface{}, error) {
	r.line++

	if r.clSet {
		r.clSet = false
		return r.clValues, r.clError
	}

	err := r.scanLine()

	if err != nil {
		return nil, err
	}

	sl := strings.TrimSpace(r.r.Text())
	if sl == "" {
		r.clValues, r.clError = r.readLine()
		if r.clError == io.EOF {
			return nil, io.EOF
		}
		r.clSet = true
	}
	line := "[" + sl + "]"

	var lv []interface{}
	err = json.Unmarshal([]byte(line), &lv)
	if err != nil {
		err = fmt.Errorf("parse error row %d : %s", r.line, err.Error())
		return nil, err
	}

	typesafe, erritem := checkCSVJTypes(lv)

	if !typesafe {
		return nil, fmt.Errorf("row %d parse error at item %d", r.line, erritem)
	}

	return lv, nil
}

func checkCSVJTypes(ar []interface{}) (bool, int) {

	for idx, el := range ar {
		if el == nil {
			continue
		}

		switch el.(type) {
		case float64:
		case string:
		case bool:
			continue

		default:
			return false, idx
		}
	}

	return true, -1
}
