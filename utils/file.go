package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

var mt sync.RWMutex

func ReadLine(fileName string, handler func(string), lock bool) error {

	if lock {
		mt.RLock()
		defer func() {
			mt.RUnlock()
		}()
	}

	f, err := os.Open(fileName)
	if err != nil {
		return err
	}

	defer f.Close()

	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		handler(line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	return nil
}

//写入一行文件
func WriteLine(fileName string, line string) error {
	mt.Lock()
	defer func() {
		mt.Unlock()
	}()

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("打开文件失败: ", err)
		return err
	}

	defer f.Close()

	_, err = f.WriteString(line)
	if err != nil {
		fmt.Println("写入文件失败: ", err)
	}

	return err
}
