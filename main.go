package main

import (
	"flag"
	"fmt"
	"git-migration/utils"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	cc       int
	filePath string
	wg       sync.WaitGroup
	outpath  string
)

type work struct {
	totalCount int
	curCount   int
	line       string
}

func getLineCount(fileName string) int {
	count := 0
	utils.ReadLine(fileName, func(string) {
		count += 1
	}, false)
	return count
}

func clone(giturl string) error {
	vs := strings.Split(giturl, "/")

	var subpath string
	if len(vs) > 2 {
		subpath = vs[len(vs)-2]
	}

	clonePath := filepath.Join(outpath, subpath)

	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(clonePath, 0777) //0777也可以os.ModePerm
	}

	cmd := fmt.Sprintf("cd %v && /usr/bin/git clone --bare %s",
		clonePath, giturl)
	err, out, e := utils.Shellout(cmd)

	fmt.Println("out:", out)
	fmt.Println("err:", e)
	if err != nil {
		return err
	}
	return nil
}

func isCloned(line string) bool {
	var cloned = false
	utils.ReadLine("fail.txt", func(line1 string) {
		vs := strings.Split(line1, ",")
		if len(vs) > 1 {
			if strings.Compare(line, vs[1]) == 0 {
				cloned = true
			}
		}
	}, true)
	if cloned {
		return cloned
	}

	utils.ReadLine("fail.txt", func(line1 string) {
		vs := strings.Split(line1, ",")
		if len(vs) > 1 {
			if strings.Compare(line, vs[1]) == 0 {
				cloned = true
			}
		}
	}, true)

	return cloned
}

func doWork(c chan work) {
	for w := range c {
		totalCount := w.totalCount
		curLine := w.curCount
		line := w.line

		fmt.Printf("[%v/%v]开始克隆: %v\n",
			totalCount, curLine, line)

		if isCloned(line) {
			fmt.Printf("[%v/%v]已经克隆过: %v\n",
				totalCount, curLine, line)
		} else {
			err := clone(line)
			if err != nil {
				//写入失败文件
				fmt.Println("克隆失败:", err)

				//时间，路径，原因
				writeLine := fmt.Sprintf("%v,%v,%v\n",
					utils.GetNowTime(), line, err)
				utils.WriteLine("fail.txt", writeLine)
			} else {
				//写入成功文件
				fmt.Println("克隆成功")

				//时间，路径
				writeLine := fmt.Sprintf("%v,%v\n", utils.GetNowTime(), line)
				utils.WriteLine("success.txt", writeLine)
			}
		}

		wg.Done()
	}
}

func main() {
	var h bool
	flag.IntVar(&cc, "c", 1, "多进程并发数量")
	flag.StringVar(&filePath, "f", "./giturl-tpi.txt", "giturl 地址配置文件")
	flag.StringVar(&outpath, "o", "/tmp", "克隆路径")
	flag.BoolVar(&h, "h", false, "使用说明")
	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	fmt.Println("并发数：", cc)
	fmt.Println("配置文件:", filePath)

	//读取到所有的行数
	totalCount := getLineCount(filePath)
	fmt.Println("总行数: ", totalCount)

	wg.Add(totalCount)
	curLine := 0

	c := make(chan work, cc)
	for i := 0; i < cc; i += 1 {
		go doWork(c)
	}

	utils.ReadLine(filePath, func(line string) {
		curLine += 1
		c <- work{totalCount, curLine, line}
	}, false)

	wg.Wait()

}
