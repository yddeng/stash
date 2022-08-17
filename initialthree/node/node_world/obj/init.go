package obj

import (
	"initialthree/node/node_world/world"
	"initialthree/pkg/golog"
	"time"
)

var logger golog.LoggerI

func Init(baseLogger golog.LoggerI, worldID uint32) {
	logger = baseLogger

	trash := &TrashCan{
		UserID: "trash",
		ID:     5555555,
		Pos: &world.Position{
			X:     0,
			Y:     0,
			Z:     0,
			Angle: 0,
		},
		Status: 0,
	}

	world.AddWorldObj(trash)

	boss := &Boss{
		UserID: "boss",
		ID:     6666666,
		Pos: &world.Position{
			X:     0,
			Y:     0,
			Z:     0,
			Angle: 0,
		},
		Status: 1,
	}

	world.AddWorldObj(boss)

	go func() {
		time.Sleep(time.Second * 5)
		world.EnterAllMap(boss, nil)
		boss.StartAi()
	}()

}
