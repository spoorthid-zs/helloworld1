package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
)

type recaptchaResponse struct {
	Success     bool     `json:"success"`
	ChallengeTs string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			renderForm(w)
			return
		}

		// Verify reCAPTCHA response
		success, err := verifyRecaptcha(r.RemoteAddr, r.FormValue("g-recaptcha-response"))
		if err != nil {
			http.Error(w, "Error verifying reCAPTCHA", http.StatusInternalServerError)
			return
		}
		if !success {
			http.Error(w, "reCAPTCHA verification failed", http.StatusUnauthorized)
			return
		}

		// reCAPTCHA verification successful, proceed with form submission
		// ...
		fmt.Fprintln(w, "reCAPTCHA verification successful")
	})

	http.ListenAndServe(":8080", nil)
}

func renderForm(w http.ResponseWriter) {
	tmpl := template.Must(template.New("form").Parse(`
		<html>
			<head>
				<title>reCAPTCHA Test Form</title>
			</head>
			<body>
				<form method="POST">
					<div class="g-recaptcha" data-sitekey="6Lew8YIlAAAAAPj4-m_jdH9ZsXsr_HA450tpPXgm"></div>
					<br>
					<input type="submit" value="Submit">
				</form>
				<script src="https://www.google.com/recaptcha/api.js"></script>
			</body>
		</html>
	`))
	tmpl.Execute(w, nil)
}

func verifyRecaptcha(remoteip, response string) (bool, error) {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {"6Lew8YIlAAAAAAikEIXjZgF-vRnL9nSOUna4A27K"},
			"remoteip": {remoteip},
			"response": {response}})

	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	var recaptchaResponse recaptchaResponse
	err = json.Unmarshal(body, &recaptchaResponse)
	if err != nil {
		return false, err
	}

	return recaptchaResponse.Success, nil
}
