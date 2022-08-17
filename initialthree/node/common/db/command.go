package db

import (
	"fmt"
	"strings"
)

/*
 * 插入数据
 * tableName:表名 fields:键值对
 */
func (c *Client) Set(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (__key__,__version__,%s)
VALUES (%s);`

	columns, values := []string{}, []string{"$1", "$2"}
	args := []interface{}{key, 1}
	var i = 3
	for k, v := range fields {
		columns = append(columns, k)
		values = append(values, fmt.Sprintf("$%d", i))
		i++
		args = append(args, v)
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(values, ","))
	//fmt.Println(sqlStatement)
	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err
}

/*
 * 更新数据
 * tableName:表名 key:key fields:键值对
 */
func (c *Client) Update(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
UPDATE "%s" 
SET %s
WHERE __key__ = '%s';`

	columns := []string{}
	args := []interface{}{}
	var i = 1
	for k, v := range fields {
		if k != "__key__" {
			columns = append(columns, fmt.Sprintf(`%s = $%d`, k, i))
			i++
			args = append(args, v)
		}
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), key)
	//fmt.Println(sqlStatement)
	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec(args...)
	return err

}

/*
 * 批量更新数据
 * tableName:表名 keyFields: key -> 键值对
 * 需控制sql语句字符串长度
 */

func (c *Client) UpdateBatch(tableName string, keyFields map[string]map[string]interface{}) error {
	sqlStr := `
UPDATE "%s"
SET %s
WHERE __key__ IN (%s)`

	keys := make([]string, 0, len(keyFields))
	filedStr := map[string]string{}
	for key, fields := range keyFields {
		keys = append(keys, fmt.Sprintf(`'%s'`, key))
		for k, v := range fields {
			s, ok := filedStr[k]
			if !ok {
				s = fmt.Sprintf(`%s = CASE __key__
`, k)
				filedStr[k] = s
			}

			s += fmt.Sprintf(`WHEN '%s' THEN %s 
`, key, ConvertValue(v))
			filedStr[k] = s

		}
	}

	caseStr := make([]string, 0, len(filedStr))
	for _, s := range filedStr {
		caseStr = append(caseStr, s+"END")
	}

	setStr := strings.Join(caseStr, ",")

	sqlStatement := fmt.Sprintf(sqlStr, tableName, setStr, strings.Join(keys, ","))
	//fmt.Println(sqlStatement)
	_, err := c.db.Exec(sqlStatement)
	return err

}

/*
 * upsert
 * 没有数据插入，有则更改。
 * tableName:表名 key:主键 fields:键值对
 */
func (c *Client) Upsert(tableName, key string, fields map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (__key__,__version__,%s)
VALUES(%s) 
ON conflict(%s) DO 
UPDATE SET %s;`

	columns, values, sets := []string{}, []string{"$1", "$2"}, []string{}
	args := []interface{}{key, 1}
	var i = 3
	for k, v := range fields {
		if !("__key__" == k || k == "__version__") {
			sets = append(sets, fmt.Sprintf(`%s = $%d`, k, i))
			columns = append(columns, k)
			values = append(values, fmt.Sprintf("$%d", i))
			i++
			args = append(args, v)
		}
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(values, ","), "__key__", strings.Join(sets, ","))
	//fmt.Println(sqlStatement, args)
	_, err := c.db.Exec(sqlStatement, args...)
	return err
}

/*
 * upsert批量操作
 * 没有数据插入，有则更改。
 * tableName:表名  keyFields: key -> 键值对
 * 需控制sql语句字符串长度
 */
func (c *Client) UpsertBatch(tableName string, keyFields map[string]map[string]interface{}) error {
	sqlStr := `
INSERT INTO "%s" (%s)
VALUES %s 
ON conflict(%s) DO 
UPDATE SET %s;`

	columns := []string{"__key__", "__version__"}
	for _, fields := range keyFields {
		for k := range fields {
			if !(k == "__key__" || k == "__version__") {
				columns = append(columns, k)
			}
		}
		break
	}

	keys := make([]string, 0, len(keyFields))
	values := make([]string, 0, len(keyFields))
	valStr := make([]string, 0, len(keyFields))

	for key, fields := range keyFields {
		keys = append(keys, fmt.Sprintf(`%s`, key))
		values = values[0:0]
		for _, k := range columns {
			if k == "__key__" {
				values = append(values, key)
			} else if k == "__version__" {
				values = append(values, "1")
			} else {
				if v, ok := fields[k]; ok {
					values = append(values, ConvertValue(v))
				} else {
					return fmt.Errorf("%s not found filed %s", key, k)
				}
			}
		}
		s := "(" + strings.Join(values, ",") + ")"
		valStr = append(valStr, s)
	}

	setStr := make([]string, 0, len(columns))
	for _, k := range columns {
		if k != "__key__" {
			setStr = append(setStr, fmt.Sprintf("%s = excluded.%s", k, k))
		}
	}

	sqlStatement := fmt.Sprintf(sqlStr, tableName, strings.Join(columns, ","), strings.Join(valStr, ","), "__key__", strings.Join(setStr, ","))
	//fmt.Println(sqlStatement)
	_, err := c.db.Exec(sqlStatement)
	return err
}

/*
 * 读取数据。
 * tableName:表名 key:主键
 * ret 返回键值对
 */
func (c *Client) Get(tableName, key string, fields ...string) (ret map[string]interface{}, err error) {
	sqlStr := `
SELECT %s FROM "%s" 
WHERE __key__ in (%s);`

	var columns []string
	var values []interface{}
	var sqlStatement string
	if len(fields) == 0 {
		sqlStatement = fmt.Sprintf(sqlStr, "*", tableName, fmt.Sprintf("'%s'", key))
	} else {
		columns = make([]string, 0, len(fields)+2)
		values = make([]interface{}, 0, len(fields)+2)
		for _, field := range fields {
			columns = append(columns, field)
			values = append(values, new(interface{}))
		}
		columns = append(columns, "__key__", "__version__")
		values = append(values, new(interface{}), new(interface{}))
		sqlStatement = fmt.Sprintf(sqlStr, strings.Join(columns, ","), tableName, fmt.Sprintf("'%s'", key))
	}

	//fmt.Println(sqlStatement)
	rows, err := c.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ret = map[string]interface{}{}

	for rows.Next() {
		if len(columns) == 0 {
			columns, err = rows.Columns()
			if err != nil {
				return nil, err
			}

			columnsLen := len(columns)
			values = make([]interface{}, 0, columnsLen)
			for i := 0; i < columnsLen; i++ {
				values = append(values, new(interface{}))
			}
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}
		for i, k := range columns {
			ret[k] = *(values[i].(*interface{}))
		}
		break
	}

	return ret, nil
}

/*
 * batch读取数据。
 * tableName:表名 key:主键
 * ret 返回键值对
 */
func (c *Client) GetBatch(tableName string, keys []string, callback func(ret map[string]interface{}) error, fields ...string) error {
	sqlStr := `
SELECT %s FROM "%s" 
WHERE __key__ in (%s);`

	_keys := make([]string, 0, len(keys))
	for _, key := range keys {
		_keys = append(_keys, fmt.Sprintf("'%s'", key))
	}

	var columns []string
	var values []interface{}
	var sqlStatement string
	if len(fields) == 0 {
		sqlStatement = fmt.Sprintf(sqlStr, "*", tableName, strings.Join(_keys, ","))
	} else {
		columns = make([]string, 0, len(fields)+2)
		values = make([]interface{}, 0, len(fields)+2)
		for _, field := range fields {
			columns = append(columns, field)
			values = append(values, new(interface{}))
		}
		columns = append(columns, "__key__", "__version__")
		values = append(values, new(interface{}), new(interface{}))
		sqlStatement = fmt.Sprintf(sqlStr, strings.Join(columns, ","), tableName, strings.Join(_keys, ","))
	}

	//fmt.Println(sqlStatement)
	rows, err := c.db.Query(sqlStatement)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if len(columns) == 0 {
			columns, err = rows.Columns()
			if err != nil {
				return err
			}

			columnsLen := len(columns)
			values = make([]interface{}, 0, columnsLen)
			for i := 0; i < columnsLen; i++ {
				values = append(values, new(interface{}))
			}
		}
		//fmt.Println(columns, values)
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		mid := map[string]interface{}{}
		for i, k := range columns {
			mid[k] = *(values[i].(*interface{}))
		}
		if e := callback(mid); e != nil {
			return e
		}
	}
	return nil

}

/*
 * 读取所有数据。
 * tableName:表名
 * ret 返回键值对的slice
 * limit 限制每次读取多少行数据.
 */
func (c *Client) GetAll2(tableName string, limit int, callback func([]map[string]interface{}) error) error {
	start := 0
	total, err := c.Count(tableName)
	if err != nil {
		return err
	}

	ret := make([]map[string]interface{}, 0, limit)
	for start < total {
		sqlStatement := fmt.Sprintf(`SELECT * FROM "%s" ORDER BY __key__ LIMIT %d OFFSET %d ;`, tableName, limit, start)
		rows, err := c.db.Query(sqlStatement)
		if err != nil {
			return err
		}

		var columns []string
		var values []interface{}

		for rows.Next() {
			start++
			if len(columns) == 0 || len(values) != len(columns) {
				columns, err = rows.Columns()
				if err != nil {
					return err
				}

				columnsLen := len(columns)
				values = make([]interface{}, 0, columnsLen)
				for i := 0; i < columnsLen; i++ {
					values = append(values, new(interface{}))
				}
			}
			//fmt.Println(columns, values)
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			mid := map[string]interface{}{}
			for i, k := range columns {
				mid[k] = *(values[i].(*interface{}))
			}
			ret = append(ret, mid)
		}

		e := callback(ret)
		ret = ret[0:0]
		_ = rows.Close()

		if e != nil {
			return e
		}
	}

	return nil
}

/*
 * 读取所有数据。
 * tableName:表名
 * ret 返回键值对的
 */
func (c *Client) GetAll(tableName string, callback func(map[string]interface{}) error) error {
	sqlStatement := fmt.Sprintf(`SELECT * FROM "%s" `, tableName)
	rows, err := c.db.Query(sqlStatement)
	if err != nil {
		return err
	}
	defer rows.Close()

	var columns []string
	var values []interface{}

	for rows.Next() {
		if len(columns) == 0 || len(values) != len(columns) {
			columns, err = rows.Columns()
			if err != nil {
				return err
			}

			columnsLen := len(columns)
			values = make([]interface{}, 0, columnsLen)
			for i := 0; i < columnsLen; i++ {
				values = append(values, new(interface{}))
			}
		}
		//fmt.Println(columns, values)
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		mid := map[string]interface{}{}
		for i, k := range columns {
			mid[k] = *(values[i].(*interface{}))
		}
		e := callback(mid)

		if e != nil {
			return e
		}
	}
	return nil
}

func (c *Client) Delete(tableName, key string) error {
	sqlStr := `
DELETE FROM "%s" 
WHERE __key__ = '%s';`

	sqlStatement := fmt.Sprintf(sqlStr, tableName, key)
	//fmt.Println(sqlStatement)

	smt, err := c.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 清空表
func (this *Client) Truncate(tableName string) error {
	sqlStr := `
TRUNCATE TABLE %s;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	//fmt.Println(sqlStatement)
	smt, err := this.db.Prepare(sqlStatement)
	if err != nil {
		return err
	}
	_, err = smt.Exec()
	return err
}

// 表行数
func (this *Client) Count(tableName string) (int, error) {
	sqlStr := `
select count(*) from %s;`

	sqlStatement := fmt.Sprintf(sqlStr, tableName)
	smt, err := this.db.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}
	row := smt.QueryRow()
	var count int
	err = row.Scan(&count)
	return count, err
}

// 执行 sql
func (this *Client) Exec(query string, args ...interface{}) error {
	_, err := this.db.Exec(query, args...)
	return err
}
