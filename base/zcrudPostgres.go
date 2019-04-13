package base

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// postgres
// 根据 tableName 生成 CRUD 的 sql

// ColumnInfo 字段属性
type ColumnInfo struct {
	Name       string `db:"column_name" json:"name" form:"name"`
	DataType   string `db:"data_type" json:"dataType" form:"dataType"`
	IsNullable string `db:"is_nullable" json:"isNullable" form:"isNullable"`
	// ColumnKey  string `db:"COLUMN_KEY" json:"columnKey" form:"columnKey"`
}

// ColumnInfoByTable 根据 tablename 获取 column
func columnInfoByTable(tableName string) ([]ColumnInfo, error) {
	var err error
	sugar := Sugar()
	db := DB()
	kv := KvCache()

	sqlCmd := kv.GetCommand("meta.postgresColumns", nil)
	sugar.Info(sqlCmd)

	var items []ColumnInfo
	err = db.Select(&items, sqlCmd, "pony_development", tableName)

	return items, err
}

// MapToName 只取得 Name
func mapToName(cols []ColumnInfo) []string {
	result := []string{}
	for _, v := range cols {
		result = append(result, v.Name)
	}
	return result
}

func printSubCode(table string) error {
	items, err := genStubCode(table)
	if err != nil {
		return err
	}
	for _, v := range items {
		fmt.Println(v)
	}
	return nil
}

func genStubCode(table string) ([]string, error) {
	sugar = Sugar()

	// 取得所有字段信息
	cols, err := columnInfoByTable(table)
	if err != nil {
		sugar.Fatalw("no column info", "table", table, "err", err)
		return nil, err
	}

	// struct
	str := structField(cols, table)

	// 取得 select
	sel := selectField(cols, table)

	// 取得 insert
	ins := insertField(cols, table)

	// 取得 update
	upd := updateField(cols, table)

	// 最终结果
	sep := "########"
	result := []string{}

	result = append(result, sep)
	result = append(result, str...)

	result = append(result, sep)
	result = append(result, sel...)

	result = append(result, sep)
	result = append(result, ins...)

	result = append(result, sep)
	result = append(result, upd...)

	result = append(result, sep)

	// 只支持单一主键
	pk := "id"
	// for _, c := range cols {
	// 	if c.ColumnKey == "PRI" {
	// 		pk = c.Name
	// 		break
	// 	}
	// }
	result = append(result, fmt.Sprintf("DELETE FROM %s WHERE %s = :%s; ", table, pk, pk))

	// excel 中各个字段的映射
	result = append(result, sep)
	excelCols := mapExcelColumns(cols)
	result = append(result, excelCols...)

	return result, nil
}

// structField 生成 struct
func structField(cols []ColumnInfo, table string) []string {
	var result []string
	result = append(result, fmt.Sprintf("type %s struct {", ucFirst(table)))

	// 数据库字段类型 和 golang 类型对照
	types := map[string]string{
		"uuid":                        "uuid.UUID",
		"character varying":           "string",
		"timestamp without time zone": "time.Time",

		"int":    "int",
		"bigint": "int",
	}

	for _, c := range cols {
		t, ok := types[c.DataType]
		if !ok {
			sugar.Fatal("missing type: %s", c.DataType)
		}

		v := c.Name
		n := fmt.Sprintf("%s %s `db:\"%s\" json:\"%s\" form:\"%s\"`", ucFirst(v), t, v, lcFirst(v), lcFirst(v))
		result = append(result, n)
	}
	result = append(result, "}")
	return result
}

// selectField 形成 select 的字段
func selectField(cols []ColumnInfo, table string) []string {
	var result []string
	result = append(result, "SELECT ")

	last := len(cols) - 1
	for i, c := range cols {
		v := c.Name
		var n string

		// 默认不是最后一行，有逗号
		lc := ", "
		if i == last {
			// 最后一列，没有逗号
			lc = ""
		}

		if c.IsNullable == "NO" {
			// 不可能为空，不需要 IFNULL 判定
			n = fmt.Sprintf("%s%s", v, lc)
		} else {
			//  可能为空
			if c.DataType == "varchar" {
				n = fmt.Sprintf("coalesce(%s, '') %s%s", v, v, lc)
			} else {
				n = fmt.Sprintf("coalesce(%s, 0) %s%s", v, v, lc)
			}
		}

		result = append(result, n)
	}
	result = append(result, fmt.Sprintf("FROM %s;  ", table))
	return result
}

// insertField 生成 insert 语句
func insertField(cols []ColumnInfo, table string) []string {
	var result []string
	result = append(result, fmt.Sprintf("INSERT INTO %s (", table))

	// :ID
	var vals []string
	last := len(cols) - 1
	for i, c := range cols {
		// if c.ColumnKey == "PRI" {
		// 	// 主键不能插入
		// 	continue
		// }

		v := c.Name
		var n string
		var n2 string
		if i == last {
			// 最后一个字段没有 逗号
			n = fmt.Sprintf(":%s ", v)
			n2 = v
		} else {
			n = fmt.Sprintf(":%s, ", v)
			n2 = fmt.Sprintf("%s, ", v)
		}
		vals = append(vals, n)
		result = append(result, n2)
	}

	result = append(result, ") VALUES ( ")

	for _, v := range vals {
		result = append(result, v)
	}

	result = append(result, ");  ")

	return result
}

// updateField 生成 update 语句
func updateField(cols []ColumnInfo, table string) []string {
	var result []string
	result = append(result, fmt.Sprintf("UPDATE %s SET ", table))

	// 只支持单一主键
	pk := "id"
	last := len(cols) - 1
	for i, c := range cols {
		// if c.ColumnKey == "PRI" {
		// 	pk = c.Name
		// 	// 主键不能更新
		// 	continue
		// }

		v := c.Name
		var n string
		if i == last {
			// 最后一个字段没有 逗号
			n = fmt.Sprintf("%s = :%s ", v, v)
		} else {
			n = fmt.Sprintf("%s = :%s, ", v, v)
		}
		result = append(result, n)
	}

	// 把主键拼接在 where 中，只支持单一主键
	result = append(result, fmt.Sprintf("WHERE %s = :%s;  ", pk, pk))

	return result
}

// ucFirst 首字符大写
func ucFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

// lcFirst 首字符小写
func lcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// lpad 在左边补 0
func lpad(v string, n int) string {
	a := strings.Repeat("0", n)
	b := fmt.Sprintf("%s%s", a, strings.TrimSpace(v))
	// 取最后的 n 位
	return b[len(b)-n:]
}

// mapExcelColumn 上载 excel 时，对应的字段
func mapExcelColumns(cols []ColumnInfo) []string {
	var result []string
	for idx, c := range cols {
		// 第一列都是 ID，跳过，Excel 中第一列是表中的第二列
		if idx == 0 {
			continue
		}
		info := fmt.Sprintf(" %s: rowData(%d), ", ucFirst(c.Name), idx-1)
		result = append(result, info)
	}

	return result
}

// ddlCreate 生成 create table
func ddlCreate(tbl string, comment string, cnt int) {
	var result []string

	result = append(result, fmt.Sprintf(" DROP TABLE %s; ", tbl))

	result = append(result, fmt.Sprintf("CREATE TABLE `%s` (", tbl))
	result = append(result, fmt.Sprintf("  `ID` bigint(20) NOT NULL AUTO_INCREMENT, "))

	for i := 0; i < cnt; i++ {
		info := fmt.Sprintf("`col%s` varchar(200) DEFAULT NULL, ", lpad(strconv.Itoa(i+1), 2))
		result = append(result, info)
	}

	result = append(result, " PRIMARY KEY (`ID`) ")
	result = append(result, fmt.Sprintf(" ) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8 COMMENT='%s'; ", comment))

	result = append(result, fmt.Sprintf(" SELECT * FROM %s; ", tbl))

	fmt.Println(strings.Join(result, "\n"))
}

func fire() {
	sugar := Sugar()

	sugar.Info("fire sugar")
}
