package array

// when no generic, that is so painful

// int64

func IndexOfInt64(arr []int64, element int64) int {
	for i, v := range arr {
		if v == element {
			return i
		}
	}
	return -1
}

func RemoveInt64(s []int64, index int) []int64 {
	return append(s[:index], s[index+1:]...)
}

// AddInt64 anti duplicate adding
func AddInt64(s []int64, element int64) []int64 {
	if IndexOfInt64(s, element) != -1 {
		return s
	} else {
		return append(s, element)
	}
}

func FilterInt64(s []int64, filter func(int64) bool) []int64 {
	arr := make([]int64, 0)
	for _, v := range s {
		if filter(v) {
			arr = append(arr, v)
		}
	}
	return arr
}
