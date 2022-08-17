package mocker

import (
	"log"

	cs_msg "initialthree/protocol/cs/message"
)

func (cu *ClientUser) actionDrawCardStoreIn() {
	resp := cu.call(&cs_msg.DrawCardStoreInToS{})
	log.Printf("client user draw card store in resp:%v", resp.GetData().(*cs_msg.DrawCardStoreInToC))
	mustBeZero(resp.GetErrCode())
}
