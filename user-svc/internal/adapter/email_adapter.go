package adapter

import (
	"be-yourmoments/user-svc/internal/helper/utils"
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailAdapter interface {
	SendEmail(email string, token string, category string) error
}

type smtpService struct {
	dialer *gomail.Dialer
	from   string
}

type emailAdapter struct {
	services    []smtpService
	backendUrl  string
	frontendUrl string
}

func NewEmailAdapter() EmailAdapter {
	services := []smtpService{}

	for i := 1; i <= 1; i++ {
		smtpPortStr := utils.GetEnv(fmt.Sprintf("SMTP_PORT%d", i))
		smtpPort, err := strconv.Atoi(smtpPortStr)
		if err != nil {
			log.Printf("Invalid SMTP_PORT%d: %v", i, err)
			continue
		}

		smtpHost := utils.GetEnv(fmt.Sprintf("SMTP_HOST%d", i))
		smtpUsername := utils.GetEnv(fmt.Sprintf("SMTP_USERNAME%d", i))
		smtpPassword := utils.GetEnv(fmt.Sprintf("SMTP_PASSWORD%d", i))
		smtpEmailFrom := utils.GetEnv(fmt.Sprintf("SMTP_EMAIL_FROM%d", i))

		if smtpHost == "" || smtpUsername == "" || smtpPassword == "" || smtpEmailFrom == "" {
			log.Printf("Incomplete SMTP configuration for service %d", i)
			continue
		}

		dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
		services = append(services, smtpService{dialer: dialer, from: smtpEmailFrom})
	}

	backendUrl := utils.GetEnv("BACKEND_URL")
	frontendUrl := utils.GetEnv("FRONTEND_URL") // Pastikan penamaan env sudah benar

	return &emailAdapter{
		services:    services,
		backendUrl:  backendUrl,
		frontendUrl: frontendUrl,
	}
}

func (a *emailAdapter) SendEmail(email string, token string, category string) error {
	emailData := struct {
		Email       string
		Token       string
		FrontendUrl string
		AppUrl      string
	}{
		Email:       email,
		Token:       token,
		FrontendUrl: a.frontendUrl,
		AppUrl:      a.backendUrl,
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}

	subject := ""
	filePath := ""
	switch category {
	case "reset password":
		subject = "Email Reset Password"
		filePath = "request_reset_password.html"
	case "new email verification":
		subject = "User Email Verification"
		filePath = "registration.html"
	default:
		return fmt.Errorf("kategori email tidak dikenali: %s", category)
	}

	templatePath := filepath.Join(cwd, "../../internal/templates", filePath)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("error parsing html template: %w", err)
	}

	var b bytes.Buffer
	if err := tmpl.Execute(&b, emailData); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	go func() {
		var lastErr error
		for i, service := range a.services {
			mailer := gomail.NewMessage()
			mailer.SetHeader("From", "hervipro@gmail.com")
			mailer.SetHeader("To", email)
			mailer.SetHeader("Subject", subject)
			mailer.SetBody("text/html", b.String())

			err := service.dialer.DialAndSend(mailer)
			if err != nil {
				log.Printf("SMTP service %d gagal: %v", i+1, err)
				fmt.Printf("SMTP service, error terakhir: %w", err)
				lastErr = err
			} else {
				break
			}

		}
		if lastErr != nil {
			log.Printf("semua layanan SMTP gagal, error terakhir: %w", lastErr)
		}

		log.Println("email berhasil dikirimkan")
	}()

	return nil
}
