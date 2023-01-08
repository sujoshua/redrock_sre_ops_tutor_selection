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

func chooseHandler(selectedStuMap SelectedStuMap, name NameConfig) func(c *gin.Context) {
	selectedStuMap.channel = make(chan consumer, config.MaxNum)
	selectedStuMap.ok = false
	go choose(name, selectedStuMap)
	return func(c *gin.Context) {
		if selectedStuMap.ok {
			c.String(200, selectedStuMap.result)
			return
		}
		token := c.GetHeader("Authorization")
		if token == "" {
			token = c.GetHeader("authorization")
			if token == "" {
				log.Printf("未提供authorization访问接口！")
				c.String(401, "无法获取到Authorization字段，请在未在访问头部添加Authorization字段！")
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
}

func choose(name NameConfig, selectedStuMap SelectedStuMap) {
	stuMap := make(map[int]bool)
	stuList := make([]UserInfo, 0, config.MaxNum)
	for !selectedStuMap.ok {
		c := <-selectedStuMap.channel
		if _, ok := stuMap[c.userInfo.ID]; !ok {
			stuList = append(stuList, c.userInfo)
			stuMap[c.userInfo.ID] = true
			c.out <- fmt.Sprintf("Hello %s！\n"+
				"你已成功选择导师：%s\n"+
				"欢迎你成为第%d/%d位学员\n"+
				"快点击链接添加导师的飞书群吧：%s\n"+
				name.SpecialWords, c.userInfo.Name, name.Name, len(stuMap), config.MaxNum, name.FeishuURL)
			close(c.out)
			if len(stuMap) >= config.MaxNum {
				break
			}
		} else {
			c.out <- "你已经成功加入啦！不用重复提交，如果数据有误想更改，请联系导师！"
			close(c.out)
		}

	}
	names := ""
	for _, u := range stuList {
		names += "\n" + u.Name
	}
	selectedStuMap.result = fmt.Sprintf("在%ds后，目前确定以下同学成功成为%s导师的学员：%s\n那本api只能遗憾的宣布该导师满员了，快去看看别的导师吧！", time.Now().Unix()-config.startTime, name.Name, names)
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
	log.Printf("%s:\n"+result, name.NickName)
	err := os.WriteFile(name.NickName+"result.txt", []byte(result), 0777)
	if err != nil {
		log.Printf("write " + name.NickName + "result.txt failed!")
		return
	}
}
