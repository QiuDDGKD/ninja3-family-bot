package tools

import "strings"

// 获取指令字符串
func GetSplits(org string) []string {
	result := []string{}
	sb := strings.Builder{}

	for _, char := range org {
		if char == ' ' {
			if sb.Len() > 0 {
				result = append(result, sb.String())
				sb.Reset()
			}
		} else {
			sb.WriteRune(char)
		}
	}

	if sb.Len() > 0 {
		result = append(result, sb.String())
	}

	return result
}
