package sms

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type BulkSMSLab struct {
	APIKey   string
	UserID   string
	Password string
	SenderID string
}

func NewBulkSMSLab(apiKey, userID, password, senderID string) *BulkSMSLab {
	return &BulkSMSLab{
		APIKey:   apiKey,
		UserID:   userID,
		Password: password,
		SenderID: senderID,
	}
}

func (b *BulkSMSLab) SendOTP(mobile string, otp string) error {
	form := url.Values{}
	form.Set("userid", b.UserID)
	form.Set("password", b.Password)
	form.Set("sendMethod", "quick")
	form.Set("mobile", "91"+mobile)
	form.Set(
		"msg",
		fmt.Sprintf("Your OTP is %s for Phone Verification.OTPSTE", otp), // EXACT TEMPLATE
	)
	form.Set("senderid", b.SenderID)
	form.Set("msgType", "text")
	form.Set("duplicatecheck", "true")
	form.Set("output", "json")
	form.Set("dltEntityId", "1701159170147084368")
	form.Set("dltTemplateId", "1707165701540733056")

	req, err := http.NewRequest(
		"POST",
		"https://sms.bulksmslab.com/SMSApi/send",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("apikey", b.APIKey)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Timeout: 10 * time.Second, // VERY IMPORTANT
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("sms http error: %d body=%s", resp.StatusCode, string(body))
	}

	var res map[string]interface{}
	if err := json.Unmarshal(body, &res); err != nil {
		return fmt.Errorf("invalid sms response: %s", string(body))
	}

	// CHECK PROVIDER-SPECIFIC SUCCESS FLAG
	if res["status"] != "success" {
		return fmt.Errorf("sms rejected: %v", res)
	}

	fmt.Println(res)

	return nil
}
