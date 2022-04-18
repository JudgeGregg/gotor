package main

import (
	"io"
	"os"

	"github.com/JudgeGregg/gotor/parser"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

const test = `[21:23:38.818] [@Zangyéf#686674938948221|(4613.04,4821.58,698.01,168.25)|(1/326800)] [] [] [AreaEntered {836045448953664}: Imperial Fleet {137438989504}] (he4002) <v7.0.0b>
[21:29:08.933] [] [@Zangyef#686674938948221|(1188.58,-4426.38,-777.19,154.51)|(337833/337833)] [Kolto Probe {814832605462528}] [RemoveEffect {836045448945478}: Kolto Probe {814832605462528}]
[21:31:11.920] [@Zangyef#686674938948221|(-14.12,4.94,0.25,-174.06)|(337833/337833)] [=] [Shoulder Cannon {3066482095292416}] [ApplyEffect {836045448945477}: Missile Loader {3066482095292718}]
[21:31:11.937] [Gamorrean Palace Guard {2470959109898240}:38983000004090|(-14.08,6.99,0.25,1.29)|(951368/964970)] [@Zangyef#686674938948221|(-14.12,4.94,0.25,-174.06)|(337833/337833)] [Stockstrike {811886257897472}] [ApplyEffect {836045448945477}: Damage {836045448945501}] (3637 ~0 kinetic {836045448940873} -shield {836045448945509} (3637 absorbed {836045448945511})) <3637>
[21:31:12.032] [@Shamiya#689102189850071|(1.00,13.61,-0.09,25.01)|(394741/394741)] [@Zangyef#686674938948221|(-14.12,4.94,0.25,-174.06)|(337833/337833)] [Revivification {808703687131136}] [ApplyEffect {836045448945477}: Heal {836045448945500}] (1690 ~0)
[21:31:12.431] [@Zangyef#686674938948221|(-14.12,4.94,0.25,-174.06)|(337833/337833)] [=] [Heat Screen {2841833830875136}] [ModifyCharges {836045448953666}: Heat Screen {2841833830875136}] (2 charges {836045448953667})
[21:31:12.451] [@Zangyef#686674938948221|(-14.12,4.94,0.25,-174.06)|(337833/337833)] [Gamorrean Bodyguard {2470954814930944}:38983000004025|(-15.25,-7.15,0.28,-147.53)|(247468/964970)] [Shocked {3960050041225514}] [ApplyEffect {836045448945477}: Damage {836045448945501}] (224 energy {836045448940874}) <2238>`

func main() {
	file, _ := os.Open("combat_2022-04-16_21_23_12_979036.txt")
	wInUTF8 := transform.NewReader(file, charmap.ISO8859_1.NewDecoder())
	str, _ := io.ReadAll(wInUTF8)
	//parser.Parse(string(test))
	parser.Parse(string(str))
	//fmt.Println("vim-go")

}
