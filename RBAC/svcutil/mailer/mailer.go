package mailer

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"net/mail"
	"path"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	gomail "gopkg.in/mail.v2"
)

type Email struct {
	mailer *gomail.Dialer

	from mail.Address
}

type EmailConfig struct {
	Server  string
	Port    int
	Timeout time.Duration

	Username, Password string
	DefaultFrom        mail.Address

	CACert       fs.File
	CACommonName string
}

func NewEmail(conf EmailConfig) (*Email, error) {
	var tlsConf *tls.Config
	if conf.CACert != nil {
		rootCAs, _ := x509.SystemCertPool()
		if rootCAs == nil {
			rootCAs = x509.NewCertPool()
		}
		certs, err := io.ReadAll(conf.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to append %q to RootCAs: %w", conf.CACommonName, err)
		}
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			return nil, errors.New("failed to append cert to system pool")
		}
		tlsConf = &tls.Config{ServerName: conf.CACommonName, RootCAs: rootCAs}
	}

	d := gomail.NewDialer(conf.Server, conf.Port, conf.Username, conf.Password)
	d.TLSConfig = tlsConf
	d.Timeout = conf.Timeout
	return &Email{mailer: d, from: conf.DefaultFrom}, nil
}

type EmailMessage struct {
	From    mail.Address
	To      mail.Address
	Subject string
	Msg     *template.Template
	Args    interface{}
	Assets  []Asset
}

type Asset struct {
	FileName string
	Data     fs.File
}

// AssetDir loads all files in the given directory as assets.
// Includes all sub-directories. Use `AssetGlob` with `fs.Sub`
// when sub-directories are not desired.
func (m *EmailMessage) AssetDir(f fs.FS, dir string) error {
	if f == nil {
		return fmt.Errorf("missing file system")
	}
	if dir == "" {
		dir = "."
	}
	return fs.WalkDir(f, dir,
		fs.WalkDirFunc(
			func(fPath string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				a, err := f.Open(fPath)
				if err != nil {
					return err
				}
				m.Assets = append(m.Assets, Asset{
					FileName: path.Base(d.Name()),
					Data:     a,
				})
				return nil
			}),
	)
}

func (m *EmailMessage) LoadAssets(dir fs.FS, filename ...string) error {
	if len(filename) == 0 {
		return nil
	}
	for _, fn := range filename {
		a, err := dir.Open(fn)
		if err != nil {
			return err
		}
		m.Assets = append(m.Assets, Asset{
			FileName: path.Base(fn),
			Data:     a,
		})
	}
	return nil
}

// AssetGlob loads all files matching the glob as assets.
func (m *EmailMessage) AssetGlob(f fs.FS, glob string) error {
	if glob == "" {
		return nil
	}
	fns, err := fs.Glob(f, glob)
	if err != nil {
		return err
	}
	return m.LoadAssets(f, fns...)
}

func (em *Email) Send(msg EmailMessage) error {
	if msg.From.Address == "" {
		msg.From = em.from
	}
	switch "" {
	case msg.From.Address:
		return fmt.Errorf(`missing "FROM" email address`)
	case msg.To.Address:
		return fmt.Errorf(`missing "TO" email address`)
	case msg.Subject:
		return fmt.Errorf(`missing email subject`)
	}
	if msg.Msg == nil {
		return fmt.Errorf(`missing message template`)
	}
	m := gomail.NewMessage()
	m.SetAddressHeader("From", msg.From.Address, msg.From.Name)
	m.SetAddressHeader("To", msg.To.Address, msg.To.Name)
	m.SetHeader("Subject", msg.Subject)

	var html bytes.Buffer
	if err := msg.Msg.Execute(&html, msg.Args); err != nil {
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
	m.SetBody("text/html", transformed)

	for _, f := range msg.Assets {
		img, err := parseFile(f.Data, f.FileName)
		if err != nil {
			return err
		}
		m.AttachReader(img.fileName, img.data, gomail.SetHeader(mimeHeader(*img, true)))
	}
	return em.mailer.DialAndSend(m)
}

func mimeHeader(img img, inline bool) map[string][]string {
	att := "attachment"
	if inline {
		att = "inline"
	}
	return map[string][]string{
		"Content-Type":              {fmt.Sprintf("%s;\n\tfilename=%q", img.mimeType, img.fileName)},
		"Content-Disposition":       {fmt.Sprintf("%s;\n\tfilename=%q", att, img.fileName)},
		"Content-Transfer-Encoding": {"base64"},
		"Content-ID":                {fmt.Sprintf("<%s>", img.fileName)},
	}
}

type img struct {
	fileName string
	data     io.Reader
	mimeType string
}

func parseFile(f fs.File, fileName string) (*img, error) {
	bf := bufio.NewReaderSize(f, 512)
	b, err := bf.Peek(512)
	if err != nil {
		return nil, err
	}
	return &img{
		fileName: fileName,
		data:     bf,
		mimeType: http.DetectContentType(b),
	}, nil
}
