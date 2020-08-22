/**
 * @author ForeverZi
 * @email txzm2018@gmail.com
 * @create date 2020-08-22 14:27:08
 * @modify date 2020-08-22 14:27:08
 * @desc [description]
 */
package confparser_test

import (
	"fmt"
	"testing"

	"github.com/ForeverZi/confparser"
)

func setupParser() confparser.Parser {
	parser := confparser.NewJSONParser("./conf")
	parser.RegisterConfMap("role.json", Role{}, func(tblname string, content map[string]interface{}) {})
	return parser
}

func BenchmarkExist(b *testing.B) {
	parser := setupParser()
	for i := 0; i < b.N; i++ {
		_ = parser.Exist("role.json", fmt.Sprint(i))
	}
}

func BenchmarkGetTable(b *testing.B) {
	parser := setupParser()
	var m map[string]interface{}
	for i := 0; i < b.N; i++ {
		m = parser.GetTable("role.json")
	}
	_ = m["1"]
}

func BenchmarkGetRecord(b *testing.B) {
	parser := setupParser()
	var role Role
	for i := 0; i < b.N; i++ {
		item := parser.GetRecord("role.json", "100")
		if item != nil {
			role = item.(Role)
		}
	}
	_ = role
}

func BenchmarkGetAllItems(b *testing.B) {
	parser := setupParser()
	for i := 0; i < b.N; i++ {
		_ = parser.GetAllItems("role.json")
	}
}
