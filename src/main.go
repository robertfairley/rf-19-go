package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	blackfriday "gopkg.in/russross/blackfriday.v2"
)

// SiteSettings - settings JSON data struct
type SiteSettings struct {
	Title    string
	Slogan   string
	Notice   string
	Hostname string
	Port     string
}

// Page - site page data
type Page struct {
	Title   string
	Content interface{}
}

// StringKeyValue - JSON key:value as string
type StringKeyValue map[string]string

func loadSettings() SiteSettings {
	cwd, _ := os.Getwd()

	settingsPath := cwd + "/site-config.json"

	settingsFileOutput, _ := ioutil.ReadFile(settingsPath)

	var result StringKeyValue
	json.Unmarshal([]byte(settingsFileOutput), &result)

	settingsOutput := SiteSettings{
		Title:    result["title"],
		Slogan:   result["slogan"],
		Notice:   result["notice"],
		Hostname: result["hostname"],
		Port:     result["port"],
	}

	return settingsOutput
}

func getContents() []byte {
	cwd, _ := os.Getwd()

	testFilePath := cwd + "/posts/2019/01/hello-world.md"
	contents, err := ioutil.ReadFile(testFilePath)

	if err != nil {
		log.Fatal(err)
	}

	output := blackfriday.Run(contents)

	return output
}

func main() {
	cwd, _ := os.Getwd()

	settings := loadSettings()

	staticPath := cwd + "/static/"
	viewsPath := cwd + "/views"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(viewsPath + "/home.html"))
		data := Page{
			Title:   settings.Title,
			Content: "This is a test home page",
		}
		tmpl.Execute(w, data)
	})

	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(viewsPath + "/about.html"))
		tmpl.Execute(w, nil)
	})

	http.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles(viewsPath + "/post.html"))

		postContents := getContents()

		data := Page{
			Title:   "Hello, world!",
			Content: template.HTML(postContents),
		}
		tmpl.Execute(w, data)
	})

	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening at %s%s", settings.Hostname, settings.Port)
	log.Fatal(http.ListenAndServe(settings.Port, nil))
}
