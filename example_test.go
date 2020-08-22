/**
 * @author ForeverZi
 * @email txzm2018@gmail.com
 * @create date 2020-08-22 11:49:14
 * @modify date 2020-08-22 11:49:14
 * @desc [description]
 */
package confparser_test

import (
	"log"

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
	parser := confparser.NewJSONParser("./conf")
	parser.RegisterConfMap("role.json", Role{}, func(tableName string, content map[string]interface{}) {
		log.Printf("配置表[%v]更新了:[%#v]\n", tableName, content)
	})
	watchEndChan, err := parser.Watch()
	if err != nil {
		log.Fatal("无法开启监听", err)
	}
	<-watchEndChan
	log.Println("监听跳出")
}
