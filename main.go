package main

import (
	"flag"
	"fmt"
	"github.com/guange2015/utils"
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
	_ = utils.ReadLine(fileName, func(string) {
		count += 1
	}, false)
	return count
}

func clone(giturl string) (error, string, string) {
	vs := strings.Split(giturl, "/")

	var subpath string
	if len(vs) > 2 {
		subpath = vs[len(vs)-2]
	}

	clonePath := filepath.Join(outpath, subpath)

	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		_ = os.Mkdir(clonePath, 0777) //0777也可以os.ModePerm
	} else {
		//如果存在已克隆过的目录，删除掉
		gitPath := filepath.Join(clonePath, vs[len(vs)-1])
		_ = os.RemoveAll(gitPath)
	}

	cmd := fmt.Sprintf("cd %v && /usr/bin/git clone --bare %s",
		clonePath, giturl)
	err, out, e := utils.Shellout(cmd)

	fmt.Println("out:", out)
	fmt.Println("err:", e)
	if err != nil {
		return err, out, e
	}
	return nil, out, e
}

func isCloned(line string) bool {
	var cloned = false
	_ = utils.ReadLine("fail.txt", func(line1 string) {
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

	_ = utils.ReadLine("success.txt", func(line1 string) {
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
			err, _, e := clone(line)
			if err != nil {
				//写入失败文件
				fmt.Println("克隆失败:", err)

				//时间，路径，原因
				writeLine := fmt.Sprintf("%v,%v,%v\n",
					utils.GetNowTime(), line, e)
				_ = utils.WriteLine("fail.txt", writeLine)
			} else {
				//写入成功文件
				fmt.Println("克隆成功")

				//时间，路径
				writeLine := fmt.Sprintf("%v,%v\n", utils.GetNowTime(), line)
				_ = utils.WriteLine("success.txt", writeLine)
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

	_ = utils.ReadLine(filePath, func(line string) {
		curLine += 1
		c <- work{totalCount, curLine, line}
	}, false)

	wg.Wait()

}
