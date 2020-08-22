/**
 * @author ForeverZi
 * @email txzm2018@gmail.com
 * @create date 2020-08-22 14:27:03
 * @modify date 2020-08-22 14:27:03
 * @desc [description]
 */
package confparser_test

import (
	"encoding/json"
	"io/ioutil"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/ForeverZi/confparser"
	"github.com/stretchr/testify/assert"
)

func TestParser_Exist(t *testing.T) {
	a := assert.New(t)
	parser := setupParser()
	a.True(parser.Exist("role.json", "1"), "1 在配置中")
	a.False(parser.Exist("role.json", "10000"), "10000 不在配置中")
}

func TestParser_GetAllItems(t *testing.T) {
	a := assert.New(t)
	parser := setupParser()
	items := parser.GetAllItems("role.json")
	a.Len(items, 1000, "长度错误,应该为1000个元素")
	sort.Slice(items, func(i, j int) bool {
		lid, _ := strconv.ParseInt(items[i].(*Role).ID, 10, 64)
		rid, _ := strconv.ParseInt(items[j].(*Role).ID, 10, 64)
		return lid < rid
	})
	for i, item := range items {
		id, _ := strconv.ParseInt(item.(*Role).ID, 10, 64)
		a.Equal(i, int(id), "元素不一致")
	}
}

func TestParser_GetRecord(t *testing.T) {
	a := assert.New(t)
	parser := setupParser()
	testCases := []struct {
		ID    string
		IsNil bool
	}{
		{"1", false},
		{"1000", true},
		{"999", false},
		{"0", false},
	}
	for _, testCase := range testCases {
		record := parser.GetRecord("role.json", testCase.ID)
		if testCase.IsNil {
			a.Nil(record, "记录应该为空")
		} else {
			a.Equal(testCase.ID, record.(Role).ID, "获取记录错误")
		}
	}
}

func TestParser_GetTable(t *testing.T) {
	a := assert.New(t)
	parser := setupParser()
	a.Nil(parser.GetTable("dada"))
	a.NotNil(parser.GetTable("role.json"))
	a.NotNil(parser.GetTable("role"))
}

func TestParser_Watch(t *testing.T) {
	a := assert.New(t)
	parser := confparser.NewJSONParser("./conf")
	times := 0
	nowText := time.Now().String()
	parser.RegisterConfMap("role.json", Role{}, func(tblname string, content map[string]interface{}) {
		times++
		a.Equal("role", tblname)
		if times > 1 {
			a.Equal(content["0"].(Role).Board, nowText)
		}
	})
	_, err := parser.Watch()
	a.Nil(err)
	content := parser.GetTable("role")
	item := content["0"].(Role)
	item.Board = nowText
	content["0"] = item
	data, _ := json.Marshal(content)
	err = ioutil.WriteFile("./conf/role.json", data, 0666)
	a.Nil(err, err)
}
