package utils

import (
	"io"
	"mime/multipart"
	"strings"
)

func GetFileContentAsString(fileHeader *multipart.FileHeader) (string, error) {
	// 打开文件
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容
	bytes, err := io.ReadAll(file) // 修改这里
	if err != nil {
		return "", err
	}

	// 将字节转换为字符串并返回
	return string(bytes), nil
}

func GetFileNameSuffix(name string) string {
	suffix := strings.Split(name, ".")
	if len(suffix) < 2 {
		return ""
	} else {
		return "." + suffix[1]
	}
}
