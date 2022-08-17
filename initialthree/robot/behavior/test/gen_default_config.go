package main

import (
	"flag"
	"initialthree/node/common/attr"
	. "initialthree/robot/behavior"
	"initialthree/robot/types"

	. "github.com/GodYY/bevtree"
)

func main() {
	loadall := flag.Bool("load-all", false, "load all trees on initailization")

	flag.Parse()

	configPath := "./config.xml"
	exporter := CreateExporter()
	exporter.SetLoadAll(*loadall)

	logintree := NewTree("登录")
	exporter.AddTree(logintree, "bevtree-login.xml")

	logintree.SetComment("登录到服务器，如果是首次登录，创建角色")
	seq := NewSequenceNode()
	logintree.Root().SetChild(seq)

	loginNode := NewBevNode(&BevLogin{ReloginDelay: 1000})
	loginNode.SetComment("登录到服务器")
	seq.AddChild(loginNode)

	createRoleNode := NewBevNode(&BevCreateRole{})
	createRoleNode.SetComment("初次登录，创建角色")
	seq.AddChild(createRoleNode)

	logictree := NewTree("逻辑")
	exporter.AddTree(logictree, "bevtree-logic.xml")

	logictree.SetComment("运行逻辑")

	selector := NewSelectorNode()
	logictree.Root().SetChild(selector)

	dungeonCondition := NewConditionNode()
	dungeonCondition.SetComment("闯关")
	selector.AddChild(dungeonCondition)
	dungeonCondition.AddCond(NewCondAttr(NewAttr(attr.CurrentFatigue, 100, GreaterEqual)))
	dungeonSelector := NewSelectorNode()
	dungeonCondition.SetChild(dungeonSelector)
	mainDungeonSelector := NewSelectorNode()
	mainDungeonSelector.SetComment("主线")
	dungeonSelector.AddChild(mainDungeonSelector)
	mainDungeonProgressSequence := NewSequenceNode()
	mainDungeonProgressSequence.SetComment("主线进度挑战")
	mainDungeonSelector.AddChild(mainDungeonProgressSequence)
	mainDungeonProgressChallenge := NewBevNode(&BevMainDungeon{Random: false})
	mainDungeonProgressSequence.AddChild(mainDungeonProgressChallenge)
	mainDungeonClaimStaraward := NewBevNode(&BevClaimStaraward{AwardType: StarAward_MainDungeon})
	mainDungeonProgressSequence.AddChild(mainDungeonClaimStaraward)
	mainDungeonRandomChallenge := NewBevNode(&BevMainDungeon{Random: true})
	mainDungeonRandomChallenge.SetComment("主线随机挑战")
	mainDungeonSelector.AddChild(mainDungeonRandomChallenge)

	otherSelector := NewSelectorNode()
	otherSelector.SetComment("其它")
	selector.AddChild(otherSelector)
	recoverFatigueCondition := NewConditionNode()
	recoverFatigueCondition.SetComment("体力低于100时，恢复100点体力")
	otherSelector.AddChild(recoverFatigueCondition)
	recoverFatigueCondition.AddCond(NewCondAttr(NewAttr(attr.CurrentFatigue, 100, Lower)))
	recoverFatigue := NewBevNode(&BevGameMaster{Desc: "恢复体力100点", Cmd: []GMCmd{{Type: GMCmd_AddAttr, ID: attr.CurrentFatigue, Count: 100}}})
	recoverFatigueCondition.SetChild(recoverFatigue)

	compTree := NewTree("综合")
	exporter.AddTree(compTree, "bevtree-comp.xml")

	compTree.SetComment("综合登录和逻辑")
	compSelector := NewSelectorNode()
	compTree.Root().SetChild(compSelector)
	condition := NewConditionNode()
	condition.SetComment("是否没有登录")
	compSelector.AddChild(condition)
	condition.AddCond(NewCondStatus(true, NewStatusBool(types.Status_IsLogin, false)))
	condition.SetChild(NewSubtreeNode(logintree, false))
	compSelector.AddChild(NewSubtreeNode(logictree, false))

	exporter.Export(configPath)
}
