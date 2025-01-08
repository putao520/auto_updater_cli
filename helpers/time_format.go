package helpers

import "time"

func GeneratePubDate() string {
	// 获取当前时间
	currentTime := time.Now()

	// 设置时区为东八区（+0800）
	location, _ := time.LoadLocation("Asia/Shanghai")
	currentTime = currentTime.In(location)

	// 格式化时间
	return currentTime.Format("Mon, 02 Jan 2006 15:04:05 -0700")
}
