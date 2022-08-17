package behavior

import (
	"initialthree/robot/types"
	"testing"

	. "github.com/GodYY/bevtree"
)

func TestMarshalUnmarshal(t *testing.T) {
	tree := NewTree("test")
	tree.SetName("登录")
	tree.SetComment("登录到服务器，如果是首次登录，创建角色")
	seq := NewSequenceNode()
	tree.Root().SetChild(seq)

	loginSelector := NewSelectorNode()
	loginCondNode := NewConditionNode()
	loginNode := NewBevNode(&BevLogin{ReloginDelay: 1000})
	loginNode.SetComment("登录到服务器")
	loginCondNode.SetNoChildSuccess(true)
	loginCondNode.SetComment("判定是否需要登录")
	loginCondNode.AddCond(NewCondStatus(false, NewStatusBool(types.Status_IsLogin, true)))
	loginSelector.AddChild(loginCondNode)
	loginSelector.AddChild(NewBevNode(&BevLogin{ReloginDelay: 1000}))

	createRoleSelector := NewSelectorNode()
	createRoleCondNode := NewConditionNode()
	createRoleNode := NewBevNode(&BevCreateRole{})
	createRoleNode.SetComment("创建角色")
	createRoleCondNode.SetNoChildSuccess(true)
	createRoleCondNode.SetComment("判定是否需要创建角色")
	createRoleCondNode.AddCond(NewCondStatus(false, NewStatusBool(types.Status_IsFirstlogin, true)))
	createRoleSelector.AddChild(createRoleCondNode)
	createRoleSelector.AddChild(createRoleNode)

	seq.AddChild(loginSelector)
	seq.AddChild(createRoleSelector)

	data, err := MarshalXMLTree(framework, tree)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))

	newTree := NewTree("test")
	if err := UnmarshalXMLTree(framework, data, newTree); err != nil {
		t.Fatal(err)
	}
}
