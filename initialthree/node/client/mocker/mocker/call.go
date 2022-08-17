package mocker

import codecs "initialthree/codec/cs"

type call struct {
	seq   uint32
	c     chan *call
	err   error
	reply *codecs.Message
}

func newCall() *call {
	return &call{c: make(chan *call, 1)}
}

func (c *call) doneWithErr(err error) {
	c.err = err
	c.done()
}

func (c *call) doneWithReply(reply *codecs.Message) {
	c.reply = reply
	c.done()
}

func (c *call) done() {
	select {
	case c.c <- c:
	default:
	}
}

func (c *call) Done() (*codecs.Message, error) {
	<-c.c
	return c.reply, c.err
}
