package main

import (
	"io"
	"os"

	"log"

	"github.com/JudgeGregg/gotor/ids"
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
	if record.Effect.ActionID == ids.ENTERCOMBATID || record.Effect.ActionID == ids.EXITCOMBATID || record.Effect.EventID == ids.AREAENTEREDID {
		handleStartStop(raid, record)
	} else if record.Effect.ActionID == ids.DAMAGEID {
		handleDamage(raid, record)
	} else if record.Effect.ActionID == ids.HEALID {
		handleHeal(raid, record)
	}
}

func handleDamage(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("NOT IN RAID", record.LineNumber)
		return
	}
	raid.CurrentPull.DamageDone[record.Actor.Name] += record.Amount.Effective
}

func handleHeal(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("HEAL NOT IN RAID", record.LineNumber)
		return
	}
	if record.Actor.Name == "" {
		return
	}
	raid.CurrentPull.HealDone[record.Actor.Name] += record.Amount.Effective
}

func handleStartStop(raid *parser.Raid, record parser.Record) {
	if record.Effect.ActionID == ids.ENTERCOMBATID {
		if raid.InPull {
			//stop pull
			raid.InPull = false
			log.Printf("Stopping fight")
		}
		raid.CurrentPull = &parser.Pull{}
		damagedone := make(map[string]uint64)
		healdone := make(map[string]uint64)
		raid.CurrentPull.DamageDone = damagedone
		raid.CurrentPull.HealDone = healdone
		//start pull
		raid.InPull = true
		log.Printf("Starting fight")
	}
	if record.Effect.ActionID == ids.EXITCOMBATID && raid.InPull {
		//stop pull
		raid.InPull = false
		log.Printf("Stopping fight")
	}
	if record.Effect.Event == ids.AREAENTEREDID && raid.InPull {
		//stop pull
		raid.InPull = false
		log.Printf("Stopping fight")
	}
}
