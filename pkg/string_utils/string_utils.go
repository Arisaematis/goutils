package string_utils

import (
	"fmt"
	"strconv"
	"unicode"
	"unsafe"
)

/*
 * @author: yeshibo
 * @date: Saturday, 2022/09/03, 2:07:56 pm
 */

// StringToByteSlice string to []byte
func StringToByteSlice(s string) []byte {

	tmp1 := (*[2]uintptr)(unsafe.Pointer(&s))

	tmp2 := [3]uintptr{tmp1[0], tmp1[1], tmp1[1]}

	return *(*[]byte)(unsafe.Pointer(&tmp2))
}

// ByteSliceToString []byte to string
func ByteSliceToString(bytes []byte) string {
	return *(*string)(unsafe.Pointer(&bytes))
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				// need日志记录一下 没有小写字母
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

// IsStartUpper 判断首字母是否大写
func IsStartUpper(s string) bool {
	return unicode.IsUpper([]rune(s)[0])
}

// Lowercase 首字母小写
func Lowercase(str string) (newString string) {
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 65 && vv[i] <= 90 {
				vv[i] += 32 // string的码表相差32位
				newString += string(vv[i])
			} else {
				// need日志记录一下 没有大写字母
				return str
			}
		} else {
			newString += string(vv[i])
		}
	}
	return newString
}

func ReplaceAtPosition(sourceText, replaceText string, start, end int) string {
	r := []rune(sourceText)
	startString := string(r[0:start])
	endString := string(r[end:])
	return startString + replaceText + endString
}

func Decimal(value float64) float64 {
	// 存在进度问题
	// return math.Trunc(value*1e2+0.5) * 1e-2
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}

// AverageFloat64 计算float64的平均值
func AverageFloat64(values []float64) float64 {
	var sum float64
	for _, v := range values {
		sum += v
	}
	return Decimal(sum / float64(len(values)))
}

// MaxFloat64 获取float64的最大值
func MaxFloat64(values []float64) float64 {
	var max float64 = 0
	// 判断是否为空
	if len(values) == 0 {
		return max
	}
	// 判断是否只有一个元素
	if len(values) == 1 {
		return values[0]
	}
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}

// ContainsString 判断 arr string数组中是否存在 s
func ContainsString(arr []string, s string) bool {
	for _, value := range arr {
		if value == s {
			return true
		}
	}
	return false
}
