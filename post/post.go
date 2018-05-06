package post

import (
	"crypto/tls"
	"fmt"
	"net"
	"strings"

	"github.com/Slava12/Computer_Market/logger"

	"github.com/Slava12/Computer_Market/config"

	"net/mail"
	"net/smtp"
)

var (
	server   string
	user     string
	password string
)

// Init инициализирует переменные для работы с почтой
func Init(configFile config.Config) {
	server = configFile.Post.Server
	user = configFile.Post.User
	password = configFile.Post.Password
}

// SendMail Отправляет письмо по указанному адресу
func SendMail(recipient string, subj string, body string) {

	from := mail.Address{"", user}
	to := mail.Address{"", recipient}

	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	servername := server

	host, _, _ := net.SplitHostPort(servername)

	stringArray := strings.Split(user, "@")

	auth := smtp.PlainAuth("", stringArray[0], password, host)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		logger.Warn(err, "Не удалось подключиться к почтовому серверу!")
		return
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		logger.Warn(err, "Не удалось создать нового клиента почтового сервера!")
		return
	}

	if err = c.Auth(auth); err != nil {
		logger.Warn(err, "Не удалось авторизоваться на почтовом сервере!")
		return
	}

	if err = c.Mail(from.Address); err != nil {
		logger.Warn(err, "Не удалось сформировать отправителя!")
		return
	}

	if err = c.Rcpt(to.Address); err != nil {
		logger.Warn(err, "Не удалось сформировать получателя!")
		return
	}

	w, err := c.Data()
	if err != nil {
		logger.Warn(err, "Не удалось создать поток!")
		return
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		logger.Warn(err, "Не удалось записать в поток!")
		return
	}

	err = w.Close()
	if err != nil {
		logger.Warn(err, "Не удалось закрыть поток!")
		return
	}
	c.Quit()
}
