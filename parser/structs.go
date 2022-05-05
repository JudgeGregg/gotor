package parser

import "time"

type Record struct {
	LineNumber int
	DateTime   time.Time
	Actor      Actor
	Target     Target
	Ability    Ability
	Effect     Effect
	Amount     Amount
	Threat     uint64
}

type Actor struct {
	Name string
	NPC  bool
	ID   string
}

type Target struct {
	Name string
	NPC  bool
	ID   string
}

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
}

type Pull struct {
	StartTime  time.Time
	StopTime   time.Time
	Target     string
	DamageDone map[Actor]*DamageDict
	HealDone   map[Actor]*HealDict
	ThreatDone map[string]uint64
}

type DamageDict struct {
	ID               string
	Name             string
	TargetDamageDict map[Target]*TargetDamageDict
}

type TargetDamageDict struct {
	ID      string
	Name    string
	Ability map[string]*AbilityDict
}

type HealDict struct {
	ID             string
	Name           string
	TargetHealDict map[Target]*TargetHealDict
}

type TargetHealDict struct {
	ID      string
	Name    string
	Ability map[string]*AbilityDict
}

type AbilityDict struct {
	ID           string
	Name         string
	Hits         uint64
	Critical     uint64
	Amount       uint64
	Miss         uint64
	Resist       uint64
	Immune       uint64
	DodgeOrParry uint64
	Shield       uint64
}
