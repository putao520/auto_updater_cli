package handle

import (
	"cli/helpers"
	"fmt"
	"log"
)

func getPackInfo(distributeOptions map[string]any, job string, pack string) map[string]any {
	releases, _ := distributeOptions["releases"].([]map[string]any)
	for _, v := range releases {
		n := v["name"]
		if n == job {
			jobs, _ := v["jobs"].([]map[string]any)
			for _, j := range jobs {
				jn, _ := j["name"]
				if jn == pack {
					p, _ := j["package"].(map[string]any)
					return p
				}
			}
		}
	}
	panic("releases:" + job + ", job:" + pack + " -> not defined")
}

func DoPublish(platform string, releaseName string, jobName string) {
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
	outputDir := distributeOptions["output"].(string)

	// 执行 build
	_, err = helpers.Exec("flutter", "build", platform, "--release")
	if err != nil {
		fmt.Errorf("error running flutter build: %w", err)
		return
	}

	// 执行 distributor
	_, err = helpers.Exec("flutter_distributor", "release", "--name", releaseName, "--jobs", jobName)
	if err != nil {
		fmt.Errorf("error running flutter_distributor: %w", err)
		return
	}

	// 执行 sign_update
	packInfo := getPackInfo(distributeOptions, releaseName, jobName)
	version, _ := pubspec["version"].(string)
	name, _ := pubspec["name"].(string)
	target, _ := packInfo["target"].(string)
	signatureFile := fmt.Sprintf("%s/%s/%s-%s-%s.%s", outputDir, version, name, version, platform, target)
	var args []string
	args = append(args, signatureFile)
	signUpdateResult, err := signUpdate(args)
	if err != nil {
		log.Fatal(err)
	}
	dsaSignature := signUpdateResult.Signature

	// 生成 appcast.xml
	helpers.GenerateAppCast(outputDir, dsaSignature, pubspec, packInfo)
}
