package util

// IndexOfString finds the index of a string in an array of strings
func IndexOfString(arr []string, value string) (index int) {
	for i, v := range arr {
		if v == value {
			return i
		}
	}
	return -1
}
