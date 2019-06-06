package base

import (
	"bytes"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"

	"pony/models"
)

// ProcessPOItems read and save po items
func ProcessPOItems() {
	var err error
	sugar := Sugar()

	fmap := map[string]func() (models.PoItems, error){
		"强化": readPoItems1,
		"句容": readPoItems2,
		"康逸": readPoItems3,
	}

	var all models.PoItems
	for k, fn := range fmap {
		items, err := fn()
		if err != nil {
			sugar.Errorw("read po items", "err", err)
			return
		}
		sugar.Debugw("load", "plant", k, "count", len(items))

		all = append(all, items...)
	}

	err = savePOItems(all)
	if err != nil {
		sugar.Errorw("save po items", "err", err)
		return
	}
	sugar.Debugw("save completed", "count", len(all))
}

// readPOItems1 load 强化 po items from excel file
func readPoItems1() (models.PoItems, error) {
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
	sheetName := "Sheet3"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 50000

	var results models.PoItems
	var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.PoItem{
			Company:    rowData(0),
			PoNum:      rowData(2),
			VendorName: rowData(3),
			MatName:    rowData(4),
			ItemUnit:   rowData(6),
			Cate2:      rowData(8),
			Cate1:      rowData(9),
			Operator:   rowData(10),
		}
		t, err = time.Parse("01-02-06", rowData(1))
		if err != nil {
			sugar.Errorw("po date parse", "line", r, "po date", rowData(1))
		}
		entry.PoDate = t

		t, err = time.Parse("2006/1/2", rowData(13))
		if err != nil {
			sugar.Errorw("inbound date parse", "line", r, "inbound date", rowData(13))
		}
		entry.InboundDate = t

		f, err = strconv.ParseFloat(rowData(5), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(5))
		}
		entry.ItemQty = f

		f, err = strconv.ParseFloat(rowData(7), 64)
		if err != nil {
			sugar.Errorw("unit price", "line", r, "unit price", rowData(7))
		}
		entry.UnitPrice = f

		f, err = strconv.ParseFloat(rowData(11), 64)
		if err != nil {
			sugar.Errorw("inbound qty", "line", r, "unit price", rowData(11))
		}
		entry.InboundQty = f

		f, err = strconv.ParseFloat(rowData(12), 64)
		if err != nil {
			sugar.Errorw("outstanding qty", "line", r, "unit price", rowData(12))
		}
		entry.OutstandingQty = f

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readPOItems2 load 句容 po items from excel file
func readPoItems2() (models.PoItems, error) {
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
	sheetName := "Sheet3"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 50000

	var results models.PoItems
	var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.PoItem{
			Company:    rowData(0),
			PoNum:      rowData(2),
			VendorName: rowData(3),
			LineNum:    rowData(4),
			MatName:    rowData(5),
			ItemUnit:   rowData(7),
			Cate2:      rowData(9),
			Cate1:      rowData(10),
			Operator:   rowData(11),
			DnNum:      rowData(12),
			DnItem:     rowData(13),
		}

		t, err = time.Parse("2006/1/2", rowData(14))
		if err != nil {
			sugar.Errorw("inbound date parse", "line", r, "inbound date", rowData(14))
		}
		entry.InboundDate = t

		f, err = strconv.ParseFloat(rowData(6), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(6))
		}
		entry.ItemQty = f

		f, err = strconv.ParseFloat(rowData(8), 64)
		if err != nil {
			sugar.Errorw("unit price", "line", r, "unit price", rowData(8))
		}
		entry.UnitPrice = f

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readPOItems3 load 康逸 po items from excel file
func readPoItems3() (models.PoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/康逸.xlsx"
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

	const MaxLoop = 50000

	var results models.PoItems
	var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.PoItem{
			VendorName:  rowData(0),
			Company:     rowData(1),
			PoNum:       rowData(2),
			ItemStatus:  rowData(3),
			Operator:    rowData(5),
			Cate2:       rowData(6),
			MatName:     rowData(7),
			MatCode:     rowData(8),
			MatSpec:     rowData(9),
			ItemUnit:    rowData(11),
			InboundUnit: rowData(13),
		}

		t, err = time.Parse("2006-01-02", rowData(4))
		if err != nil {
			sugar.Errorw("po date parse", "line", r, "inbound date", rowData(4))
		}
		entry.PoDate = t

		t, err = time.Parse("2006-01-02", rowData(17))
		if err != nil {
			sugar.Errorw("plan date parse", "line", r, "plan date", rowData(17))
		}
		entry.PlannedDate = t

		if rowData(18) != "" {
			t, err = time.Parse("2006-01-02", rowData(18))
			if err != nil {
				sugar.Errorw("inbound date parse", "line", r, "inbound date", rowData(18))
			}
			entry.InboundDate = t
		}

		f, err = strconv.ParseFloat(rowData(10), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(10))
		}
		entry.ItemQty = f

		if rowData(12) != "" {
			f, err = strconv.ParseFloat(rowData(12), 64)
			if err != nil {
				sugar.Errorw("inbound qty", "line", r, "item qty", rowData(12))
			}
			entry.InboundQty = f
		}

		f, err = strconv.ParseFloat(rowData(14), 64)
		if err != nil {
			sugar.Errorw("arrived book quantity", "line", r, "qty", rowData(14))
		}
		entry.ArriveBookQty = f

		f, err = strconv.ParseFloat(rowData(15), 64)
		if err != nil {
			sugar.Errorw("booked quantity", "line", r, "qty", rowData(15))
		}
		entry.BookedQty = f

		f, err = strconv.ParseFloat(rowData(16), 64)
		if err != nil {
			sugar.Errorw("unbooked quantity", "line", r, "qty", rowData(16))
		}
		entry.UnbookedQty = f

		if rowData(19) == "" {
			entry.DelayedDays = 0
		} else {
			f, err = strconv.ParseFloat(rowData(19), 64)
			if err != nil {
				sugar.Errorw("delayed days", "line", r, "qty", rowData(19))
			}
			entry.DelayedDays = f
		}

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// savePOItems save po items to database
func savePOItems(items models.PoItems) error {
	sugar := Sugar()

	// save to DB
	db := models.DB
	_ = db.Transaction(func(tx *pop.Connection) error {
		// delete all
		// _, err = tx.Store.Exec("truncate table po_items; ")
		// if err != nil {
		// 	sugar.Errorw("truncate error", "err", err)
		// }

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

// ProcessMOItems read and save po items
func ProcessMOItems() {
	var err error
	sugar := Sugar()

	fmap := map[string]func() (models.MoItems, error){
		"强化": readMoItems1,
		"句容": readMoItems2,
		"康逸": readMoItems3,
	}

	var all models.MoItems
	for k, fn := range fmap {
		items, err := fn()
		if err != nil {
			sugar.Errorw("read mo items", "err", err)
			return
		}
		sugar.Debugw("load", "plant", k, "count", len(items))

		all = append(all, items...)
	}

	err = saveMOItems(all)
	if err != nil {
		sugar.Errorw("save mo items", "err", err)
		return
	}

	sugar.Debugw("save completed", "count", len(all))
}

// readMOItems1 load 强化 po items from excel file
func readMoItems1() (models.MoItems, error) {
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
	sheetName := "Sheet4"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 100000

	var results models.MoItems
	var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.MoItem{
			Company:   rowData(0),
			MoNum:     rowData(1),
			WorkOrder: rowData(2),
			MatCode:   rowData(3),
			MatName:   rowData(4),
			Line:      rowData(5),
			MoType:    rowData(6),
			ItemUnit:  rowData(7),
		}
		if rowData(8) != "" {
			t, err = time.Parse("2006/1/2", rowData(8))
			if err != nil {
				sugar.Errorw("start date parse", "line", r, "start date", rowData(8))
			}
			entry.StartDate = t
		}
		if rowData(9) != "" {
			t, err = time.Parse("2006/1/2", rowData(9))
			if err != nil {
				sugar.Errorw("end date parse", "line", r, "end date", rowData(9))
			}
			entry.EndDate = t
		}

		t, err = time.Parse("2006/1/2", rowData(11))
		if err != nil {
			sugar.Errorw("mo date parse", "line", r, "mo date", rowData(11))
		}
		entry.MoDate = t

		f, err = strconv.ParseFloat(rowData(10), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(10))
		}
		entry.ItemQty = f

		entry.Source = "强化"

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readMOItems2 load 句容 po items from excel file
func readMoItems2() (models.MoItems, error) {
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
	sheetName := "Sheet4"
	// 本 worksheet 写入的 helper
	readHelper := SheetReader(xlsx, sheetName)

	const MaxLoop = 50000

	var results models.MoItems
	var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			break
		}

		entry := models.MoItem{
			Company:   rowData(0),
			Line:      rowData(1),
			MoNum:     rowData(3),
			ItemNum:   rowData(4),
			MatName:   rowData(5),
			ItemUnit:  rowData(7),
			Cate2:     rowData(8),
			Cate1:     rowData(9),
			Warehouse: rowData(11),
		}
		t, err = time.Parse("2006/1/2", rowData(2))
		if err != nil {
			sugar.Errorw("mo date parse", "line", r, "mo date", rowData(2))
		}
		entry.MoDate = t

		t, err = time.Parse("2006/1/2", rowData(10))
		if err != nil {
			sugar.Errorw("inbound date parse", "line", r, "inbound date", rowData(10))
		}
		entry.InboundDate = t

		f, err = strconv.ParseFloat(rowData(6), 64)
		if err != nil {
			sugar.Errorw("item qty", "line", r, "item qty", rowData(6))
		}
		entry.ItemQty = f

		entry.Source = "句容"

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// readMOItems3 load 康逸 po items from excel file
func readMoItems3() (models.MoItems, error) {
	var err error
	sugar := Sugar()

	// 从文件中读取
	fileName := "dump/康逸生产.xlsx"
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

	const MaxLoop = 50000

	var results models.MoItems
	// var t time.Time
	var f float64
	for r := 3; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := RowReader(readHelper, r)

		// 第四列：
		if rowData(3) == "" {
			break
		}

		entry := models.MoItem{
			Company: "康逸",
			Line:    rowData(0),
			Shift:   rowData(1),
			Step:    rowData(2),
			MatCode: rowData(3),
			MatName: rowData(4),
			MatSpec: rowData(5),
			MoNum:   rowData(6),
		}

		if rowData(7) != "" {
			f, err = strconv.ParseFloat(rowData(7), 64)
			if err != nil {
				sugar.Errorw("main qty", "line", r, "main qty", rowData(7))
			}
			entry.MainMatQty = f
		}

		if rowData(8) != "" {
			f, err = strconv.ParseFloat(rowData(8), 64)
			if err != nil {
				sugar.Errorw("Input qty1", "line", r, "input qty1", rowData(8))
			}
			entry.InputMatQty1 = f
		}

		if rowData(9) != "" {
			f, err = strconv.ParseFloat(rowData(9), 64)
			if err != nil {
				sugar.Errorw("Input qty2", "line", r, "input qty2", rowData(9))
			}
			entry.InputMatQty2 = f
		}

		if rowData(10) != "" {
			f, err = strconv.ParseFloat(rowData(10), 64)
			if err != nil {
				sugar.Errorw("Claim qty1", "line", r, "claim qty1", rowData(10))
			}
			entry.ClaimQty1 = f
		}

		if rowData(11) != "" {
			f, err = strconv.ParseFloat(rowData(11), 64)
			if err != nil {
				sugar.Errorw("Claim qty2", "line", r, "claim qty2", rowData(11))
			}
			entry.ClaimQty2 = f
		}

		entry.Source = "康逸"

		results = append(results, entry)
	}

	sugar.Infow("load formA", "file", fileName, "items count", len(results))

	return results, nil
}

// saveMOItems save mo items to database
func saveMOItems(items models.MoItems) error {
	sugar := Sugar()

	// save to DB
	db := models.DB
	_ = db.Transaction(func(tx *pop.Connection) error {
		// delete all
		// _, err = tx.Store.Exec("truncate table mo_items; ")
		// if err != nil {
		// 	sugar.Errorw("truncate error", "err", err)
		// }

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
