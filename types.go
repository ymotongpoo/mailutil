package mailutil

import (
	"io"
	"mime"
	"mime/multipart"
	"net/mail"
)

// UnsupportedMediaTypeError is specific error which covers all unsupported
// media type.
type UnsupportedMediaTypeError struct {
	mediatype string
}

func (u UnsupportedMediaTypeError) Error() string {
	return "Unsupported Media Type: " + u.mediatype
}

// MailMessage describes features all mail messages should have.
type MailMessage interface {
	String() string
	Header() mail.Header
}

type UnsupportedTransferEncodingError struct {
	transferencoding string
}

func (u UnsupportedTransferEncodingError) Error() string {
	return "Unsupported Content Transfer Encoding: " + u.transferencoding
}

// HTMLMailMessage implements MailMessage where each part of multipart is decoded
// into UTF-8 and stored as []byte in text and html fields respectively.
type HTMLMailMessage struct {
	body   io.Reader
	header mail.Header
	text   []byte
	html   []byte
}

// String returns text in text/plain part in HTML message.
func (hm *HTMLMailMessage) String() string {
	return string(hm.text)
}

// Header returns header part in HTML message.
func (hm *HTMLMailMessage) Header() mail.Header {
	return hm.header
}

// TextMailMessage implements MailMessage text message is decoded into UTF-8
// and stored as []byte.
type TextMailMessage struct {
	body   io.Reader
	header mail.Header
	text   []byte
}

func (tm *TextMailMessage) String() string {
	return string(tm.text)
}

// Header returns header part in text message.
func (tm *TextMailMessage) Header() mail.Header {
	return tm.header
}

// NewHTMLMail returns decoded HTML mail message.
func NewHTMLMail(msg *mail.Message, boundary string) (*HTMLMailMessage, error) {
	mr := multipart.NewReader(msg.Body, boundary)
	m := &HTMLMailMessage{
		body:   msg.Body,
		header: msg.Header,
	}
	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			return m, nil
		}
		if err != nil {
			return nil, err
		}
		mt, params, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
		if err != nil {
			return nil, err
		}
		switch mt {
		case "text/plain":
			text, err := parseTextPart(p, params["charset"])
			if err != nil {
				return nil, err
			}
			m.text = text
		case "text/html":
			html, err := parseHTMLPart(p, params["charset"])
			if err != nil {
				return nil, err
			}
			m.html = html
		}
	}
}

func NewTextMail(msg *mail.Message) (*TextMailMessage, error) {
	return nil, nil // TODO(ymotongpoo): stub
}
