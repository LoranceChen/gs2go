package module

type RouterName struct {
	Name int16
}

type Routes struct {
	Value map[RouterName]Router
}

func (r *Routes) RegisterRoute(router Router) {
	r.Value[router.RouterName()] = router
}

func NewRoutes() *Routes {
	var routes = &Routes{Value: make(map[RouterName]Router, 128)}
	return routes
}
