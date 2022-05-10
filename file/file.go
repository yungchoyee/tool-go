package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type ReplaceFile struct {
	FileList       []string
	Filter         map[string]string
	HandleFuncList []func(file string) error
	ProcessCh      chan string
}

//创建一个新的客户端
func NewReplaceFileClient() *ReplaceFile {
	client := &ReplaceFile{
		FileList: []string{},
		Filter:   map[string]string{},
	}
	client.ProcessCh = make(chan string, 100)
	return client
}

// 获取目录下所有文件名
func (f *ReplaceFile) GetFileList(path string) (fileList []string, err error) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	separator := string(os.PathSeparator)
	for _, fileInfo := range fileInfos {
		currentFilePath := path + separator + fileInfo.Name()
		if fileInfo.IsDir() {
			// 如果是文件夹，递归去处理
			_, err = f.GetFileList(currentFilePath)
			if err != nil {
				return
			}
		} else {
			f.FileList = append(f.FileList, currentFilePath)
		}
	}
	fileList = f.FileList
	return
}

// 获取目录下所有文件名
func (f *ReplaceFile) BatchReplace(path string) (err error) {
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	separator := string(os.PathSeparator)
	for _, fileInfo := range fileInfos {
		currentFilePath := path + separator + fileInfo.Name()
		if fileInfo.IsDir() {
			// 如果是文件夹，递归去处理
			err = f.BatchReplace(currentFilePath)
			if err != nil {
				return
			}
		} else {
			f.ProcessCh <- currentFilePath
		}
	}
	close(f.ProcessCh)
	return
}

//替换文件内容
func (f *ReplaceFile) ReplaceContent(file string) (err error) {
	in, err := os.OpenFile(file, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer in.Close()
	br := bufio.NewReader(in)
	var dataBlocks []byte
	var line []byte
	for {
		line, _, err = br.ReadLine()
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return
		}
		for oldStr, newStr := range f.Filter {
			if strings.Contains(string(line), oldStr) {
				log.Println(fmt.Sprintf("%s | %s | %s", file, oldStr, newStr))
				line = []byte(strings.Replace(string(line), oldStr, newStr, -1))
			}
		}
		line = []byte(string(line) + "\n")
		dataBlocks = append(dataBlocks, line...)
	}
	in.WriteAt(dataBlocks, 0)
	return
}

//判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}
