package slice

import "sort"

/*
 * @author: yeshibo
 * @date: Wednesday, 2022/09/07, 3:49:26 pm
 */

// RemoveRepByMap 通过map主键唯一的特性过滤重复元素 []string去重
func RemoveRepByMap(slc []string) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// RemoveRepByLoop 通过两重循环过滤重复元素 []string去重
func RemoveRepByLoop(slc []string) []string {
	var result []string // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// RemoveRep 元素去重 []string去重
func RemoveRep(slc []string) []string {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slc)
	}
}

// IsValueInList 查询数组中是否包含某个value值
func IsValueInList(value string, list []string) bool {
	sort.Strings(list)
	i := sort.SearchStrings(list, value)
	return i < len(list) && list[i] == value
}
