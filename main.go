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
	argsImpl := os.Args[1:]
	switch strings.ToLower(cmd) {
	case "version", "v":
		handle.ArgVersion()
	case "publish", "p":
		if len(argsImpl) != 3 {
			panic("publish need platform releaseName jobName")
		}
		handle.DoPublish(argsImpl[0], argsImpl[1], argsImpl[2])
	default:
		handle.DoUnknown(cmd)
	}
}
