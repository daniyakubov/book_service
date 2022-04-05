package Action

type Action struct {
	Method string
	Route  string
}

func NewAction(method string, route string) Action {
	return Action{
		Method: method,
		Route:  route,
	}
}
