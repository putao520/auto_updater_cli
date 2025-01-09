package main

import (
	"cli/handle"
	"log"
	"os"
	"runtime"
	"strings"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	if !(runtime.GOOS == "darwin" || runtime.GOOS == "windows") {
		log.Fatal("auto_updater:sign_update is not supported on this platform")
	}
	ArgFilter(os.Args)
}

func ArgFilter(args []string) {
	if len(args) == 1 {
		return
	}
	cmd := args[1]
	argsImpl := os.Args[2:]
	switch strings.ToLower(cmd) {
	case "check", "c": // 检查|初始化环境
		handle.DoInit()
	case "version", "v":
		handle.ArgVersion()
	case "sign", "s": // 签名|生成签名信息
		if len(argsImpl) != 3 {
			panic("sign need platform releaseName jobName")
		}
		handle.DoSignUpdate(handle.BuildSignArg(argsImpl[0], argsImpl[1], argsImpl[2]))
	case "build", "b": // 构建|生成包
		if len(argsImpl) != 1 {
			panic("build need platform")
		}
		handle.DoBuild(argsImpl[0])
	case "distributor", "d": // 发布|上传存储
		if len(argsImpl) != 2 {
			panic("distributor need releaseName jobName")
		}
		handle.DoDistributor(argsImpl[0], argsImpl[1])
	case "publish", "p": // 生成 appcast.xml
		if len(argsImpl) != 3 {
			panic("publish need platform releaseName jobName")
		}
		handle.DoPublish(argsImpl[0], argsImpl[1], argsImpl[2])
	case "auto", "a": // 全部流程
		if len(argsImpl) != 3 {
			panic("auto need platform releaseName jobName")
		}
		handle.DoAuto(argsImpl[0], argsImpl[1], argsImpl[2])
	default:
		handle.DoUnknown(cmd)
	}
}
