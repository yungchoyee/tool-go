package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 获取目录下所有文件名
func GetFileNames(path string) (filenames []string, err error) {
	filenames = make([]string, 0)
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	separator := string(os.PathSeparator)
	for _, fileInfo := range fileInfos {
		currentFilePath := path + separator + fileInfo.Name()
		if fileInfo.IsDir() {
			// 如果是文件夹，递归去处理
			_, err = GetFileNames(currentFilePath)
			if err != nil {
				return
			}
		} else {
			filenames = append(filenames, currentFilePath)
		}
	}
	return
}

//替换文件内容
func ReplaceContent(file string, rMap map[string]string) (err error) {
	in, err := os.Open(file)
	if err != nil {
		return
	}
	defer in.Close()
	br := bufio.NewReader(in)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		for oldStr, newStr := range rMap {
			newLine := strings.Replace(string(line), oldStr, newStr, -1)
			_, err = in.WriteString(newLine + "\n")
		}

		if err != nil {
			fmt.Println("write to file fail:", err)
			os.Exit(-1)
		}
		fmt.Println("done ", index)
		index++
	}
}
