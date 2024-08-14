package email

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
)

// Email template with embedded HTML and CSS
const emailTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f4f4f4;
            margin: 0;
            padding: 0;
            color: #333333;
        }
        .container {
            width: 80%;
            margin: 0 auto;
            background-color: #ffffff;
            padding: 20px;
            border-radius: 10px;
            box-shadow: 0px 0px 10px 0px #aaaaaa;
        }
        h1 {
            color: #4CAF50;
        }
        p {
            line-height: 1.5;
        }
        .footer {
            margin-top: 20px;
            text-align: center;
            color: #777777;
        }
    </style>
    <title>{{ .Subject }}</title>
</head>
<body>
    <div class="container">
        <h1>Hello {{ .Name }},</h1>
        <p>{{ .Body }}</p>
        <p>Best regards,</p>
        <p>The JobStar Team</p>
    </div>
    <div class="footer">
        <p>&copy; 2024 JobStar. All rights reserved.</p>
    </div>
</body>
</html>
`

// EmailData holds the data for the email template
type EmailData struct {
	Subject string
	Name    string
	Body    string
}

// SendEmail sends an email using the SMTP server
func SendEmail(to string, subject string, name string, body string) error {
	from := os.Getenv("SMTP_SENDER_EMAIL")
	password := os.Getenv("SMTP_SENDER_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	// Create the email content
	tmpl, err := template.New("email").Parse(emailTemplate)
	if err != nil {
		return fmt.Errorf("error parsing email template: %w", err)
	}

	data := EmailData{
		Subject: subject,
		Name:    name,
		Body:    body,
	}

	var emailBody bytes.Buffer
	if err := tmpl.Execute(&emailBody, data); err != nil {
		return fmt.Errorf("error executing email template: %w", err)
	}

	message := []byte(fmt.Sprintf("Subject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", subject, emailBody.String()))

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, message)
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
