package file

import "github.com/eric2788/MiraiValBot/utils/set"

// Real Type
type (
	storageData struct {
		Answers   map[string]bool   `json:"answers"`
		Responses map[string]string `json:"responses"`
		Bilibili  *bilibiliSettings `json:"bilibili"`
		Setting   *setting          `json:"setting"`
		Listening *listening        `json:"listening"`
	}

	bilibiliSettings struct {
		HighLightedUsers []int64 `json:"highLightedUsers"`
	}

	listening struct {
		Bilibili []int64  `json:"bilibili"`
		Youtube  []string `json:"youtube"`
		Twitter  []string `json:"twitter"`
	}

	setting struct {
		VerboseDelete bool  `json:"verboseDelete"`
		Verbose       bool  `json:"verbose"`
		YearlyCheck   bool  `json:"yearlyCheck"`
		LastChecked   int64 `json:"lastChecked"`
	}
)

// Wrapper Type

type (
	StorageData struct {
		Answers   map[string]bool   `json:"answers"`
		Responses map[string]string `json:"responses"`
		Bilibili  *BilibiliSettings `json:"bilibili"`
		Setting   *Setting          `json:"setting"`
		Listening *Listening        `json:"listening"`
	}

	BilibiliSettings struct {
		HighLightedUsers *set.Int64Set `json:"highLightedUsers"`
	}

	Listening struct {
		Bilibili *set.Int64Set  `json:"bilibili"`
		Youtube  *set.StringSet `json:"youtube"`
		Twitter  *set.StringSet `json:"twitter"`
	}

	Setting struct {
		VerboseDelete bool  `json:"verboseDelete"`
		Verbose       bool  `json:"verbose"`
		YearlyCheck   bool  `json:"yearlyCheck"`
		LastChecked   int64 `json:"lastChecked"`
	}
)

func (s *StorageData) toRealStorageData() *storageData {
	return &storageData{
		Answers:   s.Answers,
		Responses: s.Responses,
		Bilibili: &bilibiliSettings{
			HighLightedUsers: s.Bilibili.HighLightedUsers.ToArr(),
		},
		Setting: &setting{
			VerboseDelete: s.Setting.VerboseDelete,
			Verbose:       s.Setting.Verbose,
			YearlyCheck:   s.Setting.YearlyCheck,
			LastChecked:   s.Setting.LastChecked,
		},
		Listening: &listening{
			Bilibili: s.Listening.Bilibili.ToArr(),
			Youtube:  s.Listening.Youtube.ToArr(),
			Twitter:  s.Listening.Twitter.ToArr(),
		},
	}
}

func (s *StorageData) parse(sd *storageData) {
	s.Answers = sd.Answers
	s.Responses = sd.Responses
	s.Bilibili = &BilibiliSettings{
		HighLightedUsers: set.FromInt64Arr(sd.Bilibili.HighLightedUsers),
	}
	s.Setting = &Setting{
		VerboseDelete: sd.Setting.VerboseDelete,
		Verbose:       sd.Setting.Verbose,
		YearlyCheck:   sd.Setting.YearlyCheck,
		LastChecked:   sd.Setting.LastChecked,
	}
	s.Listening = &Listening{
		Bilibili: set.FromInt64Arr(sd.Listening.Bilibili),
		Youtube:  set.FromStrArr(sd.Listening.Youtube),
		Twitter:  set.FromStrArr(sd.Listening.Twitter),
	}
}

func makeWrapper(sd *storageData) *StorageData {
	s := &StorageData{}
	s.parse(sd)
	return s
}
