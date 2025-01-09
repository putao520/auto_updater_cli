package handle

import (
	"cli/helpers"
	"log"
)

func DoInit() {
	// 检测 pem
	if !helpers.TryGeneratePem() {
		log.Fatal("Failed to generate dsa_pub.pem file")
		return
	}

	// 检测 rc
	if !helpers.CheckDSAPEM() {
		helpers.WriteDSAPEM()
		log.Println("DSAPub DSAPEM ../../dsa_pub.pem written to windows/runner/Runner.rc")
	} else {
		log.Println("DSAPub DSAPEM ../../dsa_pub.pem already exists in windows/runner/Runner.rc")
	}
}
