package route

import "github.com/lunny/tango"

type Comment struct {
	tango.Ctx
}

func (c *Comment) Get() {
	c.WriteHeader(401)
}

func (c *Comment) Post() {
	c.WriteHeader(401)
}
