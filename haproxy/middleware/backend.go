package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all backends
func (c *Client) GetBackends() (*models.GetBackends, error) {
	url := c.base_url + "/services/haproxy/configuration/backends/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBackends{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// return single backend
func (c *Client) GetBackend(backendName string) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backendName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBackend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateBackend(transactionId string, backend models.Backend) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(backend)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Backend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateBackend(transactionId string, backendName string, backend models.Backend) (*models.Backend, error) {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backendName + "?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(backend)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Backend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteBackend(transactionId string, backendName string) error {
	url := c.base_url + "/services/haproxy/configuration/backends/" + backendName + "?transaction_id=" + transactionId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
