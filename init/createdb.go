package main

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fmt"
	"net/http"
"io/ioutil"
"encoding/json"
"sort"
)

type Champ struct {
	ChampionId int `json:"championId" bson:"championId"`
	Name string `json:"name" bson:"name"`
	Key string `json:"key" bson:"key"`
}

type Stat struct {
	Champion1 int `json:"champion1" bson:"champion1"`
	Champion2 int `json:"champion2" bson:"champion2"`
	Games int `json:"games" bson:"games"`
	Win int `json:"win" bson:"win"`
}

type ChampionDto struct {
	Id	int
	Name	string
	Key	string
}

type ChampionListDto struct {
	Data	map[string]ChampionDto
}

var path_conf_file string = "/root/go/src/who_is_my_bestie/conf.ini"

type Conf struct {
	Api_key string
	First_id string
}

var conf Conf

func fill_conf() {
	file, err := ioutil.ReadFile(path_conf_file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
}


//le rendre automatique, par la requete au static data
//stocker id, name et key

func main() {
	// A modifier pour recuperer champions
	if conf.Api_key == "" {
		fill_conf()
	}
	resp, err := http.Get("https://global.api.pvp.net/api/lol/static-data/euw/v1.2/champion?api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var champion ChampionListDto //tu me fais une bonne classe STP ??
	err = json.Unmarshal(body, &champion)
	if err != nil {
		panic(err)
	}

	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	/// MDWDFWFW
	session.SetMode(mgo.Strong, true)
	str, err := session.DatabaseNames()
	fmt.Println(str)
	c := session.DB("champ").C("fiche")
	//we had a fiche for each champion (Helpful to get name and key for request img)
	for key,elem := range champion.Data {
		fmt.Println(key)
		fmt.Println(elem)
		champ := Champ{ChampionId: elem.Id, Name: elem.Name, Key: elem.Key}
		result, err := c.Find(champ).Count()
		if err != nil {
			panic(err)
		}
		if result <= 0 {
			c.Insert(champ)
			fmt.Println(elem)
		}
	}
	var keys []string
	for k := range champion.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	c = session.DB("champ").C("stat")
	for _, k := range keys { //Tu dois remplir les fiche de stat donc un autre boucle for pour parcourir tout les autres elem
		save := false
		for _, tmpk := range keys {
			if save == true {
				stat := Stat{Champion1: champion.Data[k].Id, Champion2: champion.Data[tmpk].Id, Games: 0, Win: 0}
				query := c.Find(bson.M{"champion1": champion.Data[k].Id, "champion2": champion.Data[tmpk].Id})
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
			}
			if save == false && k == tmpk {
				save = true
			}
		}
	}
	session.Close()
// On ajoute maintenant les fiche pour les differentes relation entre champions
}
