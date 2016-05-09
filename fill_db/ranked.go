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
	ParticipantId string
	Tier string
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

func	give_next_id_player(p []ParticipantIdentity) (string){
	str := ""
	for key, elem := range p {
		if key == 0 {
			str = "" + strconv.FormatInt(elem.Player.SummonerId, 10)
		} else {
			str = str + "," + strconv.FormatInt(elem.Player.SummonerId, 10)
		}
	}
	resp, err := http.Get("https://euw.api.pvp.net/api/lol/euw/v2.5/league/by-summoner/" + str + "?api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var league map[string][]LeagueDto
	err = json.Unmarshal(body, &league)
	if err != nil {
		panic(err)
	}
	var lst_id []string
	for _, elem := range league {
		if String_in_array(elem[0].Tier, rank) {
			lst_id = append(lst_id, elem[0].ParticipantId)
		}
	}
	random_id := rand.Intn(len(lst_id))
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("fill").C("id_player")
	c.Remove(bson.M{"last": true})
	id, _ := strconv.ParseInt(lst_id[random_id], 10, 64)

	err = c.Insert(bson.M{"summonerid": id, "timestamp": time.Now().Format(time.RFC850), "last" : true})
	if err != nil {
		panic(err)
	}
	session.Close()
	return lst_id[random_id]


}
