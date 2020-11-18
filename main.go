package main

import (
	"github.com/bobbaicloudwithpants/simple-web-application/controllers"
	"github.com/codegangsta/negroni"
)

func main(){
	n := negroni.Classic()
	n.UseHandler(controllers.InitHandlers())
	n.Run(":3000")
}