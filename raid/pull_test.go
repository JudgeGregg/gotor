package raid

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/JudgeGregg/gotor/parser"
)

func TestGetPull(t *testing.T) {
	var pullTestData = []struct {
		file string
		pull parser.Pull
	}{{"testdata/sample1.txt", parser.Pull{
		Target:    "",
		StartTime: time.Date(1, 1, 1, 21, 31, 01, 391000000, time.UTC),
		StopTime:  time.Date(1, 1, 1, 21, 31, 37, 979000000, time.UTC),
		DamageDone: map[parser.Actor]*parser.DamageDict{
			{Name: "Gamorrean Palace Guard", ID: "2470959109898240", UID: "38983000004090", NPC: true}: {TargetDamageDict: map[string]*parser.TargetDamageDict{
				"Zangyef":     {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 6391}}},
				"Tenna Aiken": {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 4000}}},
			}},
			{Name: "Zangyef", ID: "686674938948221"}: {TargetDamageDict: map[string]*parser.TargetDamageDict{
				"Gamorrean Palace Guard": {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 1000}}},
			}},
		},
	}}}
	for _, pullTest := range pullTestData {
		content, err := ioutil.ReadFile(pullTest.file)
		if err != nil {
			t.Logf(err.Error())
		}
		records := parser.Parse(string(content))
		raid_ := &parser.Raid{PlayersNumber: 1}
		for _, record := range records {
			HandleRecord(raid_, record)
		}
		if raid_.Pulls[0].Target != pullTest.pull.Target {
			t.Logf("Invalid pull Target: %s is not %s", raid_.Pulls[0].Target, pullTest.pull.Target)
			t.Fail()
		}
		if raid_.Pulls[0].StartTime != pullTest.pull.StartTime {
			t.Logf("Invalid pull StartTime: %s is not %s", raid_.Pulls[0].StartTime, pullTest.pull.StartTime)
			t.Fail()
		}
		if raid_.Pulls[0].StopTime != pullTest.pull.StopTime {
			t.Logf("Invalid pull StopTime: %s is not %s", raid_.Pulls[0].StopTime, pullTest.pull.StopTime)
			t.Fail()
		}
		for actor, dmgDict := range raid_.Pulls[0].DamageDone {
			for target, targetDmgDict := range dmgDict.TargetDamageDict {
				for ability, abilityDict := range targetDmgDict.Ability {
					amount := pullTest.pull.DamageDone[actor].TargetDamageDict[target].Ability[ability].Amount
					if abilityDict.Amount != amount {
						t.Logf("Invalid pull Amount: %d is not %d", abilityDict.Amount, amount)
						t.Fail()
					}
				}
			}
		}
	}
}
