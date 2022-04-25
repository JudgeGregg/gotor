package main

import (
	"io"
	"log"
	"os"
	"path"

	"github.com/JudgeGregg/gotor/globals"
	"github.com/JudgeGregg/gotor/parser"
	"github.com/JudgeGregg/gotor/raid"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	argFile := os.Args[1]
	_, filename := path.Split(argFile)
	globals.RaidStartDate = raid.GetRaidStartDate(filename)
	file, err := os.Open(argFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	wInUTF8 := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	str, _ := io.ReadAll(wInUTF8)
	records := parser.Parse(string(str))
	raid_ := &parser.Raid{PlayersNumber: 1}
	for _, record := range records {
		raid.HandleRecord(raid_, record)
	}
}
