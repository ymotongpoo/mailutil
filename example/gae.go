package gae

import (
	"net/http"

	"appengine"
	gaemail "appengine/mail"

	"mailutil"
)

func init() {
	// Replace "appid" with your application ID.
	http.HandleFunc("/_ah/mail/string@appid.appspotmail.com", inboundMailHandler)
}

// Google App Engine for Go runtime routes mail messages received on "string@appid.appspotmail.com"
// to the endpoint "/_ah/mail/string@appid.appspotmail.com" as *http.Request, and mail header and body
// are stored as Body field in the http.Request.
//
// See details on official docs.
// https://cloud.google.com/appengine/docs/go/mail/#Go_Receiving_mail
func inboundMailHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	defer r.Body.Close()

	// Parse your mail message.
	reqmsg, err := mailutil.ParseMail(r.Body)

	addrs, err := reqmsg.Header().AddressList("from")
	if err != nil {
		c.Errorf("Error parsing mail header: %v", err)
	}
	addrStr := make([]string, len(addrs))
	for i, a := range addrs {
		addrStr[i] = a.String()
	}

	// Echo back your message.
	respmsg := &gaemail.Message{
		Sender:  "Service sender <info@appid.appspotmail.com>",
		To:      addrStr,
		Subject: "Your message -- " + reqmsg.Header().Get("Subject"),
		Body:    reqmsg.String(),
	}
	if err := gaemail.Send(c, respmsg); err != nil {
		c.Errorf("Couldn't send email: %v", err)
	}
}
