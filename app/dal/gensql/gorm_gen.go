package main

import (
	"flag"
	"fmt"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/rawsql"
)

// gorm.io/gen@v0.3.23/internal/model/base.go
var dataMap = map[string]func(gorm.ColumnType) (dataType string){
	"int":     func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	"json":    func(columnType gorm.ColumnType) string { return "json.RawMessage" },
	"tinyint": func(columnType gorm.ColumnType) string { return "int32" }, // Map bool to int
}

// 官方文档 gorm gen https://gorm.io/gen/create.html
// 文档 gorm gen工具教程，更安全更友好的ORM工具 https://blog.ysboke.cn/archives/772.html
// go run gensql.go # 生成所有db代码
func main() {
	flag.Parse()

	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithDefaultQuery, // generate mode

		FieldWithIndexTag: true,  // generate index tag for field
		FieldWithTypeTag:  true,  // generate type tag for field
		FieldNullable:     false, // generate nullable tag for field 为可能不存在的字段生成为指针
		FieldCoverable:    true,  // 如果有配置默认值则生成(没用上)
		FieldSignable:     false, // 如果数据里配置的是int(11) unsigned则生成uint64
	})

	// gormdb, err := gorm.Open(sqlite.Open("../../../forgeturl.sqlite"), &gorm.Config{})
	gormdb, err := gorm.Open(rawsql.New(rawsql.Config{
		FilePath: []string{
			"./page.sql",
			"./user.sql",
			"./user_page.sql",
		},
	}))
	if err != nil {
		panic(fmt.Errorf("open sql fail: %w", err))
	}
	g.UseDB(gormdb) // reuse your gorm db
	g.WithDataTypeMap(dataMap)
	g.ApplyBasic(g.GenerateAllTable()...)

	// 调用GenerateAllTable内部生成方法
	//tableList, err := gormdb.Migrator().GetTables()
	//fmt.Println("tableList:", tableList)
	//if err != nil {
	//	panic(fmt.Errorf("get all tables fail: %w", err))
	//}

	//tableModels := make([]interface{}, len(tableList))
	//for i, tableName := range tableList {
	//	if strings.HasPrefix(tableName, "sqlite_") {
	//		continue
	//	}
	//
	//	fmt.Println("Generating model for table:", tableName)
	//	tableModels[i] = g.GenerateModel(tableName)
	//}
	//g.ApplyBasic(tableModels...)

	// Generate the code
	g.Execute()
}
