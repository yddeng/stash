package inoutput

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"initialthree/node/common/enumType"
	"initialthree/node/table/excel/DataTable/PlayerCharacter"
	"initialthree/protocol/cs/message"
)

var (
	ErrInputNotEnough  = errors.New("inoutput input not enough")
	ErrCfgNotFound     = errors.New("inoutput output config is not found")
	ErrSpaceNotEnough  = errors.New("inoutput space is not enough")
	ErrInvalidResType  = errors.New("inoutput input invalid res type")
	ErrInvalidResCount = errors.New("inoutput invalid res count")
)

/*
 * 资源的描述
 */
type ResDesc struct {
	Type  int   //资源类型
	ID    int32 //资源ID
	Count int32 //数量
	// 实例ID，输出资源是装备、武器时，用于实例化ID。若输出资源时角色，用于实例化武器ID
	// 若消耗资源是道具，用于道具的格子ID
	InsID interface{}
}

type RollBackFunc func() // 回滚函数
type CommitFunc func()   // 整个流程执行完，要做的。不可撤销的操作

type Transaction struct {
	rollbacks []RollBackFunc
	commits   []CommitFunc
}

func (this *Transaction) removeRes(c Container, r ResDesc) error {
	err, rb, cmt := c.RemoveRes(r)
	if nil == err {
		if nil != rb {
			this.rollbacks = append(this.rollbacks, rb)
		}
		if nil != cmt {
			this.commits = append(this.commits, cmt)
		}
	}

	return err
}

func (this *Transaction) addRes(c Container, r ResDesc) (error, []ResDesc) {
	err, rb, cmt, transfer := c.AddRes(r)
	if nil == err {
		if nil != rb {
			this.rollbacks = append(this.rollbacks, rb)
		}
		if nil != cmt {
			this.commits = append(this.commits, cmt)
		}
	}
	return err, transfer
}

func (this *Transaction) rollbackAll() {
	i := len(this.rollbacks) - 1
	for ; i >= 0; i-- {
		this.rollbacks[i]()
	}
	this.rollbacks = nil
}

func (this *Transaction) doCommit() {
	for i := 0; i < len(this.commits); i++ {
		this.commits[i]()
	}
	this.commits = nil
}

/*
 *  存放资源的容器要实现的接口
 *  在游戏中这个容器包括物品管理器，装备管理器，角色管理器，属性管理器等等
 */
type Container interface {
	RemoveRes(ResDesc) (error, RollBackFunc, CommitFunc)
	AddRes(ResDesc) (error, RollBackFunc, CommitFunc, []ResDesc)
}

//容器管理器，根据资源类型获取正确的管理容器
type ContainerMgr interface {
	GetContainer(int) Container
	RegisterContainer(int, Container)
}

/*
 *  资源输入输出算法流程
 *
 *  根据输入描述从正确的Container中扣除资源
 *  根据输出描述向正确的Container添加资源
 *
 *  需要保证输入资源数量足够以及所有输出均可被添加到输出容器中（例如有容量限制）
 */

func DoInputOutput(mgr ContainerMgr, input []ResDesc, output []ResDesc) error {

	trans := Transaction{
		rollbacks: make([]RollBackFunc, 0, len(input)+len(output)),
		commits:   []CommitFunc{},
	}

	if err := check(input, output); err != nil {
		return err
	}

	for _, v := range input {
		c := mgr.GetContainer(v.Type)
		if nil == c {
			trans.rollbackAll()
			return ErrInvalidResType
		} else {
			err := trans.removeRes(c, v)
			if nil != err {
				trans.rollbackAll()
				return err
			}
		}
	}

	l := len(output)
	for i := 0; i < l; i++ {
		v := output[i]
		c := mgr.GetContainer(v.Type)
		if nil == c {
			trans.rollbackAll()
			return ErrInvalidResType
		} else {
			err, transfer := trans.addRes(c, v)
			if nil == err {
				output = append(output, transfer...)
			} else {
				trans.rollbackAll()
				return err
			}
		}
		l = len(output)
	}

	trans.doCommit()

	return nil

}

func check(input []ResDesc, output []ResDesc) error {
	for _, v := range input {
		if v.Count <= 0 {
			return ErrInvalidResCount
		}
	}
	for _, v := range output {
		if v.Count <= 0 {
			return ErrInvalidResCount
		}
	}
	return nil
}

func OutputIns2Award(out []ResDesc) *message.Award {
	award := &message.Award{
		AwardInfos: make([]*message.AwardInfo, 0, len(out)),
	}

	for _, v := range out {
		if v.Count == 1 && v.InsID != nil {
			switch v.Type {
			case enumType.IOType_Equip, enumType.IOType_Weapon:
				dropType := enumType.DropType_Equip
				if v.Type == enumType.IOType_Weapon {
					dropType = enumType.DropType_Weapon
				}

				award.AwardInfos = append(award.AwardInfos, &message.AwardInfo{
					Type:  proto.Int32(int32(dropType)),
					ID:    proto.Int32(v.ID),
					Count: proto.Int32(1),
					InsID: proto.Uint32(v.InsID.(uint32)),
				})
			case enumType.IOType_Character:
				award.AwardInfos = append(award.AwardInfos, &message.AwardInfo{
					Type:  proto.Int32(int32(enumType.DropType_Character)),
					ID:    proto.Int32(v.ID),
					Count: proto.Int32(1),
				})
				def := PlayerCharacter.GetID(v.ID)
				award.AwardInfos = append(award.AwardInfos, &message.AwardInfo{
					Type:  proto.Int32(int32(enumType.DropType_Weapon)),
					ID:    proto.Int32(def.DefaultWeapon),
					Count: proto.Int32(1),
					InsID: proto.Uint32(v.InsID.(uint32)),
				})
			default:
			}
		} else {
			award.AwardInfos = append(award.AwardInfos, &message.AwardInfo{
				Type:  proto.Int32(int32(IOType2DropType(v.Type))),
				ID:    proto.Int32(v.ID),
				Count: proto.Int32(v.Count),
			})
		}
	}
	return award
}

func IOType2DropType(ioType int) int {
	var dropType int
	switch ioType {
	case enumType.IOType_Equip:
		dropType = enumType.DropType_Equip
	case enumType.IOType_Weapon:
		dropType = enumType.DropType_Weapon
	case enumType.IOType_Character:
		dropType = enumType.DropType_Character
	case enumType.IOType_UsualAttribute:
		dropType = enumType.DropType_UsualAttribute
	case enumType.IOType_Item:
		dropType = enumType.DropType_Item
	default:
		panic(fmt.Sprintf("inoutput ioType trans to dropType %d not found", ioType))
	}
	return dropType
}
