package main

func Contains(s []string, ele string) bool {
	for _, v := range s {
		if v == ele {
			return true
		}
	}
	return false
}

func indexOf[T comparable](element T, data []T) int {
	for index, value := range data {
		if element == value {
			return index
		}
	}
	return -1
}
