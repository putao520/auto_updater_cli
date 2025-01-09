package handle

import (
	"cli/helpers"
	"fmt"
)

func DoDistributor(releaseName string, jobName string) {
	_, err := helpers.Exec("flutter_distributor", "release", "--name", releaseName, "--jobs", jobName)
	if err != nil {
		fmt.Errorf("error running flutter_distributor: %w", err)
		return
	}
}
