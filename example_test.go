/**
 * @author ForeverZi
 * @email txzm2018@gmail.com
 * @create date 2020-08-22 11:49:14
 * @modify date 2020-08-22 11:49:14
 * @desc [description]
 */
package confparser_test

import (
	"io/ioutil"
	"log"
	"sort"
	"strconv"

	"github.com/ForeverZi/confparser"
)

type Role struct {
	ID     string
	Name   string
	Age    uint8
	Gender uint8
	Board  string
}

func Example() {
	// 设置日志输出
	confparser.SetLoggerOutput(ioutil.Discard)
	parser := confparser.NewJSONParser("./conf")
	parser.RegisterConfMap("role.json", Role{}, func(tableName string, content map[string]interface{}) {
		log.Printf("配置表[%v]更新了:[%#v]\n", tableName, content)
	})
	watchEndChan, err := parser.Watch()
	if err != nil {
		log.Fatal("无法开启监听", err)
	}
	// 获取配置记录
	item := parser.GetRecord("role.json", "1")
	log.Printf("角色[%v]留言[%v]\n", item.(Role).ID, item.(Role).Board)
	// 获取全表
	table := parser.GetTable("role.json")
	log.Printf("角色[%v]留言[%v]\n", table["1"].(Role).ID, table["1"].(Role).Board)
	// 以列表的形式获取配置
	// Notice:列表中的是指针
	items := parser.GetAllItems("role.json")
	sort.Slice(items, func(i, j int) bool {
		lid, _ := strconv.ParseInt(items[i].(*Role).ID, 10, 64)
		rid, _ := strconv.ParseInt(items[j].(*Role).ID, 10, 64)
		return lid < rid
	})
	log.Printf("第[%v]个角色的ID是[%v]\n", 2, items[1].(*Role).ID)
	<-watchEndChan
	log.Println("监听跳出")
}
