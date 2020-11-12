package heidou

type Loader interface {
	// 从项目数据库元信息表中获取表定义数据
	LoadMetaTable() ([]*MetaTable, error)
}
