package parser

import (
	"sort"
	"time"
)

type Record struct {
	DateTime time.Time
	Actor    Entity
	Target   Entity
	Ability  Ability
	Effect   Effect
	Amount   Amount
	Threat   uint64
}

type Entity struct {
	Name string
	NPC  bool
	ID   string
}

type Actor Entity
type Target Entity

type Ability struct {
	Name string
	ID   string
}

type Effect struct {
	Event    string
	EventID  string
	Action   string
	ActionID string
	Spec     string
	SpecID   string
}

type Amount struct {
	Altered      bool
	Mitigated    bool
	Mitigation   uint64
	DamageType   string
	DamageTypeID string
	Absorbed     bool
	Critical     bool
	Amount       uint64
	Effective    uint64
}

type Raid struct {
	Pulls       []Pull
	InPull      bool
	CurrentPull *Pull
	Difficulty  string
	BubblerMap  BubblerMap
}

type Pull struct {
	StartTime  time.Time
	StopTime   time.Time
	Target     string
	DamageDone map[Entity]*DamageDict
	HealDone   map[Entity]*HealDict
	ThreatDone map[string]float64
	AbsDone    map[Entity]float64
}

type DamageDict struct {
	ID               string
	Name             string
	TargetDamageDict map[Entity]*TargetDamageDict
}

type TargetDamageDict struct {
	ID      string
	Name    string
	Ability map[string]*AbilityDict
}

type HealDict struct {
	ID             string
	Name           string
	TargetHealDict map[Entity]*TargetHealDict
}

type TargetHealDict struct {
	ID      string
	Name    string
	Ability map[string]*AbilityDict
}

type AbilityDict struct {
	ID             string
	Name           string
	DamageType     string
	Hits           uint64
	Critical       uint64
	Amount         uint64
	MitigationDict MitigationDict
}

type MitigationDict struct {
	Miss              uint64
	Resist            uint64
	DodgeParryDeflect uint64
	Shield            uint64
	Immune            uint64
}

type BubblerMap map[Entity]Bubbler

type Bubbler struct {
	CurrentBubbler  Entity
	PreviousBubbler Entity
}

type StatsMap map[Entity]float64

func (sm *StatsMap) Sort() []Entity {
	sortedValues := []Entity{}
	map_ := map[Entity]float64(*sm)
	for name := range *sm {
		sortedValues = append(sortedValues, name)
	}
	sort.SliceStable(sortedValues, func(i, j int) bool {
		return map_[sortedValues[i]] > map_[sortedValues[j]]
	})
	return sortedValues
}
