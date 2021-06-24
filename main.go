package main

import (
	"bytes"
	"fmt"
	"github.com/robfig/cron"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


const (
	ServiceSideTestGroup = "https://open.feishu.cn/open-apis/bot/v2/hook/08b50b8b-cbc8-4d5b-9af8-4daf83d9d8d6"
	TestGroup = "https://open.feishu.cn/open-apis/bot/v2/hook/08b50b8b-cbc8-4d5b-9af8-4daf83d9d8d6"
	fixBugSpec = "0 30 10 ? ? 1,3" //cron表达式，周一每天早上10点30
	meetingHostSpec = "0 30 10 ? ? 1" //cron表达式，周三每天早上10点30
	testSpec = "0/3 * * * * ? *" //cron表达式，从0秒开始每3秒执行一次
	fixBugFormat = `{
					"chat_id":"6975356370782552092",
					"msg_type":"interactive",
					"card":{"config":{"wide_screen_mode":true},
					"header":{"title":{"tag":"plain_text","content":"❤️同学，这周该你值班了: "},
					"template": "indigo"},
					"elements":[{"tag":"div","fields":[{"is_short":false,"text":{"tag":"lark_md","content":"<at id=%s></at>\n"}},
					{"is_short":false,"text":{"tag":"lark_md","content":"**请调查本周失败用例，详情请戳👉去OnCall**"}}]},
					{"tag":"action","actions":[{"tag":"button","text":{"tag":"plain_text","content":"去OnCall"},"type":"primary",
					"url":"https://bytedance.feishu.cn/docs/doccnqU5xG0fOuAQLwH33YWI4qd#Zbv06s"}]}]}}`

	meetingHostFormat = `{
						"chat_id":"6975356370782552092",
						"msg_type":"interactive",
						"card":{"config":{"wide_screen_mode":true},
						"header":{"title":{"tag":"plain_text","content":"❤️同学，这周该你值班了: "},
						"template": "indigo"},
						"elements":[{"tag":"div","fields":[{"is_short":false,"text":{"tag":"lark_md","content":"<at id=%s></at>\n"}},
						{"is_short":false,"text":{"tag":"lark_md","content":"**请准备主持周会以及更新上周工作进展，详情请戳👉去OnCall**"}}]},
						{"tag":"action","actions":[{"tag":"button","text":{"tag":"plain_text","content":"去OnCall"},"type":"primary",
						"url":"https://bytedance.feishu.cn/docs/doccnqU5xG0fOuAQLwH33YWI4qd#Zbv06s"}]}]}}`
	MONDAY = "Monday"
	WEDNESDAY = "Wednesday"
)

func OnCallHelper(roster1 []string, targetDay string, hookUrl string, chatFormat string ) {
	var jsonStr []byte
	var roster1Index int32

	fixBugIndex := module(roster1Index, roster1)
	fixBugIndex = updateIndex(fixBugIndex,targetDay)

	jsonStr = []byte( fmt.Sprintf(chatFormat,
							roster1[fixBugIndex]))

	req, err := http.NewRequest("POST", hookUrl, bytes.NewBuffer(jsonStr))

	if err != nil{
		log.Fatal("An error occurred",err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("An error occurred",err)
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	hea := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	fmt.Println(statusCode)
	fmt.Println(hea)
}

func CurrentTime() string {
	t := time.Now()
	weekDay := t.Weekday().String()
	log.Println("今天是",weekDay)
	return weekDay
}


func module(numerator int32, denominator []string) int32 {
	result := numerator % int32(len(denominator))
	return result
}

func updateIndex(index int32, targetDay string) int32{
	weekday := CurrentTime()
	if weekday == targetDay{
		index += 1
	}
	return index
}


func main() {

	fixBugMembers := []string{"6882565187921084418","6927445841556799516","6926702869835939841","6967159379506249730","6972359042111193089","6866971850510008321"}
	//meetingHostMembers := []string{"6882565187921084418","6908529743885516802","6926702869835939841","6927445841556799516","6941182598773227547","6967159379506249730","6972359042111193089","6866971850510008321"}

	c := cron.New() //fixBugMembers
	//fixBugMembers定时任务
	c.AddFunc(testSpec, func() {
		OnCallHelper(fixBugMembers,MONDAY,TestGroup, fixBugFormat)

	})
	c.Start()


	//c1 := cron.New() //meetingHostMembers
	////meetingHostMembers定时任务
	//c1.AddFunc(meetingHostSpec, func() {
	//	OnCallHelper(meetingHostMembers, WEDNESDAY,TestGroup, meetingHostFormat )
	//})
	//
	//c1.Start()






}