package apiproxy

import (
	"fmt"
	"net/smtp"
	"sync/atomic"

	"go.uber.org/zap"
)

func (svc *APIProxyGRPCService) sendEmail(toEmail string, subject, content string) error {
	go func() {
		curIndex := atomic.AddUint32(&svc.emailSendIndex, 1)
		emails := svc.rm.getEmailAddrs()
		for i := 0; i < len(emails); i++ {
			email := emails[(int(curIndex)+i)%len(emails)]
			auth := smtp.PlainAuth("", email.Addr, email.Passwd, email.Host)
			fmt.Println(email)
			msg := []byte(fmt.Sprintf("To: %s \r\nSubject: %s\r\nMIME-version: 1.0;\r\nContent-Type: text/html;\r\n\r\n%s\r\n", toEmail, subject, content))
			err := smtp.SendMail(email.Host+":"+email.Port, auth, email.Addr, []string{toEmail}, msg)
			if err == nil {
				break
			}
			svc.logger.Error("send email failed ", zap.Error(err))
		}
	}()
	return nil
}

func (svc *APIProxyGRPCService) sendEmailBindCode(toEmail string, bindCode string) error {
	subject := "Register Mirror Realms Account With Email Verification Code"
	content := fmt.Sprintf(`<html><body><p>Your verification code is:</p>
<p style="font-size: 20px; font-weight: bold;">%s</p>
<p>This code will be used for Mirror Realms account registration, and is valid for 10 minutes.</p>
<p>For your security, please don't share this code with anyone else. If you did not make this request, ignore this message.</p>
<p>Thanks!</p>
<p>Mirror Realms Team</p></body></html>`, bindCode)
	return svc.sendEmail(toEmail, subject, content)
}

func (svc *APIProxyGRPCService) sendEmailResetPasswordValidationCode(toEmail string, code string) error {
	subject := "Change Your Password With Email Verification Code"
	content := fmt.Sprintf(`<html><body><p>Your verification code is:</p>
<p style="font-size: 20px; font-weight: bold;">%s</p>
<p>This verification code is used to change the password, and is valid for 10 minutes.</p>
<p>For your security, please don't share this code with anyone else. If you did not make this request, ignore this message.</p>
<p>Thanks!</p>
<p>Mirror Realms Team</p></body></html>`, code)
	return svc.sendEmail(toEmail, subject, content)
}
