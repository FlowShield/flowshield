package util

import "os"

func InArray(in string, array []string) bool {
	for k := range array {
		if in == array[k] {
			return true
		}
	}
	return false
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
