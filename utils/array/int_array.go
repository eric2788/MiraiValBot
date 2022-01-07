package array

// int

func IndexOfInt(arr []int, element int) int {
	for i, v := range arr {
		if v == element {
			return i
		}
	}
	return -1
}

func RemoveInt(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

// AddInt anti duplicate adding
func AddInt(s []int, element int) []int {
	if IndexOfInt(s, element) != -1 {
		return s
	} else {
		return append(s, element)
	}
}

func FilterInt(s []int, filter func(int) bool) []int {
	arr := make([]int, 0)
	for _, v := range s {
		if filter(v) {
			arr = append(arr, v)
		}
	}
	return arr
}
