package email

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/domodwyer/mailyak/v3"
	"github.com/vanng822/go-premailer/premailer"
)

const (
	assetPath = "integrations/email/assets/"
)

type MailSender struct {
	server   string
	port     int
	username string
	password string
	fromMail string
	fromName string
	cmsURL   string
}

func New(server string, port int, username, password, fromMail, fromName, cmsURL string) *MailSender {
	return &MailSender{
		server:   server,
		port:     port,
		username: username,
		password: password,
		fromMail: fromMail,
		fromName: fromName,
		cmsURL:   cmsURL,
	}
}

func (ms *MailSender) mailYak() *mailyak.MailYak {
	host := fmt.Sprintf("%s:%d", ms.server, ms.port)
	m := mailyak.New(host, smtp.PlainAuth("", ms.username, ms.password, ms.server))
	return m
}

func (ms *MailSender) sendMail(
	fromName, email, subject string,
	tmpl *template.Template,
	data interface{},
	imgs []string,
) error {
	yak := ms.mailYak()
	yak.From(ms.fromMail)
	yak.FromName(fromName)
	yak.Subject(subject)
	yak.To(email)

	var html bytes.Buffer

	if err := tmpl.Execute(&html, data); err != nil {
		return err
	}

	prem, err := premailer.NewPremailerFromBytes(html.Bytes(), premailer.NewOptions())
	if err != nil {
		return err
	}

	transformed, err := prem.Transform()
	if err != nil {
		return err
	}

	yak.HTML().Set(transformed)

	for _, i := range imgs {
		img, err := getImage(i)
		if err != nil {
			return err
		}
		yak.AttachInlineWithMimeType(img.fileName, img.data, img.mimeType)
	}
	return yak.Send()
}

type onboardingReminderForm struct {
	RedirectURL string
}

func (ms *MailSender) OnboardingReminder(email string, orgID string, userID string) error {
	onboardingReminder := template.Must(
		template.New("onboarding-reminder.html").
			ParseFiles(assetPath + "onboarding-reminder.html"))

	const subj = "Please complete your onboarding"
	fullURL := fmt.Sprintf("%s/register/businessinfo?org_id=%s&user_id=%s", ms.cmsURL, orgID, userID)
	f := onboardingReminderForm{
		RedirectURL: fullURL,
	}
	imgs := []string{"white-logo.png"}
	return ms.sendMail(ms.fromName, email, subj, onboardingReminder, f, imgs)
}

type img struct {
	fileName string
	data     io.Reader
	mimeType string
}

func getImage(fileName string) (*img, error) {
	path := fmt.Sprintf("%s/%s", assetPath, fileName)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 512)
	if _, err = f.Read(b); err != nil {
		return nil, err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return nil, err
	}
	ct := http.DetectContentType(b)
	return &img{
		fileName: fileName,
		data:     f,
		mimeType: ct,
	}, nil
}

type DsaServiceRequestNotificationForm struct {
	Email        string
	Status       string
	ServiceName  string
	Remark       string
	PartnerNames string
}

func (ms *MailSender) DsaServiceRequestNotification(req DsaServiceRequestNotificationForm) error {
	dsaServiceRequestNotification := template.Must(
		template.New("dsa-service-request.html").
			ParseFiles(assetPath + "dsa-service-request.html"))

	req.ServiceName = strings.Title(strings.ToLower(req.ServiceName))
	const subj = "PETNET - Service Request Status"
	imgs := []string{"white-logo.png"}
	return ms.sendMail(ms.fromName, req.Email, subj, dsaServiceRequestNotification, req, imgs)
}
