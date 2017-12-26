package lib

import "strconv"

func Utoa(u uint) string { return strconv.Itoa(int(u)) }
func Itoa(i int) string  { return strconv.Itoa(i) }

func Atou(s string) uint {
	if i, err := strconv.ParseUint(s, 10, 0); err == nil {
		return uint(i)
	}
	return 0
}

func Atoi(s string) int {
	if i, err := strconv.ParseUint(s, 10, 0); err == nil {
		return int(i)
	}
	return 0
}
