package actions

import (
	"bytes"
	"pony/base"
	"pony/models"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// itemTypeA each column for Type A form
type itemTypeA struct {
	Code         string
	Name         string
	QuestionType string
	Unit         string
	Remark       string
	Value        string
	ItemType     string
}

// MatrixOpenTypeA for type A form
func MatrixOpenTypeA(c buffalo.Context) error {
	// read from querystring
	var err error
	sugar := base.Sugar()

	info := struct {
		MatrixNum string `form:"num" json:"num" db:"matrix"`
		ID        string `form:"id" json:"id"`
		Period    string `form:"period" json:"period" db:"period"`
		Company   string `form:"company" json:"company" db:"company"`
		Version   string `form:"version" json:"version" db:"version"`
	}{}
	err = c.Bind(&info)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(info.Period) == 0 {
		info.Period = "2018"
	}
	if len(info.Company) == 0 {
		info.Company = "PD"
	}
	if len(info.Version) == 0 {
		info.Version = "ACTUAL"
	}

	sugar.Debugw("parse parm", "param", info)

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	kv := base.KvCache()
	db := tx.TX

	// 查询之前保存的结果
	results := models.Matrices{}
	cmdSQL := kv.GetCommand("matrix.findBy", nil)
	nstmt, err := db.PrepareNamed(cmdSQL)
	if err != nil {
		sugar.Errorw("preparedNamed failed", "err", err)
		return errors.WithStack(err)
	}

	err = nstmt.Select(&results, info)
	if err != nil {
		sugar.Errorw("failed to query index", "param", info, "err", err)
		return errors.WithStack(err)
	}

	// 加载表单
	title, items, err := loadTypeA(info.MatrixNum)
	if err != nil {
		return errors.WithStack(err)
	}

	// 与已保存的值做匹配
	for idx, entry := range items {
		if entry.ItemType == "L1" {
			continue
		}

		for _, r := range results {
			if entry.Code == r.Code {
				items[idx].Value = r.Value
				break
			}
		}
	}

	p := struct {
		Period         string
		Company        string
		Version        string
		MatrixNum      string
		ID             string
		Title          string
		Items          []itemTypeA
		PeriodOptions  []string
		CompanyOptions []string
		VersionOptions []string
	}{
		Period:         info.Period,
		Company:        info.Company,
		Version:        info.Version,
		MatrixNum:      info.MatrixNum,
		ID:             info.ID,
		Title:          title,
		Items:          items,
		PeriodOptions:  []string{"2016", "2017", "2018"},
		CompanyOptions: []string{"IBM", "PD", "TMX"},
		VersionOptions: []string{"ACTUAL", "PLAN"},
	}
	c.Set("p", p)
	return c.Render(200, r.HTML("matrix/openTypeA", "surveys/simple.html"))
}

// loadTypeA read from excel
func loadTypeA(num string) (string, []itemTypeA, error) {
	var err error
	sugar := base.Sugar()

	// 从文件中读取
	fileName := "matrix/" + num + ".xlsx"
	source, err := base.Box.Find(fileName)
	if err != nil {
		return "", nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return "", nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := base.SheetReader(xlsx, sheetName)

	const MaxLoop = 500

	results := []itemTypeA{}
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := base.RowReader(readHelper, r)

		// 第一列：
		if len(rowData(0)) == 0 {
			break
		}

		entry := itemTypeA{
			Code:         rowData(0),
			Name:         rowData(1),
			QuestionType: rowData(2),
			Unit:         rowData(3),
			Remark:       rowData(4),
			ItemType:     "L2",
		}
		// 默认都是行项目，如果 第二列(name) 没有，那么表示是章节
		if len(entry.Name) == 0 {
			entry.ItemType = "L1"
		}

		results = append(results, entry)
	}

	// 固定位置：第二行，第一列
	title := base.RowReader(readHelper, 2)(0)
	sugar.Debugw("load formA", "file", fileName, "title", title, "items count", len(results))

	return title, results, nil
}

// MatrixSubmitTypeA for type A form
func MatrixSubmitTypeA(c buffalo.Context) error {
	// read from querystring
	var err error
	sugar := base.Sugar()

	info := struct {
		MatrixNum string `form:"num" json:"num"`
		ID        string `form:"id" json:"id"`
		Period    string `form:"period" json:"period"`
		Company   string `form:"company" json:"company"`
		Version   string `form:"version" json:"version"`
	}{}
	err = c.Bind(&info)
	if err != nil {
		return errors.WithStack(err)
	}
	if len(info.Period) == 0 {
		info.Period = "2018"
	}
	if len(info.Company) == 0 {
		info.Company = "PD"
	}
	if len(info.Version) == 0 {
		info.Version = "ACTUAL"
	}
	sugar.Debugw("parse parm", "param", info)

	title, items, err := loadTypeA(info.MatrixNum)
	if err != nil {
		return errors.WithStack(err)
	}

	// 获取 form data
	err = c.Request().ParseForm()
	if err != nil {
		return errors.WithStack(err)
	}
	formData := c.Request().Form

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	kv := base.KvCache()
	db := tx.TX

	// 按照 (问卷编号 + 提交人) 删除
	cmdSQL := kv.GetCommand("matrix.checkExist", nil)
	nstmt, err := db.PrepareNamed(cmdSQL)
	if err != nil {
		sugar.Errorw("preparedNamed failed", "err", err)
		return errors.WithStack(err)
	}

	updSQL := kv.GetCommand("matrix.updateValue", nil)
	now := time.Now()

	for _, entry := range items {
		if entry.ItemType == "L1" {
			continue
		}

		v := formData[entry.Code]
		record := models.Matrix{
			Company:    info.Company,
			Version:    info.Version,
			Period:     info.Period,
			Matrix:     info.MatrixNum,
			Code:       entry.Code,
			Value:      v[0],
			SubmitUser: info.ID,
		}

		var ids models.IDList
		// 按照关键字段查找
		err = nstmt.Select(&ids, record)
		if err != nil {
			sugar.Errorw("failed to query index", "entry", record, "err", err)
			return errors.WithStack(err)
		}

		if len(ids) == 0 {
			// 新增
			// sugar.Infow("create", "i", m)

			_, err := tx.ValidateAndCreate(&record)
			if err != nil {
				// return errors.WithStack(err)
				sugar.Errorw("failed to save index", "entry", record, "err", err)
				continue
			}
		} else {
			// 更新
			record.ID = ids[0].ID
			record.UpdatedAt = now
			// sugar.Infow("update", "i", m)

			_, err = db.NamedExec(updSQL, record)
			if err != nil {
				sugar.Errorw("update failed", "err", err, "entry", record)
			}
		}
	}

	p := struct {
		Title  string
		Period string
		From   string
	}{
		Title:  title,
		Period: info.Period,
		From:   c.Request().RequestURI,
	}

	c.Set("p", p)
	return c.Render(200, r.HTML("matrix/success", "surveys/simple.html"))
}
