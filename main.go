package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func main() {
	defer timeCost()()

	// work path
	dirPath := "E:\\java_workspace\\mongodb-practice\\"

	outputPath := "mvn-lib"

	mvnCommandResult, err := doExportLib(dirPath)
	if err != nil {
		fmt.Printf("出现异常:%s", err)
		return
	}
	fmt.Printf("执行结果 %s. \nlib=>path%slib\n", mvnCommandResult, dirPath)
	//读取目录下的lib文件 到mvn仓库进行搜索

	fileNames, err := getAllSearchFileName(dirPath + "lib")
	if err != nil {
		fmt.Printf("出现异常:%s", err)
		return
	}
	mvnPath, err := getMvnLocalRepositoryPath()
	fmt.Printf("获取到mvn Local地址 =>%s\n", mvnPath)
	dirs, err := searchFileDir(mvnPath, fileNames)
	if err != nil {
		fmt.Printf("出现异常:%s", err)
		return
	}
	for _, v := range dirs {
		directory := getParentDirectory(v)
		target := dirPath + outputPath + strings.TrimPrefix(directory, mvnPath)
		fmt.Printf("%s => %s \n", directory, target)
		copyDir(directory, target)
	}
}

func searchFileDir(dir string, names []string) ([]string, error) {

	var result []string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%v", err)
			}
			if !info.IsDir() && contains(names, info.Name()) {
				result = append(result, path)
			}
			return nil
		})
	return result, err
}

func contains(names []string, name string) bool {
	if strings.Contains(name,"parent") {
		return true
	}
	for _, v := range names {
		if v == name {
			return true
		}
	}
	return false
}

func getAllSearchFileName(dirPath string) ([]string, error) {
	rd, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return nil, err
	}
	var result []string
	for _, fi := range rd {
		if !fi.IsDir() {
			result = append(result, fi.Name())
		}
	}
	return result, nil
}

// 获取所有异常文件  *.lastUpdated 文件 和 _remote.repositories
func getAllLastUpdateFile(dir string) ([]string, error) {
	var result []string
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%v", err)
			}
			if !info.IsDir() && (strings.Contains(path, "lastUpdated") || strings.Contains(path, "_remote.repositories")) {
				result = append(result, path)
			}
			return nil
		})
	return result, err
}

func doExportLib(dir string) (string, error) {
	//根据命令 mvn
	//获取本地仓库地址
	//清除无效文件
	fmt.Printf("mvn dependency:copy-dependencies -DoutputDirectory=lib\n")
	cmd := exec.Command("mvn", "dependency:copy-dependencies", "-DoutputDirectory=lib")
	cmd.Dir = dir
	stdout, err := cmd.StdoutPipe()

	//获取输出对象，可以从该对象中读取输出结果
	if err != nil {
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()

	// 运行命令
	if err := cmd.Start(); err != nil {
		return "", err
	}
	// 读取输出结果
	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		return "", err
	} else {
		return string(opBytes), err
	}

}

func getMvnLocalRepositoryPath() (string, error) {
	//根据命令 mvn
	//获取本地仓库地址
	//清除无效文件
	fmt.Printf("mvn help:evaluate -Dexpression=settings.localRepository | grep -v '\\[INFO\\]'\n")
	cmd := exec.Command("mvn", "help:evaluate", "-Dexpression=settings.localRepository")

	stdout, err := cmd.StdoutPipe()

	//获取输出对象，可以从该对象中读取输出结果
	if err != nil {
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()

	// 运行命令
	if err := cmd.Start(); err != nil {
		return "", err
	}
	// 读取输出结果
	if opBytes, err := ioutil.ReadAll(stdout); err != nil {
		return "", err
	} else {
		return parserLocalRepositoryPath(string(opBytes))
	}

}

func parserLocalRepositoryPath(content string) (string, error) {
	lineList := strings.Split(content, "\n")
	for i := range lineList {
		if strings.Index(lineList[i], "[INFO]") == -1 {
			result := strings.TrimRight(lineList[i], "\r")
			return result, nil
		}
	}
	return "", errors.New(fmt.Sprintf("没有找到maven Local Repository maven command result: \n%s", content))
}

//@brief：耗时统计函数
func timeCost() func() {
	start := time.Now()
	return func() {
		tc := time.Since(start)
		fmt.Printf("\ntime cost = %v\n", tc)
	}
}

func FormatPath(s string) string {
	switch runtime.GOOS {
	case "windows":
		return strings.Replace(s, "/", "\\", -1)
	case "darwin", "linux":
		return strings.Replace(s, "\\", "/", -1)
	default:
		fmt.Println("only support linux,windows,darwin, but os is " + runtime.GOOS)
		return s
	}
}

func copyDir(src string, dest string) {
	src = FormatPath(src)
	dest = FormatPath(dest)
	log.Println(src)
	log.Println(dest)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("xcopy", src, dest, "/I", "/E")
	case "darwin", "linux":
		cmd = exec.Command("cp", "-R", src, dest)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		fmt.Printf("%s,error => %s\n",err.Error(),stderr.String())
		return
	}
	fmt.Printf("commmand result =>%s\n",out.String())
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "\\"))
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
