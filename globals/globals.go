package globals

import "time"

var Debug = false

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
	//Asation Veteran
	"3010424182145024": "The Writhing Horror",
	"2994837745827840": "Operator IX",
	"3013327580037120": "The Dread Guards",
	"3013121421606912": "Kephess the Undying",
	"3025220344479744": "The Terror from Beyond",
	//Misc
	"2857785339412480": "Operations Training Dummy",
	//Athiss Veteran
	"172496007036928":  "Professor Ley'arsha",
	"2263129937412096": "The Beast of Vodall Kressh",
	"1172951273570304": "Prophet of Vodal",
	//False Emperor Veteran
	"770963809501184":  "Darth Malgus",
	"770955219566592":  "HK-47",
	"1690314444111872": "Tregg the Destroyer",
	"770959514533888":  "Jindo Krey",
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
	//Taral V Master
	"2531994890141696": "General Edikar",
	//Cademimu Master
	"3210672737353728": "Officer Xander",
	"3210677032321024": "Captain Grimyk",
	"3210659852451840": "General Ortol",
	//Traitor Amongst the Chiss Master
	"4128382694457344": "Guardian Droid",
	"4135314771673088": "Syndic Zenta",
	"4132209510318080": "Vaiss",
	//Crisis on Umbara Master
	"4111512062918656": "Shadow Assassin Elli-Vaa & Technician Canni",
	"4111507767951360": "Shadow Assassin Elli-Vaa & Technician Canni",
	"4112779078270976": "Umbaran Spider Tank",
	"4112173487882240": "Vixian Mauler",
	//The Battle of Ilum Master
	"2511473536401408": "Drinda-Zel & Velasu Graege",
	"2511490716270592": "Drinda-Zel & Velasu Graege",
	"2511486421303296": "Krel Thak",
	"2511469241434112": "Darth Serevin",
	//The Red Reaper Master
	"2482594176303104": "Lord Kherus",
	"2482692960550912": "SV-3 Eradictor",
	"1478800189685760": "Darth Ikoral",
	//Depths of Manaan Master
	"3505251659284480": "Saisiri",
	"3505260249219072": "Ortuno",
	"3505268839153664": "Stivastin",
	"3506922401562624": "M2-AUX Foreman",
	//Nathema Conspiracy Master
	//Ruins of Nul Master
	//Spirit of Vengeance Master
}

var RaidStartDate = time.Time{}
