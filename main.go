package main

import (
	"code.build.gee/day2-context/gee"
	"net/http"
)

func main(){
	r := gee.New()
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK,"<h1>Hello World</h1>")
	})
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK,"hello %s",c.Query("name"))
	})
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK,gee.H{
			"username":c.PostForm("username"),
			"password":c.PostForm("password"),
		})
	})

	r.Run(":9999")
}
