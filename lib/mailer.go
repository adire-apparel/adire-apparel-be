package lib

import (
	"adire-apparel/config"
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"sync"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	templates    map[string]*template.Template
	mu           sync.RWMutex
	dialer       *gomail.Dialer
	templatesDir string
}

type EmailDto struct {
	To       []string
	Subject  string
	Template string
	Data     interface{}
}

var (
	emailService *EmailService
	once         sync.Once
)

func GetEmailService() *EmailService {
	once.Do(func() {
		emailService = &EmailService{
			templates:    make(map[string]*template.Template),
			templatesDir: "templates",
			dialer: gomail.NewDialer(
				config.AppConfig.SmtpHost,
				config.AppConfig.SmtpPort,
				config.AppConfig.SmtpUser,
				config.AppConfig.SmtpPassword,
			),
		}
	})
	return emailService
}

func (es *EmailService) getTemplate(templateName string) (*template.Template, error) {
	es.mu.RLock()
	tmpl, exists := es.templates[templateName]
	es.mu.RUnlock()

	if exists {
		return tmpl, nil
	}

	es.mu.Lock()
	defer es.mu.Unlock()

	if tmpl, exists = es.templates[templateName]; exists {
		return tmpl, nil
	}

	fullPath := filepath.Join(es.templatesDir, templateName+".html")
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		return nil, err
	}

	es.templates[templateName] = tmpl
	return tmpl, nil
}

func (es *EmailService) renderTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := es.getTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (es *EmailService) SendEmail(payload EmailDto) error {
	html, err := es.renderTemplate(payload.Template, payload.Data)
	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.AppConfig.AppEmail)
	msg.SetHeader("To", payload.To...)
	msg.SetHeader("Subject", payload.Subject)
	msg.SetBody("text/html", html)

	return es.dialer.DialAndSend(msg)
}

func SendEmail(payload EmailDto) error {
	return GetEmailService().SendEmail(payload)
}

type TestEmailDto struct {
	Name  string `json:"name" validate:"required,name"`
	Email string `json:"email" validate:"required,email"`
}

func TestEmailConfig(payload TestEmailDto) error {
	otp := GenerateOtp()
	link := GenerateUrl()
	err := SendEmail(EmailDto{
		To:       []string{payload.Email},
		Subject:  "Test Email",
		Template: "test",
		Data: map[string]interface{}{
			"name": payload.Name,
			"otp":  otp,
			"link": link,
		},
	})
	if err != nil {
		log.Printf("Email test failed: %v", err)
		return err
	}

	log.Println("Email sent succesfully")
	return nil
}
