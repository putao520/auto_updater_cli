package handle

import "fmt"

var GlobalVer = "0.0.1"

func ArgVersion() {
	fmt.Printf("auto updater cli version:%s", GlobalVer)
}
