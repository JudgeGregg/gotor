package main

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/JudgeGregg/gotor/globals"
	"github.com/JudgeGregg/gotor/parser"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func main() {

	argFile := os.Args[1]
	_, filename := path.Split(argFile)
	globals.RaidStartDate = getRaidStartDate(filename)
	file, err := os.Open(argFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	wInUTF8 := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	str, _ := io.ReadAll(wInUTF8)
	records := parser.Parse(string(str))
	raid := &parser.Raid{PlayersNumber: 1}
	for _, record := range records {
		handleRecord(raid, record)
	}
}

func getRaidStartDate(filename string) time.Time {
	filename = strings.TrimPrefix(filename, "combat_")
	filename = strings.TrimSuffix(filename, ".txt")
	last := strings.LastIndex(filename, "_")
	filenameRunes := []rune(filename)
	filenameRunes[last] = rune('.')
	raidStartDate, _ := time.Parse("2006-01-02_15_04_05.000000", string(filenameRunes))
	return raidStartDate
}

func handleRecord(raid *parser.Raid, record parser.Record) {
	if record.Effect.ActionID == globals.ENTERCOMBATID || record.Effect.ActionID == globals.EXITCOMBATID || record.Effect.EventID == globals.AREAENTEREDID || record.Effect.ActionID == globals.DEATHID || record.Effect.ActionID == globals.REVIVEID {
		handleStartStop(raid, record)
	} else if record.Effect.ActionID == globals.DAMAGEID {
		handleDamage(raid, record)
	} else if record.Effect.ActionID == globals.HEALID {
		handleHeal(raid, record)
	}
	handleThreat(raid, record)
}

func handleDamage(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		log.Println("DAMAGE NOT IN RAID", record.LineNumber)
		return
	}
	if raid.CurrentPull.Target == "" {
		checkPullTarget(raid, record)
	}
	actorPointer := record.Actor
	actorName := record.Actor.Name
	actorID := record.Actor.ID
	targetName := record.Target.Name
	targetID := record.Target.ID
	abilityName := record.Ability.Name
	abilityID := record.Ability.ID
	abilityAmount := record.Amount.Effective
	// Do we already know this actor ?
	if actor, ok := raid.CurrentPull.DamageDone[actorPointer]; ok {
		targetDamageDict := actor.TargetDamageDict
		// Do we already know this target ?
		if target, ok := targetDamageDict[targetName]; ok {
			// Do we already know this ability ?
			if ability, ok := target.Ability[abilityName]; ok {
				ability.Amount += abilityAmount
			} else {
				ability = &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
				target.Ability[abilityName] = ability
			}
		} else {
			ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
			target = &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
			target.Ability[abilityName] = ability
			actor.TargetDamageDict[targetName] = target
		}
	} else {
		ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
		target := &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
		actor := &parser.DamageDict{Name: actorName, ID: actorID, TargetDamageDict: make(map[string]*parser.TargetDamageDict)}
		target.Ability[abilityName] = ability
		actor.TargetDamageDict[targetName] = target
		raid.CurrentPull.DamageDone[actorPointer] = actor
	}
}

func handleThreat(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		//log.Println("THREAT NOT IN RAID", record.LineNumber)
		return
	}
	raid.CurrentPull.ThreatDone[record.Actor.Name] += record.Threat
}

func handleHeal(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		//log.Println("HEAL NOT IN RAID", record.LineNumber)
		return
	}
	if record.Actor.Name == "" {
		return
	}
	raid.CurrentPull.HealDone[record.Actor.Name] += record.Amount.Effective
}

func handleStartStop(raid *parser.Raid, record parser.Record) {
	if record.Effect.ActionID == globals.ENTERCOMBATID {
		if raid.InPull {
			//stop pull
			return
		}
		raid.CurrentPull = &parser.Pull{}
		raid.CurrentPull.StartTime = record.DateTime
		damageDone := make(map[parser.Actor]*parser.DamageDict)
		healDone := make(map[string]uint64)
		threatDone := make(map[string]uint64)
		raid.CurrentPull.DamageDone = damageDone
		raid.CurrentPull.HealDone = healDone
		raid.CurrentPull.ThreatDone = threatDone
		raid.AlivePlayersNumber = raid.PlayersNumber
		//start pull
		raid.InPull = true
		log.Printf("==============================")
		log.Printf("Starting fight %s", raid.CurrentPull.StartTime)
	}
	if record.Effect.ActionID == globals.EXITCOMBATID && raid.InPull {
		//stop pull
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		showDamage(raid.CurrentPull)
		log.Printf("Stopping fight exited %s", raid.CurrentPull.StopTime)
	}
	if record.Effect.EventID == globals.AREAENTEREDID {
		switch record.Effect.SpecID {
		case globals.FOURPLAYERVETERAN:
			raid.PlayersNumber = 4
			raid.Difficulty = globals.VETERAN
		case globals.FOURPLAYERMASTER:
			raid.PlayersNumber = 4
			raid.Difficulty = globals.MASTER
		case globals.EIGHTPLAYERSTORY:
			raid.PlayersNumber = 8
			raid.Difficulty = globals.STORY
		case globals.EIGHTPLAYERVETERAN:
			raid.PlayersNumber = 8
			raid.Difficulty = globals.VETERAN
		case globals.EIGHTPLAYERMASTER:
			raid.PlayersNumber = 8
			raid.Difficulty = globals.MASTER
		}
	}
	if raid.InPull && record.Effect.ActionID == globals.DEATHID && !record.Target.NPC {
		//log.Printf("%s DEAD at %s\n", record.Target.Name, record.DateTime)
		raid.AlivePlayersNumber -= 1
		if raid.AlivePlayersNumber == 0 {
			raid.InPull = false
			raid.CurrentPull.StopTime = record.DateTime
			raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
			showDamage(raid.CurrentPull)
			log.Printf("Stopping fight WIPE %s", raid.CurrentPull.StopTime)
		}
	}
	if raid.InPull && record.Effect.ActionID == globals.REVIVEID && !record.Target.NPC {
		//log.Printf("%s REVIVED at %s\n", record.Target.Name, record.DateTime)
		raid.AlivePlayersNumber += 1
	}
}

func showDamage(pull *parser.Pull) {
	if pull.Target == "" {
		return
	}
	duration := pull.StopTime.Sub(pull.StartTime)
	log.Println(duration)
	log.Println(pull.Target)
	log.Println("------------------------------")
	seconds := duration.Seconds()

	for player, dmgDict := range pull.DamageDone {
		if !player.NPC {
			totalDamage := float64(0)
			for _, targetDmgDict := range dmgDict.TargetDamageDict {
				for _, abilityDmgDict := range targetDmgDict.Ability {
					totalDamage += float64(abilityDmgDict.Amount)
				}
			}
			log.Println(player.Name, totalDamage, totalDamage/seconds)
		}
	}
}

func checkPullTarget(raid *parser.Raid, record parser.Record) {
	if name, ok := globals.Targets[record.Target.ID]; ok {
		raid.CurrentPull.Target = name
	}
}
