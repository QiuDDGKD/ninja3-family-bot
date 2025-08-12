package tools

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

// 读取 xlsx 文件
func ReadXLSX(filePath string) ([][]string, error) {
	// 打开文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer f.Close()

	// 获取所有工作表名称
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("文件中没有工作表")
	}

	// 读取第一个工作表的内容
	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("无法读取工作表内容: %v", err)
	}

	return rows, nil
}
