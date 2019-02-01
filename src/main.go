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
	Title     string
	Slogan    string
	Notice    string
	Hostname  string
	Port      string
	PostExt   string
	PostsDir  string
	ViewsDir  string
	StaticDir string
}

// Page - site page data
type Page struct {
	Title   string
	Content interface{}
}

// StringKeyValue - JSON key:value as string
type StringKeyValue map[string]string

// Globals
var cwd, _ = os.Getwd()
var settings = loadSettings()
var staticPath = cwd + settings.StaticDir
var viewsPath = cwd + settings.ViewsDir
var postsPath = cwd + settings.PostsDir

// loadSettings - Load the site settins from a JSON to prevent
// the need to recompile if certain global settings change.
func loadSettings() SiteSettings {
	settingsPath := cwd + "/site-config.json"

	settingsFileOutput, _ := ioutil.ReadFile(settingsPath)

	var result StringKeyValue
	json.Unmarshal([]byte(settingsFileOutput), &result)

	settingsOutput := SiteSettings{
		Title:     result["title"],
		Slogan:    result["slogan"],
		Notice:    result["notice"],
		Hostname:  result["hostname"],
		Port:      result["port"],
		PostExt:   result["postExtension"],
		PostsDir:  result["postsDir"],
		ViewsDir:  result["viewsDir"],
		StaticDir: result["staticDir"],
	}

	return settingsOutput
}

// getContents - Get the markdown contents of a post and
// convert it to renderable HTML
func getContents(path string) []byte {
	testFilePath := postsPath + path
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
	filename := path.Base(filepath)
	postName := strings.Split(filename, settings.PostExt)

	return postName[0]
}

// getPostFilename - concatenate the standard blog post
// file extension and the provided path to the post.
func getPostFilename(postURL string) string {
	return postURL + settings.PostExt
}

// HomeRouteHandler - Response for the home page
func HomeRouteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(viewsPath + "/home.html"))
	data := Page{
		Title:   settings.Title,
		Content: "This is a test home page",
	}
	tmpl.Execute(w, data)
}

// AboutRouteHandler - Response for the about page
func AboutRouteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(viewsPath + "/about.html"))
	tmpl.Execute(w, nil)
}

// CvRouteHandler - Response for the CV page
func CvRouteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(viewsPath + "/cv.html"))
	tmpl.Execute(w, nil)
}

// PostRouteHandler - Response for any blog post. Parses the markdown file
// and injects it into the page based on the `year/month/post-name`
// path beneath `/posts/` sub-routes`
func PostRouteHandler(w http.ResponseWriter, r *http.Request) {
	postPath := strings.TrimPrefix(r.URL.Path, settings.PostsDir)

	tmpl := template.Must(template.ParseFiles(viewsPath + "/post.html"))

	postContents := getContents(getPostFilename(postPath))

	data := Page{
		Content: template.HTML(postContents),
	}
	tmpl.Execute(w, data)
}

func main() {

	fmt.Print(settings)

	http.HandleFunc("/", HomeRouteHandler)
	http.HandleFunc("/about", AboutRouteHandler)
	http.HandleFunc("/cv", CvRouteHandler)
	http.HandleFunc(settings.PostsDir, PostRouteHandler)

	fs := http.FileServer(http.Dir(staticPath))
	http.Handle(settings.StaticDir, http.StripPrefix(settings.StaticDir, fs))

	fmt.Printf("Listening at %s%s\n", settings.Hostname, settings.Port)
	log.Fatal(http.ListenAndServe(settings.Port, nil))
}
