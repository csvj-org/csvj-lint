// Copyright CSVJ.org. All rights reserved.
// Use of this source code is governed by
// MIT license that can be found in the LICENSE file.

// Package lint verifies CSVJ file to be good, gives some warnings too
package lint

import (
	"fmt"
	"github.com/csvj-org/gocsvj"
	"io"
	"reflect"
)

// Message is a minimal element of lint result
type Message struct {
	Level   string // Level - either Info, Warning or Error
	Message string // Message - human readable string messages
}

// Info, Warning, Error - different message levels
const (
	Info    = "Info"
	Warning = "Warning"
	Error   = "Error"
)

func aplog(ll []Message, l string, m ...interface{}) []Message {
	ll = append(ll, Message{Level: l, Message: fmt.Sprint(m...)})
	return ll
}

// Do actually does CSVJ lint
func Do(reader *gocsvj.Reader) []Message {

	var lintLog []Message

	headers, err := reader.Headers()

	if err != nil {
		lintLog = aplog(lintLog, Error, err.Error())
	}

	lintLog = aplog(lintLog, Info, "header contains ", len(headers), " columns")

	rown := 1
	row1, err := reader.Read()
	lintLog = checkRowHeader(lintLog, 1, headers, row1)
	if err != nil {
		lintLog = aplog(lintLog, Error, "reading first row: ", err.Error())
	}

	for {
		row, err := reader.Read()
		if io.EOF == err {
			break
		}
		rown++
		if err != nil {
			lintLog = aplog(lintLog, Error, "reading ", rown, " row: ", err.Error())
		}
		lintLog = checkRowHeader(lintLog, rown, headers, row)
		lintLog = checkRowTypes(lintLog, rown, row1, row)
	}

	lintLog = aplog(lintLog, Info, "read ", rown, " rows")

	return lintLog
}

func checkRowTypes(lintLog []Message, rown int, row1 []interface{}, row []interface{}) []Message {
	for i, item := range row1 {
		if i >= len(row) {
			lintLog = aplog(lintLog, Warning, "row ", rown,
				" contains less number of elements ", len(row), " than first row ", len(row1))
			return lintLog
		}

		oitem := row[i]
		row1t := reflect.TypeOf(item)
		rownt := reflect.TypeOf(oitem)
		if row1t != rownt {
			lintLog = aplog(lintLog, Warning, "row ", rown, " column ", i+1, " type ", rownt,
				" differs from first row type ", row1t)
			return lintLog
		}
	}
	return lintLog
}

func checkRowHeader(lintLog []Message, rown int, headers []string, row []interface{}) []Message {
	lenh := len(headers)
	lenr := len(row)
	if lenh != lenr {
		lintLog = aplog(lintLog, Warning, "row ", rown,
			" contains different number of items ", lenr, " then header ", lenh)
	}
	return lintLog
}
