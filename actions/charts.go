package actions

import (
	"pony/base"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// ChartFirst demo chart
func ChartFirst(c buffalo.Context) error {
	return c.Render(200, r.HTML("charts/first.html", "layout/empty.html"))
}

func isStrSliceContainStr(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// readSkipWords get all skip words
func readSkipWords() []string {
	var err error
	sugar := base.Sugar()

	dir := envy.Get("WordDir", "./config/words")
	file := dir + "/skip.txt"
	lines, err := base.ReadLines(file)
	if err != nil {
		sugar.Errorw("failed to get skipword", "file", file)
		return []string{}
	}
	return lines
}

// WordCloudHandle return word frequencies
func WordCloudHandle(c buffalo.Context) error {
	var err error

	type resultT struct {
		Word    string `db:"word" json:"word"`
		Count   int    `db:"count" json:"count"`
		DocName string `db:"doc_name" json:"doc_name"`
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

	cmdSQL := kv.GetCommand("chart.wordCloud", nil)

	items := []resultT{}
	err = db.Select(&items, cmdSQL)
	if err != nil {
		sugar.Errorf("wordCloud data failed", "sql", cmdSQL, "err", err)
		return errors.WithStack(err)
	}

	// skip word
	filtered := []resultT{}
	excludes := readSkipWords()
	for _, i := range items {
		if isStrSliceContainStr(excludes, i.Word) {
			continue
		}
		filtered = append(filtered, i)
	}

	return c.Render(200, r.JSON(filtered))
}

// WordFreqHandle return word frequencies
func WordFreqHandle(c buffalo.Context) error {
	var err error

	type resultT struct {
		Level string `db:"level" json:"level"`
		Count int    `db:"count" json:"count"`
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

	cmdSQL := kv.GetCommand("chart.wordFreq", nil)

	items := []resultT{}
	err = db.Select(&items, cmdSQL)
	if err != nil {
		sugar.Errorf("wordFreq data failed", "sql", cmdSQL, "err", err)
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(items))
}

// wordDistHandle return distribution
func wordDistHandle(c buffalo.Context) error {
	var err error

	type resultT struct {
		Word      string `db:"word" json:"word"`
		WordCount int    `db:"wc_count" json:"wc_count"`
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

	cmdSQL := kv.GetCommand("chart.wordDist", nil)

	items := []resultT{}
	err = db.Select(&items, cmdSQL)
	if err != nil {
		sugar.Errorf("wordFreq data failed", "sql", cmdSQL, "err", err)
		return errors.WithStack(err)
	}

	return c.Render(200, r.JSON(items))
}
