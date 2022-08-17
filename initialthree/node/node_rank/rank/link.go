package rank

import "fmt"

type LinkNode struct {
	rankId int32
	next   *LinkNode
}

// offset 起始偏移值
// len 获取的长度， if len <= 0 , 则为偏移值之后所有节点。
// offset 超出范围，返回 nil
// len 超出范围，返回范围偏移值之后所有节点。
func (l *LinkNode) getNodes(offset, length int) []int32 {
	nodes := make([]int32, 0)
	ofc, count := 0, 0
	ll := l
	for ll != nil {
		if ofc == offset {
			nodes = append(nodes, ll.rankId)
			count++
			if count == length {
				return nodes
			}
		} else {
			ofc++
		}
		ll = ll.next
	}
	return nodes
}

func (l *LinkNode) show() {
	ll := l
	str := ""
	for ll != nil {
		str += fmt.Sprintf("-(%d)-", ll.rankId)
		ll = ll.next
	}
	fmt.Println(str)
}

/*
 * 构建链条。
 * 不能有环。不能有多个后续。可以多个接同一个后续
 * 返回多条线性链
 */
func makeLink(idToLastId map[int32]int32) map[int32]*LinkNode {
	nodes := map[int32]*LinkNode{}
	headNodes := map[int32]*LinkNode{} // 最后有多少个headNode 就有多少条 link

	for id, lastId := range idToLastId {
		var idN, lastIdN *LinkNode
		var ok bool
		if idN, ok = nodes[id]; !ok {
			idN = &LinkNode{rankId: id}
			nodes[id] = idN
			headNodes[id] = idN
		}
		if lastIdN, ok = nodes[lastId]; !ok {
			lastIdN = &LinkNode{rankId: lastId}
			nodes[lastId] = lastIdN
		}

		if idN.next != nil {
			panic(fmt.Sprintf("link id %d 已经存在后续", id))
		}
		idN.next = lastIdN
		// 如果在 headNodes， 移除
		if _, ok = headNodes[lastId]; ok {
			delete(headNodes, lastId)
		}
	}

	return headNodes
}

// 查找, 返回头 linkNode，当前 linkNode
func findLink(id int32, links map[int32]*LinkNode) (*LinkNode, *LinkNode) {
	for _, l := range links {
		ll := l
		for ll != nil {
			if ll.rankId == id {
				return l, ll
			}
			ll = ll.next
		}
	}

	return nil, nil
}
