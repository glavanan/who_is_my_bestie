package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"strconv"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"math/rand"
)

type LeagueDto struct {
	Entries []Entry
	ParticipantId string
	Tier string
}

type Entry struct {
	PlayerOrTeamId string
}

func String_in_array(value string, ref []string) (bool){
	for _, elem := range ref {
		if value == elem {
			return true
		}
	}
	return false
}

var rank []string

func push_id_player(id_player string) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("fill").C("id_player")
	c.Remove(bson.M{"last": true})
	id, _ := strconv.ParseInt(id_player, 10, 64)

	err = c.Insert(bson.M{"summonerid": id, "timestamp": time.Now().Format(time.RFC850), "last" : true})
	if err != nil {
		panic(err)
	}
	session.Close()
}

func get_rdm_player_in_ligue(ligue map[string][]LeagueDto, id string)(string) {
	random_id := rand.Intn(len(ligue[id][0].Entries))
	push_id_player(ligue[id][0].Entries[random_id].PlayerOrTeamId)
	return ligue[id][0].Entries[random_id].PlayerOrTeamId
}

//will seek all league of all player ine "P"
//If the player is >= gold we had himi in the lst
//puis resortir un mec random de cette lst
//Save id, et si seulement une personne dans liste, piocher au hasard dans la league
func	give_next_id_player(p []ParticipantIdentity) (string){
	str := ""
	for key, elem := range p {
		if key == 0 {
			str = "" + strconv.FormatInt(elem.Player.SummonerId, 10)
		} else {
			str = str + "," + strconv.FormatInt(elem.Player.SummonerId, 10)
		}
	}
	time.Sleep(200 * time.Millisecond)
	resp, err := http.Get("https://euw.api.pvp.net/api/lol/euw/v2.5/league/by-summoner/" + str + "?api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var league map[string][]LeagueDto
	err = json.Unmarshal(body, &league)
	if err != nil {
		panic(err)
	}
	save := ""
	var lst_id []string
for id, elem := range league {
		if String_in_array(elem[0].Tier, rank) {
			lst_id = append(lst_id, elem[0].ParticipantId)
			save = id
		}
	}
	if len(lst_id) == 1 {
		return get_rdm_player_in_ligue(league, save)
	}
	random_id := rand.Intn(len(lst_id))
	push_id_player(lst_id[random_id])
	return lst_id[random_id]
}

func get_new_player_id(sumid string) (string){
	time.Sleep(200 * time.Millisecond)
	resp, err := http.Get("https://euw.api.pvp.net/api/lol/euw/v2.5/league/by-summoner/" + sumid + "?api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var league map[string][]LeagueDto
	err = json.Unmarshal(body, &league)
	if err != nil {
		panic(err)
	}
	return league[sumid][0].Entries[rand.Intn(len(league[sumid][0].Entries))].PlayerOrTeamId
}

