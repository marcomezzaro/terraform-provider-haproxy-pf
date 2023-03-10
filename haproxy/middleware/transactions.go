package middleware

import (
	"net/http"
	"strconv"
	"terraform-provider-haproxy-pf/haproxy/models"
)

func (c *Client) CreateTransaction(version int) (*models.Transaction, error) {
	url := c.base_url + "/services/haproxy/transactions?version=" + strconv.Itoa(version)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.Transaction{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) CommitTransaction(transactionId string) (*models.Transaction, error) {
	url := c.base_url + "/services/haproxy/transactions/" + transactionId
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.Transaction{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
