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
}

func main() {
	r := gin.Default()
	r.GET("getToken", GetTokenHandler())
	for _, name := range config.Names {
		r.GET("choose_"+name.NickName, chooseHandler(SelectedStuMap{}, name))
	}
	err := r.Run(config.Addr)
	if err != nil {
		log.Println("start server fail!")
		return
	}
}
