package array

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

// AddString anti duplicate adding
func AddString(s []string, element string) []string {
	if IndexOfString(s, element) != -1 {
		return s
	} else {
		return append(s, element)
	}
}

func FilterString(s []string, filter func(string) bool) []string {
	arr := make([]string, 0)
	for _, v := range s {
		if filter(v) {
			arr = append(arr, v)
		}
	}
	return arr
}
