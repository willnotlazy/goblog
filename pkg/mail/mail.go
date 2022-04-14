package mail

import (
	"goblog/pkg/logger"
	"gopkg.in/gomail.v2"
)

type EmailInfo struct {
}

func SendMail(subject, body string, to ...string) error {
	m := gomail.NewMessage()

	m.SetHeader("Subject", "重置密码")
	m.SetHeader("To", to...)
	m.SetHeader("From", m.FormatAddress("519478450@qq.com", "goblog@qq.com"))

	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.qq.com", 465, "519478450@qq.com", "vqmypcczikhxbhad")
	if err := d.DialAndSend(m); err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}