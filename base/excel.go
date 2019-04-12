package base

import (
	"fmt"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// DistString 字符串去重
func DistString(items []string) []string {
	distInvs := []string{}
	for _, i := range items {
		flag := true
		for _, j := range distInvs {
			if i == j {
				flag = false
				break
			}
		}
		if flag {
			distInvs = append(distInvs, i)
		}
	}

	return distInvs
}

// 根据坐标，返回单元格名字
// r 从 1 开始，也即：1 代表 第1行，2 代表 第2行； A1, A2
// c 从 0 开始，也即 0 -> A, 1 -> B
func cellAxis(r int) func(c int) string {
	return func(c int) string {
		return fmt.Sprintf("%s%d", excelize.ToAlphaString(c), r)
	}
}

// SheetWriterT sheet 的写入的辅助函数
type SheetWriterT func(axis string, value interface{})

// SheetWriter 写
func SheetWriter(workbook *excelize.File, sheetName string) SheetWriterT {
	return func(axis string, value interface{}) {
		workbook.SetCellValue(sheetName, axis, value)
	}
}

// RowWriter 直接写入一行中的某一列
// 设置 workbook，sheetName，以及 rowIndex，每次只需要再输入 colIndex，value
// rowIndex 从 1 开始； A1, B1
// colIndex 从 0 开始； 0 -> A, 1 -> B
func RowWriter(shWriter SheetWriterT, r int) func(c int, value interface{}) {
	axis := cellAxis(r)
	return func(c int, value interface{}) {
		shWriter(axis(c), value)
	}
}

// SheetReaderT sheet 的读取的辅助函数
type SheetReaderT func(axis string) string

// SheetReader 读
func SheetReader(workbook *excelize.File, sheetName string) SheetReaderT {
	return func(axis string) string {
		return workbook.GetCellValue(sheetName, axis)
	}
}

// RowReader 一行中的某一列
// 设置 workbook，sheetName，以及 rowIndex，每次只需要再输入 colIndex
// rowIndex 从 1 开始； A1, B1
// colIndex 从 0 开始； 0 -> A, 1 -> B
func RowReader(shReader SheetReaderT, r int) func(c int) string {
	axis := cellAxis(r)
	return func(c int) string {
		return strings.TrimSpace(shReader(axis(c)))
	}
}
