package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/gin-gonic/gin"
)

var (
	
)


func UserHandler() gin.HandlerFunc {

	user_service_path _ := url.Parse("http://localhost:8080")
	proxy := httputil.NewSingleHostReverseProxy(user_service_path)

	return func(c *gin.Context){
		
		
	}
}