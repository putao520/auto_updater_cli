package handle

import (
	"cli/helpers"
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type SignUpdateResult struct {
	Signature string
	Length    int
}

func signUpdate(args []string) (SignUpdateResult, error) {
	var executable string
	if runtime.GOOS == "darwin" { // macOS
		executable = filepath.Join(".", "macos", "Pods", "Sparkle", "bin", "sign_update")
	} else if runtime.GOOS == "windows" { // Windows
		executable = filepath.Join(
			".", "windows", "flutter", "ephemeral", ".plugin_symlinks", "auto_updater_windows",
			"windows", "WinSparkle-0.8.1", "bin", "sign_update.bat")
	}

	// Ensure the arguments are passed correctly
	arguments := append([]string(nil), args...) // copy the args
	if runtime.GOOS == "windows" && len(arguments) == 1 {
		arguments = append(arguments, filepath.Join(".", "dsa_priv.pem"))
	}

	// Run the command
	outputs, err := helpers.Exec(executable, arguments...)
	if err != nil {
		return SignUpdateResult{}, fmt.Errorf("error running sign update: %w", err)
	}

	// Process output
	signUpdateOutput := strings.Join(outputs, "")
	if runtime.GOOS == "windows" {
		signUpdateOutput = strings.TrimSpace(strings.ReplaceAll(signUpdateOutput, "\r\n", ""))
		signUpdateOutput = strings.TrimSpace(strings.ReplaceAll(signUpdateOutput, "\r", ""))
		signUpdateOutput = strings.TrimSpace(strings.ReplaceAll(signUpdateOutput, "\n", ""))
		signUpdateOutput = fmt.Sprintf(`sparkle:dsaSignature="%s" length="0"`, signUpdateOutput)
	}

	// Print the result
	fmt.Println(signUpdateOutput)

	// Regular expression to capture the signature and length
	re := regexp.MustCompile(`sparkle:(dsa|ed)Signature="([^"]+)" length="(\d+)"`)
	matches := re.FindStringSubmatch(signUpdateOutput)
	if matches == nil {
		return SignUpdateResult{}, fmt.Errorf("failed to sign update")
	}

	// Parse the signature and length
	signature := matches[2]
	length, err := strconv.Atoi(matches[3])
	if err != nil {
		return SignUpdateResult{}, fmt.Errorf("invalid length value")
	}

	return SignUpdateResult{
		Signature: signature,
		Length:    length,
	}, nil
}

func DoSignUpdate(args []string) {
	result, err := signUpdate(args)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Output the result (you can use result as needed)
	fmt.Printf("Signature: %s, Length: %d\n", result.Signature, result.Length)
}

func BuildSignArg(platform string, releaseName string, jobName string) []string {
	pubspec, _ := helpers.GetPubSpecYaml()
	distributeOptions, _ := helpers.GetDistributeOptions()
	outputDir, _ := distributeOptions["output"].(string)
	packInfo := helpers.GetReleaseItemInfo(distributeOptions, releaseName, jobName, "package")
	version, _ := pubspec["version"].(string)
	name, _ := pubspec["name"].(string)
	target, _ := packInfo["target"].(string)
	signatureFile := fmt.Sprintf("%s%s/%s-%s-%s.%s", outputDir, version, name, version, platform, target)
	var args []string
	return append(args, signatureFile)
}
