package globals

import "time"

var Debug = false

const (
	ENTERCOMBATID = "836045448945489"
	EXITCOMBATID  = "836045448945490"
	AREAENTEREDID = "836045448953664"

	DAMAGEID = "836045448945501"
	HEALID   = "836045448945500"

	FOURPLAYERSTORY   = "836045448953658"
	FOURPLAYERVETERAN = "836045448953657"
	FOURPLAYERMASTER  = "836045448953659"

	EIGHTPLAYERSTORY   = "836045448953650"
	EIGHTPLAYERVETERAN = "836045448953652"
	EIGHTPLAYERMASTER  = "836045448953655"

	REFLECTEDID = "836045448953649"

	STORY   = "Story"
	VETERAN = "Veteran"
	MASTER  = "Master"

	DEATHID  = "836045448945493"
	REVIVEID = "836045448945494"

	PARRYID   = "836045448945503"
	DODGEID   = "836045448945505"
	DEFLECTID = "836045448945508"
	MISSID    = "836045448945502"
	RESISTID  = "836045448945507"
	IMMUNEID  = "836045448945506"

	ABSORBID = "836045448945511"
	SHIELDID = "836045448945509"

	//Mitigations

	DODGE_PARRY_DEFLECT uint64 = iota
	MISS
	RESIST
	IMMUNE
	SHIELD

	//Damage Types
	KINETICID   = "836045448940873"
	ENERGYID    = "836045448940874"
	ELEMENTALID = "836045448940875"
	INTERNALID  = "836045448940876"
)

var Targets = map[string]string{

	/// Misc
	"2857785339412480": "Operations Training Dummy",

	/// Flashpoints

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
	//Mandalorian Raiders Master
	"3158428755165184": "Braxx the Bloodhound",
	"3172567787503616": "Mavrix Varad",
	"3250538623795200": "Gil",
	"3158402985361408": "Ezeraline",
	//False Emperor Master
	"2511688284766208": "Darth Malgus",
	"2511701169668096": "HK-47",
	"2511709759602688": "Jindo Krey",
	"2511761299210240": "Tregg the Destroyer",
	//Czerka Core Meltdown Master
	"3279293429841920": "Enhanced Duneclaw",
	"3279289134874624": "Enhanced Vrblther",
	"3279310609711104": "The Vigilant",
	//Czerka XXX
	//The Foundry
	//The Esseles
	//Nathema Conspiracy Master
	//Ruins of Nul Master
	//Spirit of Vengeance Master
	//Blood Hunt Master
	"3507051250581504": "Shae Vizla",
	"3507029775745024": "Valk & Jos",
	"3507034070712320": "Valk & Jos",
	//Boarding Party Master
	"2514497193377792": "HXI-54 Juggernaut",
	"2514522963181568": "Sakan Do'nair",
	"2514531553116160": "Commander Jorland",

	/// Operations

	//Eternity Vault 8 Veteran
	"2289823159156736": "Soa",
	"2034526008115200": "Gharj",
	"2034573252755456": "Annihilation Droid XRR-3",
	//Karagga 8 Veteran
	"2761191524925440": "Karagga the Unyielding",
	"2760637474144256": "Foreman Crusher",
	"2748401112317952": "G4-B3 Heavy Fabricator",
	"2760482855321600": "Jarg & Sorno",
	"2760487150288896": "Jarg & Sorno",
	"2624474125959168": "Bonethrasher",
	//Asation 8 Veteran
	"3010424182145024": "The Writhing Horror",
	"2994837745827840": "Operator IX",
	"3013327580037120": "The Dread Guards",
	"3013121421606912": "Kephess the Undying",
	"3025220344479744": "The Terror from Beyond",
	//Darvannis 8 Veteran
	//The Dread Fortress 8 Master
	"3303031714086912": "Nefra, Who Bars the Way",
	"3153558262251520": "Dash'Roode",
	"3152458750623744": "Titan 6",
	"3154563284598784": "Thrasher",
	"3157548286869504": "Operations Chief",
	"3050663730741248": "Olok the Shadow",
	"3156895451840512": "Cartel Warlords",
	"3054400352288768": "Cartel Warlords",
	"3054408942223360": "Cartel Warlords",
	"3054404647256064": "Cartel Warlords",
	"3067057620910080": "Dread Master Styrak",
	"3152407211016192": "Dread Master Styrak",
}

var RaidStartDate = time.Time{}
var MainPlayerName = ""
