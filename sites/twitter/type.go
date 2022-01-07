package twitter

const (
	Tweet           = "TWEET"
	ReTweet         = "RETWEET"
	Delete          = "DELETE"
	Reply           = "REPLY"
	ReTweetWithText = "RETWEET_WITH_TEXT"
)

// if it has *, then it may be nullable

type TweetUser struct {
	CreatedAt                      string  `json:"created_at"`
	DefaultProfile                 bool    `json:"default_profile"`
	Url                            *string `json:"url"`
	DefaultProfileImage            bool    `json:"default_profile_image"`
	Description                    string  `json:"description"`
	FavouritesCount                int64   `json:"favourites_count"`
	FollowersCount                 int64   `json:"followers_count"`
	FriendsCount                   int64   `json:"friends_count"`
	Id                             int64   `json:"id"`
	IdStr                          string  `json:"id_str"`
	ListedCount                    int64   `json:"listed_count"`
	Location                       string  `json:"location"`
	Name                           string  `json:"name"`
	ProfileBackgroundColor         string  `json:"profile_background_color"`
	ProfileBackgroundImageUrl      string  `json:"profile_background_image_url"`
	ProfileBackgroundImageUrlHttps string  `json:"profile_background_image_url_https"`
	ProfileBackgroundTile          bool    `json:"profile_background_tile"`
	ProfileBannerUrl               string  `json:"profile_banner_url"`
	ProfileImageUrl                string  `json:"profile_image_url"`
	ProfileImageUrlHttps           string  `json:"profile_image_url_https"`
	ProfileLinkColor               string  `json:"profile_link_color"`
	ProfileSidebarBorderColor      string  `json:"profile_sidebar_border_color"`
	ProfileSidebarFillColor        string  `json:"profile_sidebar_fill_color"`
	ProfileTextColor               string  `json:"profile_text_color"`
	ProfileUseBackgroundImage      bool    `json:"profile_use_background_image"`
	Protected                      bool    `json:"protected"`
	ScreenName                     string  `json:"screen_name"`
	StatusesCount                  int64   `json:"statuses_count"`
	Verified                       bool    `json:"verified"`
}

type TweetStreamData struct {
	CreatedAt        string         `json:"created_at"`
	Entities         TweetEntities  `json:"entities"`
	ExtendedEntities *TweetEntities `json:"extended_entities"`
	FavouriteCount   int64          `json:"favourite_count"`
	Id               int64          `json:"id"`
	IdStr            string         `json:"id_str"`
	Lang             string         `json:"lang"`
	QuoteCount       int64          `json:"quote_count"`
	ReplyCount       int64          `json:"reply_count"`
	RetweetCount     int64          `json:"retweet_count"`
	Source           string         `json:"source"`
	Text             string         `json:"text"`
	TimestampMs      string         `json:"timestamp_ms"`
	Truncated        bool           `json:"truncated"`
	User             TweetUser      `json:"user"`

	//retweet with text
	IsQuoteStatus     bool             `json:"is_quote_status"`
	QuotedStatusId    int64            `json:"quoted_status_id"`
	QuotedStatusIdStr string           `json:"quoted_status_id_str"`
	QuotedStatus      *TweetStreamData `json:"quoted_status"`

	//reply
	InReplyToStatusId    int64   `json:"in_reply_to_status_id"`
	InReplyToStatusIdStr string  `json:"in_reply_to_status_id_str"`
	InReplyToUserId      int64   `json:"in_reply_to_user_id"`
	InReplyToUserIdStr   string  `json:"in_reply_to_user_id_str"`
	InReplyToScreenName  *string `json:"in_reply_to_screen_name"`

	PossiblySensitive bool `json:"possibly_sensitive"`

	RetweetedStatus *TweetStreamData `json:"retweeted_status"`

	Delete *struct {
		Status struct {
			Id        int64  `json:"id"`
			IdStr     string `json:"id_str"`
			UserId    int64  `json:"user_id"`
			UserIdStr string `json:"user_id_str"`
		} `json:"status"`
		TimestampMs int64 `json:"timestamp_ms"`
	} `json:"delete"`
}

type TweetEntities struct {
	HashTags []interface{} `json:"hash_tags"`
	Symbols  []interface{} `json:"symbols"`
	Urls     []struct {
		DisplayUrl  *string `json:"display_url"`
		ExpandedUrl string  `json:"expanded_url"`
		Url         string  `json:"url"`
	} `json:"urls"`
	UserMentions []struct {
		Id         int64  `json:"id"`
		IdStr      string `json:"id_str"`
		Name       string `json:"name"`
		ScreenName string `json:"screen_name"`
	} `json:"user_mentions"`
	Media *[]struct {
		Id            int64  `json:"id"`
		IdStr         string `json:"id_str"`
		MediaUrl      string `json:"media_url"`
		MediaUrlHttps string `json:"media_url_https"`
		Url           string `json:"url"`
		DisplayUrl    string `json:"display_url"`
		ExpandUrl     string `json:"expand_url"`
		Type          string `json:"type"`
	} `json:"media"`
}

func (t TweetStreamData) IsDeleteTweet() bool {
	return t.Delete != nil
}

func (t TweetStreamData) IsRetweet() bool {
	return t.RetweetedStatus != nil
}

func (t TweetStreamData) IsRetweetWithText() bool {
	return !t.IsRetweet() && t.QuotedStatus != nil
}

func (t TweetStreamData) IsReply() bool {
	return t.InReplyToScreenName != nil
}

func (t TweetStreamData) GetCommand() string {
	switch {
	case t.IsDeleteTweet():
		return Delete
	case t.IsRetweet():
		return ReTweet
	case t.IsRetweetWithText():
		return ReTweetWithText
	case t.IsReply():
		return Reply
	default:
		return Tweet
	}
}
