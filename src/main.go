package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

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
	Title    string
	Greeting interface{}
	Content  interface{}
}

// Post - post metadata
type Post struct {
	Title     string
	Date      string
	Image     string
	Excerpt   string
	URL       string
	LocalPath string
	HTML      interface{}
}

// DirInfo - directory info including the original path
type DirInfo struct {
	children []os.FileInfo
	info     os.FileInfo
	path     string
}

// PostInfo - ... tbd
type PostInfo struct {
	info []byte
	path string
	data interface{}
}

// StringKeyValue - JSON key:value as string
type StringKeyValue map[string]string

// Globals
var cwd, _ = os.Getwd()
var settings = loadSettings()
var staticPath = cwd + settings.StaticDir
var viewsPath = cwd + settings.ViewsDir
var postsPath = cwd + settings.PostsDir

//var postPaths, postNames = getPostList()

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
		PostExt:   result["postExt"],
		PostsDir:  result["postsDir"],
		ViewsDir:  result["viewsDir"],
		StaticDir: result["staticDir"],
	}

	return settingsOutput
}

// getContents - Get the markdown contents of a post and
// convert it to renderable HTML
func getContents(path string) []byte {
	postFilePath := postsPath + path
	contents, err := ioutil.ReadFile(postFilePath)

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

//
func getPostMeta(path string) Post {
	postMeta := Post{}

	postMeta.Title = getPostFilename(path)
	postMeta.Date = time.Now().String()
	postMeta.Excerpt = "Test item"
	postMeta.LocalPath = path
	postMeta.HTML = ""

	return postMeta
}

// getPostList - Get and return a list of absolute paths
// to every available blog post and as a second value return
// a list of every available post filename
func getPostList() ([]string, []string) {
	var yearsDirs []os.FileInfo
	//var monthsDirs []os.FileInfo
	var years []DirInfo
	var months []DirInfo
	var posts []PostInfo
	//var postFiles []os.FileInfo
	//var postPaths []string
	//var postNames []string

	yearsDirs, _ = ioutil.ReadDir(postsPath)

	for _, year := range yearsDirs {
		thisYear := DirInfo{
			info: year,
			path: postsPath + year.Name() + "/",
		}
		thisYear.children, _ = ioutil.ReadDir(thisYear.path)
		years = append(years, thisYear)
	}

	for i := 0; i < len(years); i++ {
		for j := 0; j < len(years[i].children); j++ {
			thisMonthPath := years[i].path + years[i].children[j].Name() + "/"
			thisMonth := DirInfo{
				info: years[i].children[j],
				path: thisMonthPath,
			}
			thisMonth.children, _ = ioutil.ReadDir(thisMonth.path)
			months = append(months, thisMonth)
		}
	}

	for _, month := range months {

		for i := 0; i < len(month.children); i++ {
			postName := month.children[i].Name()
			postData, _ := ioutil.ReadFile(postName)
			thisPost := PostInfo{
				info: postData,
				path: month.path + postName,
			}
			posts = append(posts, thisPost)
		}
	}

	for _, test := range posts {
		fmt.Println(test.path)
	}

	//for _, postFile := range postFiles {
	//	relPath := "/" + year.Name() + "/" + month.Name() + "/" + postFile.Name()
	//	postURLPath := strings.Replace(relPath, ".md", "", 1)
	//	postProperName := strings.Replace(postFile.Name(), ".md", "", 1)
	//	postPaths = append(postPaths, postURLPath)
	//	postNames = append(postNames, postProperName)
	//}

	return nil, nil
}

/*
// HomeRouteHandler - Response for the home page
func HomeRouteHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(viewsPath + "/home.html"))

	var postLinks string

	for i := 0; i < len(postNames); i++ {
		postLinks += `<li><a href="/posts/` + postPaths[i] + `">` + postNames[i] + `</a></li>`
	}

	md := "```\n" +
		"#include <stdio.h>\n" +
		"int main(void) {\n" +
		"  printf(\"Hello, world!\");\n" +
		"}```"

	greet := []byte(md)

	greeting := blackfriday.Run(greet)

	data := Page{
		Title:    settings.Title,
		Greeting: template.HTML(string(greeting)),
		Content:  template.HTML(`<div class="container"><ul>` + postLinks + `</ul></div>`),
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
*/
func main() {

	//http.HandleFunc("/", HomeRouteHandler)
	//http.HandleFunc("/about", AboutRouteHandler)
	//http.HandleFunc("/cv", CvRouteHandler)
	//http.HandleFunc(settings.PostsDir, PostRouteHandler)

	//fs := http.FileServer(http.Dir(staticPath))
	//http.Handle(settings.StaticDir, http.StripPrefix(settings.StaticDir, fs))

	//fmt.Printf("Listening at %s%s\n", settings.Hostname, settings.Port)
	//log.Fatal(http.ListenAndServe(settings.Port, nil))

	getPostList()
}
