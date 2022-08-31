package framework

type IGroup interface {
	Get(string,ControllerHandler)
	Post(string,ControllerHandler)
	Put(string,ControllerHandler)
	Delete(string,ControllerHandler)
}

type Group struct {
	core *Core
	prefix string
}

func NewGroup(core *Core,prefix string) *Group {
	return &Group{
		core: core,
		prefix: prefix,
	}
}

func (g *Group)Get(url string,handler ControllerHandler)  {
	url = g.prefix + url
	g.core.Get(url,handler)
}

func (g *Group)Post(url string,handler ControllerHandler)  {
	url = g.prefix + url
	g.core.Get(url,handler)
}

func (g *Group)Put(url string,handler ControllerHandler)  {
	url = g.prefix + url
	g.core.Get(url,handler)
}

func (g *Group)Delete(url string,handler ControllerHandler)  {
	url = g.prefix + url
	g.core.Get(url,handler)
}

