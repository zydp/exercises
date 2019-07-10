package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	postElectricArcEvent = "http://ip:port/jhsq/api/Server/postElectricArcEvent.shtml"		//电气灭弧预警事件
)

func main() {

	URL, err := url.Parse(postElectricArcEvent)
	if err != nil {
		fmt.Println("Error in parsing:", err)
		return
	}

	c := &http.Client{
		Timeout: 15 * time.Second,
	}
	data := make(url.Values)
	var state map[string]interface{} = make(map[string]interface{})
	var jsonList []interface{}

	data.Set("token","eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9")
	data.Set("reportDate", time.Now().Format("2006-01-02 15:04:05"))

	//postElectricArcEvent
	state["id"] ="111111"
	state["communityNum"] = "0001"
	state["state"] = "true"
	state["message"] = "测试电气灭弧预警事件！"
	state["dateTime"] = data.Get("reportDate")
	jsonList = append(jsonList, state)

	jsonObj, _ := json.Marshal(jsonList)
	data.Set("jsonObj",string(jsonObj))

	fmt.Printf("%s\n", jsonObj)

	request, err := http.NewRequest("POST", URL.String(), strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("Get:", err)
		return
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpData, err := c.Do(request)
	if err != nil {
		fmt.Println("Error in Do():", err)
		return
	}

	fmt.Println("Status code:", httpData.Status)

	length := 0
	var buffer [1024]byte
	r := httpData.Body
	for {
		n, err := r.Read(buffer[0:])
		if err ==io.EOF {
			break
		}else if err!=nil{
			fmt.Println(err)
			break
		}
		length = length + n
	}
	fmt.Println("Calculated response data length:", length)
	fmt.Println("-----------------------------")
	fmt.Printf("[%s]\n", buffer)
	fmt.Println("+++++++++++++++++++++++++++++")
}
