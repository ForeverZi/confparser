/**
 * @author ForeverZi
 * @email txzm2018@gmail.com
 * @create date 2020-08-22 11:49:22
 * @modify date 2020-08-22 11:49:22
 * @desc [description]
 */
package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
)

var (
	//ErrUnregisteredTable 表示该配置表对应的配置类型未注册
	ErrUnregisteredTable = fmt.Errorf("未注册的配置表")
	//ErrWatcherAbort  监听协程跳出
	ErrWatcherAbort = fmt.Errorf("监听协程跳出")
)

//OnTableChanged 配表变更响应(包括初始化)
type OnTableChanged = func(tableName string, content map[string]interface{})

//GameTables  存储结构
type GameTables = map[string]map[string]interface{}

//CacheTables 数组缓存结构
type CacheTables = map[string][]interface{}

//JSONParser 支持并发的游戏配置解析器，仅支持第一层级的.json文件
type JSONParser struct {
	tables       atomic.Value
	mutex        sync.Mutex
	watcher      *fsnotify.Watcher
	watchEndChan chan error
	confDir      string
	confType     map[string]reflect.Type
	itemsCache   atomic.Value
	changeFuncs  map[string]OnTableChanged
}

//NewJSONParser 创建方法
func NewJSONParser(confDir string) *JSONParser {
	parser := &JSONParser{
		confDir:      confDir,
		confType:     make(map[string]reflect.Type),
		changeFuncs:  make(map[string]OnTableChanged),
		watchEndChan: make(chan error, 1),
	}
	parser.tables.Store(make(GameTables))
	parser.itemsCache.Store(make(CacheTables))
	return parser
}

//RegisterConfMap 注册指定配置表名称的类型
func (parser *JSONParser) RegisterConfMap(filename string, confType interface{}, onTableChanged OnTableChanged) {
	tableName := parser.getTableName(filename)
	parser.confType[tableName] = reflect.TypeOf(confType)
	parser.changeFuncs[tableName] = onTableChanged
	parser.onConfFileModify(filename)
}

//Watch 开启监视协程
func (parser *JSONParser) Watch() (endchan <-chan error, err error) {
	parser.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	err = parser.watcher.Add(parser.confDir)
	if err != nil {
		return
	}
	go parser.watch()
	endchan = parser.watchEndChan
	return
}

//Exist 查看游戏配置项是否存在
func (parser *JSONParser) Exist(tableName, id string) bool {
	table := parser.GetTable(tableName)
	if table == nil {
		return false
	}
	_, ok := table[id]
	return ok
}

//GetAllItems 获取配置表所有的条目的复制
func (parser *JSONParser) GetAllItems(tableName string) []interface{} {
	tableName = parser.getTableName(tableName)
	parser.mutex.Lock()
	defer parser.mutex.Unlock()
	cache := parser.itemsCache.Load().(CacheTables)
	if items, ok := cache[tableName]; ok {
		res := make([]interface{}, len(items))
		copy(res, items)
		return res
	}
	content := parser.GetTable(tableName)
	items := make([]interface{}, 0, len(content))
	for _, item := range content {
		ptrV := reflect.New(parser.confType[tableName])
		ptrV.Elem().Set(reflect.ValueOf(item))
		items = append(items, ptrV.Interface())
	}
	cache[tableName] = items
	parser.itemsCache.Store(cache)
	res := make([]interface{}, len(items))
	copy(res, items)
	return res
}

//GetTable 获取游戏表, 首先需要判断返回的是否为nil
func (parser *JSONParser) GetTable(tableName string) map[string]interface{} {
	//读取最新的配置表
	tableName = parser.getTableName(tableName)
	tables := parser.tables.Load().(GameTables)
	return tables[tableName]
}

//GetRecord 获取游戏配置行，在使用的时候首先判断是否为nil，然后转换为对应的类型再使用
func (parser *JSONParser) GetRecord(tableName, id string) interface{} {
	table := parser.GetTable(tableName)
	if table == nil {
		return nil
	}
	return table[id]
}

func (parser *JSONParser) watch() {
	defer func() {
		logger.Println("结束监听游戏配置")
		close(parser.watchEndChan)
	}()
	registeredOp := []fsnotify.Op{fsnotify.Create, fsnotify.Write, fsnotify.Rename}
	for {
		select {
		case event, ok := <-parser.watcher.Events:
			logger.Printf("监听到游戏配置变更:%#v\n", event)
			if !ok {
				parser.watchEndChan <- ErrWatcherAbort
				return
			}
			for _, targetOp := range registeredOp {
				if event.Op&targetOp == targetOp {
					err := parser.onConfFileModify(event.Name)
					if err != nil && err != ErrUnregisteredTable {
						logger.Printf("更新游戏配置失败:[%#v]event:[%#v]\n", err, event)
					}
					break
				}
			}
		case err, ok := <-parser.watcher.Errors:
			if !ok {
				parser.watchEndChan <- ErrWatcherAbort
				return
			}
			logger.Printf("游戏配置监听发生错误:[%#v]\n", err)
		}
	}
}

func (parser *JSONParser) readTable(fileName string) (table map[string]interface{}, err error) {
	tableName := parser.getTableName(fileName)
	confType, ok := parser.confType[tableName]
	if !ok {
		logger.Printf("没有注册的配置表:[%v]\n", tableName)
		err = ErrUnregisteredTable
		return
	}
	file, err := os.Open(filepath.Join(parser.confDir, fileName))
	if err != nil {
		return
	}
	defer file.Close()
	mapType := reflect.MapOf(reflect.TypeOf(tableName), confType)
	mapPtr := reflect.New(mapType)
	mapPtr.Elem().Set(reflect.MakeMap(mapType))
	err = json.NewDecoder(file).Decode(mapPtr.Interface())
	if err != nil {
		return
	}
	table = make(map[string]interface{})
	iter := mapPtr.Elem().MapRange()
	for iter.Next() {
		k := iter.Key().String()
		v := iter.Value().Interface()
		table[k] = v
	}
	return
}

func (parser *JSONParser) updateTable(tableName string, content map[string]interface{}) {
	parser.mutex.Lock()
	tables := parser.tables.Load().(GameTables)
	tables[tableName] = content
	parser.tables.Store(tables)
	cache := parser.itemsCache.Load().(CacheTables)
	delete(cache, tableName)
	parser.itemsCache.Store(cache)
	parser.mutex.Unlock()
	// 防止注册方法耗时太久，提前释放锁
	parser.changeFuncs[tableName](tableName, content)
}

func (parser *JSONParser) getTableName(fileName string) string {
	fileName = filepath.Base(fileName)
	tableName := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	return strings.ToLower(tableName)
}

func (parser *JSONParser) onConfFileModify(fileName string) (err error) {
	tableName := parser.getTableName(fileName)
	logger.Printf("配置变更:[%v]\n", tableName)
	table, err := parser.readTable(filepath.Base(fileName))
	if err != nil {
		return
	}
	parser.updateTable(tableName, table)
	return
}
