# 使用指南

## 配置项说明

- 数据表配置项 : tables

| 名称    | 说明                   | 类型   | 默认值                                                | 例子               |
| ------- | ---------------------- | ------ | ----------------------------------------------------- | ------------------ |
| name    | 数据表名               | string | -                                                     | product            |
| isSkip  | 是否生成相应代码       | bool   | false                                                 | false              |
| extra   | 扩展配置               | map    | -                                                     |
| fields  | 字段数组               | array | -                                                     | -                  |
| errorCodes | 错误码 | array | - |
| methods | 支持的 Api 方法 | array  | ["list", "update", "delete", "bulkGet", "bulkDelete"] | ["list", "update"] |

- 字段配置项 : fields:

| 名称           | 说明                                                      | 类型   | 默认值 | 例子                      |
| -------------- | --------------------------------------------------------- | ------ | ------ | ------------------------- |
| name           | 字段名                                                    | string | -      | id                        |
| alias           | 别名                                                    | string | -      | nameAlias                        |
| isSkip     | 是否忽略此字段                                              | bool   | false  | true                      |
| isRequired     | 是否必填字段                                              | bool   | false  | true                      |
| isSortable   | 是否可按此字段排序                                            | bool   | false  | true                      |
| isFilterable   | 是否可按此字段过滤                                            | bool   | false  | true                      |
| operations   | 排序时的可用操作                                            | array   | -  | true                      |
| tags           | 扩展 struct tags                                          | string | ""     | binding:"required,max=64" |
| joinType       | 关联类型,取值 None, BelongTo, HasOne, HasMany, ManyToMany | enum   | None   | ManyToMany                |
| tableName      | 指定关联表表名                                            | string | ""     | category                  |
| joinTableName  | joinType 为ManyToMany时，指定连接表表名                   | string | ""     | product_category_relation |
| foreignKey     | 指定外键，可选，默认使用拥有者的类型名加上主字段名        | string | ""     | CategoryID                |
| references     | 指定引用                                                  | string | ""     | name                      |
| joinForeignKey | 指定连接表的外键                                          | string | ""     | productReferID            |
| joinReferences | 指定连接表的引用外键                                      | string | ""     | productRefer              |


关联关系相关字段 (foreignKey, references, joinForeignKey, joinReferences）的配置与 Gorm 保持一致，详见：[Gorm](https://gorm.io/zh_CN/docs) 的关联说明