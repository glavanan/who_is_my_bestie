package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"strconv"
)

//store information of match
//Regarde si tu as bien toutes tes info pour pouvoir fill t'es fiches champ
type Match struct {
	ParticipantsIdentities []ParticipantIdentity
	Participants []Participant
	Teams []Team
}

type Team struct {
	TeamId int
	winner bool
}

type Participant struct {
	ChampionId int
	TeamId int
	ParticipantId int
}

type ParticipantIdentity struct {
	ParticipantId int
	Player PlayerStat
}

type PlayerStat struct {
	SummonerId	int64
}

// Struct for matchlist by id user
type Match_lst struct {
	Matches []Matches_info
	TotalGames int
	StartIndex int
	EndIndex int
}

type Matches_info struct {
	Timestamp int
	Champion int
	Region string
	Queue string
	Season string
	MatchId int
	Role string
	PlatformId string
	Lane string
}

type Matches_db struct {
	MatchId int `json:"matchid" bson:"matchid"`
	Timestamp int `json:"timestamp" bson:"timestamp"`
}

// Get the initial conf
type Conf struct {
	Api_key string
	First_id string
}

var conf Conf

var path_conf_file string = "/root/go/src/who_is_my_bestie/conf.ini"
//

func fill_conf() {
	file, err := ioutil.ReadFile(path_conf_file)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file))
	err = json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
}

func get_first_id() (sumid string) {
	if conf.First_id == "" {
		fill_conf()
	}
	return conf.First_id
}

func get_id_match(sumid string) (matchid string) {
	if conf.Api_key == "" {
		fill_conf()
	}
	resp, err := http.Get("https://euw.api.pvp.net/api/lol/euw/v2.2/matchlist/by-summoner/" + sumid + "?api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var match Match_lst
	err = json.Unmarshal(body, &match)
	if err != nil {
		panic(err)
	}
	i := 0
	session, err := mgo.Dial("127.0.0.1:27017")
	session.SetMode(mgo.Strong, true)
	c := session.DB("fill").C("id_match")
	match_db := Matches_db{MatchId: match.Matches[i].MatchId, Timestamp: match.Matches[i].Timestamp}
	result, err := c.Find(match_db).Count()
	for result != 0 {
		i++;
		match_db := Matches_db{MatchId: match.Matches[i].MatchId, Timestamp: match.Matches[i].Timestamp}
		result, err = c.Find(match_db).Count()
	}
	c.Insert(match_db)
	fmt.Println(match.Matches[i].MatchId)
	return strconv.Itoa(match.Matches[i].MatchId)
}

func get_match(matchid string) {
	if conf.Api_key == "" {
		fill_conf()
	}
	resp, err := http.Get ("https://euw.api.pvp.net/api/lol/euw/v2.2/match/" + matchid + "?includeTimeline=false&api_key=" + conf.Api_key)
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
}

func main() {
	get_id_match(get_first_id())
}


