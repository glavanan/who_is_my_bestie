package main

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fmt"
)

type Champ struct {
	ChampionId int `json:"championId" bson:"championId"`
	Name string `json:"name" bson:"name"`
}

type Stat struct {
	Champion1 int `json:"champion1" bson:"champion1"`
	Champion2 int `json:"champion2" bson:"champion2"`
	Games int `json:"games" bson:"games"`
	Win int `json:"win" bson:"win"`
}


func main() {
	lst_champ := []Champ{
	{1, "Annie"},
	{2, "Olaf"},
	{3, "Galio"},
	{4, "TwistedFate"},
	{5, "XinZhao"},
	{6, "Urgot"},
	{7, "Leblanc"},
	{8, "Vladimir"},
	{9, "FiddleSticks"},
	{10, "Kayle"},
	{11, "MasterYi"},
	{12, "Alistar"},
	{13, "Ryze"},
	{14, "Sion"},
	{15, "Sivir"},
	{16, "Soraka"},
	{17, "Teemo"},
	{18, "Tristana"},
	{19, "Warwick"},
	{20, "Nunu"},
	{21, "MissFortune"},
	{22, "Ashe"},
	{23, "Tryndamere"},
	{24, "Jax"},
	{25, "Morgana"},
	{26, "Zilean"},
	{27, "Singed"},
	{28, "Evelynn"},
	{29, "Twitch"},
	{30, "Karthus"},
	{31, "Chogath"},
	{32, "Amumu"},
	{33, "Rammus"},
	{34, "Anivia"},
	{35, "Shaco"},
	{36, "DrMundo"},
	{37, "Sona"},
	{38, "Kassadin"},
	{39, "Irelia"},
	{40, "Janna"},
	{41, "Gangplank"},
	{42, "Corki"},
	{43, "Karma"},
	{44, "Taric"},
	{45, "Veigar"},
	{48, "Trundle"},
	{50, "Swain"},
	{51, "Caitlyn"},
	{53, "Blitzcrank"},
	{54, "Malphite"},
	{55, "Katarina"},
	{56, "Nocturne"},
	{57, "Maokai"},
	{58, "Renekton"},
	{59, "JarvanIV"},
	{60, "Elise"},
	{61, "Orianna"},
	{62, "MonkeyKing"},
	{63, "Brand"},
	{64, "LeeSin"},
	{67, "Vayne"},
	{68, "Rumble"},
	{69, "Cassiopeia"},
	{72, "Skarner"},
	{74, "Heimerdinger"},
	{75, "Nasus"},
	{76, "Nidalee"},
	{77, "Udyr"},
	{78, "Poppy"},
	{79, "Gragas"},
	{80, "Pantheon"},
	{81, "Ezreal"},
	{82, "Mordekaiser"},
	{83, "Yorick"},
	{84, "Akali"},
	{85, "Kennen"},
	{86, "Garen"},
	{89, "Leona"},
	{90, "Malzahar"},
	{92, "Riven"},
	{91, "Talon"},
	{96, "KogMaw"},
	{98, "Shen"},
	{99, "Lux"},
	{101, "Xerath"},
	{102, "Shyvana"},
	{103, "Ahri"},
	{104, "Graves"},
	{105, "Fizz"},
	{106, "Volibear"},
	{107, "Rengar"},
	{110, "Varus"},
	{111, "Nautilus"},
	{112, "Viktor"},
	{113, "Sejuani"},
	{114, "Fiora"},
	{115, "Ziggs"},
	{117, "Lulu"},
	{119, "Draven"},
	{120, "Hecarim"},
	{121, "Khazix"},
	{122, "Darius"},
	{126, "Jayce"},
	{127, "Lissandra"},
	{131, "Diana"},
	{133, "Quinn"},
	{134, "Syndra"},
	{136, "AurelionSol"},
	{143, "Zyra"},
	{150, "Gnar"},
	{154, "Zac"},
	{157, "Yasuo"},
	{161, "Velkoz"},
	{201, "Braum"},
	{202, "Jhin"},
	{203, "Kindred"},
	{222, "Jinx"},
	{223, "TahmKench"},
	{236, "Lucian"},
	{238, "Zed"},
	{245, "Ekko"},
	{254, "Vi"},
	{266, "Aatrox"},
	{267, "Nami"},
	{268, "Azir"},
	{412, "Thresh"},
	{429, "Kalista"},
	{421, "RekSai"},
	{420, "Illaoi"},
	{432, "Bard"},
	}
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Strong, true)
	str, err := session.DatabaseNames()
	fmt.Println(str)
	c := session.DB("champ").C("fiche")
	for _,elem := range lst_champ {
		result, err := c.Find(elem).Count()
		if result <= 0 {
			err = c.Insert(elem)
			if err != nil {
				panic(err)
			}
		} else if result > 1 {
			err = c.Remove(elem)
			if err != nil {
				panic(err)
			}
		}
		aatrox := Champ{}
		err = c.Find(elem).One(&aatrox)
		if err != nil {
			panic(err)
		}
		fmt.Println(aatrox)
	}
	c = session.DB("champ").C("stat")
	for key, elem := range lst_champ {
		key++
		for len(lst_champ) > key {
			stat := Stat{Champion1: elem.ChampionId, Champion2: lst_champ[key].ChampionId, Games: 0, Win: 0}
			query := c.Find(bson.M{"champion1": elem.ChampionId, "champion2": lst_champ[key].ChampionId})
			result, _ := query.Count()

			if err != nil {
				panic(err)
}
			if result > 0 {
				fmt.Println("%i, %i find", stat.Champion1, stat.Champion2)
			} else {
				fmt.Println("%i, %i create", stat.Champion1, stat.Champion2)
				c.Insert(stat)
			}
			key++
		}
	}
// On ajoute maintenant les fiche pour les differentes relation entre champions

}
