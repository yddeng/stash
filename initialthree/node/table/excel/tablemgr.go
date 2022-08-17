package excel

import (
	path2 "path"
)

type Table interface {
	Load(path string)
}

type BeforeLoad interface {
	BeforeLoad()
}

type AfterLoad interface {
	AfterLoad()
}

type AfterLoadAll interface {
	AfterLoadAll()
}

var (
	tables []Table
)

func Reload(path string) {
	Load(path)
}

func Load(path string) {
	for _, v := range tables {
		if bl, ok := v.(BeforeLoad); ok {
			bl.BeforeLoad()
		}

		v.Load(path2.Join(path))

		if al, ok := v.(AfterLoad); ok {
			al.AfterLoad()
		}
	}

	for _, v := range tables {
		if al, ok := v.(AfterLoadAll); ok {
			al.AfterLoadAll()
		}
	}
}

func AddTable(t Table) {
	tables = append(tables, t)
}

func init() {
	tables = make([]Table, 0)
}
