package helpers

import (
	"encoding/xml"
	"fmt"
	"log"
	"os"
)

func templateXml() []byte {
	return []byte(`<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0" xmlns:sparkle="http://www.andymatuschak.org/xml-namespaces/sparkle">
    <channel>
        <title></title>
        <description></description>
        <item>
            <title>Version 1.1.0</title>
            <pubDate>Mon, 08 Jan 2025 10:15:00 +0800</pubDate>
            <enclosure url="ai_app-1.0.0+1-windows.zip"
                       sparkle:dsaSignature="MDwCHH42lx5clQKtXCjlB3aq9IYpWE15FUdWgtaQmFcCHF/TaLLajijwGUyqfkH5BUcM8q88xtrlHyJQodE="
                       sparkle:version="1.0.0+1"
                       sparkle:os="windows"
                       length="0"
                       type="application/octet-stream" />
        </item>
    </channel>
</rss>`)
}

func writeXml(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}

func getAppCast(rootDir string) (map[string]interface{}, bool, error) {
	appCastXml := "./" + rootDir + "/appcast.xml"
	_, err := os.Stat(appCastXml)
	isCreated := false

	if os.IsNotExist(err) {
		// 生成
		err := writeXml(appCastXml, templateXml())
		if err != nil {
			return nil, false, err
		}
		isCreated = true
	}

	data, err := os.ReadFile(appCastXml)
	if err != nil {
		log.Fatal(err)
	}

	// 使用 map 动态解析
	var result map[string]interface{}
	err = xml.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err)
	}

	return result, isCreated, nil
}

func buildAppCastItem(dsaSignature string, spec map[string]any, pack map[string]any) map[string]any {
	var info map[string]any
	version, _ := spec["version"].(string)
	platform, _ := pack["platform"].(string)
	name, _ := spec["name"].(string)
	target, _ := pack["target"].(string)
	info["title"] = "Version " + version
	info["pubDate"] = GeneratePubDate()
	url := fmt.Sprintf("%s-%s-%s.%s", name, version, platform, target)
	var enclosure map[string]any
	enclosure["url"] = url
	enclosure["sparkle:dsaSignature"] = dsaSignature
	enclosure["sparkle:version"] = version
	enclosure["sparkle:os"] = platform
	enclosure["length"] = 0
	enclosure["type"] = "application/octet-stream"
	info["enclosure"] = enclosure
	return info
}

func GenerateAppCast(rootDir string, dsaSignature string, spec map[string]interface{}, pack map[string]interface{}) {
	result, isCreated, err := getAppCast(rootDir)
	if err != nil {
		log.Fatal(err)
		return
	}
	// 初始化
	channel, _ := result["rss"].(map[string]interface{})["channel"].(map[string]interface{})
	if isCreated {
		channel["title"] = spec["name"]
		channel["description"] = spec["description"]
	}

	// 构造 item
	var items []map[string]interface{}
	item := buildAppCastItem(dsaSignature, spec, pack)
	channel["item"] = append(items, item)

	// 输出 xml
	xmlData, err := xml.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling to XML: %v", err)
		return
	}

	appCastXml := "./" + rootDir + "/appcast.xml"
	err = writeXml(appCastXml, xmlData)
	if err != nil {
		log.Fatalf("Error marshaling to XML: %v", err)
		return
	}
}
