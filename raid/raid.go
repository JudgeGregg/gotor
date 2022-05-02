package raid

import (
	"log"
	"strings"
	"time"

	"github.com/JudgeGregg/gotor/globals"
	"github.com/JudgeGregg/gotor/parser"
)

func GetRaidStartDate(filename string) time.Time {
	filename = strings.TrimPrefix(filename, "combat_")
	filename = strings.TrimSuffix(filename, ".txt")
	last := strings.LastIndex(filename, "_")
	filenameRunes := []rune(filename)
	filenameRunes[last] = rune('.')
	raidStartDate, _ := time.Parse("2006-01-02_15_04_05.000000", string(filenameRunes))
	return raidStartDate
}

func HandleRecord(raid *parser.Raid, record parser.Record) {
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
		//log.Println("DAMAGE NOT IN RAID", record.LineNumber)
		return
	}
	if raid.CurrentPull.Target == "" {
		checkPullTarget(raid, record)
	}
	actor := record.Actor
	actorName := record.Actor.Name
	actorID := record.Actor.ID
	target := record.Target
	targetName := record.Target.Name
	targetID := record.Target.ID
	abilityName := record.Ability.Name
	abilityID := record.Ability.ID
	abilityAmount := record.Amount.Effective
	// Do we already know this actor ?
	if actorDmgDict, ok := raid.CurrentPull.DamageDone[actor]; ok {
		targetDamageDict := actorDmgDict.TargetDamageDict
		// Do we already know this target ?
		if targetDmgDict, ok := targetDamageDict[target]; ok {
			// Do we already know this ability ?
			if ability, ok := targetDmgDict.Ability[abilityName]; ok {
				ability.Amount += abilityAmount
				ability.Hits += 1
				if record.Amount.Mitigated {
					switch record.Amount.Mitigation {
					case globals.IMMUNE:
						{
							ability.Immune += 1
						}
					case globals.RESIST:
						{
							ability.Resist += 1
						}
					case globals.MISS:
						{
							ability.Miss += 1
						}
					case globals.DODGE_OR_PARRY:
						{
							ability.DodgeOrParry += 1
						}
					case globals.SHIELD:
						{
							ability.Shield += 1
						}
					}
				}
			} else {
				ability = &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
				targetDmgDict.Ability[abilityName] = ability
			}
		} else {
			ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
			targetDmgDict = &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
			targetDmgDict.Ability[abilityName] = ability
			actorDmgDict.TargetDamageDict[target] = targetDmgDict
		}
	} else {
		ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
		targetDmgDict := &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
		actorDmgDict := &parser.DamageDict{Name: actorName, ID: actorID, TargetDamageDict: make(map[parser.Target]*parser.TargetDamageDict)}
		targetDmgDict.Ability[abilityName] = ability
		actorDmgDict.TargetDamageDict[target] = targetDmgDict
		raid.CurrentPull.DamageDone[actor] = actorDmgDict
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
	actor := record.Actor
	actorName := record.Actor.Name
	actorID := record.Actor.ID
	target := record.Target
	targetID := record.Target.ID
	targetName := record.Target.ID
	abilityName := record.Ability.Name
	abilityID := record.Ability.ID
	abilityAmount := record.Amount.Effective
	// Do we already know this actor ?
	if healDict, ok := raid.CurrentPull.HealDone[actor]; ok {
		targetHealDict := healDict.TargetHealDict
		// Do we already know this target ?
		if targetHeaDict, ok := targetHealDict[target]; ok {
			// Do we already know this ability ?
			if ability, ok := targetHeaDict.Ability[abilityName]; ok {
				ability.Amount += abilityAmount
				ability.Hits += 1
			} else {
				ability = &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
				targetHeaDict.Ability[abilityName] = ability
			}
		} else {
			ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
			targetHeaDict = &parser.TargetHealDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
			targetHeaDict.Ability[abilityName] = ability
			healDict.TargetHealDict[target] = targetHeaDict
		}
	} else {
		ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective}
		targetHeaDict := &parser.TargetHealDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
		heaDict := &parser.HealDict{Name: actorName, ID: actorID, TargetHealDict: make(map[parser.Target]*parser.TargetHealDict)}
		targetHeaDict.Ability[abilityName] = ability
		heaDict.TargetHealDict[target] = targetHeaDict
		raid.CurrentPull.HealDone[actor] = heaDict
	}
}

func handleStartStop(raid *parser.Raid, record parser.Record) {
	if record.Effect.ActionID == globals.ENTERCOMBATID {
		if raid.InPull {
			//stop pull
			return
		}
		//start pull
		raid.CurrentPull = &parser.Pull{}
		raid.CurrentPull.StartTime = record.DateTime
		damageDone := make(map[parser.Actor]*parser.DamageDict)
		healDone := make(map[parser.Actor]*parser.HealDict)
		threatDone := make(map[string]uint64)
		raid.CurrentPull.DamageDone = damageDone
		raid.CurrentPull.HealDone = healDone
		raid.CurrentPull.ThreatDone = threatDone
		raid.AlivePlayersNumber = raid.PlayersNumber
		raid.InPull = true
		log.Printf("%d Starting fight %s", record.LineNumber, raid.CurrentPull.StartTime)
	}
	if record.Effect.ActionID == globals.EXITCOMBATID && raid.InPull {
		//stop pull
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		showDamage(raid.CurrentPull)
		log.Printf("%d Stopping fight exited %s", record.LineNumber, raid.CurrentPull.StopTime)
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
			//stop pull, WIPE
			raid.InPull = false
			raid.CurrentPull.StopTime = record.DateTime
			raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
			showDamage(raid.CurrentPull)
			//log.Printf("%d Stopping fight WIPE %s", record.LineNumber, raid.CurrentPull.StopTime)
		}
	}
	if raid.InPull && record.Effect.ActionID == globals.REVIVEID && !record.Target.NPC {
		//log.Printf("%s REVIVED at %s\n", record.Target.Name, record.DateTime)
		raid.AlivePlayersNumber += 1
	}
}

func showDamage(pull *parser.Pull) {
	if pull.Target == "" && !globals.Debug {
		return
	}
	log.Printf("==============================")
	log.Printf("STARTING FIGHT %s", pull.StartTime)
	duration := pull.StopTime.Sub(pull.StartTime)
	seconds := duration.Seconds()
	log.Println(duration)
	log.Println(pull.Target)
	log.Println("------------------------------")
	log.Println("DAMAGE DONE")

	receivedMap := make(map[parser.Target]float64)
	for player, dmgDict := range pull.DamageDone {
		if !player.NPC && player.Name != "" {
			totalDamage := float64(0)
			for _, targetDmgDict := range dmgDict.TargetDamageDict {
				for _, abilityDmgDict := range targetDmgDict.Ability {
					totalDamage += float64(abilityDmgDict.Amount)
				}
			}
			log.Println(player.Name, totalDamage, "Total", totalDamage/seconds, "DPS")
		} else {
			for target, targetDmgDict := range dmgDict.TargetDamageDict {
				if !target.NPC {
					for _, abilityDmgDict := range targetDmgDict.Ability {
						receivedMap[target] += float64(abilityDmgDict.Amount)
					}
				}
			}
		}
	}
	log.Println("------------------------------")
	log.Println("DAMAGE RECEIVED")
	for target, amount := range receivedMap {
		log.Println(target.Name, amount, "Total", amount/seconds, "DPS")
	}
	log.Println("------------------------------")
	log.Println("HEAL DONE")
	for player, healDict := range pull.HealDone {
		if !player.NPC {
			totalHeal := float64(0)
			for _, targetHeaDict := range healDict.TargetHealDict {
				for _, abilityDmgDict := range targetHeaDict.Ability {
					totalHeal += float64(abilityDmgDict.Amount)
				}
			}
			log.Println(player.Name, totalHeal, "Total", totalHeal/seconds, "HPS")
		}
	}
	log.Println("------------------------------")
	log.Printf("STOPPING FIGHT %s", pull.StopTime)
	log.Printf("==============================")
}

func checkPullTarget(raid *parser.Raid, record parser.Record) {
	if name, ok := globals.Targets[record.Target.ID]; ok {
		raid.CurrentPull.Target = name
	}
}
