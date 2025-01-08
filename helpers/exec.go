package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

func Exec(executable string, arguments ...string) ([]string, error) {
	var outs []string
	cmd := exec.Command(executable, arguments...)

	// 获取命令的标准输出和标准错误
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("获取标准输出管道失败: %v", err)
		return outs, err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatalf("获取标准错误管道失败: %v", err)
		return outs, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Fatalf("启动命令失败: %v", err)
		return outs, err
	}

	go func() {
		// 将标准输出实时读取并追加到 outs
		buf := make([]byte, 4096)
		for {
			n, err := stdoutPipe.Read(buf)
			if err != nil && err.Error() != "EOF" {
				log.Printf("标准输出读取失败: %v", err)
				break
			}
			if n > 0 {
				outs = append(outs, string(buf[:n]))
				fmt.Print(string(buf[:n]))
			}
			if err == io.EOF {
				break
			}
		}
	}()

	go func() {
		_, err := io.Copy(os.Stderr, stderrPipe)
		if err != nil {
			log.Printf("标准错误打印失败: %v", err)
		}
	}()

	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		log.Fatalf("命令执行失败: %v", err)
		return outs, err
	}

	fmt.Println("命令执行完毕")
	return outs, nil
}
