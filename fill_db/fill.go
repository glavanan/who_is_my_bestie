package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"encoding/json"
	"strconv"
)

//store information of match

//struct summid + timestamp

type Summoner struct {
	Summonerid int64 `json:"summonerid" bson:"summonerid"`
	Timestamp int64 `json:"timestamp" bson:"timestamp"`
	Last bool `json:"last" bason:"last"`
}

//Struct des fiches champ
type Stat struct {
	Champion1 int `json:"champion1" bson:"champion1"`
	Champion2 int `json:"champion2" bson:"champion2"`
	Games int `json:"games" bson:"games"`
	Win int `json:"win" bson:"win"`
}

//Regarde si tu as bien toutes tes info pour pouvoir fill t'es fiches champ
type Match struct {
	ParticipantIdentities []ParticipantIdentity
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

func get_id_player() (sumid string) {
	if conf.Api_key == "" {
		fill_conf()
	}
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("fill").C("id_player")
	var result Summoner
	c.Find(bson.M{"last": true}).One(&result)
	session.Close()
	fmt.Println(result)
	if result.Summonerid == 0 {
		return conf.First_id
	}
	return strconv.FormatInt(result.Summonerid, 10)
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
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("fill").C("id_match")
	c.EnsureIndex(mgo.Index{Key: []string{"matchid"}, Unique:true})
	result := 1
	for len(match.Matches) > i && (result != 0 || (match.Matches[i].Queue != "TEAM_BUILDER_DRAFT_RANKED_5x5" || match.Matches[i].Season != "SEASON2016")) {
		result, err = c.Find(bson.M{"matchid" : match.Matches[i].MatchId}).Count()
		i++;
	}
	fmt.Println(i, " game on :", len(match.Matches))
	if result > 0 || i >= len(match.Matches) {
		return strconv.Itoa(match.Matches[0].MatchId)
	}
	c.Insert(bson.M{"matchid" : match.Matches[i - 1].MatchId})
	session.Close()
	return strconv.Itoa(match.Matches[i - 1].MatchId)
}

func fill_db(teamw []int, teaml []int) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("champ").C("stat")
	var stat Stat
	for i, elem := range teamw {
		for j := i + 1; j < len(teamw); j++ {
			if elem < teamw[j] {
				c.Find(bson.M{"champion1": elem, "champion2": teamw[j]}).One(&stat)
				c.Update(bson.M{"champion1": elem,"champion2": teamw[j]}, bson.M{"$set": bson.M{"games": stat.Games + 1,"win": stat.Win + 1}})
			} else {
				c.Find(bson.M {	"champion1": teamw[j],"champion2": elem}).One(&stat)
				c.Update(bson.M{"champion1": teamw[j],"champion2": elem},bson.M{"$set": bson.M{"games": stat.Games + 1,"win": stat.Win + 1}})
			}
		}
	}
	for i, elem := range teaml {
		for j := i + 1; j < len(teaml); j++ {
			if elem < teaml[j] {
				c.Find(bson.M{"champion1": elem, "champion2": teaml[j]}).One(&stat)
				c.Update(bson.M{"champion1": elem, "champion2": teaml[j]}, bson.M{"$set": bson.M{"games": stat.Games + 1, "win": stat.Win}})
			} else {
				c.Find(bson.M{"champion1": teaml[j], "champion2": elem}).One(&stat)
				c.Update(bson.M{"champion1": teaml[j], "champion2": elem}, bson.M{"$set": bson.M{"games": stat.Games + 1, "win": stat.Win}})
			}
		}
	}
	session.Close()
}

func get_match(matchid string) {
	fmt.Println("Match id : ", matchid)
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
	var match Match
	err = json.Unmarshal(body, &match)
	var team1, team2 []int
	for i := 0; i < len(match.Participants); i++ {
		if match.Participants[i].TeamId == 100 {
			team1 = append(team1, match.Participants[i].ChampionId)
		} else {
			team2 = append(team2, match.Participants[i].ChampionId)
		}
	}
	if match.Teams[0].winner {
		fill_db(team1, team2)
	} else {
		fill_db(team2, team1)
	}
	give_next_id_player(match.ParticipantIdentities)
}

func main() {
	rank = []string {"CHALLENGER", "MASTER", "DIAMOND", "PLATINUM", "GOLD"}
	for i := 0; i < 3; i++ {
		get_match(get_id_match(get_id_player()))
	}
}


