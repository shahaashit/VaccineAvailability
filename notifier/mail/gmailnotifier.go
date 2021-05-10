package mail

import (
	"VaccineAvailability/config"
	"gopkg.in/gomail.v2"
)

var (
	mailClientConfig config.MailClientConf
	dialer           *gomail.Dialer
)

func init() {
	mailClientConfig = config.AppConfiguration.MailClientConfig
	dialer = gomail.NewDialer(mailClientConfig.Host, mailClientConfig.Port, mailClientConfig.UserName, mailClientConfig.Password)
}

func SendMail(mailTo string, subject string, body string) error {
	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("From", mailClientConfig.UserName)
	goMailMessage.SetHeader("To", mailTo)
	goMailMessage.SetHeader("Subject", subject)
	goMailMessage.SetBody("text/plain", body)

	if err := dialer.DialAndSend(goMailMessage); err != nil {
		return err
	}
	return nil
}
