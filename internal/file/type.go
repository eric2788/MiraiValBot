package file

import "github.com/eric2788/common-utils/set"

// Real Type
type (
	storageData struct {
		Answers    map[string]bool            `json:"answers"`
		Responses  map[string]string          `json:"responses"`
		WordCounts map[string]map[int64]int64 `json:"word_counts"`
		Bilibili   *bilibiliSettings          `json:"bilibili"`
		Youtube    *youtubeSettings           `json:"youtube"`
		Twitter    *twitterSettings           `json:"twitter"`
		Setting    *setting                   `json:"setting"`
		Listening  *listening                 `json:"listening"`
	}

	youtubeSettings struct {
		BroadcastIdle bool `json:"broadcastIdle"`
		AntiDuplicate bool `json:"antiDuplicate"`
	}

	bilibiliSettings struct {
		HighLightedUsers []int64 `json:"highLightedUsers"`
	}

	twitterSettings struct {
		ShowReply bool `json:"showReply"`
	}

	listening struct {
		Bilibili []int64  `json:"bilibili"`
		Youtube  []string `json:"youtube"`
		Twitter  []string `json:"twitter"`
		Valorant []string `json:"valorant"`
	}

	setting struct {
		VerboseDelete    bool    `json:"verboseDelete"`
		Verbose          bool    `json:"verbose"`
		YearlyCheck      bool    `json:"yearlyCheck"`
		LastChecked      int64   `json:"lastChecked"`
		MsgSeqAfter      int64   `json:"msgSeqAfter"`
		TimesPerNotify   int     `json:"timesPerNotify"`
		TagClassifyLimit float64 `json:"tagClassifyLimit"`
	}
)

// Wrapper Type

type (
	StorageData struct {
		Answers    map[string]bool
		Responses  map[string]string
		WordCounts map[string]map[int64]int64
		Bilibili   *BilibiliSettings
		Youtube    *YoutubeSettings
		Twitter    *TwitterSettings
		Setting    *Setting
		Listening  *Listening
	}

	YoutubeSettings struct {
		BroadcastIdle bool
		AntiDuplicate bool
	}

	BilibiliSettings struct {
		HighLightedUsers *set.Int64Set
	}

	TwitterSettings struct {
		ShowReply bool
	}

	Listening struct {
		Bilibili *set.Int64Set
		Youtube  *set.StringSet
		Twitter  *set.StringSet
		Valorant *set.StringSet
	}

	Setting struct {
		VerboseDelete    bool
		Verbose          bool
		YearlyCheck      bool
		LastChecked      int64
		MsgSeqAfter      int64
		TimesPerNotify   int
		TagClassifyLimit float64
	}
)

func (s *StorageData) toRealStorageData() *storageData {
	return &storageData{
		Answers:    s.Answers,
		Responses:  s.Responses,
		WordCounts: s.WordCounts,
		Youtube: &youtubeSettings{
			BroadcastIdle: s.Youtube.BroadcastIdle,
			AntiDuplicate: s.Youtube.AntiDuplicate,
		},
		Bilibili: &bilibiliSettings{
			HighLightedUsers: s.Bilibili.HighLightedUsers.ToArr(),
		},
		Twitter: &twitterSettings{
			ShowReply: s.Twitter.ShowReply,
		},
		Setting: &setting{
			VerboseDelete:    s.Setting.VerboseDelete,
			Verbose:          s.Setting.Verbose,
			YearlyCheck:      s.Setting.YearlyCheck,
			LastChecked:      s.Setting.LastChecked,
			MsgSeqAfter:      s.Setting.MsgSeqAfter,
			TimesPerNotify:   s.Setting.TimesPerNotify,
			TagClassifyLimit: s.Setting.TagClassifyLimit,
		},
		Listening: &listening{
			Bilibili: s.Listening.Bilibili.ToArr(),
			Youtube:  s.Listening.Youtube.ToArr(),
			Twitter:  s.Listening.Twitter.ToArr(),
			Valorant: s.Listening.Valorant.ToArr(),
		},
	}
}

func (s *StorageData) parse(sd *storageData) {
	s.Answers = sd.Answers
	s.Responses = sd.Responses
	s.WordCounts = sd.WordCounts
	s.Bilibili = &BilibiliSettings{
		HighLightedUsers: set.FromInt64Arr(sd.Bilibili.HighLightedUsers),
	}
	s.Youtube = &YoutubeSettings{
		BroadcastIdle: sd.Youtube.BroadcastIdle,
		AntiDuplicate: sd.Youtube.AntiDuplicate,
	}
	s.Twitter = &TwitterSettings{
		ShowReply: sd.Twitter.ShowReply,
	}
	s.Setting = &Setting{
		VerboseDelete:    sd.Setting.VerboseDelete,
		Verbose:          sd.Setting.Verbose,
		YearlyCheck:      sd.Setting.YearlyCheck,
		LastChecked:      sd.Setting.LastChecked,
		MsgSeqAfter:      sd.Setting.MsgSeqAfter,
		TimesPerNotify:   sd.Setting.TimesPerNotify,
		TagClassifyLimit: sd.Setting.TagClassifyLimit,
	}
	s.Listening = &Listening{
		Bilibili: set.FromInt64Arr(sd.Listening.Bilibili),
		Youtube:  set.FromStrArr(sd.Listening.Youtube),
		Twitter:  set.FromStrArr(sd.Listening.Twitter),
		Valorant: set.FromStrArr(sd.Listening.Valorant),
	}
}

func makeWrapper(sd *storageData) *StorageData {
	s := &StorageData{}
	s.parse(sd)
	return s
}
