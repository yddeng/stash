package IDGen

import (
	"fmt"
	"time"
)

const beginTimeStamp int64 = 1514736000 //2018-1-1
const AreaIDMax uint16 = 4095

var genError error = fmt.Errorf("beyond gen limit")
var timeError error = fmt.Errorf("time rollback")

type IDGenerator struct {
	areaID   uint16 //12区ID支持4095个区
	serverID byte   //8位服务器ID支持255台服务
	counter  uint16 //14位自增值，支持每秒16383个ID生成
	time     int64  //30位自beginTimeStamp以来的秒数，支持34年
	prefix   uint64
	start    time.Time
}

/*

注意：

进程运行过程中出现时间回调不会影响ID生成算法的正确性(保证无重复)
但是如果进程运行一段时间关闭之后，将机器时间调到之前的一个时间点
再重新运行进程将会产生重复ID。为了避免这种情况，应该将每个生成器
实例的start序列化保存下来(例如数据库)，每次创建生成器之后将实例
的start与保存下来的start做比较，如果发现时间出现了回调及时终止
进程并报告错误。

*/

func (this *IDGenerator) GetStartTimeStamp() time.Time {
	return this.start
}

func (this *IDGenerator) getSecond() int64 {
	now := time.Now()
	elapse := now.Sub(this.start) //Monotonic elapse
	return this.start.Unix() + (int64(elapse / time.Second))
}

func New(areaID uint16, serverID byte) *IDGenerator {
	if areaID > AreaIDMax {
		return nil
	}

	now := time.Now()

	if now.Unix() < beginTimeStamp {
		return nil
	}

	generator := &IDGenerator{}
	generator.areaID = areaID
	generator.serverID = serverID
	generator.start = now
	generator.time = generator.getSecond() - beginTimeStamp
	generator.prefix = (uint64(areaID) << 52) | (uint64(serverID) << 44)
	return generator
}

func (this *IDGenerator) Gen() (uint64, error) {
	timestamp := this.getSecond() - beginTimeStamp

	if timestamp < this.time {
		return 0, timeError
	} else if timestamp > this.time {
		this.time = timestamp
		this.counter = 0
	} else if this.counter+1 > uint16(0x3FFF) {
		return 0, genError
	}

	this.counter++
	id := this.prefix | (uint64(this.time) << 14) | uint64(this.counter)
	return id, nil
}
