package main

import (
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	gomail "gopkg.in/mail.v2"
)

type ConfigJSON struct {
	From         string `json:"from"`
	To           string `json:"to"`
	Password     string `json:"password"`
	SMTPUsername string `json:"smtp_username"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	Subject      string `json:"subject"`
	Interval     int    `json:"interval"`
}

var (
	Config ConfigJSON
	lastIP = "1.2.3.4"
)

func parseConfig() {
	conf := os.Getenv("MAIL_CONFIG")
	if conf == "" {
		conf = "config.json"
	}
	file, err := os.Open(conf)
	if err != nil {
		log.Fatalf("Read config fiel err: %v\n", err)
	}

	j := json.NewDecoder(file)
	err = j.Decode(&Config)
	if err != nil {
		log.Fatalf("Parse config failed: %v\n", err)
	}

	defer file.Close()
}

func httpGet(url string) string {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Read HTML Body failed: %v\n", err)
	}
	return string(content)
}

func getIP() string {
	ctx := httpGet("https://ip.tool.lu/")
	reg := regexp.MustCompile(`(\d+)\.(\d+).(\d+).(\d+)`)
	ips := reg.FindStringSubmatch(ctx)
	for _, ip := range ips {
		return ip
	}
	return ""
}

func sendMail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", Config.From)
	m.SetHeader("To", Config.To)
	m.SetHeader("Subject", Config.Subject)
	m.SetBody("text/plain", getIP())
	d := gomail.NewDialer(Config.SMTPHost, Config.SMTPPort, Config.SMTPUsername, Config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func main() {
	parseConfig()
	for {
		if strings.EqualFold(lastIP, getIP()) {
			log.Printf("IP has not changed: %v\n", getIP())
			time.Sleep(time.Duration(Config.Interval) * time.Second)
		} else {
			sendMail()
			log.Printf("IP send successfully: %v\n", getIP())
			lastIP = getIP()
			time.Sleep(time.Duration(Config.Interval) * time.Second)
		}
	}
}
