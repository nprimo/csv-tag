package main

import (
	"encoding/csv"
	"fmt"
	"os"

	csvtag "github.com/nprimo/csv-tag"
)

func main() {
	f, err := os.Open("example.csv")
	if err != nil {
		panic(err)
	}
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	type A struct {
		Id   int     `csv:"a"`
		Word string  `csv:"b"`
		Num  float32 `csv:"c"`
	}
	var data []A
	if err := csvtag.Unmarshal(rows, &data); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", data)
}
