package main

import (
	"fmt"
	"strings"
)

type ContactForm struct{ Name, Email, Message string }
type BuildEvent struct{ Commit, Status string }
type ReleaseOperation struct{ Version, Environment string }
type Diagnostic struct{ Code, Detail string }

func routeContact(form ContactForm, inbox string) (EmailResult, error) {
	if strings.TrimSpace(form.Email) == "" || strings.TrimSpace(form.Message) == "" {
		return EmailResult{}, fmt.Errorf("email and message are required")
	}
	payload := map[string]string{"to": inbox, "subject": "Developer tools contact: " + form.Name, "body": fmt.Sprintf("From: %s <%s>\n\n%s", form.Name, form.Email, form.Message)}
	return sendEmail(payload, "contact-"+form.Email+"-"+form.Message)
}

func buildDiagnostic(event BuildEvent) Diagnostic {
	if event.Status == "failed" {
		return Diagnostic{Code: "BUILD_FAILED", Detail: "Inspect compiler output for commit " + event.Commit}
	}
	return Diagnostic{Code: "BUILD_OK", Detail: "Build passed for commit " + event.Commit}
}
