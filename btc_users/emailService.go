package main
import (
	"github.com/sendgrid/sendgrid-go"
	"fmt"
)

func emailPassword(email string, password string) {
	body := "Your password has been reset successfully. Your new password is: <br><br><h2>" + password + "</h2><br><br>"
	body = body + "If you believe this is a mistake, do not reply to this email. Please contact us at btcpredictions@gmail.com"
	sg := sendgrid.NewSendGridClientWithApiKey("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")

	message := sendgrid.NewMail()
	message.AddTo(email)
	message.SetFrom("NoReply@BTCPredictions.com")
	message.SetSubject("Successfully reset password")
	message.SetHTML(body)
	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}
}

func sendEmail(userAction UserAction, token string) {
	link := "http://localhost:8080/api/doaction/" + token

	var body string
	var subject string
	if userAction.Action == "resetPassword" {
		body = "<h3>To reset your password, please click the following link. You will then be sent an email containing your new password.</h3>"
		body = body + "<br><br>" + link
		subject = "Request to reset password"
	} else if userAction.Action == "verifyEmail" {
		body = "<h3>Please click the following link to verify your account</h3>"
		body = body + "<br><br>" + link
		subject = "Please verify your account"
	}

	sg := sendgrid.NewSendGridClientWithApiKey("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")

	message := sendgrid.NewMail()
	message.AddTo(userAction.Email)
	message.SetFrom("NoReply@BTCPredictions.com")
	message.SetSubject(subject)
	message.SetHTML(body)
	if r := sg.Send(message); r == nil {
		fmt.Println("Email sent!")
	} else {
		fmt.Println(r)
	}
}