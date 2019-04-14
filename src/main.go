package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
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
	Title    string
	Date     string
	Greeting interface{}
	Content  interface{}
}

// DirInfo - directory info including the original path
type DirInfo struct {
	Children []os.FileInfo
	Info     os.FileInfo
	Path     string
}

// PostInfo - ... tbd
type PostInfo struct {
	Info []byte
	Path string
	Meta PostMeta
	Data interface{}
}

// PostDate - structure for holding parsed dates
type PostDate struct {
	Year  string
	Month string
	Day   string
}

// PostMeta - Basic post metadata
type PostMeta struct {
	Title   string
	DateStr string
	Date    PostDate
	Excerpt string
	URL     string
	image   string
}

// StringKeyValue - JSON key:value as string
type StringKeyValue map[string]string

// Globals
var cwd, _ = os.Getwd()
var settings = loadSettings()
var staticPath = cwd + settings.StaticDir
var viewsPath = cwd + settings.ViewsDir
var postsPath = cwd + settings.PostsDir
var postList, getPostListErr = getPostList()

// loadSettings - Load the site settins from a JSON to prevent
// the need to recompile if certain global settings change.
func loadSettings() SiteSettings {
	settingsPath := cwd + "/site-config.json"

	settingsFileOutput, _ := ioutil.ReadFile(settingsPath)

	var result StringKeyValue
	json.Unmarshal(settingsFileOutput, &result)

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
	fullContents, err := ioutil.ReadFile(postFilePath)
	// Remove the metadata header. Maybe extract this process later on.
	contents := strings.Split(string(fullContents), "---")[1]

	if err != nil {
		log.Fatal(err)
	}

	output := blackfriday.Run([]byte(contents))

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

// getPostDate - format a date as a structure for convenient use
func getPostDate(dateStr string) PostDate {
	dateSpl := strings.Split(dateStr, "-")

	return PostDate{
		Year:  dateSpl[0],
		Month: dateSpl[1],
		Day:   dateSpl[2],
	}
}

// getPostURL - get a formatted post url string
func getPostURL(path string, date PostDate) string {
	var postURL string

	postName := getPostName(path)
	postURL = "/posts/" + date.Year + "/" + date.Month + "/" + postName

	return postURL
}

// getPostMeta - generate a struct containing basic post metadata
func getPostMeta(path string) PostMeta {
	postFile, _ := ioutil.ReadFile(path)
	dat := string(postFile)

	headerDelim := "---"
	infoDelim := ": "

	postHeader := strings.Split(dat, headerDelim)[0]

	postHeaderSpl := strings.Split(postHeader, "\n")

	postTitle := strings.Split(postHeaderSpl[0], infoDelim)[1]
	postDateStr := strings.Split(postHeaderSpl[1], infoDelim)[1]
	postDate := getPostDate(postDateStr)
	postExcerpt := strings.Split(postHeaderSpl[2], infoDelim)[1]
	postURL := getPostURL(path, postDate)
	postImage := strings.Split(postHeaderSpl[3], infoDelim)[1]

	postMeta := PostMeta{
		Title:   postTitle,
		Date:    postDate,
		DateStr: postDateStr,
		Excerpt: postExcerpt,
		URL:     postURL,
		image:   postImage,
	}

	return postMeta
}

// getPostList - Get and return a list of absolute paths
// to every available blog post and as a second value return
// a list of every available post filename
func getPostList() ([]PostInfo, error) {
	var yearsDirs []os.FileInfo
	var years []DirInfo
	var months []DirInfo
	var posts []PostInfo
	var err error

	yearsDirs, err = ioutil.ReadDir(postsPath)

	if err != nil {
		return nil, err
	}

	for _, year := range yearsDirs {
		thisYear := DirInfo{
			Info: year,
			Path: postsPath + year.Name() + "/",
		}
		thisYear.Children, err = ioutil.ReadDir(thisYear.Path)

		if err != nil {
			return nil, err
		}

		years = append(years, thisYear)
	}

	for i := 0; i < len(years); i++ {
		for j := 0; j < len(years[i].Children); j++ {
			thisMonthPath := years[i].Path + years[i].Children[j].Name() + "/"
			thisMonth := DirInfo{
				Info: years[i].Children[j],
				Path: thisMonthPath,
			}
			thisMonth.Children, err = ioutil.ReadDir(thisMonth.Path)

			if err != nil {
				return nil, err
			}

			months = append(months, thisMonth)
		}
	}

	for _, month := range months {

		for i := 0; i < len(month.Children); i++ {
			postName := month.Children[i].Name()
			postPath := month.Path + postName
			postData, err := ioutil.ReadFile(postPath)

			if err != nil {
				return nil, err
			}

			fmt.Println(postPath)

			postMeta := getPostMeta(postPath)
			thisPost := PostInfo{
				Info: postData,
				Path: postPath,
				Meta: postMeta,
			}
			posts = append(posts, thisPost)
		}
	}

	return posts, nil
}

// HomeRouteHandler - Response for the home page
func HomeRouteHandler(w http.ResponseWriter, r *http.Request) {
	pageTemplate := template.Must(template.ParseFiles(viewsPath + "/home.html"))
	postImageBasePath := "/static/images/"

	var postLinks string
	var noScriptPostLinks string
	var postImage string

	for i := 0; i < len(postList); i++ {

		if postList[i].Meta.image != "" {
			postImage = "url('" + postImageBasePath + postList[i].Meta.image + "')"
		} else {
			r := strconv.Itoa(rand.Intn(255))
			g := strconv.Itoa(rand.Intn(255))
			b := strconv.Itoa(rand.Intn(255))
			postImage = "rgb(" + r + "," + g + "," + b + "); filter: saturate(0.5)"
		}

		postLinks += `<div class="column col-6 col-sm-12">
		<div class="parallax post-card rounded" onClick="window.location.href = '` + postList[i].Meta.URL + `'">
			<div class="parallax-top-left" tabindex="1"></div>
			<div class="parallax-top-right" tabindex="2"></div>
			<div class="parallax-bottom-left" tabindex="3"></div>
			<div class="parallax-bottom-right" tabindex="4"></div>
			<div class="parallax-content">
				<div class="post-card__front parallax-front">
					<div class="post-card__label title rounded">` + postList[i].Meta.Title + `</div>
					<div class="post-card__label"><small>` + postList[i].Meta.DateStr + `</small></div>
					
					<div class="post-card__label rounded">
						<a href="` + postList[i].Meta.URL + `">Read More</a>
					</div>
				</div>
				<div class="parallax-back" style="background: ` + postImage + `; width: 100%; height: 400px; clear: both;">
				</div>
			</div>
		</div>
	</div>`

		noScriptPostLinks += `<div>
		<h5><a href="` + postList[i].Meta.URL + `">` + postList[i].Meta.Title + `</a>	-	` + postList[i].Meta.DateStr + `</h5>
	</div>`
	}

	data := Page{
		Title: settings.Title,
		Content: template.HTML(`<div class="container">
			<div class="columns" style="display:none;" id="post-cards">` + postLinks + `</div>
			<noscript>
			<nav>
				<h2>Post</h2>
				<hr />
				` + noScriptPostLinks + `
			</nav>
		</noscript>
		<script>
			document.querySelector("#post-cards").style = null;
		</script>
		</div>`),
	}
	pageTemplate.Execute(w, data)
}

// AboutRouteHandler - Response for the about page
func AboutRouteHandler(w http.ResponseWriter, r *http.Request) {
	pageTemplate := template.Must(template.ParseFiles(viewsPath + "/about.html"))
	pageTemplate.Execute(w, nil)
}

// CvRouteHandler - Response for the CV page
func CvRouteHandler(w http.ResponseWriter, r *http.Request) {
	pageTemplate := template.Must(template.ParseFiles(viewsPath + "/cv.html"))
	pageTemplate.Execute(w, nil)
}

// PostRouteHandler - Response for any blog post. Parses the markdown file
// and injects it into the page based on the `year/month/post-name`
// path beneath `/posts/` sub-routes`
func PostRouteHandler(w http.ResponseWriter, r *http.Request) {
	postPath := strings.TrimPrefix(r.URL.Path, settings.PostsDir)

	pageTemplate := template.Must(template.ParseFiles(viewsPath + "/post.html"))

	postContents := getContents(getPostFilename(postPath))

	data := Page{
		Content: template.HTML(postContents),
	}
	pageTemplate.Execute(w, data)
}

func main() {

	if getPostListErr != nil {
		fmt.Print("SERVER EXITED DUE TO ERROR READING/PARSING POST LIST:\n")
		log.Fatal(getPostListErr)
	}

	http.HandleFunc("/", HomeRouteHandler)
	http.HandleFunc("/about", AboutRouteHandler)
	http.HandleFunc("/cv", CvRouteHandler)
	http.HandleFunc(settings.PostsDir, PostRouteHandler)

	fs := http.FileServer(http.Dir(staticPath))
	http.Handle(settings.StaticDir, http.StripPrefix(settings.StaticDir, fs))

	fmt.Printf("Listening at %s%s\n", settings.Hostname, settings.Port)
	log.Fatal(http.ListenAndServe(settings.Port, nil))

	getPostList()
}
