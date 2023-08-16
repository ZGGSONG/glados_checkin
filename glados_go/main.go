package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	gladosURL = "https://glados.rocks"
)

var log = logrus.New()

type Config struct {
	Cookies []string `yaml:"cookies"`
	BarkURL string   `yaml:"bark_url"`
}

func sendBark(url, title, text string) (string, error) {
	if url == "" {
		return "bark: 未配置，无法进行消息推送.", nil
	}

	log.Info("=================================================================")
	log.Info("Bark: 开始推送消息！")

	uri := fmt.Sprintf("%s/%s/%s", url, title, text)
	resp, err := http.Get(uri)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	message, _, _, err := jsonparser.Get(body, "message")
	if err != nil {
		return "", err
	}

	return string(message), nil
}

func checkin(cookie string) (string, string, error) {
	client := &http.Client{}

	checkinURL := fmt.Sprintf("%s/api/user/checkin", gladosURL)
	statusURL := fmt.Sprintf("%s/api/user/status", gladosURL)

	payload := map[string]interface{}{
		"token": "glados.one",
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	reqCheckin, err := http.NewRequest("POST", checkinURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", "", err
	}

	reqStatus, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return "", "", err
	}

	reqCheckin.Header.Set("cookie", cookie)
	reqCheckin.Header.Set("referer", gladosURL+"/console/checkin")
	reqCheckin.Header.Set("origin", gladosURL)
	reqCheckin.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")
	reqCheckin.Header.Set("content-type", "application/json;charset=UTF-8")

	reqStatus.Header.Set("cookie", cookie)
	reqStatus.Header.Set("referer", gladosURL+"/console/checkin")
	reqStatus.Header.Set("origin", gladosURL)
	reqStatus.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.1 Safari/605.1.15")

	respCheckin, err := client.Do(reqCheckin)
	if err != nil {
		return "", "", err
	}
	defer respCheckin.Body.Close()

	respStatus, err := client.Do(reqStatus)
	if err != nil {
		return "", "", err
	}
	defer respStatus.Body.Close()

	checkinBody, err := ioutil.ReadAll(respCheckin.Body)
	if err != nil {
		return "", "", err
	}

	statusBody, err := ioutil.ReadAll(respStatus.Body)
	if err != nil {
		return "", "", err
	}

	// Parse JSON responses
	checkinMessage, _ := jsonparser.GetString(checkinBody, "message")
	statusLeftDays, _ := jsonparser.GetString(statusBody, "data", "leftDays")
	statusEmail, _ := jsonparser.GetString(statusBody, "data", "email")

	timeNow := time.Now().Format("2006-01-02 15:04:05")

	if checkinCode, _ := jsonparser.GetInt(checkinBody, "code"); checkinCode == -2 {
		return "", "", fmt.Errorf(checkinMessage)
	}

	return checkinMessage, fmt.Sprintf("现在时间是：%s\nemail: %s\ncheckin: %d | state: %d\n%s\n剩余天数：%s天", timeNow, statusEmail, respCheckin.StatusCode, respStatus.StatusCode, checkinMessage, statusLeftDays), nil
}

func main() {
	yamlPath := "config.yml"
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Fatalf("无法读取配置文件：%s", err)
	}

	var conf Config
	if err := yaml.Unmarshal(yamlFile, &conf); err != nil {
		log.Fatalf("无法解析配置文件：%s", err)
	}

	if len(conf.Cookies) == 0 || conf.Cookies[0] == "" {
		log.Fatal("没有配置cookie！")
	}

	for _, cookie := range conf.Cookies {
		title := ""
		text := ""
		var message string

		title, text, err = checkin(cookie)
		if err != nil {
			log.Error("程序出错！")
			title = "程序出错！"
			if len(err.Error()) > 0 {
				text = "网络信息: " + err.Error()
			} else {
				text = "没有获取到网络信息"
			}
		} else {
			log.Info("签到成功！")
		}

		log.Info(text)

		if conf.BarkURL != "" {
			message, err = sendBark(conf.BarkURL, url.QueryEscape(title), url.QueryEscape(text))
			if err != nil {
				log.Errorf("Bark推送出错：%s", err)
			} else {
				log.Info(message)
			}
		}
	}
}
