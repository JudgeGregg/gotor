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
	UID  string
}

type Target struct {
	Name string
	NPC  bool
	ID   string
	UID  string
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
	Mitigation   string
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
}

type Pull struct {
	StartTime  time.Time
	EndTime    time.Time
	Players    []Player
	DamageDone map[string]uint64
	HealDone   map[string]uint64
}

type Player struct {
	Name string
	ID   string
}
