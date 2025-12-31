package sms

type Sender interface {
	SendOTP(mobile string, message string) error
}
