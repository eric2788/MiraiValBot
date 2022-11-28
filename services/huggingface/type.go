package huggingface

import (
	"encoding/json"
	"fmt"
)

type (

	// Inference API

	FaceParam struct {
		Inputs  interface{}  `json:"inputs"`
		Options *FaceOptions `json:"options"`
	}

	FaceOptions struct {
		WaitForModel bool `json:"wait_for_model"`
		UseCache     bool `json:"use_cache"`
	}

	Option func(*FaceParam)

	InferenceApi struct {
		model string
		param *FaceParam
	}

	// Space API

	SpaceApiHandler interface {
		Handle(*SpaceApi) (*SpaceResp, error)
	}

	ProcessStatus string

	SpaceWssResp struct {
		Msg ProcessStatus `json:"msg"`

		// process completed
		Output map[string]interface{} `json:"output,omitempty"`

		// estimation
		AvgEventConcurrentProcessTime float64 `json:"avg_event_concurrent_process_time,omitempty"`
		AvgEventProcessTime           float64 `json:"avg_event_process_time,omitempty"`
		Rank                          int     `json:"rank,omitempty"`
		RankETA                       float64 `json:"rank_eta,omitempty"`
		Queue                         int     `json:"queue,omitempty"`
		QueueETA                      int     `json:"queue_eta,omitempty"`
	}

	SpaceWssPush struct {
		SessionHash string        `json:"session_hash"`
		FnIndex     int           `json:"fn_index"`
		Data        []interface{} `json:"data,omitempty"`
	}

	SpaceApi struct {
		Id       string
		Data     []interface{}
		endpoint string
		Hash     string
		handler  SpaceApiHandler
	}

	// Common Resp
	SpaceResp struct {
		Data            []interface{} `json:"data"`
		Duration        float64       `json:"duration"`
		AverageDuration float64       `json:"average_duration"`
		IsGenerating    bool          `json:"is_generating"`

		Durations        []float64 `json:"durations"`
		AverageDurations []float64 `json:"average_durations"`
	}

	SpaceLabelTag struct {
		Label       string `json:"label"`
		Confidences []struct {
			Label      string  `json:"label"`
			Confidence float64 `json:"confidence"`
		} `json:"confidences"`
	}
)

func (sp *SpaceResp) GetDuration() float64 {
	if sp.Duration > 0 {
		return sp.Duration
	}
	if len(sp.Durations) > 0 {
		return sp.Durations[0]
	}
	return 0
}

func (sp *SpaceResp) GetAverageDuration() float64 {
	if sp.AverageDuration > 0 {
		return sp.AverageDuration
	}
	if len(sp.AverageDurations) > 0 {
		return sp.AverageDurations[0]
	}
	return 0
}

func (sp *SpaceResp) ParseData(index int, arg interface{}) error {
	if index >= len(sp.Data) {
		return fmt.Errorf("index out of range: %d", index)
	}
	m, ok := sp.Data[index].(map[string]interface{})
	if !ok {
		return fmt.Errorf("cannot parse data: Data[%d] is not map type", index)
	}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, arg)
}
