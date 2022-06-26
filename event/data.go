package event

type Type int

type Data struct {
	Type Type
	Data interface{}
}
