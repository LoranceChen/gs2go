package module

type ActionName struct {
	Name string
}

type Router interface {
	WsPbActionHandler(action ActionName) error
	RouterName() RouterName
}
