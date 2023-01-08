package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type TokenClaims struct {
	UserInfo
	jwt.StandardClaims
}

type UserInfo struct {
	ID   int
	Name string
	QQ   string
}

// roster 允许获取token的学员
var roster map[string]struct{}

// 选择对应导师对应的接口信息
var urlInfo string

func GetTokenHandler() func(c *gin.Context) {
	roster = make(map[string]struct{})
	// 解析学员名单，之后放入roster的map中
	f, err := os.Open(config.RosterFile)
	if err != nil {
		log.Panicf("parese roster file failed! err: %s", err.Error())
	}
	r := bufio.NewReader(f)

	for {
		line, err := r.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Printf("read roaster string fail, err: %s", err.Error())
		}
		roster[strings.TrimSpace(line)] = struct{}{}
		if err == io.EOF {
			break
		}
	}

	// 生成导师对应的接口
	for _, name := range config.Names {
		urlInfo += fmt.Sprintf("\n%s：%s", name.Name, config.serverURL.JoinPath("choose_"+name.NickName).String())
	}
	return getTokenHandler
}

// 返回token
func getTokenHandler(c *gin.Context) {
	if time.Now().Unix()-config.startTime < 0 {
		c.String(401, "禁止抢跑！")
		return
	}
	values := c.Request.URL.Query()
	id, ok := values["id"]
	if !ok {
		id, ok = values["ID"]
		if !ok {
			log.Println("未输入ID字段")
			c.String(401, "未能在请求中找到id字段，请输入你的学号并检查你的请求哦！")
			return
		}
	}

	// 判断有无资格
	if _, ok := roster[id[0]]; !ok {
		log.Printf("未发现学号：%s, 有考核资格！", id[0])
		c.String(401, "根据你提供的学号信息，未发现你有考核资格哦！如有疑问，请私聊学长！")
		return
	}

	idInt, err := strconv.Atoi(id[0])
	if err != nil {
		log.Printf("输入ID: %s, 非法", id[0])
		c.String(401, "id字段输入非法，确保输入的是你的学号哦！")
		return
	}

	name, ok := values["name"]
	if !ok {
		c.String(401, "未能在请求中找到name字段，请输入你的名字并检查你的请求哦！")
		return
	}
	if name[0] == "" {
		c.String(401, "name字段啥都没有输入，请确保输入了自己的姓名哦！")
		return
	}

	qq, ok := values["qq"]
	if !ok {
		c.String(400, "未能在请求中找到qq字段，请输入你的QQ号并检查你的请求哦！")
		return
	}
	if name[0] == "" {
		c.String(400, "qq字段啥都没有输入，请确保输入了自己的qq哦！")
		return
	}
	token, err := generateToken(UserInfo{idInt, name[0], qq[0]})
	if err != nil {
		log.Printf("generate token error:%s", err.Error())
		c.String(500, "系统好像出现错误啦，不能预期的生成token（恼")
		return
	}

	c.String(200, token+"\n"+urlInfo)
}

// generate token with the info
func generateToken(info UserInfo) (string, error) {
	nowTime := time.Now()
	expiredTime := nowTime.Add(time.Hour)
	claims := TokenClaims{info, jwt.StandardClaims{ExpiresAt: expiredTime.Unix(), IssuedAt: nowTime.Unix(), Issuer: "joshua"}}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(config.tokenKey)
	return token, err
}

// 解析token
func parseToken(tokenStr string) (info UserInfo, err error) {
	tokenClaims, err := jwt.ParseWithClaims(tokenStr, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return config.tokenKey, nil
	})
	if err != nil {
		log.Printf("parse token: %s,error: %s", tokenStr, err.Error())
		return info, errors.New("解析Token失败，请检查传入的token是否合法！")
	}
	if !tokenClaims.Valid {
		return info, errors.New("Token无效，请检查传入的token是否有错误或是超过了token有效时间！")
	}
	claims, ok := tokenClaims.Claims.(*TokenClaims)
	if !ok {
		log.Println("断言token失败")
		return info, errors.New("token解析失败")
	}
	return claims.UserInfo, nil
}
