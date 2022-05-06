package parser

import (
	"testing"
	"time"

	"github.com/JudgeGregg/gotor/globals"
)

var recordTestMap = map[string]Record{
	"[21:39:27.720] [@Kiss Assoka#689409546916090|(-148.01,-3.75,-23.43,166.41)|(319754/376614)] [=] [Ball Lightning {3408941312638976}] [Event {836045448945472}: AbilityActivate {836045448945479}]":                                                                                                                                      {Ability: Ability{Name: "Ball Lightning", ID: "3408941312638976"}, Actor: Actor{Name: "Kiss Assoka", ID: "689409546916090"}, Target: Target{Name: "Kiss Assoka", ID: "689409546916090"}, Effect: Effect{Event: "Event", EventID: "836045448945472", Action: "AbilityActivate", ActionID: "836045448945479"}},
	"[21:39:27.842] [@Yyuukkii#689994892081602|(-154.48,-15.22,-23.43,-151.35)|(303424/363102)] [Vicious Manka Cat {2624461241057280}:38983000128695|(-149.66,-7.14,-23.43,-140.64)|(360016/1113496)] [Battering Assault {807720139620352}] [ApplyEffect {836045448945477}: Damage {836045448945501}] (343 energy {836045448940874}) <343>": {Ability: Ability{Name: "Battering Assault", ID: "807720139620352"}, Actor: Actor{Name: "Yyuukkii", ID: "689994892081602"}, Target: Target{Name: "Vicious Manka Cat", ID: "2624461241057280", NPC: true}, Effect: Effect{Event: "ApplyEffect", EventID: "836045448945477", Action: "Damage", ActionID: "836045448945501"}, Threat: 343},
	"[21:48:33.385] [@Tenna Aiken#689371682814222|(-169.84,19.73,22.55,58.53)|(391728/391728)] [=] [Tactical Marker: Damage {3322956067373056}] [ApplyEffect {836045448945477}: Tactical Marker: Damage {3322956067373056}]":                                                                                                                {Ability: Ability{Name: "Tactical Marker: Damage", ID: "3322956067373056"}, Actor: Actor{Name: "Tenna Aiken", ID: "689371682814222"}, Target: Target{Name: "Tenna Aiken", ID: "689371682814222"}, Effect: Effect{Event: "ApplyEffect", EventID: "836045448945477", Action: "Tactical Marker: Damage", ActionID: "3322956067373056"}},
	`[21:49:01.896] [@Zangyef#686674938948221|(-181.56,19.07,22.55,50.57)|(337833/337833)] [@Aurraxx#689419779008511|(-181.60,20.42,22.56,-9.39)|(377207/377207)] [] [Event {836045448945472}: TargetSet {836045448953668}]`:                                                                                                                {Actor: Actor{Name: "Zangyef", ID: "686674938948221"}, Target: Target{Name: "Aurraxx", ID: "689419779008511"}, Effect: Effect{Event: "Event", EventID: "836045448945472", Action: "TargetSet", ActionID: "836045448953668"}},
	`[21:23:38.818] [@Zangyef#686674938948221|(4613.04,4821.58,698.01,168.25)|(1/326800)] [] [] [DisciplineChanged {836045448953665}: Powertech {16141007401395916385}/Shield Tech {2031339142381604}]`:                                                                                                                                     {Actor: Actor{Name: "Zangyef", ID: "686674938948221"}, Effect: Effect{Event: "DisciplineChanged", EventID: "836045448953665", Action: "Powertech", ActionID: "16141007401395916385", Spec: "Shield Tech", SpecID: "2031339142381604"}},
}

var amountTestMap = map[string]Amount{
	`5612 ~0 kinetic {836045448940873} -shield {836045448945509} (5612 absorbed {836045448945511})`: {Altered: true, Mitigated: true, Mitigation: globals.SHIELD, DamageType: "kinetic", DamageTypeID: globals.KINETICID, Amount: 5612, Effective: 0},
	`247 energy {836045448940874} -shield {836045448945509} (3450 absorbed {836045448945511})`:      {Mitigated: true, Mitigation: globals.SHIELD, DamageType: "energy", DamageTypeID: globals.ENERGYID, Amount: 247, Effective: 247},
	`0 -dodge {836045448945505}`:   {Mitigated: true, Mitigation: globals.DODGE_PARRY_DEFLECT, Amount: 0, Effective: 0},
	`0 -deflect {836045448945508}`: {Mitigated: true, Mitigation: globals.DODGE_PARRY_DEFLECT, Amount: 0, Effective: 0},
	`525 ~0 energy {836045448940874} -shield {836045448945509} (525 absorbed {836045448945511})`: {Altered: true, Mitigated: true, Mitigation: globals.SHIELD, DamageType: "energy", DamageTypeID: "836045448940874", Amount: 525, Effective: 0},
	`0 -immune {836045448945506}`:                                  {Mitigated: true, Mitigation: globals.IMMUNE, Amount: 0, Effective: 0},
	`11053 kinetic {836045448940873}`:                              {Amount: 11053, Effective: 11053, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
	`24906 kinetic {836045448940873}(reflected {836045448953649})`: {Amount: 24906, Effective: 24906, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
}

var timeTestMap = map[string]time.Time{"[21:39:27.720": time.Date(1, 1, 1, 21, 39, 27, 720000000, time.UTC)}

var areaEnteredTestMap = map[string]Record{
	`[23:25:48.788] [@Zangyef#686674938948221|(125.67,35.44,5.02,55.10)|(338278/338278)] [] [] [AreaEntered {836045448953664}: Karagga's Palace {833571547775669} 8 Player Veteran {836045448953652}] (he4002) <v7.0.0b>`: {Actor: Actor{Name: "Zangyef", ID: "686674938948221"}, Effect: Effect{Event: "AreaEntered", EventID: "836045448953664", Action: "Karagga's Palace", ActionID: "833571547775669", Spec: "8 Player Veteran", SpecID: "836045448953652"}}}

func TestGetRecord(t *testing.T) {
	for line, result := range recordTestMap {
		record := getRecord(line)
		if record.Ability != result.Ability {
			t.Logf("Invalid record ability: %v is not %v", record.Ability, result.Ability)
			t.Fail()
		}
		if record.Actor != result.Actor {
			t.Logf("Invalid record actor: %v is not %v", record.Actor, result.Actor)
			t.Fail()
		}
		if record.Target != result.Target {
			t.Logf("Invalid record target: %v is not %v", record.Target, result.Target)
			t.Fail()
		}
		if record.Effect != result.Effect {
			t.Logf("Invalid record effect: %v is not %v", record.Effect, result.Effect)
			t.Fail()
		}
		if record.Threat != result.Threat {
			t.Logf("Invalid record threat: %v is not %v", record.Threat, result.Threat)
			t.Fail()
		}
	}
}

func TestGetAmount(t *testing.T) {
	for line, result := range amountTestMap {
		amount := getAmount(line)
		if amount != result {
			t.Logf("Invalid amount: %v is not %v", amount, result)
			t.Fail()
		}
	}
}

func TestGetTime(t *testing.T) {
	for line, result := range timeTestMap {
		time_ := getTime(line)
		if time_ != result {
			t.Logf("Invalid time: %v is not %v", time_, result)
			t.Fail()
		}
	}
}

func TestGetAreaEntered(t *testing.T) {
	for line, result := range areaEnteredTestMap {
		record := getRecord(line)
		if record.Ability != result.Ability {
			t.Logf("Invalid record ability: %v is not %v", record.Ability, result.Ability)
			t.Fail()
		}
		if record.Actor != result.Actor {
			t.Logf("Invalid record actor: %v is not %v", record.Actor, result.Actor)
			t.Fail()
		}
		if record.Target != result.Target {
			t.Logf("Invalid record target: %v is not %v", record.Target, result.Target)
			t.Fail()
		}
		if record.Effect != result.Effect {
			t.Logf("Invalid record effect: %v is not %v", record.Effect, result.Effect)
			t.Fail()
		}
		if record.Threat != result.Threat {
			t.Logf("Invalid record threat: %v is not %v", record.Threat, result.Threat)
			t.Fail()
		}
	}
}
