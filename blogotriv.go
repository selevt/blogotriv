package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Post struct {
	ID    string
	Date  string
	Title string
	Text  template.HTML
}

type Data struct {
	Title string
	Posts []Post
}

var content = make(map[string]string)

var templates = template.Must(template.ParseFiles("index.html"))

func handlerHome(w http.ResponseWriter, r *http.Request) {
	var data *Data

	if r.FormValue("id") != "" {
		data = loadDataForPost(r.FormValue("id"))
	} else if r.FormValue("p") != "" {
		data = loadDataForPage(r.FormValue("p"))
	} else {
		data = loadData()
	}

	err := templates.ExecuteTemplate(w, "index.html", data)
	if err != nil {
		fmt.Println("Error loading template: ", err)
	}
}

func loadData() *Data {
	data := &Data{}
	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		fmt.Println("Error reading data: ", err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		fmt.Println("Error parsing data to JSON: ", err)
	}

	return data
}

func loadDataForPost(postId string) *Data {
	data := loadData()
	if postId != "" {
		for _, val := range data.Posts {
			if val.ID == postId {
				data.Posts = []Post{val}
			}
		}
	}
	return data
}

func loadDataForPage(pageId string) *Data {
	data := loadData()

	// avoid loading of external files
	if strings.Contains(pageId, "..") {
		return data
	}

	content, err := ioutil.ReadFile("pages/" + pageId + ".html")
	if err != nil {
		return data
	}

	data.Posts = []Post{Post{Text: template.HTML(content)}}
	return data
}

func main() {
	http.HandleFunc("/", handlerHome)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))
	http.ListenAndServe(":8080", nil)
}
