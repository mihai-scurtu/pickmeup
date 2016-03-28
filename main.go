package main

import (
	"net/http"
	"os"
	"time"

	"encoding/json"
	"fmt"
	"github.com/mihai-scurtu/reddit-go/reddit"
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"strings"
)

type Post struct {
	Title string
	Url   string
	Embed string
}

type embedlyResponse struct {
	Url  string `json:"url"`
	Html string `json:"html"`
}

func (this *Post) updateEmbed() {
	url := fmt.Sprintf("https://api.embedly.com/1/oembed?url=%s&key=%s", url.QueryEscape(this.Url), "62b4ea459a544a7c9441f34a887b9951")
	req, _ := http.NewRequest("GET", url, nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var embedly embedlyResponse
	json.Unmarshal(body, &embedly)

	log.Printf("%+v", embedly)

	if embedly.Html != "" {
		this.Embed = embedly.Html
	} else {
		this.Embed = fmt.Sprintf(`<img src="%s">`, embedly.Url)
	}
}

var Subreddits []string
var Posts []*Post
var Client *reddit.Client

func UpdatePosts() {
	uri := "r/" + strings.Join(Subreddits, "+")

	posts := Client.GetPostListing(uri)

	for _, redditPost := range posts.GetChildren() {
		post := &Post{
			Title: redditPost.Title,
			Url:   redditPost.Url,
		}

		if !postExists(post) {
			Posts = append(Posts, post)
		}
	}

	log.Printf("Posts updated. Current count: %d", len(Posts))
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	Client = reddit.NewClient("PickMeUp cute image service")

	Subreddits = []string{
		"aww",
		"animalgifs",
		"animalsbeingderps",
		"dogpictures",
		"dogswearinghats",
		"rabbits",
		"redpandas",
		"pandagifs",
		"corgigifs",
		"puppysmiles",
	}

	rand.Seed(time.Now().UnixNano())

	// run update process
	go func() {
		for {
			UpdatePosts()
			time.Sleep(10 * time.Minute)
		}
	}()

	http.HandleFunc("/random", func(rw http.ResponseWriter, req *http.Request) {
		if len(Posts) == 0 {
			log.Println("No posts")
			return
		}

		r := rand.Intn(len(Posts))
		post := Posts[r]

		if post.Embed == "" {
			log.Printf("Embed: %s", post.Embed)
			post.updateEmbed()
		}

		rw.Write([]byte(post.Embed))
	})

	http.Handle("/", http.FileServer(http.Dir("public/")))

	http.ListenAndServe(":"+port, nil)
}

func postExists(post *Post) bool {
	for _, p := range Posts {
		if p.Url == post.Url {
			return true
		}
	}

	return false
}
