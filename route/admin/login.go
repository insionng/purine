package admin
import (
    "github.com/tango-contrib/renders"
    "github.com/lunny/tango"
)

type Login struct {
    renders.Renderer
    tango.Ctx
}

func (l *Login) Get(){
    if err := l.Render("admin/login.tmpl");err != nil{
        panic(err)
    }
}