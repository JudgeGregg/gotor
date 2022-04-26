package globals

import "time"

const (
	ENTERCOMBATID = "836045448945489"
	EXITCOMBATID  = "836045448945490"
	AREAENTEREDID = "836045448953664"

	DAMAGEID = "836045448945501"
	HEALID   = "836045448945500"

	FOURPLAYERVETERAN = "836045448953657"
	FOURPLAYERMASTER  = "836045448953659"

	EIGHTPLAYERSTORY   = "836045448953650" //FIXME
	EIGHTPLAYERVETERAN = "836045448953652"
	EIGHTPLAYERMASTER  = "836045448953654" //FIXME

	STORY   = "Story"
	VETERAN = "Veteran"
	MASTER  = "Master"

	DEATHID  = "836045448945493"
	REVIVEID = "836045448945494"

	//Mitigations
	DODGE_OR_PARRY uint64 = iota
	MISS
	RESIST
	IMMUNE
	SHIELD
)

var Targets = map[string]string{
	//Karagga Veteran
	"2761191524925440": "Karagga the Unyielding",
	"2760637474144256": "Foreman Crusher",
	"2748401112317952": "G4-B3 Heavy Fabricator",
	"2760482855321600": "Jarg & Sorno",
	"2760487150288896": "Jarg & Sorno",
	"2624474125959168": "Bonethrasher",
	//Misc
	"2857785339412480": "Operations Training Dummy",
	//Athiss Veteran
	"172496007036928":  "Professor Ley'arsha",
	"2263129937412096": "The Beast of Vodall Kressh",
	"1172951273570304": "Prophet of Vodal",
	//Athiss Master
	"3158072272879616": "Professor Ley'arsha",
	"3158119517519872": "The Beast of Vodall Kressh",
	"3158123812487168": "Prophet of Vodal",
	"3247459132243968": "Ancient Abomination",
	//Hammer Station Master
	"3172554902601728": "Battlelord Kreshan",
	"3251019660132352": "Asteroid Beast",
	"3152162397880320": "Vorgan the Volcano",
	"3152166692847616": "Vorgan the Volcano",
	"3152158102913024": "Vorgan the Volcano",
	"3152170987814912": "DN-314 Tunneler",
	//Kaon Under Siege Master
	"2762260971782144": "Commander Lk'graagth",
	"2762269561716736": "Commander Lk'graagth",
	"2762265266749440": "Commander Lk'graagth",
	"2765357643202560": "Rakghoul Behemoth",
	//False Emperor Veteran
	"770963809501184":  "Darth Malgus",
	"770955219566592":  "HK-47",
	"1690314444111872": "Tregg the Destroyer",
	"770959514533888":  "Jindo Krey",
	//Taral V Master
	"2531994890141696": "General Edikar",
	//Cademimu Master
	"3210672737353728": "Officer Xander",
	"3210677032321024": "Captain Grimyk",
	"3210659852451840": "General Ortol",
}

var RaidStartDate = time.Time{}
