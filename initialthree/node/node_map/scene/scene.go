package scene

import (
	"github.com/golang/protobuf/proto"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"initialthree/node/common/aoi"
	"initialthree/node/common/config"
	"initialthree/node/table/excel/ConstTable/Global"
	"initialthree/node/table/excel/DataTable/Map"
	"initialthree/pkg/timer"
	"initialthree/protocol/cs/message"
	"initialthree/protocol/ss/ssmessage"
	"time"

	"initialthree/pkg/event"
	"initialthree/pkg/golog"
)

var logger golog.LoggerI
var hearbeatTimeout = time.Second * 10

type SceneObject interface {
	aoi.User
	SetAoiEntity(en aoi.Entity)
	GetAoiEntity() aoi.Entity
	SetScene(*Scene)
	GetScene() *Scene
	GetAoiPosition() aoi.Position
	GetUserID() string
	GetID() uint64
	Tick(time.Time)
	ClearMoved()
	PackRoleData() *message.ViewObj
	Heartbeat() time.Time
}

type Position struct {
	X     int32
	Y     int32
	Z     int32
	Angle int32
}

type Scene struct {
	sceneIdx      int32
	aoiMgr        aoi.Manager
	objects       map[string]SceneObject
	moveAOIEntity map[aoi.Entity]bool // 移动的aoi对象
	processQueue  *event.EventQueue

	lastTime int64
}

func newScene(idx int32) *Scene {
	logger.Debugln("newScene", idx)
	s := &Scene{
		sceneIdx:      idx,
		aoiMgr:        aoi.NewOrth(AOIRadius),
		objects:       map[string]SceneObject{},
		moveAOIEntity: map[aoi.Entity]bool{},
		processQueue:  event.NewEventQueue(),
	}

	go func() {
		s.processQueue.Run()
		logger.Debugln("Scene destroy", s.sceneIdx)
	}()

	// 每秒tick两次
	cluster.RegisterTimer(time.Millisecond*500, func(i *timer.Timer, _ interface{}) {
		if nil != s.processQueue.Post(s.tick) {
			logger.Debugln("Scene destroy2", s.sceneIdx)
			return
		}
	}, nil)

	return s
}

func (this *Scene) PostTask(task func()) {
	this.processQueue.Post(task)
}

func (this *Scene) GetObjByID(id string) SceneObject {
	return this.objects[id]
}

func (this *Scene) StartAoi(o SceneObject) {
	if o.GetScene() == this {
		entity, err := this.aoiMgr.Add(o.GetUserID(), o.GetAoiPosition(), o)
		if err != nil {
			logger.Errorf("aoi add err:%s", err)
			return
		}
		o.SetAoiEntity(entity)
		logger.Infoln(o.GetUserID(), "StarAoi OK")
	}
}

func (this *Scene) EnterScene(o SceneObject) bool {
	if nil != o.GetScene() {
		logger.Debugln("EnterScene o.GetScene != nil")
		return false
	}

	o.SetScene(this)
	this.objects[o.GetUserID()] = o
	return true
}

func (this *Scene) LeaveScene(o SceneObject) bool {
	if this != o.GetScene() {
		logger.Debugln("LeaveScene this != o.GetScene")
		return false
	}

	entity := o.GetAoiEntity()
	if entity != nil {
		logger.Debugln("aoiEntity remove")
		err := this.aoiMgr.Rem(entity)
		if err != nil {
			logger.Debugln("aoi remove failed")
			return false
		}
		delete(this.moveAOIEntity, entity)
	}

	o.SetScene(nil)
	delete(this.objects, o.GetUserID())
	return true
}

func (this *Scene) Move(o SceneObject) bool {
	if this != o.GetScene() {
		return false
	}

	aoiEntity := o.GetAoiEntity()
	if o.GetAoiEntity() == nil {
		return false
	}

	this.moveAOIEntity[aoiEntity] = true
	return true

}

func (this *Scene) tick() {
	now := time.Now()

	// 心跳超时
	for _, o := range this.objects {
		if now.After(o.Heartbeat().Add(hearbeatTimeout)) {
			this.LeaveScene(o)
		}
	}

	for o := range this.moveAOIEntity {
		obj := o.User().(SceneObject)
		o.Move(obj.GetAoiPosition())
	}

	for _, o := range this.objects {
		o.Tick(now)
	}

	for _, o := range this.objects {
		o.ClearMoved()
	}

}

type SceneMgr struct {
	selfAddr   addr.Addr
	worldLogic addr.LogicAddr
	mapID      uint32
	sceneMap   []*Scene

	leftBottom aoi.Position
	rightTop   aoi.Position
}

var sceneMgr *SceneMgr

func FixPosition(pos Position) Position {

	if pos.X < sceneMgr.leftBottom.X {
		pos.X = sceneMgr.leftBottom.X
	}

	if pos.Y < sceneMgr.leftBottom.Y {
		pos.Y = sceneMgr.leftBottom.Y
	}

	if pos.X >= sceneMgr.rightTop.X {
		pos.X = sceneMgr.rightTop.X - 1
	}

	if pos.Y >= sceneMgr.rightTop.Y {
		pos.Y = sceneMgr.rightTop.Y - 1
	}

	return pos
}

func GetSelfWorld() addr.LogicAddr {
	return sceneMgr.worldLogic
}

func GetScene(idx int32) *Scene {
	return sceneMgr.sceneMap[idx]
}

func GetScenes() []*Scene {
	return sceneMgr.sceneMap
}

// 给每一个scene投递
func AllSceneDo(fn func()) {
	for _, s := range sceneMgr.sceneMap {
		scene_ := s
		scene_.PostTask(fn)
	}
}

func Init(baseLogger golog.LoggerI, conf *config.Map) error {
	logger = baseLogger

	logic, err := addr.MakeLogicAddr(conf.WorldAddr)
	if err != nil {
		return err
	}

	mapId := logic.Server()
	mapDef := Map.GetID(int32(mapId))

	leftBottom := aoi.Position{
		X: 0,
		Y: 0,
	}

	rightTop := aoi.Position{
		X: mapDef.Length * int32(Global.GetID(1).PositionRate),
		Y: mapDef.Width * int32(Global.GetID(1).PositionRate),
	}

	sceneMgr = &SceneMgr{
		worldLogic: logic,
		mapID:      mapId,
		sceneMap:   make([]*Scene, 0, SceneSliceLen),
		leftBottom: leftBottom,
		rightTop:   rightTop,
	}

	for i := 0; i < SceneSliceLen; i++ {
		sceneMgr.sceneMap = append(sceneMgr.sceneMap, newScene(int32(i)))
	}

	cluster.RegisterTimer(time.Millisecond*1000, func(timer *timer.Timer, _ interface{}) {
		msg := &ssmessage.MapToWorld{
			MapAddr: proto.String(cluster.SelfAddr().Logic.String()),
			Scenes:  make([]*ssmessage.Scene, 0, len(sceneMgr.sceneMap)),
		}
		for _, s := range sceneMgr.sceneMap {
			msg.Scenes = append(msg.Scenes, &ssmessage.Scene{
				SceneID:  proto.Int32(s.sceneIdx),
				AOICurr:  proto.Int32(int32(len(s.objects))),
				AOITotal: proto.Int32(SceneAoiMax),
			})
		}

		cluster.PostMessage(sceneMgr.worldLogic, msg)
	}, nil)

	return nil
}
