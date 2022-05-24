package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/JudgeGregg/gotor/globals"
)

func Parse(lines chan string, records chan Record) {
	index := 0
	for line := range lines {
		index++
		if ignoreLine(line) {
			continue
		}
		record := getRecord(line)
		records <- record
	}
	close(records)
}

func GetRecord(line string) Record {
	return getRecord(line)
}

func ignoreLine(line string) bool {
	//Ignore GSF
	return strings.Contains(line, "::")
}

func getRecord(line string) Record {
	//Do we have an amount and threat ?
	var threat uint64
	var amount Amount
	split := strings.Split(line, " <")
	if len(split) > 1 {
		threatField := strings.Trim(split[1], ">")
		threat = getThreat(threatField)
		line = split[0]
	}
	split = strings.Split(line, "] (")
	if len(split) > 1 {
		amountField := split[1]
		// Trim last ")"
		amountField = strings.Trim(amountField, ")")
		amount = getAmount(amountField)
		line = split[0]
	}

	fields := strings.Split(line, "] [")
	var target Entity
	firstField := fields[0]
	time_ := getTime(firstField)
	actorField := fields[1]
	actor := getActor(actorField)
	targetField := fields[2]
	if targetField == "=" {
		target = actor
	} else {
		target = getTarget(targetField)
	}
	abilityField := fields[3]
	ability := getAbility(abilityField)
	effectField := fields[4]
	effect := getEffect(effectField)
	record := Record{DateTime: time_, Actor: actor, Target: target, Ability: ability, Effect: effect, Amount: amount, Threat: threat}
	return record
}

func getThreat(threat string) uint64 {
	var threatInt uint64 = 0
	isNotDigit := func(c rune) bool { return c < '0' || c > '9' }
	if (strings.IndexFunc(threat, isNotDigit)) == -1 {
		threatInt, _ := strconv.ParseUint(threat, 10, 64)
		return threatInt
	}
	return threatInt
}

func getTime(time_ string) time.Time {
	timeField := strings.ReplaceAll(time_, "[", "")
	res, _ := time.Parse("15:04:05", timeField)
	year, month, day := globals.RaidStartDate.Date()
	// time.Parse sets month and day to 1
	res = res.AddDate(year, int(month)-1, day-1)
	if res.Before(globals.RaidStartDate) {
		res = res.AddDate(0, 0, 1)
	}
	return res
}

func getActor(actorField string) Entity {
	actor := Entity{}
	if actorField == "" {
		return actor
	}
	actorFields := strings.Split(actorField, "|")
	nameField := actorFields[0]
	if strings.Contains(nameField, "#") {
		actor.NPC = false
		name := strings.SplitN(nameField, "#", 2)[0]
		name = strings.ReplaceAll(name, "@", "")
		id := strings.SplitN(nameField, "#", 2)[1]
		name = strings.ReplaceAll(name, "#", "")
		actor.Name = name
		actor.ID = id
		if strings.Contains(nameField, "/") {
			companionName := strings.Split(nameField, "/")[1]
			companionName = strings.Split(companionName, " {")[0]
			actor.Name = name + " (" + companionName + ")"
		}
	} else {
		actor.NPC = true
		name := strings.SplitN(nameField, " {", 2)[0]
		idField := strings.SplitN(nameField, " {", 2)[1]
		id := strings.Split(idField, "}")[0]
		actor.Name = name
		actor.ID = id
	}
	return actor
}

func getTarget(targetField string) Entity {
	target := Entity{}
	if targetField == "" {
		return target
	}
	targetFields := strings.Split(targetField, "|")
	nameField := targetFields[0]
	if strings.Contains(nameField, "#") {
		target.NPC = false
		name := strings.SplitN(nameField, "#", 2)[0]
		name = strings.ReplaceAll(name, "@", "")
		id := strings.SplitN(nameField, "#", 2)[1]
		name = strings.ReplaceAll(name, "#", "")
		target.Name = name
		target.ID = id
		if strings.Contains(nameField, "/") {
			companionName := strings.Split(nameField, "/")[1]
			companionName = strings.Split(companionName, " {")[0]
			target.Name = name + " (" + companionName + ")"
		}
	} else {
		target.NPC = true
		name := strings.SplitN(nameField, " {", 2)[0]
		idField := strings.SplitN(nameField, " {", 2)[1]
		id := strings.Split(idField, "}")[0]
		target.Name = name
		target.ID = id
	}
	return target
}

func getAbility(abilityField string) Ability {
	ability := Ability{}
	if abilityField == "" {
		return ability
	}
	abilitySplit := strings.Split(abilityField, " {")
	abilityName := abilitySplit[0]
	abilityID := abilitySplit[1]
	abilityID = strings.ReplaceAll(abilityID, "}", "")
	ability.Name = abilityName
	ability.ID = abilityID
	return ability
}

func getEffect(effectField string) Effect {
	var spec string
	var specID string
	effect := Effect{}
	if effectField == "" {
		return effect
	}
	effectSplit := strings.Split(effectField, "}")
	event := effectSplit[0]
	action := effectSplit[1]
	eventSplit := strings.Split(event, "{")
	event = eventSplit[0]
	event = strings.Trim(event, " :{}[]")
	eventID := eventSplit[1]
	eventID = strings.Trim(eventID, " :{}[]")
	actionSplit := strings.Split(action, "{")
	action = actionSplit[0]
	action = strings.Trim(action, " :{}[]")
	actionID := actionSplit[1]
	actionID = strings.Trim(actionID, " :{}[]")
	effect.Event = event
	effect.EventID = eventID
	effect.Action = action
	effect.ActionID = actionID

	if len(effectSplit) == 4 {
		spec = effectSplit[2]
		specSplit := strings.Split(spec, "{")
		spec = strings.Trim(specSplit[0], " {}[]/")
		specID = strings.Trim(specSplit[1], " {}[]")
		effect.Spec = spec
		effect.SpecID = specID
	}
	return effect
}

func getAmount(amountField string) Amount {
	amount := Amount{}
	if amountField == "" {
		return amount
	}
	if strings.HasSuffix(amountField, "-") {
		// "0" effect, eg poisoning or target already dead, ignore amount
		return amount
	}
	//Crit
	if strings.Contains(amountField, "*") {
		amount.Critical = true
		amountField = strings.ReplaceAll(amountField, "*", "")
	}
	split := strings.Split(amountField, " ")
	//Heal, Charges or Energy
	if len(split) <= 2 {
		return getAmountHealChargeEnergy(amount, split)
	}
	//Reflected Damage
	if strings.Contains(amountField, globals.REFLECTEDID) {
		return getAmountDamageReflected(amount, split)
	}
	//Damage
	if len(split) == 3 {
		return getAmountDamageRegular(amount, amountField, split)
	}
	if len(split) == 4 {
		//Regular Damage Altered
		return getAmountDamageAltered(amount, split)
	}
	//Bubbled Damage
	if strings.ContainsAny(amountField, "~") && strings.Contains(amountField, globals.ABSORBID) {
		return getAmountBubbled(amount, amountField, split)
	}
	//Shield but no bubble
	if strings.Contains(amountField, globals.ABSORBID) {
		amount.Mitigated = true
		amount.Mitigation = globals.SHIELD
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		amount.Amount = quantityInt
		amount.Effective = quantityInt
		amount.DamageType = split[1]
		damageTypeID := strings.ReplaceAll(split[2], "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
		return amount
	}
	panic("Parsing Error")
}

func getAmountBubbled(amount Amount, amountField string, split []string) Amount {
	if strings.Contains(amountField, globals.SHIELDID) {
		//Shield
		amount.Altered = true
		amount.Absorbed = true
		amount.Mitigated = true
		amount.Mitigation = globals.SHIELD
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		alteredQuantity := strings.ReplaceAll(split[1], "~", "")
		alteredQuantityInt, _ := strconv.ParseUint(alteredQuantity, 10, 64)
		amount.Amount = quantityInt
		amount.Effective = alteredQuantityInt
		amount.DamageType = split[2]
		damageTypeID := strings.ReplaceAll(split[3], "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
		return amount
	} else {
		//No shield
		amount.Altered = true
		amount.Absorbed = true
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		alteredQuantity := strings.ReplaceAll(split[1], "~", "")
		alteredQuantityInt, _ := strconv.ParseUint(alteredQuantity, 10, 64)
		amount.Amount = quantityInt
		amount.Effective = alteredQuantityInt
		amount.DamageType = split[2]
		damageTypeID := strings.ReplaceAll(split[3], "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
		return amount
	}
}

func getAmountDamageAltered(amount Amount, split []string) Amount {
	amount.Altered = true
	quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
	alteredQuantity := strings.ReplaceAll(split[1], "~", "")
	alteredQuantityInt, _ := strconv.ParseUint(alteredQuantity, 10, 64)
	amount.Amount = quantityInt
	amount.Effective = alteredQuantityInt
	amount.DamageType = split[2]
	damageTypeID := strings.ReplaceAll(split[3], "{", "")
	damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
	amount.DamageTypeID = damageTypeID
	return amount
}

func getAmountDamageRegular(amount Amount, amountField string, split []string) Amount {
	if strings.Contains(amountField, globals.PARRYID) || strings.Contains(amountField, globals.DEFLECTID) || strings.Contains(amountField, globals.DODGEID) {
		amount.Mitigated = true
		amount.Mitigation = globals.DODGE_PARRY_DEFLECT
		amount.Amount = 0
		amount.Effective = 0
	} else if strings.Contains(amountField, globals.MISSID) {
		amount.Mitigated = true
		amount.Mitigation = globals.MISS
		amount.Amount = 0
		amount.Effective = 0
	} else if strings.Contains(amountField, globals.RESISTID) {
		amount.Mitigated = true
		amount.Mitigation = globals.RESIST
		amount.Amount = 0
		amount.Effective = 0
	} else if strings.Contains(amountField, globals.IMMUNEID) {
		amount.Mitigated = true
		amount.Mitigation = globals.IMMUNE
		amount.Amount = 0
		amount.Effective = 0
	} else if strings.Contains(amountField, globals.SHIELDID) {
		amount.Mitigated = true
		amount.Mitigation = globals.SHIELD
		amount.Amount = 0
		amount.Effective = 0
	} else {
		//Regular Damage
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		amount.Amount = quantityInt
		amount.Effective = quantityInt
		amount.DamageType = split[1]
		damageTypeID := strings.ReplaceAll(split[2], "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
	}
	return amount
}

func getAmountDamageReflected(amount Amount, split []string) Amount {
	if len(split) == 4 {
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		amount.Amount = quantityInt
		amount.Effective = quantityInt
		amount.DamageType = split[1]
		damageTypeID := strings.Split(split[2], "(")[0]
		damageTypeID = strings.ReplaceAll(damageTypeID, "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
		return amount
	}
	if len(split) == 5 {
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		effective := strings.ReplaceAll(split[1], "~", "")
		effectiveInt, _ := strconv.ParseUint(effective, 10, 64)
		amount.Amount = quantityInt
		amount.Effective = effectiveInt
		amount.Altered = true
		amount.DamageType = split[2]
		damageTypeID := strings.Split(split[3], "(")[0]
		damageTypeID = strings.ReplaceAll(damageTypeID, "{", "")
		damageTypeID = strings.ReplaceAll(damageTypeID, "}", "")
		amount.DamageTypeID = damageTypeID
		return amount
	}
	return amount
}

func getAmountHealChargeEnergy(amount Amount, split []string) Amount {
	if len(split) == 1 {
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		amount.Amount = quantityInt
		amount.Effective = quantityInt
	} else {
		amount.Altered = true
		quantityInt, _ := strconv.ParseUint(split[0], 10, 64)
		amount.Amount = quantityInt
		effective := strings.ReplaceAll(split[1], "~", "")
		effectiveInt, _ := strconv.ParseUint(effective, 10, 64)
		amount.Effective = effectiveInt
	}
	return amount
}
