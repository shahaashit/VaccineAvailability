package notifier

import (
	"VaccineAvailability/config"
	"gopkg.in/gomail.v2"
)

var (
	mailClientConfig config.EmailConf
	dialer           *gomail.Dialer
)

func init() {
	mailClientConfig = config.AppConfiguration.EmailConfig
	dialer = gomail.NewDialer(mailClientConfig.Host, mailClientConfig.Port, mailClientConfig.UserName, mailClientConfig.Password)
}

func sendMail(mailTo string, subject string, body string) error {
	goMailMessage := gomail.NewMessage()
	goMailMessage.SetHeader("From", mailClientConfig.UserName)
	goMailMessage.SetHeader("To", mailTo)
	goMailMessage.SetHeader("Subject", subject)
	goMailMessage.SetBody("text/plain", body)
	//for _, file := range attachmentPaths {
	//	goMailMessage.Attach(file)
	//}

	if err := dialer.DialAndSend(goMailMessage); err != nil {
		return err
	}
	return nil
}
