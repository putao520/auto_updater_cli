package helpers

import (
	"fmt"
	"github.com/putao520/gscXml"
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

func getAppCast(rootDir string) (*gscXml.XmlDocument, bool, error) {
	appCastXml := "./" + rootDir + "appcast.xml"
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
	doc := gscXml.NewXmlDocument(string(data))
	return &doc, isCreated, nil
}

func buildAppCastItem(dsaSignature string, spec map[string]any, pack map[string]any) []*gscXml.XmlNode {
	d := make(map[gscXml.XmlName][]*gscXml.XmlNode)

	version, _ := spec["version"].(string)
	platform, _ := pack["platform"].(string)
	fileName := BuildBaseFileName(spec, pack)

	d[gscXml.Name("title", "")] = []*gscXml.XmlNode{
		{
			Val:  "Version " + version,
			Attr: nil,
		},
	}
	d[gscXml.Name("pubDate", "")] = []*gscXml.XmlNode{
		{
			Val:  GeneratePubDate(),
			Attr: nil,
		},
	}
	d[gscXml.Name("enclosure", "")] = []*gscXml.XmlNode{
		{
			Val: "",
			Attr: []gscXml.XmlAttr{
				{
					Name:  gscXml.Name("url", ""),
					Value: fmt.Sprintf("%s/%s", version, fileName),
				},
				{
					Name:  gscXml.Name("dsaSignature", "sparkle"),
					Value: dsaSignature,
				},
				{
					Name:  gscXml.Name("version", "sparkle"),
					Value: version,
				},
				{
					Name:  gscXml.Name("os", "sparkle"),
					Value: platform,
				},
				{
					Name:  gscXml.Name("length", ""),
					Value: "0",
				},
				{
					Name:  gscXml.Name("type", ""),
					Value: "application/octet-stream",
				},
			},
		},
	}

	return []*gscXml.XmlNode{
		{
			Val:  d,
			Attr: nil,
		},
	}
}

func GenerateAppCast(rootDir string, dsaSignature string, spec map[string]interface{}, pack map[string]interface{}) string {
	doc, isCreated, err := getAppCast(rootDir)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	// 初始化
	rss := gscXml.GetNode(doc.Body, gscXml.Name("rss", "")).Val.(map[gscXml.XmlName][]*gscXml.XmlNode)
	channel := gscXml.GetNode(rss, gscXml.Name("channel", "")).Val.(map[gscXml.XmlName][]*gscXml.XmlNode)
	if isCreated {
		gscXml.GetNode(channel, gscXml.Name("title", "")).Val = spec["name"]
		gscXml.GetNode(channel, gscXml.Name("description", "")).Val = spec["description"]
	}
	// 构造 item
	items := buildAppCastItem(dsaSignature, spec, pack)
	channel[gscXml.Name("item", "")] = items
	// 输出 xml
	xmlData := doc.String()
	appCastXml := "./" + rootDir + "/appcast.xml"
	err = writeXml(appCastXml, []byte(xmlData))
	if err != nil {
		log.Fatalf("Error marshaling to XML: %v", err)
		return ""
	}
	return appCastXml
}
