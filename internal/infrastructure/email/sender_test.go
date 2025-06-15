package email

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestConsoleEmailSender_SendAgreementEmail(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name         string
		email        string
		loanID       string
		agreementURL string
		expectError  bool
	}{
		{
			name:         "Valid email parameters",
			email:        "investor@example.com",
			loanID:       "LOAN123",
			agreementURL: "https://example.com/agreements/LOAN123.pdf",
			expectError:  false,
		},
		{
			name:         "Empty email address",
			email:        "",
			loanID:       "LOAN123",
			agreementURL: "https://example.com/agreements/LOAN123.pdf",
			expectError:  false,
		},
		{
			name:         "Empty loan ID",
			email:        "investor@example.com",
			loanID:       "",
			agreementURL: "https://example.com/agreements/LOAN123.pdf",
			expectError:  false, // Current implementation doesn't validate inputs
		},
		{
			name:         "Empty agreement URL",
			email:        "investor@example.com",
			loanID:       "LOAN123",
			agreementURL: "",
			expectError:  false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			logger.SetLevel(logrus.InfoLevel)

			sender := NewConsoleEmailSender(logger)

			err := sender.SendAgreementEmail(tc.email, tc.loanID, tc.agreementURL)

			// Assert expectations
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Verify log entry was created
			assert.Equal(t, 1, len(hook.Entries))
			assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
			assert.Equal(t, "Sending agreement email", hook.LastEntry().Message)

			// Verify log fields
			assert.Equal(t, "email", hook.LastEntry().Data["layer"])
			assert.Equal(t, "SendAgreementEmail", hook.LastEntry().Data["function"])
			assert.Equal(t, tc.email, hook.LastEntry().Data["email"])
			assert.Equal(t, tc.loanID, hook.LastEntry().Data["loan_id"])
			assert.Equal(t, tc.agreementURL, hook.LastEntry().Data["agreement_url"])

			// Clear log entries for next test
			hook.Reset()
		})
	}
}
