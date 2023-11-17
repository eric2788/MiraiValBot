package twitter

import (
	"html"
	"strings"
	"time"
)

const (
	Tweet           = "TWEET"
	ReTweet         = "RETWEET"
	Delete          = "DELETE"
	Reply           = "REPLY"
	ReTweetWithText = "RETWEET_WITH_TEXT"
)

type (
	TweetContent struct {
		Tweet    *TweetData `json:"tweet"`
		Profile  *Profile   `json:"profile"`
		NickName string     `json:"nick_name"`
	}

	TweetData struct {
		ConversationID    string
		GIFs              []Media
		Hashtags          []string
		HTML              string
		ID                string
		InReplyToStatus   *TweetData
		InReplyToStatusID string
		IsQuoted          bool
		IsPin             bool
		IsReply           bool
		IsRetweet         bool
		IsSelfThread      bool
		Likes             int
		Name              string
		Mentions          []Mention
		PermanentURL      string
		Photos            []Photo
		Place             *Place
		QuotedStatus      *TweetData
		QuotedStatusID    string
		Replies           int
		Retweets          int
		RetweetedStatus   *TweetData
		RetweetedStatusID string
		Text              string
		Thread            []*TweetData
		TimeParsed        time.Time
		Timestamp         int64
		URLs              []string
		UserID            string
		Username          string
		Videos            []Media
		Views             int
		SensitiveContent  bool
	}

	Profile struct {
		Avatar         string
		Banner         string
		Biography      string
		Birthday       string
		FollowersCount int
		FollowingCount int
		FriendsCount   int
		IsPrivate      bool
		IsVerified     bool
		Joined         *time.Time
		LikesCount     int
		ListedCount    int
		Location       string
		Name           string
		PinnedTweetIDs []string
		TweetsCount    int
		URL            string
		UserID         string
		Username       string
		Website        string
	}

	Media struct {
		ID      string
		Preview string // preview image url
		URL     string // video url
	}

	Mention struct {
		ID       string
		Username string
		Name     string
	}

	Photo struct {
		ID  string
		URL string
	}

	Place struct {
		ID          string `json:"id"`
		PlaceType   string `json:"place_type"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		CountryCode string `json:"country_code"`
		Country     string `json:"country"`
		BoundingBox struct {
			Type        string        `json:"type"`
			Coordinates [][][]float64 `json:"coordinates"`
		} `json:"bounding_box"`
	}
)

func (t TweetData) IsDeleteTweet() bool {
	return t.ID == ""
}

func (t TweetData) IsRetweetWithText() bool {
	return !t.IsRetweet && t.QuotedStatus != nil
}

func (t TweetData) GetCommand() string {
	switch {
	case t.IsDeleteTweet():
		return Delete
	case t.IsRetweet:
		return ReTweet
	case t.IsRetweetWithText():
		return ReTweetWithText
	case t.IsReply:
		return Reply
	default:
		return Tweet
	}
}

// with html unescaped string
func (t TweetData) UnEsacapedText() string {
	i := strings.IndexByte(t.Text, '&')

	if i < 0 {
		return t.Text
	}

	return html.UnescapeString(t.Text)
}
