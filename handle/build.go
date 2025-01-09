package handle

import (
	"cli/helpers"
	"fmt"
)

func DoBuild(platform string) {
	_, err := helpers.Exec("flutter", "build", platform, "--release")
	if err != nil {
		fmt.Errorf("error running flutter build: %w", err)
		return
	}
}
