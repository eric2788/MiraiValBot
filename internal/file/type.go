package file

import (
	mapset "github.com/deckarep/golang-set/v2"
)

// Real Type
type (
	storageData struct {
		Answers    map[string]bool            `json:"answers"`
		Responses  map[string]string          `json:"responses"`
		WordCounts map[string]map[int64]int64 `json:"word_counts"`
		Points     map[int64]int64            `json:"points"`
		Bilibili   *bilibiliSettings          `json:"bilibili"`
		Youtube    *youtubeSettings           `json:"youtube"`
		Twitter    *twitterSettings           `json:"twitter"`
		Setting    *setting                   `json:"setting"`
		Listening  *listening                 `json:"listening"`
		AiDraw     *aiDrawSettings            `json:"aiDraw"`
	}

	aiDrawSettings struct {
		SexyAISession string `json:"sexy_ai_session"`
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
		Points     map[int64]int64
		Bilibili   *BilibiliSettings
		Youtube    *YoutubeSettings
		Twitter    *TwitterSettings
		Setting    *Setting
		Listening  *Listening
		AiDraw     *AIDrawSettings
	}

	AIDrawSettings struct {
		SexyAISession string
	}

	YoutubeSettings struct {
		BroadcastIdle bool
		AntiDuplicate bool
	}

	BilibiliSettings struct {
		HighLightedUsers mapset.Set[int64]
	}

	TwitterSettings struct {
		ShowReply bool
	}

	Listening struct {
		Bilibili mapset.Set[int64]
		Youtube  mapset.Set[string]
		Twitter  mapset.Set[string]
		Valorant mapset.Set[string]
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
		Points:     s.Points,
		Youtube: &youtubeSettings{
			BroadcastIdle: s.Youtube.BroadcastIdle,
			AntiDuplicate: s.Youtube.AntiDuplicate,
		},
		Bilibili: &bilibiliSettings{
			HighLightedUsers: s.Bilibili.HighLightedUsers.ToSlice(),
		},
		Twitter: &twitterSettings{
			ShowReply: s.Twitter.ShowReply,
		},
		AiDraw: &aiDrawSettings{
			SexyAISession: s.AiDraw.SexyAISession,
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
			Bilibili: s.Listening.Bilibili.ToSlice(),
			Youtube:  s.Listening.Youtube.ToSlice(),
			Twitter:  s.Listening.Twitter.ToSlice(),
			Valorant: s.Listening.Valorant.ToSlice(),
		},
	}
}

func (s *StorageData) parse(sd *storageData) {
	s.Answers = sd.Answers
	s.Responses = sd.Responses
	s.WordCounts = sd.WordCounts
	s.Points = sd.Points
	s.Bilibili = &BilibiliSettings{
		HighLightedUsers: mapset.NewSet(sd.Bilibili.HighLightedUsers...),
	}
	s.Youtube = &YoutubeSettings{
		BroadcastIdle: sd.Youtube.BroadcastIdle,
		AntiDuplicate: sd.Youtube.AntiDuplicate,
	}
	s.Twitter = &TwitterSettings{
		ShowReply: sd.Twitter.ShowReply,
	}
	s.AiDraw = &AIDrawSettings{
		SexyAISession: sd.AiDraw.SexyAISession,
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
		Bilibili: mapset.NewSet(sd.Listening.Bilibili...),
		Youtube:  mapset.NewSet(sd.Listening.Youtube...),
		Twitter:  mapset.NewSet(sd.Listening.Twitter...),
		Valorant: mapset.NewSet(sd.Listening.Valorant...),
	}
}

func makeWrapper(sd *storageData) *StorageData {
	s := &StorageData{}
	s.parse(sd)
	return s
}
