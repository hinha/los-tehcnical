// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package provider is a generated GoMock package.
package provider

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	loan "github.com/hinha/los-technical/internal/domain/loan"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// SendAgreementEmail mocks base method.
func (m *MockEmailSender) SendAgreementEmail(email, loanID, agreementURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAgreementEmail", email, loanID, agreementURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAgreementEmail indicates an expected call of SendAgreementEmail.
func (mr *MockEmailSenderMockRecorder) SendAgreementEmail(email, loanID, agreementURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAgreementEmail", reflect.TypeOf((*MockEmailSender)(nil).SendAgreementEmail), email, loanID, agreementURL)
}

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddInvestment mocks base method.
func (m *MockService) AddInvestment(id, investorID, email string, amount float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddInvestment", id, investorID, email, amount)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddInvestment indicates an expected call of AddInvestment.
func (mr *MockServiceMockRecorder) AddInvestment(id, investorID, email, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddInvestment", reflect.TypeOf((*MockService)(nil).AddInvestment), id, investorID, email, amount)
}

// ApproveLoan mocks base method.
func (m *MockService) ApproveLoan(id, validatorID, proofURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ApproveLoan", id, validatorID, proofURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// ApproveLoan indicates an expected call of ApproveLoan.
func (mr *MockServiceMockRecorder) ApproveLoan(id, validatorID, proofURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ApproveLoan", reflect.TypeOf((*MockService)(nil).ApproveLoan), id, validatorID, proofURL)
}

// CreateLoan mocks base method.
func (m *MockService) CreateLoan(borrowerID string, principal, rate, roi float64) (*loan.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLoan", borrowerID, principal, rate, roi)
	ret0, _ := ret[0].(*loan.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateLoan indicates an expected call of CreateLoan.
func (mr *MockServiceMockRecorder) CreateLoan(borrowerID, principal, rate, roi interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLoan", reflect.TypeOf((*MockService)(nil).CreateLoan), borrowerID, principal, rate, roi)
}

// DisburseLoan mocks base method.
func (m *MockService) DisburseLoan(id, fieldOfficerID, signedAgreement string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DisburseLoan", id, fieldOfficerID, signedAgreement)
	ret0, _ := ret[0].(error)
	return ret0
}

// DisburseLoan indicates an expected call of DisburseLoan.
func (mr *MockServiceMockRecorder) DisburseLoan(id, fieldOfficerID, signedAgreement interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisburseLoan", reflect.TypeOf((*MockService)(nil).DisburseLoan), id, fieldOfficerID, signedAgreement)
}

// GenerateAgreementLetter mocks base method.
func (m *MockService) GenerateAgreementLetter(id, letterURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateAgreementLetter", id, letterURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateAgreementLetter indicates an expected call of GenerateAgreementLetter.
func (mr *MockServiceMockRecorder) GenerateAgreementLetter(id, letterURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateAgreementLetter", reflect.TypeOf((*MockService)(nil).GenerateAgreementLetter), id, letterURL)
}

// GetLoan mocks base method.
func (m *MockService) GetLoan(id string) (*loan.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoan", id)
	ret0, _ := ret[0].(*loan.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoan indicates an expected call of GetLoan.
func (mr *MockServiceMockRecorder) GetLoan(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoan", reflect.TypeOf((*MockService)(nil).GetLoan), id)
}

// GetLoans mocks base method.
func (m *MockService) GetLoans(page, limit int) ([]*loan.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoans", page, limit)
	ret0, _ := ret[0].([]*loan.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoans indicates an expected call of GetLoans.
func (mr *MockServiceMockRecorder) GetLoans(page, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoans", reflect.TypeOf((*MockService)(nil).GetLoans), page, limit)
}

// GetLoansByBorrower mocks base method.
func (m *MockService) GetLoansByBorrower(borrowerID string) ([]*loan.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoansByBorrower", borrowerID)
	ret0, _ := ret[0].([]*loan.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoansByBorrower indicates an expected call of GetLoansByBorrower.
func (mr *MockServiceMockRecorder) GetLoansByBorrower(borrowerID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoansByBorrower", reflect.TypeOf((*MockService)(nil).GetLoansByBorrower), borrowerID)
}

// GetLoansByState mocks base method.
func (m *MockService) GetLoansByState(state loan.LoanState) ([]*loan.Loan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoansByState", state)
	ret0, _ := ret[0].([]*loan.Loan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoansByState indicates an expected call of GetLoansByState.
func (mr *MockServiceMockRecorder) GetLoansByState(state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoansByState", reflect.TypeOf((*MockService)(nil).GetLoansByState), state)
}
