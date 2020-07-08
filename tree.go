package mygin

func NewTree() *tree {
	return &tree{
		root: &routerGroup{},
	}
}

type tree struct {
	root *routerGroup
}

func (t *tree) match(path, method string) HandlerFunc {
	nodes := getNodeFromPath(path)
	g := t.root
	var handler HandlerFunc
	for i, node := range nodes {
		temp := g.getRouterGroup(node)
		if temp == nil {
			handler = g.GetHandler(method)
			break
		}
		g = temp
		if i == len(nodes)-1 {
			handler = g.GetHandler(method)
		}
	}
	if handler == nil {
		return bodyWriter(default404Body)
	}
	return handler
}
