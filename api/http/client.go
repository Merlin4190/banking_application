package http

import (
	"banking_application/api/http/models"
	"banking_application/api/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClientAction interface {
	PostPayment(url string, request models.PaymentRequest) (*models.PaymentResponse, error)
	GetPayment(url string) (*models.PaymentResponse, error)
}

type Client struct {
	Client *Client
}

const baseURL = "/third-party/payments"
const env = "test"

func NewClient() *Client {
	return &Client{}
}

func (s *Client) PostPayment(url string, request models.PaymentRequest) (*models.PaymentResponse, error) {
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	if env == "test" {
		mockResponse := models.PaymentResponse{
			AccountID: request.AccountID,
			Reference: util.GenerateUniqueAlphaNumeric(15),
			Amount:    request.Amount,
		}
		return &mockResponse, nil
	}

	resp, err := http.Post(baseURL+url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentResponse models.PaymentResponse
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}

func (s *Client) GetPayment(url string) (*models.PaymentResponse, error) {
	resp, err := http.Get(baseURL + url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var paymentResponse models.PaymentResponse
	err = json.Unmarshal(body, &paymentResponse)
	if err != nil {
		return nil, err
	}

	return &paymentResponse, nil
}
