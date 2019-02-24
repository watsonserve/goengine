package mail

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"log"
	"net/smtp"
)

type Mailer struct {
	Name   string
	Prefix string
	Domain string
}

type MailContent struct {
	Subject string
	From    *Mailer
	To      *Mailer
	Cc      []*Mailer
	Body    string
}

func (this *MailContent) SendMail(server string, password string) error {
	msg := ""
	from := this.From.Prefix + "@" + this.From.Domain
	to := this.To.Prefix + "@" + this.To.Domain

	// 建立安全链接
	conn, err := tls.Dial("tcp", server + ":465", nil)
	if nil != err {
		return errors.New("on Connect: " + err.Error())
	}

	// 初始化smtp客户端
	client, err := smtp.NewClient(conn, this.From.Domain)
	if nil != err {
		return errors.New("on Init: " + err.Error())
	}

	defer client.Close()

	auth := smtp.PlainAuth("", from, password, this.From.Domain)

	err = client.Auth(auth)
	if nil != err {
		return errors.New("on Auth: " + err.Error())
	}

	err = client.Mail(from)
	if nil != err {
		return errors.New("on Mail: " + err.Error())
	}

	err = client.Rcpt(to)
	if nil != err {
		return errors.New("on Rcpt: " + err.Error())
	}

	writer, err := client.Data()
	if nil != err {
		return errors.New("on Data: " + err.Error())
	}
	body := []byte(this.Body)
	this.Body = base64.StdEncoding.EncodeToString(body)

	msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Transfer-Encoding: base64\r\n"
	msg += "From: " + this.From.Name + "<" + from + ">\r\n"
	msg += "To: " + this.To.Name + "<" + to + ">\r\n"
	msg += "Subject: " + this.Subject + "\r\n"
	msg += "\r\n" + this.Body + "\r\n"
	_, err = writer.Write([]byte(msg))
	if nil != err {
		return errors.New("on Write: " + err.Error())
	}

	err = client.Quit()
	if nil != err {
		msg := err.Error()
		if "250" != msg[0:3] {
			return errors.New("on Quit: " + msg)
		}
	}
	return nil
}

func Test() {
	mailContent := &MailContent{
		From: &Mailer{
			Name:   "JamesWatson",
			Prefix: "sender",
			Domain: "watsonserve.com",
		},
		To: &Mailer{
			Name:   "somebody",
			Prefix: "recver",
			Domain: "watsonserve.com",
		},
		Subject: "a test email",
		Body:    "<!DOCTYPE html><html><body><p>Dear my baby:</p><p>this is a test email.</p></body></html>",
	}

	err := mailContent.SendMail("smtp.watsonserve.com", "password")
	if nil != err {
		log.Println(err)
	}
}
