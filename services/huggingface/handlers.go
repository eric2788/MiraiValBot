package huggingface

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/eric2788/common-utils/request"
	"github.com/gorilla/websocket"
)

const (
	SpaceApiUrl = "https://%s.hf.space/%s"
	SpaceWssUrl = "wss://%s.hf.space/queue/join"

	ProcessCompleted ProcessStatus = "process_completed"
	ProcessStarts    ProcessStatus = "process_starts"
	Estimation       ProcessStatus = "estimation"
	ProcessError     ProcessStatus = "process_error"
	SendData         ProcessStatus = "send_data"
	SendHash         ProcessStatus = "send_hash"
)

type (
	httpRequestHandler struct {
	}

	websocketHandler struct {
	}
)

func (w *websocketHandler) Handle(s *SpaceApi) (*SpaceResp, error) {
	url := fmt.Sprintf(SpaceWssUrl, s.Id)
	logger.Debugf("Requesting URL: %s", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// max 20 mins
	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			return nil, fmt.Errorf("timeout")
		default:
			_, message, err := conn.ReadMessage()
			if err != nil {
				return nil, err
			}
			var resp SpaceWssResp
			err = json.Unmarshal(message, &resp)
			if err != nil {
				return nil, err
			}
			switch resp.Msg {
			case ProcessCompleted:
				logger.Debugf("Process Completed: %v", resp.Output)
				return &SpaceResp{
					Data:            resp.GetStringData(),
					Duration:        resp.AvgEventProcessTime,
					AverageDuration: resp.AvgEventConcurrentProcessTime,
				}, nil
			case ProcessError:
				return nil, fmt.Errorf("error: %+v", resp)
			case ProcessStarts:
				logger.Debugf("Process starts")
			case Estimation:
				logger.Debugf("Estimation: %d", resp.QueueETA)
			case SendData:
				logger.Debugf("Send data: %v", s.Data)
				if err = conn.WriteJSON(SpaceWssPush{
					FnIndex:     1,
					Data:        s.Data,
					SessionHash: s.Hash,
				}); err != nil {
					logger.Errorf("write json error on send data: %s", err)
					return nil, err
				}
			case SendHash:
				logger.Debugf("Send hash: %v", s.Hash)
				if err = conn.WriteJSON(SpaceWssPush{
					SessionHash: s.Hash,
					FnIndex:     1,
				}); err != nil {
					logger.Errorf("write json error on send hash: %s", err)
					return nil, err
				}
			default:
				logger.Debugf("Unknown msg: %s", resp.Msg)
			}
		}
	}
}

func (h *httpRequestHandler) Handle(s *SpaceApi) (*SpaceResp, error) {
	body, err := json.Marshal(map[string]interface{}{
		"data": s.Data,
	})
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf(SpaceApiUrl, s.Id, s.endpoint)
	logger.Debugf("Requesting URL: %s", url)
	res, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err == nil && res.StatusCode != 200 {
		defer res.Body.Close() // only close when no err but non 200 code
		if b, berr := io.ReadAll(res.Body); berr == nil {
			err = fmt.Errorf(string(b))
		} else {
			err = fmt.Errorf(res.Status)
		}
		return nil, err
	}
	defer res.Body.Close()
	var resp SpaceResp
	err = request.Read(res, &resp)
	return &resp, err
}
