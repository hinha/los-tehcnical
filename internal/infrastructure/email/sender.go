package email

import (
	"github.com/sirupsen/logrus"
)

// ConsoleEmailSender is a simple implementation of the EmailSender interface that logs emails to the console
type ConsoleEmailSender struct {
	logger *logrus.Logger
}

// NewConsoleEmailSender creates a new console email sender
func NewConsoleEmailSender(logger *logrus.Logger) *ConsoleEmailSender {
	return &ConsoleEmailSender{
		logger: logger,
	}
}

// SendAgreementEmail sends an agreement email to an investor
func (s *ConsoleEmailSender) SendAgreementEmail(email, loanID, agreementURL string) error {
	s.logger.WithFields(logrus.Fields{
		"layer":         "email",
		"function":      "SendAgreementEmail",
		"email":         email,
		"loan_id":       loanID,
		"agreement_url": agreementURL,
	}).Info("Sending agreement email")

	// In a real implementation, this would send an actual email
	// For now, just log the email details

	return nil
}

// PDFGenerator generates PDF agreement letters
//type PDFGenerator struct {
//	logger *logrus.Logger
//}
//
//// NewPDFGenerator creates a new PDF generator
//func NewPDFGenerator(logger *logrus.Logger) *PDFGenerator {
//	return &PDFGenerator{
//		logger: logger,
//	}
//}
//
//// GenerateAgreementLetter generates a PDF agreement letter for a loan
//func (g *PDFGenerator) GenerateAgreementLetter(loanID string) (string, error) {
//	// In a real implementation, this would generate an actual PDF
//	// For now, just return a fake URL
//
//	pdfURL := fmt.Sprintf("https://asia-southeast2-dummy.cloudfunctions.net/agreements/%s.pdf", loanID)
//	g.logger.WithFields(logrus.Fields{
//		"layer":    "email",
//		"function": "GenerateAgreementLetter",
//		"loan_id":  loanID,
//		"pdf_url":  pdfURL,
//	}).Info("Generated agreement letter")
//
//	return pdfURL, nil
//}
