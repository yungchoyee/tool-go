package main

import (
	"fmt"
	"os"
	"sync"
	"time"
	"tool-go/excel"
	"tool-go/file"
)

func main() {
	fmt.Println("开始执行....")

	start := time.Now()
	curPath, _ := os.Getwd() //当前目录
	args := os.Args
	if len(args) < 3 {
		fmt.Println("缺少参数")
		return
	}
	projectDir := args[1]
	path := curPath + "/" + projectDir
	exists, err := file.PathExists(path)
	if !exists || err != nil {
		fmt.Println("项目目录不存在")
		return
	}

	fileClient := file.NewReplaceFileClient()

	excelPath := curPath + "/" + args[2]
	exists, err = file.PathExists(excelPath)
	if !exists || err != nil {
		fmt.Println("excel不存在")
		return
	}

	excelObj, err := excel.ReadExcel(excelPath)
	if err != nil {
		fmt.Println("读取Excel错误")
		return
	}

	list := excelObj.GetSheetList()
	if len(list) <= 1 {
		fmt.Println("excel表没有sheet")
		return
	}
	sheet := list[0]

	rows, err := excelObj.Rows(sheet)
	if err != nil {
		fmt.Println("没有读取到Excel行数")
		return
	}
	replaceMap := make(map[string]string, 0)
	rowNum := 1
	for rows.Next() {
		bValue, _ := excelObj.GetCellValue(sheet, fmt.Sprintf("B%d", rowNum))
		cValue, _ := excelObj.GetCellValue(sheet, fmt.Sprintf("C%d", rowNum))
		oldStr, _ := excelObj.GetCellValue(sheet, fmt.Sprintf("D%d", rowNum))
		rowNum++

		if _, ok := replaceMap[oldStr]; !ok {
			replaceStr := "I18n.of(FrogApp().currentContext)." + bValue + cValue
			replaceMap["\""+oldStr+"\""] = replaceStr
			replaceMap["'"+oldStr+"'"] = replaceStr
		} else {
			fmt.Println(fmt.Sprintf("有重复翻译的内容:rows:%d\n", rowNum))
		}
	}
	fileClient.Filter = replaceMap
	go func() {
		err = fileClient.BatchReplace(path)
		if err != nil {
			fmt.Println("扫描文件夹出错" + err.Error())
			return
		}
	}()
	fileClient.HandleFuncList = append(fileClient.HandleFuncList, fileClient.ReplaceContent)
	wg := sync.WaitGroup{}
	for filename := range fileClient.ProcessCh {
		wg.Add(1)
		go func(filePath string) {
			fmt.Println(filePath)
			for _, fn := range fileClient.HandleFuncList {
				err = fn(filePath)
				if err != nil {
					fmt.Println("替换文件内容发生错误,err:" + err.Error() + ",filename:" + filePath)
				}
			}
			wg.Done()
		}(filename)
	}
	wg.Wait()
	end := time.Now()
	fmt.Println(fmt.Sprintf("开始执行时间:%s", start.Format("2006-01-02 15:04:05")))
	fmt.Println(fmt.Sprintf("结束时间:%s", end.Format("2006-01-02 15:04:05")))
	fmt.Println(fmt.Sprintf("执行耗时:%d 毫秒", end.Sub(start)/time.Millisecond))
}
