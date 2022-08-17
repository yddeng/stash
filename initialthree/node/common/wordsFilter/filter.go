package wordsFilter

import (
	"io/ioutil"
	"strings"
	"sync"
	"unicode/utf8"
)

var filter_once sync.Once

var wordsFilter *Filter
var rLock = sync.RWMutex{}

type Filter struct {
	Root   *FilterNode
	Ignore string
}

type FilterNode struct {
	Children map[rune]*FilterNode
	End      bool
}

func NewFilter(ignore, words, sep string) *Filter {
	filter := &Filter{
		Root:   newFilterNode(),
		Ignore: ignore,
	}

	filter.NewFilter(strings.Split(words, sep))

	return filter
}

func newFilterNode() *FilterNode {
	n := &FilterNode{}
	n.Children = map[rune]*FilterNode{}
	return n
}

func (this *Filter) NewFilter(s []string) {
	for _, _s := range s {
		if _s != "" {
			this.Inster(_s)
		}
	}
}

func (this *Filter) Inster(txt string) {
	if len(txt) < 1 {
		return
	}
	node := this.Root
	key := []rune(txt)
	for i := 0; i < len(key); i++ {
		if _, b := node.Children[key[i]]; !b {
			node.Children[key[i]] = newFilterNode()
		}
		node = node.Children[key[i]]
	}

	node.End = true
}

func eliminate(str string, e string) string {
	str_ := []byte(str)
	e_ := []byte(e)
	out := []byte{}

	f := func(c byte) bool {
		for _, v := range e_ {
			if c == v {
				return false
			}
		}
		return true
	}

	for _, v := range str_ {
		if f(v) {
			out = append(out, v)
		}
	}
	return string(out)
}

func (this *Filter) Check(txt string) bool {
	if len(txt) < 1 {
		return false
	}

	txt = eliminate(txt, this.Ignore)
	node := this.Root
	key := []rune(txt)

	for i := 0; i < len(key); i++ {
		for j := i; j < len(key); j++ {
			if node, b := node.Children[key[j]]; !b {
				break
			} else {
				if node.End == true {
					return true
				}
			}
			node = node.Children[key[j]]
		}
		node = this.Root
	}
	return false
}

func Check(txt string) bool {
	rLock.RLock()
	defer rLock.RUnlock()

	if wordsFilter == nil {
		return false
	}
	return wordsFilter.Check(txt)
}

func (this *Filter) Replace(txt string) string {
	if len(txt) < 1 {
		return txt
	}

	node := this.Root
	key := []rune(txt)
	var chars []rune = key
	c, _ := utf8.DecodeRuneInString("*")

	for i := 0; i < len(key); i++ {
		if node1_, b := node.Children[key[i]]; b {
			if node1_.End {
				chars[i] = c
				continue
			}
			node = node1_
			for j := i + 1; j < len(key); j++ {
				key_ := eliminate(string(key[j]), this.Ignore)
				if key_ == "" {
					continue
				} else {
					if node2_, bb := node.Children[[]rune(key_)[0]]; bb {
						if node2_.End {
							for t := i; t <= j; t++ {
								chars[t] = c
							}
							i = j
							//break
						}
						node = node2_
					} else {
						break
					}
				}
			}
			node = this.Root
		}
	}

	return string(chars)
}

func Replace(txt string) string {
	rLock.RLock()
	defer rLock.RUnlock()
	if wordsFilter == nil {
		return txt
	}
	return wordsFilter.Replace(txt)
}

//func init() {
//	ignore := "!@#$%^&*()_+/*-="
//	words := "习近平;习主席;毛泽东;毛主席;周恩来;刘少奇;江泽民;温家宝;赵云;诸葛亮;周瑜;狗屎;臭狗屎;fuck;fuck you;"
//	wordsFilter = NewFilter(ignore, words, ";")
//}

//var ignore = "!@#$%^&*()_+/*-="
var ignore = ""

func readFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	// 忽略掉 windows \r
	dataStr := eliminate(string(data), "\r")
	return dataStr, nil
}

func Init(path string) error {
	dataStr, err := readFile(path)
	if err != nil {
		return err
	}

	rLock.Lock()
	wordsFilter = NewFilter(ignore, dataStr, "\n")
	rLock.Unlock()
	return nil
}

func ReLoad(path string) error {
	dataStr, err := readFile(path)
	if err != nil {
		return err
	}

	tmp := NewFilter(ignore, dataStr, "\n")

	rLock.Lock()
	wordsFilter = tmp
	rLock.Unlock()
	return nil

}
