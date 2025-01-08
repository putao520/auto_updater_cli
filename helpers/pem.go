package helpers

import (
	"log"
	"os"
	"time"
)

// 判断当前目录是否存在 dsa_pub.pem 文件
func checkPubPem() bool {
	_, err := os.Stat("dsa_pub.pem")
	return !os.IsNotExist(err)
}

// 运行 dsa 命令 dart run auto_updater:generate_keys
func generatePem() {
	outputs, err := Exec("dart", "run", "auto_updater:generate_keys")
	if err != nil {
		log.Fatalf("Error running dart run auto_updater:generate_keys: %v", err)
	}
	for _, l := range outputs {
		log.Println(l)
	}
}

func TryGeneratePem() bool {
	for i := 0; i < 3; i++ {
		if checkPubPem() {
			return true
		}
		generatePem()
		// sleep 1000
		time.Sleep(1000 * time.Millisecond)
	}
	return false
}
