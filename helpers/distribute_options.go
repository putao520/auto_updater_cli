package helpers

import (
	"github.com/putao520/yaml"
	"log"
	"os"
)

func GetDistributeOptions() (map[string]interface{}, error) {
	data, err := os.ReadFile("distribute_options.yaml")
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

func GetReleaseItemInfo(distributeOptions map[string]any, prod string, job string, itemName string) map[string]any {
	releases, _ := distributeOptions["releases"].([]any)
	for _, v_ := range releases {
		v := v_.(map[string]any)
		n := v["name"]
		if n == prod {
			jobs, _ := v["jobs"].([]any)
			for _, j_ := range jobs {
				j := j_.(map[string]any)
				jn, _ := j["name"]
				if jn == job {
					p, _ := j[itemName].(map[string]any)
					return p
				}
			}
		}
	}
	panic("releases:" + prod + ", job:" + job + "#" + itemName + " -> not defined")
}
