package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"log"
	//"mime"
	"net/smtp"
	"strings"
	"time"
)

var (
	recvUsers = []string{"daiping@inesadt.com", "daiping_zy@139.com", "yanglaoxj@inesadt.com"}
	user      = "yanglaoxj@inesadt.com"
	passwd    = "INEsa2024"
	smtpS     = "mail.inesadt.com"
	smtpP     = 25
)

func init() {
	recvUsers = formatEmails(recvUsers)
	//fmt.Println(recvUsers)
	//fmt.Println(toEmailString(recvUsers))
}

func formatEmails(emails []string) []string {
	formattedEmails := make([]string, len(emails))
	for i, email := range emails {
		localPart := strings.Split(email, "@")[0]
		domain := strings.Split(email, "@")[1]
		formattedEmails[i] = fmt.Sprintf("%s<%s@%s>", strings.Title(localPart), localPart, domain)
	}
	return formattedEmails
}

func toEmailString(emails []string) string {
	var r string = ""
	for i, v := range emails {
		if 0 == i {
			r = v
			continue
		}
		r = r + ";" + v
	}
	return r
}

func main() {
	SendEmail("Hello, My name is DaiPing. Here is a test message. 1234567890!@#$%^&*()_+不要回答！不要回答！不要回答！")
}

func SendEmail(content string) {
	// 收件人和抄送人列表
	recipient := recvUsers[0]
	ccRecipients := recvUsers[1:]
	log.Println(1)
	// 构建邮件内容
	subject := "[自动巡检] 告警触发 " + time.Now().Format("2006年01月02日 15时04分05秒")
	subjectBase := base64.StdEncoding.EncodeToString([]byte(subject))

	// 连接到SMTP服务器
	log.Println(smtpS, ":", smtpP)
	connection, err := smtp.Dial(fmt.Sprintf("%s:%d", smtpS, smtpP))
	if err != nil {
		log.Println("[WARN]Dial: ", err)
		return
	}
	defer connection.Close()
	log.Println(2)

	// 创建TLS配置
	tlsConfig := &tls.Config{
		ServerName:         smtpS,            // 设置服务器名称
		MinVersion:         tls.VersionTLS10, // 设置最小支持的TLS版本
		MaxVersion:         tls.VersionTLS13, // 设置最大支持的TLS版本
		InsecureSkipVerify: true,             // 跳过证书验证，仅用于测试或信任的服务器
	}
	// 启动TLS加密
	if err := connection.StartTLS(tlsConfig); err != nil {
		log.Fatalf("smtp.StartTLS error: %s", err)
		return
	}
	log.Println(3)

	// 进行身份验证
	if err = connection.Auth(smtp.PlainAuth("", user, passwd, smtpS)); err != nil {
		log.Println("[WARN]Auth: ", err)
	}
	log.Println(4)

	// 设置发件人
	if err = connection.Mail(user); err != nil {
		log.Println("[WARN]Mail: ", err)
		return
	}
	log.Println(5)

	// 设置收件人
	if err = connection.Rcpt(recipient); err != nil {
		log.Println("[WARN]Rcpt: ", err)
		return
	}
	log.Println(6)

	// 设置抄送人
	for _, cc := range ccRecipients {
		if err = connection.Rcpt(cc); err != nil {
			log.Println("[WARN]Rcpt: ", err)
			return
		}
	}
	log.Println(8)

	// 开始写入邮件内容
	w, err := connection.Data()
	if err != nil {
		log.Println("[WARN]Data: ", err)
		return
	}
	log.Println(9)
	var from = fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte("养老服务自动巡检")))
	//	mime.QEncoding.Encode("UTF-8", "养老服务自动巡检"),
	// 写入邮件头部
	_, err = fmt.Fprintf(w, "From: %s\r\nTo: %s\r\nCc: %s\r\nSubject: =?UTF-8?B?%s?=\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		from+"<"+user+">",
		recipient,
		toEmailString(ccRecipients),
		subjectBase)
	if err != nil {
		log.Println("[WARN]Fprintf: ", err)
	}

	// 写入邮件正文
	_, err = fmt.Fprintf(w, "%s\r\n", content)
	if err != nil {
		log.Println("[WARN]Fprintf: ", err)
	}

	// 关闭写入器
	err = w.Close()
	if err != nil {
		log.Println("[WARN]Close: ", err)
	}

	log.Println("Email sent successfully!")
}
