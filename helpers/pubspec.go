package helpers

import (
	"github.com/putao520/yaml"
	"log"
	"os"
)

func GetPubSpecYaml() (map[string]interface{}, error) {
	data, err := os.ReadFile("pubspec.yaml")
	if err != nil {
		log.Fatalf("error reading file: %v", err)
		return nil, err
	}

	// 解析 YAML 数据
	var result map[string]interface{}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("error unmarshalling YAML: %v", err)
		return nil, err
	}
	return result, nil
}
