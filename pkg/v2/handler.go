package webv2

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}

type handlerFunc func(c *Context)