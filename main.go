package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 推送token
var pushToken = []string{""}

func main() {
	p := NewXBPush()
	p.getXB()

	fmt.Println("http serve run :10098")
	http.ListenAndServe(":10098", nil)
}

type XBPusher interface {
	getXB()
}

type XBPush struct {
	Cron *gocron.Scheduler
}

func NewXBPush() XBPusher {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	return &XBPush{
		Cron: gocron.NewScheduler(timezone),
	}
}

func (p *XBPush) start() {
	p.Cron.StartAsync()
}

// 定时获取线报
func (p *XBPush) getXB() {
	go p.start()
	p.Cron.Every(15).Second().Do(p.xb)
}

type XBResponse struct {
	Id            int         `json:"id"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	Datetime      string      `json:"datetime"`
	Shorttime     string      `json:"shorttime"`
	Shijianchuo   int64       `json:"shijianchuo"`
	Cateid        string      `json:"cateid"`
	Catename      string      `json:"catename"`
	Comments      int         `json:"comments"`
	Louzhu        string      `json:"louzhu"`
	Louzhuregtime interface{} `json:"louzhuregtime"`
	Url           string      `json:"url"`
	Yuanurl       string      `json:"yuanurl"`
}

func (p *XBPush) xb() {
	var response []XBResponse
	if err := fasthttpGet("http://new.xianbao.fun/plus/json/push.json?230406", &response); err != nil {
		fmt.Println("获取线报错误：", err.Error())
	}
	for _, xbResponse := range response {
		go p.push(xbResponse)
	}
}

// 推送
func (p *XBPush) push(result XBResponse) {
	if ok := p.pushRule(result); !ok {
		return
	}
	domain := "http://www.pushplus.plus/send"
	for _, token := range pushToken {
		pushUrl := fmt.Sprintf("%s?%s", domain, p.pushParam(token, result))
		fasthttpGet(pushUrl, nil)
	}
}

// 推送参数
func (p *XBPush) pushParam(token string, result XBResponse) string {
	xbTime := time.Unix(result.Shijianchuo, 0).Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(`<a href="http://new.xianbao.fun/%s">%s</a>
	<span>时间：%s</span>`, result.Url, result.Title, xbTime)
	val := url.Values{}
	val.Set("token", token)
	val.Set("title", result.Title)
	val.Set("content", content)
	return val.Encode()
}

// 配置推送规则 需要的内容才推送
func (p *XBPush) pushRule(result XBResponse) bool {
	//fmt.Println(result)

	substr := []string{"bug", "猫超", "工行"}

	if !p.containsRule(result.Content, substr) {
		return false
	}

	if !p.pushTimeRule() {
		return false
	}
	return true
}

// 字符串包含
func (p *XBPush) containsRule(s string, substr []string) bool {
	if len(substr) == 0 {
		return false
	}

	for _, v := range substr {
		if strings.Contains(s, v) {
			return true
		}
	}
	return false
}

// 推送时间
func (p *XBPush) pushTimeRule() bool {
	begin := 8
	end := 23

	if time.Now().Hour() <= end && time.Now().Hour() >= begin {
		return true
	}
	return false
}

func fasthttpGet(url string, response interface{}) error {
	status, resp, err := fasthttp.Get(nil, url)

	if err != nil {
		fmt.Println("请求失败:", err.Error())
		return err
	}

	if status != fasthttp.StatusOK {
		fmt.Println("请求没有成功:", status)
		return err
	}

	if response == nil {
		return nil
	}
	if err := json.Unmarshal(resp, &response); err != nil {
		return err
	}
	return nil
}
