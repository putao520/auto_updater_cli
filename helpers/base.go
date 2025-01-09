package helpers

import "fmt"

func BuildBaseFileName(spec map[string]interface{}, pack map[string]interface{}) string {
	version, _ := spec["version"].(string)
	platform, _ := pack["platform"].(string)
	name, _ := spec["name"].(string)
	target, _ := pack["target"].(string)
	return fmt.Sprintf("%s-%s-%s.%s", name, version, platform, target)
}
