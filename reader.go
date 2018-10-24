package gocsvj

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type CSVJValue interface{}

type Reader struct {
	line       int
	headerRead bool
	header     []string
	r          *bufio.Scanner
	clSet      bool
	clValues   []CSVJValue
	clError    error
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: bufio.NewScanner(r),
	}
}

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

func (r *Reader) Read() ([]CSVJValue, error) {
	if r.headerRead == false {
		r.Headers()
	}
	return r.readLine()
}

func valuesAsStrings(vs []CSVJValue) ([]string, error) {
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

func (r *Reader) readLine() ([]CSVJValue, error) {
	r.line++

	if r.clSet {
		r.clSet = false
		return r.clValues, r.clError
	}

	if !r.r.Scan() {
		err := r.r.Err()
		if err == nil {
			return nil, io.EOF
		}
		return nil, err
	}

	if err := r.r.Err(); err != nil {
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

	var lv []CSVJValue
	err := json.Unmarshal([]byte(line), &lv)
	if err != nil {
		err = errors.New(fmt.Sprint("parse error row ", r.line, ": ", err.Error()))
		return nil, err
	}

	for idx, el := range lv {
		if el == nil {
			continue
		}

		switch el.(type) {
		case float64:
		case string:
		case bool:
			continue

		default:
			return nil, errors.New(fmt.Sprintf("row %d parse error at item %d", r.line, idx))
		}
	}

	return lv, nil
}
