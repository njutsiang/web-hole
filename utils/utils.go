package utils

// 在字符串右边填充指定长度的字符
func StrPadRight(s string, p string, length int) string {
	if len(s) >= length {
		return s
	}
	diffLength := length - len(s)
	diffString := ""
	for {
		diffString += p
		if len(diffString) >= diffLength {
			diffString = diffString[0:diffLength]
			break
		}
	}
	return s + diffString
}
