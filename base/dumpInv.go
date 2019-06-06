package base

import (
	"bytes"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"pony/models"
)

// ProcessInvItems read and save inv items
func ProcessInvItems() {
	sugar := Sugar()

	fmap := map[string]func() (models.InvItems, error){
		"强化": readInvItems1,
		"句容": readInvItems2,
		// "康逸": readInvItems3,
	}

	var all models.InvItems
	for k, fn := range fmap {
		items, err := fn()
		if err != nil {
			sugar.Errorw("read inv items", "err", err)
			return
		}
		sugar.Debugw("load", "plant", k, "count", len(items))

		all = append(all, items...)
	}

	err := saveInvItems(all)
	if err != nil {
		sugar.Errorw("save po items", "err", err)
		return
	}
	sugar.Debugw("save completed", "count", len(all))
}

// readInvItems1 load 强化 inv items from excel file
func readInvItems1() (models.InvItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/强化.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet5"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.InvItems
	// var t time.Time
	var f float64
	for r := 2; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.InvItem{
			Company:   rowData(0),
			Warehouse: rowData(1),
			Year:      "2018",
			Month:     rowData(2),
			MatCode:   rowData(3),
			MatName:   rowData(4),
			MatGrade:  rowData(6),
			Cate2:     rowData(7),
			Cate1:     rowData(8),
		}

		f, err = strconv.ParseFloat(rowData(5), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(5))
		}
		entry.MatQty = f
		entry.Source = "强化"
		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readInvItems2 load 句容 inv items from excel file
func readInvItems2() (models.InvItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/句容.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet5"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.InvItems
	// var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.InvItem{
			Company:   rowData(0),
			Warehouse: rowData(1),
			Year:      "2018",
			Month:     rowData(2),
			MatName:   rowData(3),
			Cate2:     rowData(6),
			Cate1:     rowData(7),
		}

		f, err = strconv.ParseFloat(rowData(4), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(4))
		}
		entry.MatQty = f
		entry.Source = "句容"
		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// saveInvItems save inv items to database
func saveInvItems(items models.InvItems) error {
	sugar := Sugar()

	// save to DB
	db := models.DB
	_ = db.Transaction(func(tx *pop.Connection) error {
		// each doc
		for _, i := range items {
			verr, err := tx.ValidateAndSave(&i)
			if err != nil {
				sugar.Errorw("save failed", "item", i, "err", err)
				return errors.Wrap(err, "save failed")
			}

			if len(verr.Errors) > 0 {
				sugar.Errorw("save err", "item", i, "err", verr.Errors)
				return errors.Wrap(err, "save with error")
			}
		}
		return nil
	})
	return nil
}
