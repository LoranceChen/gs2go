package service

type RouterName int16

type Services struct {
	Value map[RouterName]Router
}

func (r *Services) RegisterRoute(router Router) {
	r.Value[router.RouterName()] = router
}

func NewRoutes() *Services {
	var routes = &Services{Value: make(map[RouterName]Router, 128)}
	return routes
}
