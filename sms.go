package gotiniyo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// SmsResponse is returned after a text/sms message is posted to Tiniyo
type SmsResponse struct {
	Sid         string  `json:"sid"`
	DateCreated string  `json:"date_created"`
	DateUpdate  string  `json:"date_updated"`
	DateSent    string  `json:"date_sent"`
	AccountSid  string  `json:"account_sid"`
	To          string  `json:"to"`
	From        string  `json:"from"`
	MediaUrl    string  `json:"media_url"`
	Body        string  `json:"body"`
	Status      string  `json:"status"`
	Direction   string  `json:"direction"`
	ApiVersion  string  `json:"api_version"`
	Price       *string `json:"price,omitempty"`
	Url         string  `json:"uri"`
}

// DateCreatedAsTime returns SmsResponse.DateCreated as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateCreatedAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateCreated)
}

// DateUpdateAsTime returns SmsResponse.DateUpdate as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateUpdateAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateUpdate)
}

// DateSentAsTime returns SmsResponse.DateSent as a time.Time object
// instead of a string.
func (sms *SmsResponse) DateSentAsTime() (time.Time, error) {
	return time.Parse(time.RFC1123Z, sms.DateSent)
}


// SendSMS uses Tiniyo to send a text message.
// See http://www.tiniyo.com/docs/api/rest/sending-sms for more information.
func (tiniyo *Tiniyo) SendSMS(from, to, text, callbackUrl, applicationSid string) (smsResponse *SmsResponse, exception *Exception, err error) {
	payload := initPayloadValues(to, from, text, callbackUrl)
	smsResponse, exception, err = tiniyo.sendMessage(payload)
	return
}

// Core method to send message
func (tiniyo *Tiniyo) sendMessage(payload *strings.Reader) (smsResponse *SmsResponse, exception *Exception, err error) {
	tiniyoUrl := tiniyo.BaseUrl + "Accounts/" + tiniyo.AuthID + "/Messages"

	res, err := tiniyo.post(payload, tiniyoUrl)
	if err != nil {
		return smsResponse, exception, err
	}
	defer res.Body.Close()

	responseBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return smsResponse, exception, err
	}

	if res.StatusCode != http.StatusCreated {
		exception = new(Exception)
		err = json.Unmarshal(responseBody, exception)

		// We aren't checking the error because we don't actually care.
		// It's going to be passed to the client either way.
		return smsResponse, exception, err
	}

	smsResponse = new(SmsResponse)
	err = json.Unmarshal(responseBody, smsResponse)
	return smsResponse, exception, err
}

// Payload initialization
func initPayloadValues(to, from, text,callbackUrl string) *strings.Reader {
	payload := strings.NewReader("{\n  \"dst\": \""+ to +"\",\n  \"src\": \"TINIQE\",\n  \"method\": \"post\",\n  \"text\": \""+ text +"\",\n  \"url\": \""+ callbackUrl +"\"\n}")
	return payload
}
