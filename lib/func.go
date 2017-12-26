package lib

func SafeRunFunc(f func()) {
	if f != nil {
		f()
	}
}
