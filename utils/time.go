package utils

import (
	"fmt"
	"time"
)

func GetNowTime() string {
	t1 := time.Now().Year()   //年
	t2 := time.Now().Month()  //月
	t3 := time.Now().Day()    //日
	t4 := time.Now().Hour()   //小时
	t5 := time.Now().Minute() //分钟
	t6 := time.Now().Second() //秒
	// t7 := time.Now().Nanosecond() //纳秒

	return fmt.Sprintf("%v-%02d-%02d %02d:%02d:%02d",
		t1,
		t2,
		t3,
		t4,
		t5,
		t6)
}
