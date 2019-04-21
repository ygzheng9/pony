package actions

import (
	"bytes"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gobuffalo/buffalo"

	"pony/base"
)

// scatterShow show scatter for tasks
func scatterShow(c buffalo.Context) error {
	return c.Render(200, r.HTML("games/scatter.html", "games/layout.html"))
}

// scatterData get data for scatter
func scatterData(c buffalo.Context) error {
	var err error
	items, err := scatterLoad()
	if err != nil {
		return c.Render(200, r.JSON(H{
			"status":  10,
			"message": "can not load scatter data",
		}))
	}
	return c.Render(200, r.JSON(H{
		"status": 0,
		"items":  items,
	}))
}

// scatterItem each line for scatter data
type scatterItem struct {
	Group string `json:"group"`
	Item  string `json:"item"`
	IndA  string `json:"ind_a"`
	IndB  string `json:"ind_b"`
}

// scatterLoad load data from file
func scatterLoad() ([]scatterItem, error) {
	var err error
	sugar := base.Sugar()

	var results []scatterItem

	fileName := "scatter/S01.xlsx"
	source, err := base.Box.Find(fileName)
	if err != nil {
		sugar.Errorw("no file", "file", fileName, "err", err)
		return results, err
	}

	xlsx, err := excelize.OpenReader(bytes.NewReader(source))
	if err != nil {
		return results, err
	}

	//  上载的文件必须有这个 tab
	sheetName := "Sheet1"
	// 本 worksheet 写入的 helper
	readHelper := base.SheetReader(xlsx, sheetName)

	const MaxLoop = 2000
	accEmpty := 0
	for r := 2; r < MaxLoop; r++ {
		// 从第三行开始
		rowData := base.RowReader(readHelper, r)

		// 第一列：
		if rowData(0) == "" {
			accEmpty++
			// 连续空 5 行，表示结束
			if accEmpty > 5 {
				break
			} else {
				continue
			}
		}
		accEmpty = 0

		entry := scatterItem{
			Group: rowData(0),
			Item:  rowData(1),
			IndA:  rowData(2),
			IndB:  rowData(3),
		}

		results = append(results, entry)
	}
	sugar.Debugw("item count", "count", len(results))

	return results, nil
}
