package raid

import (
	"log"
	"math"
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
	if record.Effect.EventID == globals.AREAENTEREDID {
		handleAreaEntered(raid, record)
	}
	switch record.Effect.ActionID {
	case globals.ENTERCOMBATID, globals.EXITCOMBATID, globals.DEATHID:
		handleStartStop(raid, record)
	case globals.DAMAGEID:
		handleDamage(raid, record)
	case globals.HEALID:
		handleHeal(raid, record)
	case globals.FORCE_ARMOURID, globals.STATIC_BARRIERID, globals.EMERGENCY_POWERID, globals.ABSORB_SHIELDID, globals.SABER_WARD_GID, globals.SABER_WARD_JID, globals.SONIC_BARRIERID, globals.SHIELD_PROBEID, globals.BALLISTIC_DAMPERSID, globals.ENERGY_REDOUBT_PID, globals.ENERGY_REDOUBT_VID, globals.BLADE_BARRIERID:
		handleBubble(raid, record)
	}
	handleThreat(raid, record)
}

func handleAreaEntered(raid *parser.Raid, record parser.Record) {
	switch record.Effect.SpecID {
	case globals.FOURPLAYERSTORY:
		raid.Difficulty = globals.STORY
	case globals.FOURPLAYERVETERAN:
		raid.Difficulty = globals.VETERAN
	case globals.FOURPLAYERMASTER:
		raid.Difficulty = globals.MASTER
	case globals.EIGHTPLAYERSTORY:
		raid.Difficulty = globals.STORY
	case globals.EIGHTPLAYERVETERAN:
		raid.Difficulty = globals.VETERAN
	case globals.EIGHTPLAYERMASTER:
		raid.Difficulty = globals.MASTER
	}
}

func handleBubble(raid *parser.Raid, record parser.Record) {
	currentBubbler := record.Actor
	bubblee := record.Target
	switch record.Effect.EventID {
	case globals.APPLY_EFFECTID:
		bubbler := parser.Bubbler{CurrentBubbler: currentBubbler}
		if bubb, ok := raid.BubblerMap[bubblee]; ok {
			bubbler.PreviousBubbler = bubb.CurrentBubbler
		}
		raid.BubblerMap[bubblee] = bubbler
		//log.Printf("%s SET A BUBBLE ON %s, at %s! ", record.Actor.Name, record.Target.Name, record.DateTime)
		//log.Printf("%v", raid.BubblerMap)
	case globals.REMOVE_EFFECTID:
		bubbler := parser.Bubbler{}
		if bubb, ok := raid.BubblerMap[bubblee]; ok {
			if bubb.CurrentBubbler.Name == "" {
				bubbler.PreviousBubbler = bubb.PreviousBubbler
			} else {
				bubbler.PreviousBubbler = bubb.CurrentBubbler
			}
		}
		raid.BubblerMap[bubblee] = bubbler
		//log.Printf("BUBBLE SET BY %s ON %s EXPIRED, at %s! ", record.Actor.Name, record.Target.Name, record.DateTime)
		//log.Printf("%#v", raid.BubblerMap)
	}
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
				if ability.DamageType == "" {
					ability.DamageType = record.Amount.DamageType
				}
				ability.Hits += 1
				handleMitigation(ability, record)
			} else {
				ability = &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective, DamageType: record.Amount.DamageType, Hits: 1}
				handleMitigation(ability, record)
				targetDmgDict.Ability[abilityName] = ability
			}
		} else {
			ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective, DamageType: record.Amount.DamageType, Hits: 1}
			handleMitigation(ability, record)
			targetDmgDict = &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
			targetDmgDict.Ability[abilityName] = ability
			actorDmgDict.TargetDamageDict[target] = targetDmgDict
		}
	} else {
		ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: record.Amount.Effective, DamageType: record.Amount.DamageType, Hits: 1}
		handleMitigation(ability, record)
		targetDmgDict := &parser.TargetDamageDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
		actorDmgDict := &parser.DamageDict{Name: actorName, ID: actorID, TargetDamageDict: make(map[parser.Entity]*parser.TargetDamageDict)}
		targetDmgDict.Ability[abilityName] = ability
		actorDmgDict.TargetDamageDict[target] = targetDmgDict
		raid.CurrentPull.DamageDone[actor] = actorDmgDict
	}
	if record.Amount.Absorbed {
		handleAbsorb(raid, record)
	}
}

func handleAbsorb(raid *parser.Raid, record parser.Record) {
	if math.Abs(float64(record.Amount.Effective-record.Amount.Amount)) <= 1.0 {
		//Weird +1 delta, ignore BUG ?
		return
	}
	target := record.Target
	bubbler := raid.BubblerMap[target].CurrentBubbler
	if bubbler.Name == "" {
		bubbler = raid.BubblerMap[target].PreviousBubbler
	}
	if bubbler.Name == "" {
		bubbler = raid.BubblerMap[target].PreviousBubbler
		//log.Printf("%s took %d,  absorbed %d, at %s\n", record.Target.Name, record.Amount.Amount, record.Amount.Amount-record.Amount.Effective, record.DateTime)
	}
	absDone := raid.CurrentPull.AbsDone
	absDone[bubbler] += float64(record.Amount.Amount - record.Amount.Effective)
}

func handleMitigation(ability *parser.AbilityDict, record parser.Record) {
	if record.Amount.Mitigated {
		switch record.Amount.Mitigation {
		case globals.IMMUNE:
			ability.MitigationDict.Immune += 1
		case globals.RESIST:
			ability.MitigationDict.Resist += 1
		case globals.MISS:
			ability.MitigationDict.Miss += 1
		case globals.DODGE_PARRY_DEFLECT:
			ability.MitigationDict.DodgeParryDeflect += 1
		case globals.SHIELD:
			ability.MitigationDict.Shield += 1
		}
	}
}

func handleThreat(raid *parser.Raid, record parser.Record) {
	if !raid.InPull {
		//log.Println("THREAT NOT IN RAID", record.LineNumber)
		return
	}
	raid.CurrentPull.ThreatDone[record.Actor.Name] += float64(record.Threat)
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
				ability = &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: abilityAmount}
				targetHeaDict.Ability[abilityName] = ability
			}
		} else {
			ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: abilityAmount}
			targetHeaDict = &parser.TargetHealDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
			targetHeaDict.Ability[abilityName] = ability
			healDict.TargetHealDict[target] = targetHeaDict
		}
	} else {
		ability := &parser.AbilityDict{Name: abilityName, ID: abilityID, Amount: abilityAmount}
		targetHeaDict := &parser.TargetHealDict{Name: targetName, ID: targetID, Ability: make(map[string]*parser.AbilityDict)}
		heaDict := &parser.HealDict{Name: actorName, ID: actorID, TargetHealDict: make(map[parser.Entity]*parser.TargetHealDict)}
		targetHeaDict.Ability[abilityName] = ability
		heaDict.TargetHealDict[target] = targetHeaDict
		raid.CurrentPull.HealDone[actor] = heaDict
	}
}

func handleStartStop(raid *parser.Raid, record parser.Record) {
	if record.Effect.ActionID == globals.ENTERCOMBATID {
		//start pull
		raid.CurrentPull = &parser.Pull{}
		raid.CurrentPull.StartTime = record.DateTime
		damageDone := make(map[parser.Entity]*parser.DamageDict)
		healDone := make(map[parser.Entity]*parser.HealDict)
		absDone := make(map[parser.Entity]float64)
		threatDone := make(map[string]float64)
		raid.CurrentPull.DamageDone = damageDone
		raid.CurrentPull.HealDone = healDone
		raid.CurrentPull.ThreatDone = threatDone
		raid.CurrentPull.AbsDone = absDone
		raid.InPull = true
		//log.Printf("%d Starting fight %s", record.LineNumber, raid.CurrentPull.StartTime)
	}
	if record.Effect.ActionID == globals.EXITCOMBATID && raid.InPull {
		//stop pull
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		showDamage(raid.CurrentPull)
		showDetails(raid.CurrentPull)
		//log.Printf("%d Stopping fight exited %s", record.LineNumber, raid.CurrentPull.StopTime)
	}
	if raid.InPull && record.Effect.ActionID == globals.DEATHID && record.Target.Name == globals.MainPlayerName {
		//stop pull, dead
		//log.Printf("%d %s DEAD at %s\n", record.LineNumber, record.Target.Name, record.DateTime)
		raid.InPull = false
		raid.CurrentPull.StopTime = record.DateTime
		raid.Pulls = append(raid.Pulls, *raid.CurrentPull)
		showDamage(raid.CurrentPull)
		showDetails(raid.CurrentPull)
		//log.Printf("%d Stopping fight DEAD %s", record.LineNumber, raid.CurrentPull.StopTime)
	}
}

func showDamage(pull *parser.Pull) {
	if pull.Target == "" && !globals.Debug {
		return
	}
	duration := pull.StopTime.Sub(pull.StartTime)
	seconds := duration.Seconds()

	doneMap := make(parser.StatsMap)
	receivedMap := make(parser.StatsMap)
	healMap := make(parser.StatsMap)
	for player, dmgDict := range pull.DamageDone {
		if !player.NPC && player.Name != "" {
			// Damage Done
			totalDamage := float64(0)
			for _, targetDmgDict := range dmgDict.TargetDamageDict {
				for _, abilityDmgDict := range targetDmgDict.Ability {
					totalDamage += float64(abilityDmgDict.Amount)
				}
			}
			doneMap[player] = totalDamage
		} else {
			// Damage Received
			for target, targetDmgDict := range dmgDict.TargetDamageDict {
				if !target.NPC {
					for _, abilityDmgDict := range targetDmgDict.Ability {
						receivedMap[target] += float64(abilityDmgDict.Amount)
					}
				}
			}
		}
	}
	for player, healDict := range pull.HealDone {
		if !player.NPC {
			totalHeal := float64(0)
			for _, targetHeaDict := range healDict.TargetHealDict {
				for _, abilityDmgDict := range targetHeaDict.Ability {
					totalHeal += float64(abilityDmgDict.Amount)
				}
			}
			healMap[player] = totalHeal
		}
	}
	log.Printf("==============================")
	log.Printf("STARTING FIGHT %s", pull.StartTime)
	log.Println(duration)
	log.Println(pull.Target)
	log.Println("------------------------------")
	log.Println("DAMAGE DONE")
	sorted := doneMap.Sort()
	for _, player := range sorted {
		amount := doneMap[player]
		log.Printf("%s: %.1f DPS", player.Name, amount/seconds)
	}
	log.Println("------------------------------")
	log.Println("DAMAGE RECEIVED")
	sorted = receivedMap.Sort()
	for _, player := range sorted {
		amount := receivedMap[player]
		log.Printf("%s: %.1f DPS", player.Name, amount/seconds)
	}
	log.Println("------------------------------")
	log.Println("HEAL DONE")
	sorted = healMap.Sort()
	for _, player := range sorted {
		amount := healMap[player]
		log.Printf("%s: %.1f HPS", player.Name, amount/seconds)
	}
	log.Println("------------------------------")
	log.Println("ABSORB DONE")
	maps := parser.StatsMap(pull.AbsDone)
	absDone := maps.Sort()
	for _, player := range absDone {
		amount := pull.AbsDone[player]
		log.Printf("%s: %.1f APS", player.Name, amount/seconds)
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

func showDetails(pull *parser.Pull) {
	if pull.Target == "" && !globals.Debug {
		return
	}
	mainPlayerName := globals.MainPlayerName
	log.Printf("Damage taken details for %s:", mainPlayerName)
	for player, dmgDict := range pull.DamageDone {
		for target, targetDmgDict := range dmgDict.TargetDamageDict {
			if target.Name == mainPlayerName {
				for _, abilityDmgDict := range targetDmgDict.Ability {
					log.Printf("%s %s %s %s %s Hits: %d, Amounts: %d Miss: %d, Dodge/Parry/Deflect: %d Shield: %d Resist: %d, Immune: %d\n", player.Name, player.ID, abilityDmgDict.ID, abilityDmgDict.Name, abilityDmgDict.DamageType, abilityDmgDict.Hits, abilityDmgDict.Amount, abilityDmgDict.MitigationDict.Miss, abilityDmgDict.MitigationDict.DodgeParryDeflect, abilityDmgDict.MitigationDict.Shield, abilityDmgDict.MitigationDict.Resist, abilityDmgDict.MitigationDict.Immune)
				}
			}

		}
	}
}
