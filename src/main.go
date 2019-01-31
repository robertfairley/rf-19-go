package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

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

//
var cwd, _ = os.Getwd()

//
var settings = loadSettings()

//
var staticPath = cwd + "/static/"

//
var viewsPath = cwd + "/views"

// loadSettings - Load the site settins from a JSON to prevent
// the need to recompile if certain global settings change.
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

// getContents - Get the markdown contents of a post and
// convert it to renderable HTML
func getContents(path string) []byte {
	cwd, _ := os.Getwd()

	testFilePath := cwd + "/posts/" + path
	contents, err := ioutil.ReadFile(testFilePath)

	if err != nil {
		log.Fatal(err)
	}

	output := blackfriday.Run(contents)

	return output
}

// getPostName - Parse the provided string and return
// the full post filename.
func getPostName(filepath string) string {
	var extension string = ".md"

	filename := path.Base(filepath)
	postName := strings.Split(filename, extension)

	return postName[0]
}

//
func getPostFilename(postURL string) string {
	var extension string = ".md"

	return postURL + extension
}

//
func PostRouteHandler(w http.ResponseWriter, r *http.Request) {
	postPath := strings.TrimPrefix(r.URL.Path, "/posts/")

	tmpl := template.Must(template.ParseFiles(viewsPath + "/post.html"))

	postContents := getContents(getPostFilename(postPath))

	data := Page{
		Content: template.HTML(postContents),
	}
	tmpl.Execute(w, data)
}

func main() {

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

	http.HandleFunc("/posts/", PostRouteHandler)

	fs := http.FileServer(http.Dir(staticPath))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	fmt.Printf("Listening at %s%s", settings.Hostname, settings.Port)
	log.Fatal(http.ListenAndServe(settings.Port, nil))
}
