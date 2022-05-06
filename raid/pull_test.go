package raid

import (
	"bufio"
	"os"
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
			{Name: "Gamorrean Palace Guard", ID: "2470959109898240", NPC: true}: {TargetDamageDict: map[parser.Target]*parser.TargetDamageDict{
				{Name: "Zangyef", ID: "686674938948221"}:     {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 6391}, "Close Attack": {Amount: 4000}}},
				{Name: "Tenna Aiken", ID: "689371682814222"}: {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 4000}}},
			}},
			{Name: "Zangyef", ID: "686674938948221"}: {TargetDamageDict: map[parser.Target]*parser.TargetDamageDict{
				{Name: "Gamorrean Palace Guard", ID: "2470959109898240", NPC: true}: {Ability: map[string]*parser.AbilityDict{"Ranged Attack": {Amount: 1000}}},
			}},
		},
	}}}
	for _, pullTest := range pullTestData {
		content, err := os.Open(pullTest.file)
		fileScanner := bufio.NewScanner(content)
		fileScanner.Split(bufio.ScanLines)
		lines := make(chan string)
		go func() {
			for fileScanner.Scan() {
				line := fileScanner.Text()
				lines <- line
			}
			close(lines)
		}()
		if err != nil {
			t.Logf(err.Error())
		}
		records := make(chan parser.Record)
		go parser.Parse(lines, records)
		raid_ := &parser.Raid{}
		for record := range records {
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
					t.Log(actor, target, ability)
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
