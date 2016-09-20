package interfaces

type Model interface {
}

type Entity interface {
	Build()
}

type Response interface {
	Data()
	Status()
}
