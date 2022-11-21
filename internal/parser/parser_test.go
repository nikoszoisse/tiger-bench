package parser

import (
	"log"
	"os"
	"testing"
)

func Test(t *testing.T) {

	outputRecords := ReadCsvFile("./data/query_params.csv")

	for {
		select {
		case result, ok := <-outputRecords:
			if ok {
				log.Printf("%+v \n", result)
			} else {
				return
			}
		}
	}

	os.Exit(1)
}
