package helper

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

/**
 * 字符串拼接
 */
func Append(str1 string, strs ...string) string {
	buf := bytes.NewBufferString(str1)
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf.String()
}

func Write2Buffer(buf *bytes.Buffer, strs ...string) *bytes.Buffer {
	for _, str := range strs {
		buf.WriteString(str)
	}
	return buf
}

/**
 * 构建url
 */
func BuildUrl(url string, params map[string]string) string {
	param := "?"
	if strings.Index(url, "?") >= 0 {
		param = "&"
	}

	buf := bytes.NewBufferString("")
	Write2Buffer(buf, strings.TrimSuffix(url, "&"), param)

	for key, value := range params {
		Write2Buffer(buf, key, "=", value, "&")
	}

	return strings.TrimSuffix(buf.String(), "&")
}

/**
 * 过滤emoji
 */
func FilterEmoji(content string) string {
	newContent := ""
	buf := bytes.NewBufferString(newContent)

	for _, value := range content {
		if _, size := utf8.DecodeRuneInString(string(value)); size <= 3 {
			buf.WriteRune(value)
		}
	}
	return buf.String()
}

/**
 * 是否含有emoji
 */
func HasEmoji(content string) bool {
	for _, value := range content {
		if _, size := utf8.DecodeRuneInString(string(value)); size > 3 {
			return true
		}
	}
	return false
}

func LowerFirst(str string) []byte {
	buf := bytes.NewBufferString(strings.ToLower(str[0:1]))
	buf.WriteString(str[1:])
	return buf.Bytes()
}

func UpperFirst(str string) []byte {
	buf := bytes.NewBufferString(strings.ToUpper(str[0:1]))
	buf.WriteString(str[1:])
	return buf.Bytes()
}

func BigCase2Line(str string, lineStr string) string {
	line := bytes.NewBufferString("")
	for _, rn := range str {
		if rn >= 'A' && rn <= 'Z' {
			line.WriteString(lineStr)
			rn += 32
		}
		line.WriteRune(rn)
	}

	return strings.TrimPrefix(line.String(), lineStr)
}

func Line2BigCase(str string, lineStr string) string {
	caseStr := bytes.NewBufferString("")
	tmpStrs := strings.Split(str, lineStr)
	for _, strItem := range tmpStrs {
		if len(strItem) > 0 {
			rn := strItem[0]
			if rn >= 'a' && rn <= 'z' {
				rn -= 32
			}
			Write2Buffer(caseStr, string(rn), strItem[1:])
		}
	}

	return caseStr.String()
}
