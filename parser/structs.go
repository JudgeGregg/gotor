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
