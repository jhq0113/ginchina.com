package helper

import (
	"io"
	"os"
	"strings"
	"syscall"
)

func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func DirName(fileName string) string {
	end := strings.LastIndex(fileName, "/")
	if end > -1 {
		return fileName[0:end]
	}
	return fileName
}

func IsFile(path string) bool {
	return !IsDir(path)
}

/**
 *实现fopen
 * openType支持解析：
 * "r"	只读方式打开，将文件指针指向文件头。
 * "r+"	读写方式打开，将文件指针指向文件头。
 * "w"	写入方式打开，将文件指针指向文件头并将文件大小截为零。如果文件不存在则尝试创建之。
 * "w+"	读写方式打开，将文件指针指向文件头并将文件大小截为零。如果文件不存在则尝试创建之。
 * "a"	写入方式打开，将文件指针指向文件末尾。如果文件不存在则尝试创建之。
 * "a+"	读写方式打开，将文件指针指向文件末尾。如果文件不存在则尝试创建之。
 */
func FOpen(fileName string, openType string, fileMode os.FileMode) (*os.File, error) {
	switch openType {
	case "r":
		return os.OpenFile(fileName, syscall.O_RDONLY, fileMode)
	case "r+":
		return os.OpenFile(fileName, syscall.O_RDWR, fileMode)
	case "w":
		return os.OpenFile(fileName, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, fileMode)
	case "w+":
		return os.OpenFile(fileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_TRUNC, fileMode)
	case "a":
		return os.OpenFile(fileName, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_APPEND, fileMode)
	case "a+":
		return os.OpenFile(fileName, syscall.O_RDWR|syscall.O_CREAT|syscall.O_APPEND, fileMode)
	}

	return os.OpenFile(fileName, syscall.O_RDONLY, fileMode)
}

/**
 * 分块读取文件
 */
func ReadAllWithChunk(fileName string) (string, error) {
	f, err := FOpen(fileName, "r", 0)
	if err != nil {
		return "", err
	}

	defer f.Close()

	chunk := make([]byte, 0)
	buf := make([]byte, 1024)

	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 {
			break
		}

		chunk = append(chunk, buf[:n]...)
	}

	return string(chunk), nil
}

func Unlink(fileName string) error {
	if FileExists(fileName) {
		return os.Remove(fileName)
	}

	return nil
}
