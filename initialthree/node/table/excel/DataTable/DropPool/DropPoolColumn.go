package DropPool

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"initialthree/node/table/excel"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

type ColumnTable struct {
	loadPath string
	indexID  atomic.Value // map[int32]*DropPool.DropPool
}

func getColumnIDMap() map[int32]*DropPool {
	return columnTable.indexID.Load().(map[int32]*DropPool)
}

var columnTable ColumnTable

func init() {
	columnTable.loadPath = "/DataTable/掉落/DropPool纵表"
	excel.AddTable(&columnTable)
}

func (this *ColumnTable) Load(readPath string) {
	tmpIDMap := map[int32]*DropPool{}

	loadPath := path.Join(readPath, "/DataTable/掉落/DropPool纵表")
	if err := filepath.Walk(loadPath, func(path string, f os.FileInfo, err error) error {
		if f != nil && !f.IsDir() {
			n := strings.Split(f.Name(), ".")
			if len(n) != 2 {
				return nil
			}
			nn := strings.Split(n[0], "_")
			if len(nn) != 2 {
				return nil
			}

			id, err := strconv.Atoi(nn[1])
			if err != nil {
				return err
			}
			xlFile, err := excelize.OpenFile(path)
			if err != nil {
				return err
			}

			tableRows := xlFile.GetRows(xlFile.GetSheetName(xlFile.GetActiveSheetIndex()))

			constRow := tableRows[4]
			pool := &DropPool{
				ID:         int32(id),
				Type:       constRow[1],
				MinCount:   excel.ReadInt32(constRow[2]),
				MaxCount:   excel.ReadInt32(constRow[3]),
				Repeatable: excel.ReadBool(constRow[4]),
				TypeEnum:   excel.ReadEnum(constRow[1]),
			}

			rows := tableRows[9:]
			pool.DropList = make([]*DropList_, 0, len(rows))
			for _, row := range rows {
				if row[0] != "" {
					p := &DropList_{
						Type:   excel.ReadEnum(row[0]),
						ID:     excel.ReadInt32(row[1]),
						Count:  excel.ReadInt32(row[2]),
						Wave:   excel.ReadInt32(row[3]),
						Weight: excel.ReadInt32(row[4]),
					}
					pool.DropList = append(pool.DropList, p)
				}
			}

			tmpIDMap[pool.ID] = pool
			//fmt.Println(pool.ID, pool)
		}
		return nil
	}); err != nil {
		panic(err)
	}

	this.indexID.Store(tmpIDMap)
}
