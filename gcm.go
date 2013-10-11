package gcm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

const (
	GCMSendEndpoint    = "https://android.googleapis.com/gcm/send"
	BackoffInitalDelay = time.Second
)

// Sender
type Sender struct {
	// Client specifies the http.Client used for sending
	Client *http.Client

	// Authorization key for sending messages to the GCM endpoint
	Key string
}

// Message represents the data that will be send to the server
type Message struct {
	CollapseKey           string                 `json:"collapse_key,omitempty"`
	DelayWhileIdle        bool                   `json:"delay_while_idle,omitempty"`
	TimeToLive            int                    `json:"time_to_live,omitempty"`
	RestrictedPackageName string                 `json:"restricted_package_name,omitempty"`
	Data                  map[string]interface{} `json:"data"`
	RegistrationIds       []string               `json:"registration_ids"`
}

// MulticastResult is the response of a Message with multiple registration IDs
type MulticastResult struct {
	Success      int      `json:"success"`
	Failure      int      `json:"failure"`
	CanonicalIds int      `json:"canonical_ids"`
	MulticastId  int64    `json:"multicast_id"`
	Results      []Result `json:"results"`
	//RetryMulticastIds []int64  `json:"collapse_key"`
}

// Result is the response of a Message with a single registration ID
type Result struct {
	MessageId               string `json:"message_id"`
	CanonicalRegistrationId string `json:"registration_id"`
	ErrorCode               string `json:"error"`
}

// NewSender creates a Sender that uses the default HTTP client
func NewSender(key string) *Sender {
	return &Sender{
		Client: &http.Client{},
		Key:    key,
	}
}

// NewMessage creates a empty Message
func NewMessage(registrationIds ...string) *Message {
	return &Message{
		Data:            make(map[string]interface{}),
		RegistrationIds: registrationIds,
	}
}

// Add adds a new key value pair to the message
func (m *Message) Add(key string, value interface{}) {
	if m.Data == nil {
		m.Data = make(map[string]interface{})
	}
	m.Data[key] = value
}

// Send tries to send the message with retries if the sending faild
func (s *Sender) Send(msg *Message, retries int) (*MulticastResult, error) {
	for i := 0; i < retries; i++ {
		if r, err := s.SendNoRetry(msg); err == nil {
			return r, nil
		}
		time.Sleep(time.Duration(i*2) * BackoffInitalDelay)
	}
	return nil, errors.New("Could not send message after " + strconv.Itoa(retries) + " attempts")
}

// Send tries to send the message without retries
func (s *Sender) SendNoRetry(msg *Message) (*MulticastResult, error) {
	if len(msg.RegistrationIds) == 0 {
		return nil, errors.New("RegistrationIds cannot be empty")
	}
	b, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	res, err := s.post(GCMSendEndpoint, "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	return res, err
}

func (s *Sender) post(url string, bodyType string, body io.Reader) (*MulticastResult, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", bodyType)
	req.Header.Set("Authorization", "key="+s.Key)

	res, err := s.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var v MulticastResult
	err = json.Unmarshal(b, &v)
	if err != nil {
		return nil, errors.New("Error parsing JSON response: " + string(b))
	}
	return &v, nil
}
