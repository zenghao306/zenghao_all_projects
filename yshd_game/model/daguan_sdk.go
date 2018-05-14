package model

import (
	"encoding/json"

	"github.com/yshd_game/common"

	"net/http"
	"net/url"

	//"github.com/liudng/godump"
	"io/ioutil"
	"strings"
	//"strconv"
	"fmt"
	"strconv"
)

type Err struct {
	Errors string `json:"errors"`
}

type DaGuanRes struct {
	Status   string `json:"status"`
	Errors   string `json:"errors"`
	ReuestId string `json:"request_id"`
}

func ReportAnchorDate(cmd string, rid string, uid int, score int, title string, item_modify_time int64) {
	content := make([]map[string]interface{}, 0)
	data_list := make(map[string]interface{})
	data_list["cmd"] = cmd

	fields := make(map[string]interface{})
	fields["itemid"] = strconv.Itoa(uid)
	fields["score"] = score
	fields["title"] = title
	fields["item_tags"] = rid
	fields["item_modify_time"] = item_modify_time

	data_list["fields"] = fields

	content = append(content, data_list)

	b, err := json.Marshal(content)
	if err != nil {
		return
	}
	v := url.Values{}
	v.Set("appid", "5215044")
	v.Set("table_name", "item")
	v.Set("table_content", string(b))

	bytes_req := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", "http://datareportapi.datagrand.com/data/17wan", bytes_req)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;  param=value ; charset=utf-8") //

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		common.Log.Err("client error")
		return
	}

	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		common.Log.Err("read weixin resp err, is:%s", err2.Error())
		return
	}

	var s DaGuanRes
	json.Unmarshal(body, &s)

	if s.Status == "ok" {
		return
	} else {
		common.Log.Err(string(body))
	}
}

func ReportActionDate(anchorid int, action_type string, uid int) {

	content := make([]map[string]interface{}, 0)
	data_list := make(map[string]interface{})
	data_list["cmd"] = "add"

	fields := make(map[string]interface{})
	fields["userid"] = uid
	fields["itemid"] = anchorid
	fields["action_type"] = action_type

	data_list["fields"] = fields

	content = append(content, data_list)

	b, err := json.Marshal(content)
	if err != nil {
		return
	}
	v := url.Values{}
	v.Set("appid", "5215044")
	v.Set("table_name", "user_action")
	v.Set("table_content", string(b))

	bytes_req := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("POST", "http://datareportapi.datagrand.com/data/17wan", bytes_req)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;  param=value ; charset=utf-8") //

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		common.Log.Err("client error")
		return
	}

	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		common.Log.Err("read weixin resp err, is:%s", err2.Error())
		return
	}

	var s DaGuanRes
	json.Unmarshal(body, &s)

	if s.Status == "ok" {
		return
	} else {
		common.Log.Err(string(body))
	}

}

type Recdata struct {
	Itemid string `json:"itemid"`
}
type DaGuanRecommandRes struct {
	Status   string    `json:"status"`
	Recdata  []Recdata `json:"recdata"`
	ReuestId string    `json:"request_id"`
	Errors   string    `json:"errors"`
}

func GetDaGuanRecommand(uid int) []string {
	/*
		v := url.Values{}
		v.Set("appid", "5215044")
		v.Set("cnt",64)
		v.Set("userid", strconv.Itoa(3))
	*/
	ret := make([]string, 0)
	url := fmt.Sprintf("http://recapi.datagrand.com/personal/17wan?appid=%s&cnt=64&userid=%d", "5215044", uid)

	//bytes_req := ioutil.NopCloser(strings.NewReader(v.Encode()))
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", url, nil)

	reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded;  param=value ; charset=utf-8") //

	response, err := client.Do(reqest)
	defer response.Body.Close()
	if err != nil {
		common.Log.Err("client error")
		return nil
	}

	body, err2 := ioutil.ReadAll(response.Body)
	if err2 != nil {
		common.Log.Err("read weixin resp err, is:%s", err2.Error())
		return nil
	}

	var s DaGuanRecommandRes
	s.Recdata = make([]Recdata, 0)
	err = json.Unmarshal(body, &s)
	if err != nil {
		common.Log.Errf("json err err is %s", err.Error())
		return nil
	}

	if s.Status == "OK" {
		for _, v := range s.Recdata {
			ret = append(ret, v.Itemid)
		}
		return ret

	} else {
		common.Log.Err(string(body))
	}
	return nil
}
