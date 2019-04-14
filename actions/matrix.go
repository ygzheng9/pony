package actions

import (
	"fmt"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"pony/base"
	"pony/models"
)

// MatrixOpen open matrix sheet for input
func MatrixOpen(c buffalo.Context) error {
	// 从 querystring 中获取参数
	var err error
	info := struct {
		MatrixNum  string `form:"num" json:"num" db:"matrix"`
		SubmitUser string `form:"user" json:"user"`
		Period     string `form:"period" json:"period" db:"period"`
		Company    string `form:"company" json:"company" db:"company"`
		Version    string `form:"version" json:"version" db:"version"`
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

	// 根据 matrix，period，company 找到之前提交的数据
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	kv := base.KvCache()
	sugar := base.Sugar()
	db := tx.TX

	cmdSQL := kv.GetCommand("matrix.findBy", nil)
	nstmt, err := db.PrepareNamed(cmdSQL)
	if err != nil {
		sugar.Errorw("preparedNamed failed", "err", err)
		return errors.WithStack(err)
	}

	// 查询参数
	items := models.Matrices{}
	err = nstmt.Select(&items, info)
	if err != nil {
		sugar.Errorw("failed to query index", "param", info, "err", err)
		return errors.WithStack(err)
	}

	s, err := loadMatrix(info.MatrixNum)
	if err != nil {
		return errors.WithStack(err)
	}

	// 根据数据库中数据，更新
	// for i, _ := range s.Sections {
	// 	for j, _ := range s.Sections[i].Indexes {
	// 		for _, v := range items {
	// 			if v.Code == s.Sections[i].Indexes[j].Code {
	// 				s.Sections[i].Indexes[j].Value = v.Value
	// 				continue
	// 			}
	// 		}
	// 	}
	// }

	for _, section := range s.Sections {
		for j, idx := range section.Indexes {
			for _, v := range items {
				if v.Code == idx.Code {
					section.Indexes[j].Value = v.Value
					continue
				}
			}
		}
	}

	p := struct {
		Period        string
		MatrixNum     string
		ID            string
		Title         string
		Sections      []matrixSectionT
		PeriodOptions []string
	}{
		Period:        info.Period,
		MatrixNum:     info.MatrixNum,
		ID:            info.SubmitUser,
		Title:         s.Title,
		Sections:      s.Sections,
		PeriodOptions: []string{"2016", "2017", "2018"},
	}
	c.Set("p", p)

	return c.Render(200, r.HTML("matrix/open.html", "surveys/simple.html"))
}

type matrixT struct {
	Title    string           `yaml:"title"`
	Sections []matrixSectionT `yaml:"sections"`
}

type matrixSectionT struct {
	Section string   `yaml:"section"`
	Indexes []indexT `yaml:"indexes"`
}

type indexT struct {
	Name        string `yaml:"name"`
	Code        string `yaml:"code"`
	Description string `yaml:"description"`
	Unit        string `yaml:"unit"`
	Value       string `yaml:"value"`
}

func loadMatrix(num string) (matrixT, error) {
	var s matrixT

	// dir := envy.Get("MatrixDir", "")
	// fileName := dir + "/" + num + ".yaml"
	// source, err := ioutil.ReadFile(fileName)

	fileName := "matrix/" + num + ".yaml"
	source, err := base.Box.Find(fileName)
	if err != nil {
		return s, err
	}

	err = yaml.Unmarshal(source, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

// MatrixSubmit submit results
func MatrixSubmit(c buffalo.Context) error {
	// 从 querystring 中获取 问卷编号
	var err error
	// 解析 querystring
	info := struct {
		MatrixNum  string `form:"num" json:"num"`
		SubmitUser string `form:"id" json:"id"`
		Period     string `form:"period" json:"period"`
		Company    string `form:"company" json:"company"`
		Version    string `form:"version" json:"version" db:"version"`
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

	s, err := loadMatrix(info.MatrixNum)
	if err != nil {
		return errors.WithStack(err)
	}

	// 获取 form data
	err = c.Request().ParseForm()
	if err != nil {
		return errors.WithStack(err)
	}

	var matrices models.Matrices
	for _, m := range s.Sections {
		for _, i := range m.Indexes {
			v := c.Request().Form[i.Code]
			// f, err := strconv.ParseFloat(v[0], 64)
			// if err != nil {
			// 	c.LogField("raw", v)
			// 	f = 0.0
			// }
			fmt.Printf("[%s]: %s\n", i.Code, v)

			a := models.Matrix{
				Company:    info.Company,
				Version:    info.Version,
				Matrix:     info.MatrixNum,
				Period:     info.Period,
				SubmitUser: info.SubmitUser,
				Code:       i.Code,
				Value:      v[0],
			}

			matrices = append(matrices, a)
		}
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.WithStack(errors.New("no transaction found"))
	}

	kv := base.KvCache()
	sugar := base.Sugar()
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

	for _, m := range matrices {
		var ids models.IDList

		// 按照关键字段查找
		err = nstmt.Select(&ids, m)
		if err != nil {
			sugar.Errorw("failed to query index", "i", m, "err", err)
			return errors.WithStack(err)
		}

		if len(ids) == 0 {
			// 新增
			// sugar.Infow("create", "i", m)

			_, err := tx.ValidateAndCreate(&m)
			if err != nil {
				// return errors.WithStack(err)
				sugar.Errorw("failed to save index", "i", m, "err", err)
				continue
			}
		} else {
			// 更新
			m.ID = ids[0].ID
			m.UpdatedAt = now
			// sugar.Infow("update", "i", m)

			_, err = db.NamedExec(updSQL, m)
			if err != nil {
				sugar.Errorw("update failed", "err", err, "entry", m)
			}
		}
	}

	// 提交表单
	p := struct {
		Title  string
		Period string
		From   string
	}{
		Title:  s.Title,
		Period: info.Period,
		From:   c.Request().RequestURI,
	}
	c.Set("p", p)
	return c.Render(200, r.HTML("matrix/success.html", "surveys/simple.html"))
}
