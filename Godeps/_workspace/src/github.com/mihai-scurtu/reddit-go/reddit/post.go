package reddit

type Post struct {
	Id           string  `json:id`
	Name         string  `json:name`
	Author       string  `json:author`
	Created      float64 `json:created_utc`
	Domain       string  `json:domain`
	IsSelf       bool    `json:is_self`
	IsSticky     bool    `json:stickied`
	Permalink    string  `json:permalink`
	Score        float64 `json:score`
	SelfText     string  `json:selftext`
	SelfTextHtml string  `json:selftext_html`
	Subreddit    string  `json:subreddit`
	SubredditId  string  `json:subreddit_id`
	Thumbnail    string  `json:thumbnail`
	Title        string  `json:title`
	Url          string  `json:url`
	Upvotes      float64 `json:ups`
	Downvotes    float64 `json:downs`
}

type PostListing struct {
	Data struct {
		Children []struct {
			Data Post
		} `json:children`
		After  string `json:after`
		Before string `json:before`
	}
}

func (pl *PostListing) GetChildren() []Post {
	var list []Post

	for _, p := range pl.Data.Children {
		list = append(list, p.Data)
	}

	return list
}
