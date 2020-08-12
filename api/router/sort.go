package router

type sortedRoutes struct {
	routes Routes
}

func (s sortedRoutes) Len() int {
	return len(s.routes.Routes)
}

func (s sortedRoutes) Less(i, j int) bool {
	return s.routes.Routes[i].Priority < s.routes.Routes[j].Priority
}

func (s sortedRoutes) Swap(i, j int) {
	s.routes.Routes[i], s.routes.Routes[j] = s.routes.Routes[j], s.routes.Routes[i]
}
