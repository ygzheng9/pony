package base

import (
	"bytes"
	"strconv"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"github.com/360EntSecGroup-Skylar/excelize"

	"pony/models"
)

// ProcessSoItems read and save po items
func ProcessSoItems() {
	sugar := Sugar()

	fmap := map[string]func() (models.SoItems, error){
		"集团":     readSoItems1,
		"MDS2":   readSoItems2,
		"PDBS-1": readSoItems31,
		"PDBS-2": readSoItems32,
	}

	var all models.SoItems
	for k, fn := range fmap {
		items, err := fn()
		if err != nil {
			sugar.Errorw("read so items", "err", err)
			return
		}
		sugar.Debugw("load", "plant", k, "count", len(items))

		all = append(all, items...)
	}

	err := saveSoItems(all)
	if err != nil {
		sugar.Errorw("save so items", "err", err)
		return
	}
	sugar.Debugw("save completed", "count", len(all))
}

// readSOItems1 load 集团 SO items from excel file
func readSoItems1() (models.SoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/PDBS3渠道发货清单2018.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.SoItems
	var t time.Time
	var f float64
	var s string
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		if rowData(11) == "" {
			// 	正常情况
			entry := models.SoItem{
				Company:   "集团",
				OrderNum:  rowData(0),
				CustNum:   rowData(1),
				Category:  rowData(2),
				Serial:    rowData(3),
				MatName:   rowData(4),
				MatModel:  rowData(5),
				Period:    rowData(7),
				SalesType: rowData(8),
				Source:    "PDBS3",
			}

			f, err = strconv.ParseFloat(rowData(6), 64)
			if err != nil {
				sugar.Errorw("item qty", "line", r, "item qty", rowData(6))
			}
			entry.ItemQty = f

			s = rowData(9)
			if len(s) <= 8 {
				t, err = time.Parse("01-02-06", rowData(9))
				if err != nil {
					sugar.Errorw("wh date parse", "line", r, "wh date", rowData(9))
				}
				entry.WhDate = t
			} else {
				t, err = time.Parse("1/2/06 15:04", rowData(9))
				if err != nil {
					sugar.Errorw("wh date parse", "line", r, "wh date", rowData(9))
				}
				entry.WhDate = t
			}

			t, err = time.Parse("1/2/06 15:04", rowData(10))
			if err != nil {
				sugar.Errorw("doc date parse", "line", r, "doc date", rowData(10))
			}
			entry.DocDate = t

			results = append(results, entry)
		} else {
			entry := models.SoItem{
				Company:  "集团",
				OrderNum: rowData(0),
				CustNum:  rowData(1),

				Category:  rowData(3),
				Serial:    rowData(4),
				MatName:   rowData(5),
				MatModel:  rowData(6),
				Period:    rowData(7),
				SalesType: rowData(9),
				Source:    "PDBS3",
			}

			f, err = strconv.ParseFloat(rowData(7), 64)
			if err != nil {
				sugar.Errorw("item qty", "line", r, "item qty", rowData(7))
			}
			entry.ItemQty = f

			s = rowData(10)
			if len(s) <= 8 {
				t, err = time.Parse("01-02-06", s)
				if err != nil {
					sugar.Errorw("wh date parse", "line", r, "wh date", s)
				}
				entry.WhDate = t
			} else {
				t, err = time.Parse("1/2/06 15:04", s)
				if err != nil {
					sugar.Errorw("wh date parse", "line", r, "wh date", s)
				}
				entry.WhDate = t
			}

			s = rowData(11)
			t, err = time.Parse("1/2/06 15:04", s)
			if err != nil {
				sugar.Errorw("doc date parse", "line", r, "doc date", s)
			}
			entry.DocDate = t

			results = append(results, entry)
		}
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readSOItems2 load 集团 SO items from excel file
func readSoItems2() (models.SoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/2018MDS2渠道出库.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.SoItems
	var t time.Time
	var f float64
	var s string
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		// 	正常情况
		entry := models.SoItem{
			Company:   rowData(1),
			OrderNum:  rowData(0),
			CustNum:   rowData(2),
			Category:  rowData(3),
			Serial:    rowData(4),
			MatModel:  rowData(5),
			Period:    rowData(7),
			BookParty: rowData(8),
			MoveType:  rowData(9),
			SalesType: rowData(10),
			Warehouse: rowData(12),
			Remark:    rowData(14),
			Source:    "MDS2",
		}

		s = rowData(6)
		f, err = strconv.ParseFloat(s, 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", s)
		}
		entry.ItemQty = f

		s = rowData(11)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			sugar.Errorw("wh date parse", "line", r, "wh date", s)
		}
		entry.WhDate = t

		s = rowData(13)
		if len(s) >= 11 {
			t, err = time.Parse("2006-01-02 15:04:05", s)
			if err != nil {
				sugar.Errorw("doc date parse", "line", r, "doc date", s)
			}
			entry.DocDate = t
		} else {
			t, err = time.Parse("2006-01-02", s)
			if err != nil {
				sugar.Errorw("doc date parse", "line", r, "doc date", s)
			}
			entry.DocDate = t
		}

		results = append(results, entry)

	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

func readSoItems31() (models.SoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/2018pdbs渠道出库.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.SoItems
	var t time.Time
	var f float64
	var s string
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		// 	正常情况
		entry := models.SoItem{
			Company:   rowData(1),
			OrderNum:  rowData(0),
			CustNum:   rowData(2),
			Category:  rowData(3),
			Serial:    rowData(4),
			MatName:   rowData(5),
			MatModel:  rowData(6),
			Period:    rowData(8),
			BookParty: rowData(9),
			MoveType:  rowData(10),
			SalesType: rowData(11),
			Warehouse: rowData(14),
			Remark:    rowData(17),
			Source:    "PDBS-1",
		}

		s = rowData(7)
		if s == "" {
			entry.ItemQty = 0
		} else {
			f, err = strconv.ParseFloat(s, 64)
			if err != nil {
				sugar.Errorw("item qty", "line", r, "item qty", s)
			}
			entry.ItemQty = f
		}

		s = rowData(13)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			sugar.Errorw("wh date parse", "line", r, "wh date", s)
		}
		entry.WhDate = t

		s = rowData(16)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			sugar.Errorw("doc date parse", "line", r, "doc date", s)
		}
		entry.DocDate = t

		results = append(results, entry)

	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readSOItems32 load 集团 SO items from excel file
func readSoItems32() (models.SoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/2018pdbs渠道出库.xlsx"
	source, err := ABox.Find(fileName)
	if err != nil {
		return nil, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return nil, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet2"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 150000

	var results models.SoItems
	var t time.Time
	var f float64
	var s string
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		// 	正常情况
		entry := models.SoItem{
			Company:   rowData(1),
			OrderNum:  rowData(0),
			CustNum:   rowData(2),
			Category:  rowData(3),
			Serial:    rowData(4),
			MatName:   rowData(5),
			MatModel:  rowData(6),
			Period:    rowData(8),
			BookParty: rowData(9),
			MoveType:  rowData(10),
			SalesType: rowData(11),
			Warehouse: rowData(14),
			Remark:    rowData(17),
			Source:    "PDBS-2",
		}

		s = rowData(7)
		if s == "" {
			entry.ItemQty = 0
		} else {
			f, err = strconv.ParseFloat(s, 64)
			if err != nil {
				sugar.Errorw("item qty", "line", r, "item qty", s)
			}
			entry.ItemQty = f
		}

		s = rowData(13)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			sugar.Errorw("wh date parse", "line", r, "wh date", s)
		}
		entry.WhDate = t

		s = rowData(16)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			sugar.Errorw("doc date parse", "line", r, "doc date", s)
		}
		entry.DocDate = t

		results = append(results, entry)

	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// saveSoItems save mo items to database
func saveSoItems(items models.SoItems) error {
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
