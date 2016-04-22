package main

import (
	"fmt"
	"net/http"
//	"io/ioutil"
	"html/template"
	"strings"
	"log"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
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

func print_ratio(champid int, w http.ResponseWriter) {
	fmt.Fprint(w, "id : ")
	fmt.Fprint(w, champid)
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("champ").C("stat")
	var stat []Stat
	c.Find(bson.M{"champion1": champid}).All(&stat)
	for _, elem := range stat {
		fmt.Fprint(w, elem)
		fmt.Fprint(w, "\n")
	}
	c.Find(bson.M{"champion2": champid}).All(&stat)
	for _, elem := range stat {
		fmt.Fprint(w, elem)
		fmt.Fprint(w, "\n")
	}
}

func championPage(w http.ResponseWriter, r *http.Request) {
	champion := r.URL.Query().Get("champion")
	champion = strings.ToLower(champion)
	champ_query := []rune(champion)
	champ_query[0] = rune(champion[0] - 32)
	champion = string(champ_query)
	if (!strings.ContainsAny(champion, "\",|&*;=%'+-_")) {
		session, err := mgo.Dial("127.0.0.1:27017")
		if err != nil {
			panic(err)
		}
		session.SetMode(mgo.Monotonic, true)
		c := session.DB("champ").C("fiche")
		var champ Champ
		err = c.Find(bson.M{"name" : champion}).One(&champ)
		session.Close()
		if champ.ChampionId != 0 {
			print_ratio (champ.ChampionId, w)
		} else {
			fmt.Fprint(w, "You failed man")
		}
	}
	fmt.Fprint(w, champion)
}

func acceuil(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/root/go/src/who_is_my_bestie/templates/acceuil.html")
	t.Execute(w, "")
}

func main() {
	http.HandleFunc("/fiche", championPage) // set router
	http.HandleFunc("/", acceuil) // set router
	err := http.ListenAndServe(":9090", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
