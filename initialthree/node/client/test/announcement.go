package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	fnet "github.com/sniperHW/flyfish/pkg/net"
	"initialthree/cs"
	"initialthree/node/client/dispatcher"
	"initialthree/pkg/event"
	"initialthree/pkg/json"
	"initialthree/protocol/cmdEnum"
	"log"
	"os"

	codecs "initialthree/codec/cs"
	cs_msg "initialthree/protocol/cs/message"
)

type announcements struct {
	Groups []announcementGroup `json:"group"`
}

type announcementGroup struct {
	ID       int32                 `json:"id"`
	Title    string                `json:"title"`
	Contents []announcementContent `json:"contents"`
}

type announcementContent struct {
	ID      int32  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf(" addr userID\n")
		return
	}

	processQueue := event.NewEventQueue()
	go func() {
		processQueue.Run()
		fmt.Println("queue break-------------")
	}()

	userID := os.Args[2]
	addr := os.Args[1]

	dis := dispatcher.New(processQueue)
	dis.RegisterOnce("Establish", func(session *fnet.Socket) {
		session.Send(codecs.NewMessage(uint32(0), &cs_msg.AnnouncementToS{
			UserID: proto.String(userID),
		}))
	})

	dis.Register(cmdEnum.CS_Announcement, func(session *fnet.Socket, msg *codecs.Message) {

		data := msg.GetData().(*cs_msg.AnnouncementToC)

		announcements := new(announcements)
		if err := json.Unmarshal([]byte(data.GetAnnouncement()), announcements); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("AnnouncementToC %v\n", announcements)

	})

	cs.DialTcp(addr, 0, dis)

	sigStop := make(chan bool)
	_, _ = <-sigStop
}
