package confparser

import (
	"io"

	"github.com/ForeverZi/confparser/internal"
)

var (
	//ErrUnregisteredTable 表示该配置表对应的配置类型未注册
	ErrUnregisteredTable = internal.ErrUnregisteredTable
	//ErrWatcherAbort  监听协程跳出
	ErrWatcherAbort = internal.ErrWatcherAbort
)

//OnTableChanged 配表变更响应(包括初始化)
type OnTableChanged = internal.OnTableChanged

//Parser 注册解析器接口
//不要修改获取的数据（获取到的列表可以排序，内部指针不可以修改）
type Parser interface {
	// 注册对应的表的解析结构体和变更时的回调
	RegisterConfMap(tableName string, confType interface{}, onTableChanged OnTableChanged) (err error)
	// 开启监听
	Watch() (endchan <-chan error, err error)
	// 查看配置是否存在
	Exist(tableName, id string) bool
	// 以列表的形式获取所有的配置指针
	GetAllItems(tableName string) []interface{}
	// 获取
	GetTable(tableName string) map[string]interface{}
	GetRecord(tableName, id string) interface{}
}

//NewJSONParser 新建游戏配置对象
func NewJSONParser(confDir string) Parser {
	return internal.NewJSONParser(confDir)
}

//SetLoggerOutput 设置日志输出
func SetLoggerOutput(out io.Writer) {
	internal.SetLoggerOutput(out)
}
