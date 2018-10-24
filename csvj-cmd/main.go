package main

import (
	"fmt"
	"github.com/csvj-org/gocsvj"
	"github.com/csvj-org/gocsvj/lint"
	"log"
	"os"
)

func main() {

	if len(os.Args) == 1 {
		fmt.Println("usage:", os.Args[0], "<command> <file>")
		fmt.Println("")
		fmt.Println("following commands are supported")
		fmt.Println("  lint <file.csvj> - check CSVJ syntax")
		os.Exit(1)
	}

	fileName := os.Args[1]

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Cannot open file: ", err)
	}
	defer file.Close()

	reader := gocsvj.NewReader(file)

	result := lint.Do(reader)

	ecode := 0

	for _, m := range result {
		fmt.Printf("%s: %s\n", m.Level, m.Message)
		if m.Level != lint.Info {
			ecode = 1
		}
	}

	os.Exit(ecode)

}
