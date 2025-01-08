package helpers

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

// 检测是否存在DSAPEM
func CheckDSAPEM() bool {
	file, err := os.Open("./windows/runner/Runner.rc")
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "//") {
			continue
		}
		if strings.Contains(line, `DSAPub      DSAPEM      "../../dsa_pub.pem"`) {

		}
		arr := strings.Fields(line)
		if len(arr) == 3 {
			if arr[0] == "DSAPub" && arr[1] == "DSAPEM" && arr[2] == `"../../dsa_pub.pem"` {
				return true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return false
}

// 写入DSAPEM
func WriteDSAPEM() {
	file, err := os.OpenFile("./windows/runner/Runner.rc", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
		// lines.append(line)
		if strings.HasPrefix(line, "// Version") {
			lines = append(lines, `DSAPub      DSAPEM      "../../dsa_pub.pem"`)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	file.Seek(0, 0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Error writing to file: %v", err)
		}
	}
	writer.Flush()
}

func TryAppendRC() bool {
	for i := 0; i < 3; i++ {
		if CheckDSAPEM() {
			return true
		}
		WriteDSAPEM()
		time.Sleep(1000 * time.Millisecond)
	}
	return false
}
