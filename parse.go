// Copyright 2014 Yoshi Yamaguchi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//      Unless required by applicable law or agreed to in writing, software
//      distributed under the License is distributed on an "AS IS" BASIS,
//      WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//      See the License for the specific language governing permissions and
//      limitations under the License.
package mailutil

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
)

// ParseMail parse general mail messages in r and returns HTMLMailMessage or
// TextMailMessage.
func ParseMail(r io.Reader) (MailMessage, error) {
	msg, err := mail.ReadMessage(r)
	if err != nil {
		return nil, err
	}

	mt, params, err := mime.ParseMediaType(msg.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	switch {
	case strings.HasPrefix(mt, "multipart/"):
		return NewHTMLMail(msg, params["boundary"])
	case strings.HasPrefix(mt, "text/"):
		return NewTextMail(msg)
	default:
		return nil, UnsupportedMediaTypeError{mt}
	}
}

// parseTextPart parses text/plain part in multipart message.
func parseTextPart(p *multipart.Part, charset string) ([]byte, error) {
	return DecodeText(p, charset)
}

// parseHTMLPart parses text/html part in multipart message.
func parseHTMLPart(p *multipart.Part, charset string) ([]byte, error) {
	mt, _, err := mime.ParseMediaType(p.Header.Get("Content-Transfer-Encoding"))
	if err != nil {
		return nil, err
	}
	switch mt {
	case "base64":
		decoder := base64.NewDecoder(base64.StdEncoding, p)
		return DecodeText(decoder, charset)
	case "quoted-printable":
		return ioutil.ReadAll(p)
	default:
		return nil, UnsupportedTransferEncodingError{mt}
	}
}

// DecodeText decodes encoded byte array data coming from r,
// and returns decoded data.
func DecodeText(r io.Reader, charset string) ([]byte, error) {
	lowered := strings.ToLower(charset)
	var decoder transform.Transformer
	switch lowered {
	case "iso-2022-jp":
		decoder = japanese.ISO2022JP.NewDecoder()
	case "shift_jis":
		decoder = japanese.ShiftJIS.NewDecoder()
	case "euc-jp":
		decoder = japanese.EUCJP.NewDecoder()
	}
	tr := transform.NewReader(r, decoder)
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
