package main

import (
	"fmt"
	"net/http"
//	"io/ioutil"
	"html/template"
	"strings"
	"log"
	"sort"
	"strconv"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
)

var lst_champ []Champ

type Template struct {
	Champion string
	Wins int
	Games int
	Ratio string
	Pos bool
}

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

type Templates []Template

type Fiche struct {
	Temp Templates
	Champ string
}

var fiche Fiche

func (slice Templates) Len() int {
	return len(slice)
}

func (slice Templates) Less(i, j int) bool {
	ratio1 := float64(slice[i].Wins) / float64(slice[i].Games) * float64(100)
	ratio2 := float64(slice[j].Wins) / float64(slice[j].Games) * float64(100)
	return ratio1 > ratio2
}

func (slice Templates) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func Get_elem(name string, id int) (Champ){
	for _, elem := range lst_champ {
		if strings.EqualFold(elem.Name, name) || elem.ChampionId == id {
			return elem
		}
	}
	var ret Champ
	ret.ChampionId = 0
	ret.Name = ""
	return ret
}

func return_format_str(elem Stat, id int, name string) (string) {
//	name = "<div class=\"img \" style=\"height:48px; width:48px; background: url('//ddragon.leagueoflegends.com/cdn/6.9.1/img/champion/" + Get_elem("", id).Name + ".png') -0px -0px no-repeat; background-size: 100% 100%; \" data-rg-name=\"champion\" data-rg-id=" + Get_elem("", id).Name + "></div>"
	name = "<p><img src=\"http://ddragon.leagueoflegends.com/cdn/6.9.1/img/champion/" + Get_elem("", id).Name + ".png\" data-rg-name=\"champion\" data-rg-id=" + Get_elem("", id).Name + " height=\"48\" width=\"48\">"
	str := name + Get_elem("", id).Name + " -- Ratio : "
	ratio := float64(elem.Win) / float64(elem.Games) * float64(100)
	str = str + strconv.FormatFloat(ratio, 'f', 2, 64)
	str = str + " -- Games : " + strconv.Itoa(elem.Games) + "</p>"
	return str
}

func print_ratio(champid int, champname string, w http.ResponseWriter) {
//	fmt.Fprint(w, "Liste des ratio avec : ")
//	fmt.Fprint(w, champname + "\n")
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("champ").C("stat")
	var stat []Stat
	c.Find(bson.M{"champion1": champid}).All(&stat)
	var temp Templates
	for _, elem := range stat {
		if elem.Games > 9 {
			value := float64(elem.Win) / float64(elem.Games) * float64(100)
			ratio := strconv.FormatFloat(value, 'f', 2, 64)
			temp = append(temp, Template {Champion: Get_elem("", elem.Champion2).Name, Wins: elem.Win, Games: elem.Games, Ratio: ratio, Pos: (value >= 50)})
		}
	}
	c.Find(bson.M{"champion2": champid}).All(&stat)
	for _, elem := range stat {
		if elem.Games > 9 {
			value := float64(elem.Win) / float64(elem.Games) * float64(100)
			ratio := strconv.FormatFloat(value, 'f', 2, 64)
			temp = append(temp, Template {Champion: Get_elem("", elem.Champion1).Name, Wins: elem.Win, Games: elem.Games, Ratio: ratio, Pos: (value >= 50)})
		}
	}
	fmt.Println(temp.Len())
	sort.Sort(temp)
	t, _ := template.ParseFiles("/root/go/src/who_is_my_bestie/templates/fiche.html")
	fiche.Temp = temp
	t.Execute(w, &fiche)
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

		err = c.Find(nil).All(&lst_champ)
		champ := Get_elem(champion, 0)
		session.Close()
		if champ.ChampionId != 0 {
			fiche.Champ = champ.Name
			print_ratio (champ.ChampionId, champ.Name, w)
		} else {
			fmt.Fprint(w, "You failed man")
		}
	}
}


func about(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/root/go/src/who_is_my_bestie/templates/about.html")
	t.Execute(w, "")
}

func acceuil(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("/root/go/src/who_is_my_bestie/templates/acceuil.html")
	t.Execute(w, "")
}

func main() {
	r := mux.NewRouter()
	r.Host("http://178.62.52.164")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/root/go/src/who_is_my_bestie/static"))))
	r.HandleFunc("/fiche", championPage) // set router
	r.HandleFunc("/about", about) // set router
	r.HandleFunc("/", acceuil) // set router
	http.Handle("/", r)
	err := http.ListenAndServe(":80", nil) // set listen port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
