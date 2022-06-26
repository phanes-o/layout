package utils

func Throw(err error) {
	if err != nil {
		panic(err)
	}
}
