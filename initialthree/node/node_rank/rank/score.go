package rank

type Score int32

// 分数由大到小排序，相同分数按时间先后排序
func (this Score) Less(other interface{}) bool {
	return this >= other.(Score)
}
