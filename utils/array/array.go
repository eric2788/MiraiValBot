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

// string

func IndexOfString(arr []string, element string) int {
	for i, v := range arr {
		if v == element {
			return i
		}
	}
	return -1
}

func RemoveString(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}
