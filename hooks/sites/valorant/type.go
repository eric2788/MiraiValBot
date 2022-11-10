package valorant

import "github.com/eric2788/MiraiValBot/services/valorant"

type MatchMetaDataSub struct {
	DisplayName string                  `json:"display_name"`
	Data        *valorant.MatchMetaData `json:"data"`
}
