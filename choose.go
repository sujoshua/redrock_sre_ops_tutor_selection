package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

type SelectedStuMap struct {
	ok      bool
	result  string
	channel chan consumer
}

type consumer struct {
	userInfo UserInfo
	out      chan string
}

var selectedStuMap SelectedStuMap

func chooseHandler(c *gin.Context) {
	if selectedStuMap.ok {
		c.String(200, selectedStuMap.result)
		return
	}
	token := c.GetHeader("token")
	if token == "" {
		token = c.GetHeader("Token")
		if token == "" {
			log.Printf("未提供token访问接口！")
			c.String(401, "无法获取到token字段，请在未在访问头部添加token字段！")
			return
		}
	}
	userInfo, err := parseToken(token)
	if err != nil {
		c.String(400, err.Error())
		return
	}
	out := make(chan string)
	selectedStuMap.channel <- consumer{userInfo, out}
	c.String(200, <-out)
}

func choose() {
	selectedStuMap.channel = make(chan consumer, config.MaxNum)
	stuMap := make(map[int]bool)
	stuList := make([]UserInfo, 0, config.MaxNum)
	for !selectedStuMap.ok {
		c := <-selectedStuMap.channel
		if _, ok := stuMap[c.userInfo.ID]; !ok {
			stuList = append(stuList, c.userInfo)
			stuMap[c.userInfo.ID] = true
			c.out <- fmt.Sprintf("Hello %s！欢迎你成为第%d/%d位学员，在接下来的日子里我们一起进步呀！", c.userInfo.Name, len(stuMap), config.MaxNum)
			close(c.out)
			if len(stuMap) >= config.MaxNum {
				break
			}
		} else {
			c.out <- "你已经成功加入啦！不用重复提交，如果数据有误想更改，请联系我！"
			close(c.out)
		}

	}
	names := ""
	for _, u := range stuList {
		names += "\n" + u.Name
	}
	selectedStuMap.result = fmt.Sprintf("在%ds后，目前确定以下同学成功加入我的学员：%s\n那我只能遗憾的宣布满员了，快去看看别的导师吧！", time.Now().Unix()-config.startTime, names)
	selectedStuMap.ok = true
	for len(selectedStuMap.channel) > 0 {
		c := <-selectedStuMap.channel
		c.out <- selectedStuMap.result
		close(c.out)
	}

	result := ""
	for i := range stuList {
		result += fmt.Sprintf("%d %s %d %s\n", i+1, stuList[i].Name, stuList[i].ID, stuList[i].QQ)
	}
	println(result)
	err := os.WriteFile("result.txt", []byte(result), 0777)
	if err != nil {
		log.Printf("write result.txt failed!")
		return
	}
}
