package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"log"

	"github.com/JudgeGregg/gotor/ids"
	"github.com/JudgeGregg/gotor/parser"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

var raidStartDate time.Time

func main() {
	raidStartDate = time.Date(22, 4, 16, 21, 23, 12, 979036000, time.UTC)
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
	handleThreat(raid, record)
}

func handleDamage(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("DAMAGE NOT IN RAID", record.LineNumber)
		return
	}
	actorName := record.Actor.Name
	actorID := record.Actor.ID
	targetName := record.Target.Name
	targetID := record.Target.ID
	abilityName := record.Ability.Name
	abilityID := record.Ability.ID
	abilityAmount := record.Amount.Effective
	// Do we already know this actor ?
	if actor, ok := raid.CurrentPull.DamageDone[actorName]; ok {
		targetDamageDict := actor.TargetDamageDict
		// Do we already know this target ?
		if target, ok := targetDamageDict[targetName]; ok {
			// Do we already know this ability ?
			if ability, ok := target.Ability[abilityName]; ok {
				ability.Amount += abilityAmount
			} else {
				ability = parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
				target.Ability[abilityName] = ability
				actor.TargetDamageDict[targetName] = target
				raid.CurrentPull.DamageDone[actorName] = actor
			}
			ability := parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
			target.Ability[abilityName] = ability
			actor.TargetDamageDict[targetName] = target
		} else {
			ability := parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
			target = parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]parser.AbilityDict)}
			target.Ability[abilityName] = ability
			actor.TargetDamageDict[targetName] = target
		}
	} else {
		ability := parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
		target := parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]parser.AbilityDict)}
		actor := parser.DamageDict{Name: actorName, ID: actorID, TargetDamageDict: make(map[string]parser.TargetDamageDict)}
		target.Ability[abilityName] = ability
		actor.TargetDamageDict[targetName] = target
		raid.CurrentPull.DamageDone[actorName] = actor
	}
}

func handleThreat(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("THREAT NOT IN RAID", record.LineNumber)
		return
	}
	raid.CurrentPull.ThreatDone[record.Actor.Name] += record.Threat
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
		raid.CurrentPull.StartTime = record.DateTime
		damageDone := make(map[string]parser.DamageDict)
		healDone := make(map[string]uint64)
		threatDone := make(map[string]uint64)
		raid.CurrentPull.DamageDone = damageDone
		raid.CurrentPull.HealDone = healDone
		raid.CurrentPull.ThreatDone = threatDone
		//start pull
		raid.InPull = true
		log.Printf("Starting fight")
	}
	if record.Effect.ActionID == ids.EXITCOMBATID && raid.InPull {
		//stop pull
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		fmt.Println(raid.CurrentPull.StopTime.Sub(raid.CurrentPull.StartTime))
		fmt.Println(raid.CurrentPull.DamageDone)
		log.Printf("Stopping fight")
	}
	if record.Effect.Event == ids.AREAENTEREDID && raid.InPull {
		//stop pull
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		log.Printf("Stopping fight")
	}
}
