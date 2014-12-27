// Package recaptcha handles reCaptcha (http://www.google.com/recaptcha) form submissions
//
// This package is designed to be called from within an HTTP server or web framework
// which offers reCaptcha form inputs and requires them to be evaluated for correctness
//
// Edit the recaptcha_private_key constant before building and using
package recaptcha

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// https://www.google.com/recaptcha/api/siteverify?secret=your_secret&response=response_string&remoteip=user_ip_address
const recaptcha_server_name = "https://www.google.com/recaptcha/api/siteverify?"

var recaptcha_secret string

// check uses the client ip address and the client's response input to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func check(remoteip, response string) (body []byte) {
	vals := url.Values{"secret": recaptcha_secret, "remoteip": remoteip, "response": response}
	resp, err := http.Get(recaptcha_server_name + vals.Encode())
	if err != nil {
		log.Println("Get error: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read error: could not read body: %s", err)
	}
	return
}

type recaptcha_api_response struct {
	success bool
	errors  []string `json:"error-codes"`
}

// Confirm is the public interface function.
// It calls check, which the client ip address, the challenge code from the reCaptcha form,
// and the client's response input to that challenge to determine whether or not
// the client answered the reCaptcha input question correctly.
// It returns a boolean value indicating whether or not the client answered correctly.
func Confirm(remoteip, response string) bool {
	var r recaptcha_api_response
	if err := json.Unmarshal(check(remoteip, response), &r); err != nil {
		log.Println("JSON error:", err)
	}
	if len(r.errors) > 0 {
		log.Println("Recaptcha errors:", r.errors)
	}
	return r.success
}

// Init allows the webserver or code evaluating the reCaptcha form input to set the
// reCaptcha secret (string) value.
func Init(secret string) {
	recaptcha_secret = secret
}
