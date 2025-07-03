package holder

import utils "tiktok_e-commerce/common/infra/hot_key_client/model/util"

var (
	TurnKeyCollector *utils.TurnKeyCollector
)

func init() {
	TurnKeyCollector = utils.NewTurnKeyCollector()
}
