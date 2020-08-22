# confparser简单的配置解析器
## 使用
[GoDoc](https://godoc.org/github.com/foreverzi/confparser)
- 每个方法的使用可以参照confparser_test.go
## 特性
- 基于反射的注册解析
- 热更新
- 目前仅支持map[string]注册类型格式的JSON配置
- 建议配合使用validator等库做单项校验
- 并发安全
