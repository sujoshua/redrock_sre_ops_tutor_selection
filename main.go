package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

var config Config

func init() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	var err error
	config, err = GetConfig()
	if err != nil {
		os.Exit(1)
	}
	go choose()
}

func main() {
	r := gin.Default()
	r.GET("getToken", GetTokenHandler)
	r.GET("choose", chooseHandler)
	err := r.Run(config.Addr)
	if err != nil {
		log.Println("start server fail!")
		return
	}
}
