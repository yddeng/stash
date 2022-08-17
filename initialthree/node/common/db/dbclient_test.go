package db

import (
	"fmt"
	"log"
	"testing"
	"time"
)

var cli *Client

/*
DROP TABLE IF EXISTS "test";
CREATE TABLE "test" (
"__key__" varchar(255) NOT NULL,
"__version__" int8 NOT NULL DEFAULT 0,
"name" varchar NOT NULL DEFAULT '',
"age" int8 NOT NULL DEFAULT 0,
PRIMARY KEY ("__key__")
);
*/

func init() {
	c, err := NewClient("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
	if err != nil {
		panic(err)
		return
	}
	cli = c
}

func TestClient_Set(t *testing.T) {
	err := cli.Set("test", "333", map[string]interface{}{
		"name": "name",
		"age":  2,
	})

	fmt.Println(err)
}

func TestClient_Update(t *testing.T) {
	err := cli.Update("test", "123", map[string]interface{}{
		"age": 2,
	})

	fmt.Println(err)
}

func TestClient_Upsert(t *testing.T) {
	//err := cli.SetNx("test", "234", map[string]interface{}{
	//	"name": 234,
	//})
	//fmt.Println(err)

	//times := 0
	sTime := time.Now()

	for i := 1; i <= 1000000; i++ {
		id := fmt.Sprintf("%d", i)
		err := cli.Upsert("test", id, map[string]interface{}{
			"name": id, "age": i,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
	//useT := time.Now().Sub(sTime).Nanoseconds()

	fmt.Println(time.Now().Sub(sTime).String())

}

func TestClient_Get(t *testing.T) {
	ret, err := cli.Get("test", "111", "name")
	if err != nil {
		log.Printf("%s \n", err)
		return
	}
	fmt.Println(ret, err)

}

func TestClient_GetBatch(t *testing.T) {
	err := cli.GetBatch("test", []string{"111", "222", "444"}, func(ret map[string]interface{}) error {
		fmt.Println(ret)
		return nil
	}, "name")
	fmt.Println("end", err)

}

func TestClient_GetAll(t *testing.T) {

	err := cli.GetAll("test", func(i map[string]interface{}) error {
		fmt.Println(i)
		return nil
	})
	if err != nil {
		log.Printf("%s \n", err)
		return
	}
	fmt.Println("end")
}

func TestClient_DELETE(t *testing.T) {
	err := cli.Delete("test", "234")
	fmt.Println(err)
}

func TestNewClient(t *testing.T) {
	clis := make([]*Client, 0, 100)
	for i := 0; i < 100; i++ {
		cli, err := NewClient("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
		if err != nil {
			log.Printf("%s \n", err)
			fmt.Println(i)
			return
		}
		clis = append(clis, cli)
	}
}

func TestClient_UpdateBatch(t *testing.T) {
	err := cli.UpdateBatch("test", map[string]map[string]interface{}{
		"123": {"name": "123", "age": 1111},
		"345": {"name": "345", "age": 2222},
	})
	fmt.Println(err)

	ret, err := cli.Get("test", "123")
	if err != nil {
		log.Printf("%s \n", err)
		return
	}
	fmt.Println(ret)

	ret, err = cli.Get("test", "345")
	if err != nil {
		log.Printf("%s \n", err)
		return
	}
	fmt.Println(ret)
}

func TestClient_UpsertBatch(t *testing.T) {

	times := 0
	sTime := time.Now()
	batch := 1

	for i := 0; i < 1000000; {
		keyFields := make(map[string]map[string]interface{}, 500)
		for k := 0; k < batch; k++ {
			id := fmt.Sprintf("%d", i+k)
			keyFields[id] = map[string]interface{}{
				"name": id, "age": i + k,
			}
		}
		err := cli.UpsertBatch("test", keyFields)
		if err != nil {
			fmt.Println(err)
		}
		i += batch
		times++
	}
	useT := time.Now().Sub(sTime).Nanoseconds()

	fmt.Println(times, time.Now().Sub(sTime).String(), int(useT)/times)

	//err := cli.SetNxBatch("test", map[string]map[string]interface{}{
	//	"123": {"name": "1123", "age": 1123},
	//	"789": {"name": "789", "age": 789},
	//})
	//fmt.Println(err)
	//
	//ret, err := cli.Get("test", "123")
	//if err != nil {
	//	log.Printf("%s \n", err)
	//	return
	//}
	//fmt.Println(ret)
	//
	//ret, err = cli.Get("test", "789")
	//if err != nil {
	//	log.Printf("%s \n", err)
	//	return
	//}
	//fmt.Println(ret)
}

func Test_T(t *testing.T) {
	ch := make(chan int, 100)

	c, err := NewClient("pgsql", "127.0.0.1", 5432, "yidongdeng", "dbuser", "123456")
	if err != nil {
		panic(err)
		return
	}
	c.SetMaxOpenConns(100)

	for i := 100; i < 200; i++ {
		key := fmt.Sprintf("select * from test where __key__='%d'", i)
		go func(i int) {
			rows, err := c.db.Query(key)
			if err != nil {
				fmt.Println(err)
			}
			ch <- 1
			fmt.Println(i, rows.Next())
		}(i)
	}

	i := 0
	for {
		select {
		case <-ch:
			i++
			fmt.Println("---", i, len(ch))
		}
	}
}
