package reddit

type Comment struct {
	Id string `json:"id"`
	ParentId string `json:"parent_id"`
	Created float64 `json:"created"`

	SubredditId string `json:"subreddit_id"`
	Subreddit   string `json:"subreddit"`

	Author string `json:"author"`
	Score  int `json:"score"`

	Body     string `json:"body"`
	BodyHtml string `json:"body_html"`

	LinkTitle  string `json:"link_title"`
	LinkId     string `json:"link_id"`
	LinkAuthor string `json:"link_author"`
	LinkUrl    string `json:"link_url"`
}

type CommentListing struct {
	Data struct {
		Children []struct {
			Data Comment `json:"data"`
		} `json:children`
		After  string `json:after`
		Before string `json:before`
	}
}

func (this *CommentListing) GetChildren() []Comment {
	var list []Comment

	for _, p := range this.Data.Children {
		list = append(list, p.Data)
	}

	return list
}
