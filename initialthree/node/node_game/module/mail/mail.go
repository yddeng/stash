package mail

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	flyfish "github.com/sniperHW/flyfish/client"
	"github.com/sniperHW/flyfish/errcode"
	"initialthree/node/common/db"
	"initialthree/node/common/offlinemsg"
	"initialthree/node/node_game/global/mail"
	"initialthree/node/node_game/module"
	"initialthree/pkg/json"
	"initialthree/protocol/cs/message"
	"initialthree/zaplogger"
	"sort"
	"time"
)

const (
	tableName   = "mail"
	fieldData   = "data"
	fieldStatus = "status"
	slotCount   = 10
	spaceCount  = 100
)

func slotName(slotIdx int) string {
	return fmt.Sprintf("slot%d", slotIdx)
}

func calcSlotIdx(id int) int {
	return id % slotCount
}

type statusType int

const (
	statusOffline = 0
	statusGlobal  = 1
)

func genStatusID(t statusType, id int) int {
	return int(t)*1000 + id
}

type itemStatus struct {
	Idx     int        `json:"idx"`
	Type    statusType `json:"t"`
	Version int64      `json:"v"`
	Read    bool       `json:"r"`
}

type mailData struct {
	GenID uint32 `json:"genId"`
}

type Mail struct {
	userI module.UserI

	// 数据
	mails    []map[uint32]*message.Mail
	status   map[int]*itemStatus
	mailData mailData

	// 临时数据
	dbClient       offlinemsg.DB
	mailCount      int
	offlineVersion int64

	dirty map[uint32]*message.Mail
	*module.ModuleSaveBase
}

// 1）	邮件容量上限为100
// 2）	当邮件超出上限时，自动删除超出邮件
//	①	优先删除已读邮件中邮件生成时间最早的
//	②	无已读邮件可删除时，删除未读邮件中邮件生成时间最早的
func sortMail(mails []*message.Mail) []*message.Mail {
	var iRead, jRead bool
	var iTime, jTime int64
	sort.Slice(mails, func(i, j int) bool {
		iRead, jRead = mails[i].GetRead(), mails[j].GetRead()
		iTime, jTime = mails[i].GetCreateTime(), mails[j].GetCreateTime()
		if iRead && !jRead {
			return true
		} else if iRead == jRead {
			return iTime <= jTime
		}
		return false
	})
	return mails
}

func (this *Mail) AddMails(mails []*message.Mail) {
	this.addMails(mails)
}

func (this *Mail) addMails(mails []*message.Mail) {
	addLength := len(mails)
	if addLength >= spaceCount {
		// 插入的邮件数大于等于容量，移除全部邮件
		if this.mailCount > 0 {
			delIds := make([]uint32, 0, this.mailCount)
			for _, slot := range this.mails {
				for _, m := range slot {
					delIds = append(delIds, m.GetId())
				}
			}
			this.MailDelete(delIds)
		}

		if addLength != spaceCount {
			// 大于容量，仅插入排序后 容量数的邮件
			delCount := addLength - spaceCount
			mails = sortMail(mails)
			mails = mails[delCount:]
		}
	} else if addLength+this.mailCount > spaceCount {
		// 插入的邮件数与已有的大于等于容量，移除已有邮件排序后的数目
		oldMails := make([]*message.Mail, 0, this.mailCount)
		for _, slot := range this.mails {
			for _, m := range slot {
				oldMails = append(oldMails, m)
			}
		}
		oldMails = sortMail(oldMails)
		delCount := addLength + this.mailCount - spaceCount
		delIds := make([]uint32, delCount)
		for i := 0; i < delCount; i++ {
			delIds[i] = oldMails[i].GetId()
		}
		this.MailDelete(delIds)
	}

	// add
	for _, m := range mails {
		// 实例化ID
		id := this.genID()
		m.Id = proto.Uint32(id)

		slotIdx := calcSlotIdx(int(id))
		this.mails[slotIdx][id] = m
		this.mailCount++
		this.dirty[id] = m
		this.SetDirty(slotIdx)
	}

}

func (this *Mail) genID() uint32 {
	this.mailData.GenID++
	this.SetDirty(fieldData)
	return this.mailData.GenID
}

func (this *Mail) GetMail(id uint32) *message.Mail {
	slotIdx := calcSlotIdx(int(id))
	return this.mails[slotIdx][id]
}

func (this *Mail) IsExpire(m *message.Mail) bool {
	if m.GetExpireTime() != 0 && time.Now().Unix() > m.GetExpireTime() {
		return true
	}
	return false
}

func (this *Mail) MailRead(mailIds []uint32) {
	for _, id := range mailIds {
		slotIdx := calcSlotIdx(int(id))
		if m, ok := this.mails[slotIdx][id]; ok && !m.GetRead() {
			m.Read = proto.Bool(true)
			this.dirty[id] = m
			this.SetDirty(slotIdx)
		}
	}
}

func (this *Mail) MailDelete(mailIds []uint32) {
	for _, id := range mailIds {
		slotIdx := calcSlotIdx(int(id))
		slot := this.mails[slotIdx]
		if _, ok := slot[id]; ok {
			delete(slot, id)
			this.mailCount--
			this.dirty[id] = &message.Mail{
				Id:      proto.Uint32(id),
				Deleted: proto.Bool(true),
			}
			this.SetDirty(slotIdx)
		}
	}
}

func (this *Mail) itemDispose(items []*offlinemsg.Item, t statusType) {

	addMails := make([]*message.Mail, 0, len(items))
	for _, it := range items {
		id := genStatusID(t, it.Idx)
		status, ok := this.status[id]
		if !ok {
			status = &itemStatus{
				Idx:     it.Idx,
				Type:    t,
				Version: it.Version,
			}
			this.status[id] = status
			this.SetDirty(fieldStatus)
		} else if it.Version != status.Version {
			status.Version = it.Version
			status.Read = false
			this.SetDirty(fieldStatus)
		}

		if !status.Read {
			status.Read = true
			this.SetDirty(fieldStatus)

			var m *message.Mail
			if err := json.Unmarshal(it.Content, &m); err == nil {
				// todo 条件判断
				if !this.IsExpire(m) {
					addMails = append(addMails, m)
				}
			}

		}
	}

	if len(addMails) > 0 {
		this.addMails(addMails)
	}

}

func (this *Mail) PullGlobalMailData() {
	mail.GetGlobalMail(0, func(items []*offlinemsg.Item, nowVersion int64) {
		if len(items) > 0 {
			this.itemDispose(items, statusGlobal)
		}
	})
}

func (this *Mail) PullOfflineMailData() {
	offlinemsg.PullMsg(this.dbClient, this.userI.GetIDStr(), "mail", this.offlineVersion, time.Second*6, func(err errcode.Error, version int64, items []*offlinemsg.Item) {
		if err == nil {
			if len(items) > 0 {
				this.itemDispose(items, statusOffline)
			}
		} else {
			zaplogger.GetSugar().Error(err)
		}
	})
}

func (this *Mail) AfterInitAll() error {
	this.dbClient = offlinemsg.NewFlyfishDB(db.GetFlyfishClient("global"))

	this.PullGlobalMailData()
	this.PullOfflineMailData()
	return nil
}

func (this *Mail) ModuleType() module.ModuleType {
	return module.Mail
}

func (this *Mail) Tick(now time.Time) {}

func (this *Mail) ReadOut() *module.ReadOutCommand {
	out := &module.ReadOutCommand{
		Table:  tableName,
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]string, 0, slotCount+2),
	}

	out.Fields = append(out.Fields, fieldData)
	out.Fields = append(out.Fields, fieldStatus)
	for i := 0; i < slotCount; i++ {
		out.Fields = append(out.Fields, slotName(i))
	}

	return out
}

func (this *Mail) WriteCommand(fields map[interface{}]struct{}) *module.WriteBackCommand {
	cmd := &module.WriteBackCommand{
		Table:  this.ModuleType().String(),
		Key:    this.userI.GetIDStr(),
		Module: this,
		Fields: make([]*module.WriteBackFiled, 0, len(fields)),
	}

	for field := range fields {
		switch field.(type) {
		case string:
			fieldName := field.(string)
			var data []byte
			switch fieldName {
			case fieldData:
				data, _ = json.Marshal(this.mailData)
			case fieldStatus:
				data, _ = json.Marshal(this.status)
			}
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  fieldName,
				Value: data,
			})
		case int:
			idx := field.(int)
			slot := this.mails[idx]
			data, _ := json.Marshal(slot)
			cmd.Fields = append(cmd.Fields, &module.WriteBackFiled{
				Name:  slotName(idx),
				Value: data,
			})

		}
	}

	return cmd
}

func (this *Mail) FlushDirtyToClient() {
	if len(this.dirty) > 0 {
		msg := &message.MailSyncToC{
			IsAll: proto.Bool(false),
			Mails: make([]*message.Mail, 0, len(this.dirty)),
		}
		for _, m := range this.dirty {
			msg.Mails = append(msg.Mails, m)
		}

		this.userI.Post(msg)
		this.dirty = map[uint32]*message.Mail{}

	}
}

func (this *Mail) FlushAllToClient(seqNo ...uint32) {
	msg := &message.MailSyncToC{
		IsAll: proto.Bool(true),
		Mails: make([]*message.Mail, 0, this.mailCount),
	}

	for _, slot := range this.mails {
		for _, m := range slot {
			msg.Mails = append(msg.Mails, m)
		}
	}

	this.userI.Post(msg)
	this.dirty = map[uint32]*message.Mail{}
}

func (this *Mail) Init(fields map[string]*flyfish.Field) error {
	for i := 0; i < slotCount; i++ {
		fieldName := slotName(i)
		field, ok := fields[fieldName]
		this.mails[i] = map[uint32]*message.Mail{}
		if !ok || len(field.GetBlob()) == 0 {
			this.SetDirty(i)
		} else {
			if err := json.Unmarshal(field.GetBlob(), &this.mails[i]); err != nil {
				return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
			}

			// 移除已经超时的邮件
			for id, v := range this.mails[i] {
				if this.IsExpire(v) {
					delete(this.mails[i], id)
				}
			}

			this.mailCount += len(this.mails[i])
			this.SetDirty(i)
		}
	}

	field, ok := fields[fieldData]
	if ok && len(field.GetBlob()) != 0 {
		if err := json.Unmarshal(field.GetBlob(), &this.mailData); err != nil {
			return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
		}
	} else {
		this.SetDirty(fieldData)
	}

	field, ok = fields[fieldStatus]
	if ok && len(field.GetBlob()) != 0 {
		if err := json.Unmarshal(field.GetBlob(), &this.status); err != nil {
			return fmt.Errorf("unmarshal:%s %s", string(field.GetBlob()), err)
		}
	} else {
		this.SetDirty(fieldStatus)
	}

	return nil
}

func init() {
	module.RegisterModule(module.Mail, func(userI module.UserI) module.ModuleI {
		m := &Mail{
			userI:  userI,
			mails:  make([]map[uint32]*message.Mail, slotCount),
			status: map[int]*itemStatus{},
			dirty:  map[uint32]*message.Mail{},
		}
		m.ModuleSaveBase = module.NewModuleSaveBase(m)
		return m
	})
}
