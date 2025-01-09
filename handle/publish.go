package handle

import (
	"cli/helpers"
	"fmt"
	"log"
)

func DoPublish(platform string, releaseName string, jobName string) {
	// 读取 pubspec.yaml
	pubspec, err := helpers.GetPubSpecYaml()
	if err != nil {
		log.Println("pubspec.yaml was error!")
	}

	// 读取 distribute_options.yaml
	distributeOptions, err := helpers.GetDistributeOptions()
	if err != nil {
		log.Println("distribute_options.yaml was error!")
	}
	outputDir, _ := distributeOptions["output"].(string)

	// 执行 sign_update
	packInfo := helpers.GetReleaseItemInfo(distributeOptions, releaseName, jobName, "package")
	version, _ := pubspec["version"].(string)
	name, _ := pubspec["name"].(string)
	target, _ := packInfo["target"].(string)
	signatureFile := fmt.Sprintf("%s%s/%s-%s-%s.%s", outputDir, version, name, version, platform, target)
	var args []string
	args = append(args, signatureFile)
	signUpdateResult, err := signUpdate(args)
	if err != nil {
		log.Fatal(err)
	}
	dsaSignature := signUpdateResult.Signature

	// 生成 appcast.xml
	appcastPath := helpers.GenerateAppCast(outputDir, dsaSignature, pubspec, packInfo)

	// 上传到 qiniu
	err = helpers.PutWithRewrite(distributeOptions, pubspec, releaseName, jobName, appcastPath)
	if err != nil {
		log.Fatal(err)
	}
}

func DoAuto(platform string, releaseName string, jobName string) {
	// 检查
	DoInit()

	// 执行 build
	DoBuild(platform)

	// 执行 distributor
	DoDistributor(releaseName, jobName)

	// 发布
	DoPublish(platform, releaseName, jobName)
}
