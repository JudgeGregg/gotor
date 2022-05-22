package parser

import (
	"testing"
	"time"

	"github.com/JudgeGregg/gotor/globals"
)

var line1 = `[21:23:39.060] [@Zangyef#686674938948221|(4613.04,4821.58,698.01,168.25)|(326800/337833)] [=] [Hunter's Boon {4503178720575488}] [ApplyEffect {836045448945477}: Hunter's Boon {4503178720575765}]`
var line2 = `[21:23:39.063] [::PEWPEWGSF::|(4613.04,4821.58,698.01,168.25)|(326800/337833)] [=] [Sprint {810670782152704}] [ApplyEffect {836045448945477}: Sprint {810670782152704}]`

var recordTestMap = map[string]Record{
	"[21:39:27.720] [@Kiss Assoka#689409546916090|(-148.01,-3.75,-23.43,166.41)|(319754/376614)] [=] [Ball Lightning {3408941312638976}] [Event {836045448945472}: AbilityActivate {836045448945479}]":                                                                                                                                      {Ability: Ability{Name: "Ball Lightning", ID: "3408941312638976"}, Actor: Entity{Name: "Kiss Assoka", ID: "689409546916090"}, Target: Entity{Name: "Kiss Assoka", ID: "689409546916090"}, Effect: Effect{Event: "Event", EventID: "836045448945472", Action: "AbilityActivate", ActionID: "836045448945479"}},
	"[21:39:27.842] [@Yyuukkii#689994892081602|(-154.48,-15.22,-23.43,-151.35)|(303424/363102)] [Vicious Manka Cat {2624461241057280}:38983000128695|(-149.66,-7.14,-23.43,-140.64)|(360016/1113496)] [Battering Assault {807720139620352}] [ApplyEffect {836045448945477}: Damage {836045448945501}] (343 energy {836045448940874}) <343>": {Ability: Ability{Name: "Battering Assault", ID: "807720139620352"}, Actor: Entity{Name: "Yyuukkii", ID: "689994892081602"}, Target: Entity{Name: "Vicious Manka Cat", ID: "2624461241057280", NPC: true}, Effect: Effect{Event: "ApplyEffect", EventID: "836045448945477", Action: "Damage", ActionID: "836045448945501"}, Threat: 343},
	"[21:48:33.385] [@Tenna Aiken#689371682814222|(-169.84,19.73,22.55,58.53)|(391728/391728)] [=] [Tactical Marker: Damage {3322956067373056}] [ApplyEffect {836045448945477}: Tactical Marker: Damage {3322956067373056}]":                                                                                                                {Ability: Ability{Name: "Tactical Marker: Damage", ID: "3322956067373056"}, Actor: Entity{Name: "Tenna Aiken", ID: "689371682814222"}, Target: Entity{Name: "Tenna Aiken", ID: "689371682814222"}, Effect: Effect{Event: "ApplyEffect", EventID: "836045448945477", Action: "Tactical Marker: Damage", ActionID: "3322956067373056"}},
	`[21:49:01.896] [@Zangyef#686674938948221|(-181.56,19.07,22.55,50.57)|(337833/337833)] [@Aurraxx#689419779008511|(-181.60,20.42,22.56,-9.39)|(377207/377207)] [] [Event {836045448945472}: TargetSet {836045448953668}]`:                                                                                                                {Actor: Entity{Name: "Zangyef", ID: "686674938948221"}, Target: Entity{Name: "Aurraxx", ID: "689419779008511"}, Effect: Effect{Event: "Event", EventID: "836045448945472", Action: "TargetSet", ActionID: "836045448953668"}},
	`[21:23:38.818] [@Zangyef#686674938948221|(4613.04,4821.58,698.01,168.25)|(1/326800)] [] [] [DisciplineChanged {836045448953665}: Powertech {16141007401395916385}/Shield Tech {2031339142381604}]`:                                                                                                                                     {Actor: Entity{Name: "Zangyef", ID: "686674938948221"}, Effect: Effect{Event: "DisciplineChanged", EventID: "836045448953665", Action: "Powertech", ActionID: "16141007401395916385", Spec: "Shield Tech", SpecID: "2031339142381604"}},
}

var amountTestMap = map[string]Amount{
	`5612 ~0 kinetic {836045448940873} -shield {836045448945509} (5612 absorbed {836045448945511})`: {Absorbed: true, Altered: true, Mitigated: true, Mitigation: globals.SHIELD, DamageType: "kinetic", DamageTypeID: globals.KINETICID, Amount: 5612, Effective: 0},
	`247 energy {836045448940874} -shield {836045448945509} (3450 absorbed {836045448945511})`:      {Mitigated: true, Mitigation: globals.SHIELD, DamageType: "energy", DamageTypeID: globals.ENERGYID, Amount: 247, Effective: 247},
	`0 -dodge {836045448945505}`:   {Mitigated: true, Mitigation: globals.DODGE_PARRY_DEFLECT, Amount: 0, Effective: 0},
	`0 -deflect {836045448945508}`: {Mitigated: true, Mitigation: globals.DODGE_PARRY_DEFLECT, Amount: 0, Effective: 0},
	`0 -resist {836045448945507}`:  {Mitigated: true, Mitigation: globals.RESIST, Amount: 0, Effective: 0},
	`0 -immune {836045448945506}`:  {Mitigated: true, Mitigation: globals.IMMUNE, Amount: 0, Effective: 0},
	`0 -miss {836045448945502}`:    {Mitigated: true, Mitigation: globals.MISS, Amount: 0, Effective: 0},
	`525 ~0 energy {836045448940874} -shield {836045448945509} (525 absorbed {836045448945511})`: {Absorbed: true, Altered: true, Mitigated: true, Mitigation: globals.SHIELD, DamageType: "energy", DamageTypeID: "836045448940874", Amount: 525, Effective: 0},
	`11053 kinetic {836045448940873}`:                                  {Amount: 11053, Effective: 11053, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
	`11053 ~11055 kinetic {836045448940873}`:                           {Altered: true, Amount: 11053, Effective: 11055, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
	`24906 kinetic {836045448940873}(reflected {836045448953649})`:     {Amount: 24906, Effective: 24906, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
	`247 ~0 energy {836045448940874} (247 absorbed {836045448945511})`: {Absorbed: true, Altered: true, DamageType: "energy", DamageTypeID: globals.ENERGYID, Amount: 247, Effective: 0},
	`11053* kinetic {836045448940873}`:                                 {Amount: 11053, Critical: true, Effective: 11053, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
	// Critical Overheal
	`8649* ~0`: {Amount: 8649, Altered: true, Critical: true, Effective: 0},
	`76845 ~76846 kinetic {836045448940873}(reflected {836045448953649})`: {Altered: true, Amount: 76845, Effective: 76846, DamageType: "kinetic", DamageTypeID: globals.KINETICID},
}

var timeTestMap = map[string]time.Time{"[21:39:27.720": time.Date(1, 1, 1, 21, 39, 27, 720000000, time.UTC)}

var threatTestMap = map[string]uint64{"<25>": 25, "<aa>": 0}

var areaEnteredTestMap = map[string]Record{
	`[23:25:48.788] [@Zangyef#686674938948221|(125.67,35.44,5.02,55.10)|(338278/338278)] [] [] [AreaEntered {836045448953664}: Karagga's Palace {833571547775669} 8 Player Veteran {836045448953652}] (he4002) <v7.0.0b>`: {Actor: Entity{Name: "Zangyef", ID: "686674938948221"}, Effect: Effect{Event: "AreaEntered", EventID: "836045448953664", Action: "Karagga's Palace", ActionID: "833571547775669", Spec: "8 Player Veteran", SpecID: "836045448953652"}}}

var actorTestMap = map[string]Entity{"[Gamorrean Palace Guard {2470959109898240}:38983000004090|(-13.81,7.33,0.24,20.99)|(56242/964970)]": {Name: "Gamorrean Palace Guard", ID: "2470959109898240"}}

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

func TestParse(t *testing.T) {
	records := make(chan Record)
	lines := make(chan string)
	go func() { lines <- line1; lines <- line2; close(lines) }()
	go Parse(lines, records)
	for range records {
		//Consume records
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

func TestGetThreat(t *testing.T) {
	for line, result := range threatTestMap {
		threat := getThreat(line)
		if threat != result {
			t.Logf("Invalid threat: %v is not %v", threat, result)
		}
	}
}

func TestGetActor(t *testing.T) {
	for line, result := range actorTestMap {
		actor := getActor(line)
		if actor != result {
			t.Logf("Invalid actor: %v is not %v", actor, result)
		}
	}
}
