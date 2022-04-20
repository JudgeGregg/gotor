package main

import (
	"io"
	"os"

	"log"

	"github.com/JudgeGregg/gotor/parser"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {
	file, _ := os.Open("samples/sample1.txt")
	//file, _ := os.Open("combat_2022-04-16_21_23_12_979036.txt")
	wInUTF8 := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	str, _ := io.ReadAll(wInUTF8)
	records := parser.Parse(string(str))
	raid := &parser.Raid{}
	for _, record := range records {
		handleRecord(raid, record)
	}
}

func handleRecord(raid *parser.Raid, record parser.Record) {
	if record.Effect.Action == "EnterCombat" || record.Effect.Action == "ExitCombat" || record.Effect.Event == "AreaEntered" {
		handleStartStop(raid, record)
	} else if record.Effect.Action == "Damage" {
		handleDamage(raid, record)
	}
}

func handleDamage(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("NOT IN RAID", record.LineNumber)
	}
	log.Printf("%s hit %s for %d damage\n", record.Actor.Name, record.Target.Name, record.Amount.Effective)
}

func handleStartStop(raid *parser.Raid, record parser.Record) {
	if record.Effect.Action == "EnterCombat" {
		if raid.InPull {
			//stop pull
			raid.InPull = false
			log.Printf("Stopping fight")
		}
		//start pull
		raid.InPull = true
		log.Printf("Starting fight")
	}
	if record.Effect.Action == "ExitCombat" && raid.InPull {
		//stop pull
		raid.InPull = false
		log.Printf("Stopping fight")
	}
	if record.Effect.Event == "AreaEntered" && raid.InPull {
		//stop pull
		raid.InPull = false
		log.Printf("Stopping fight")
	}
}
