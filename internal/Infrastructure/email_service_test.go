package infrastructure

import (
	"testing"
)

func TestEmailService_WrapperCalls(t *testing.T) {
	// Use bogus config so sendEmail will fail fast; we only assert non-panics
	m := NewMailtrapService("localhost", "0", "", "", "from@example.com")
	e := NewEmailService(m, "http://localhost")
	_ = e.SendVerificationEmail("to@example.com", "code")
	_ = e.SendPasswordReset("to@example.com", "token")
	_ = e.SendUpdateVerificationEmail("to@example.com", "code")
}
