package excel

import (
	"github.com/xuri/excelize/v2"
)

func ReadExcel(filename string) (file *excelize.File, err error) {
	file, err = excelize.OpenFile(filename)
	return
}
