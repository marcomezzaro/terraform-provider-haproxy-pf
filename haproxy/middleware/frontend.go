package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all frontends
func (c *Client) GetFrontends() (*models.GetFrontends, error) {
	url := c.base_url + "/services/haproxy/configuration/frontends/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetFrontends{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// return single frontend
func (c *Client) GetFrontend(frontendName string) (*models.Frontend, error) {
	url := c.base_url + "/services/haproxy/configuration/frontends/" + frontendName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetFrontend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateFrontend(transactionId string, frontend models.Frontend) (*models.Frontend, error) {
	url := c.base_url + "/services/haproxy/configuration/frontends?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(frontend)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Frontend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateFrontend(transactionId string, frontendName string, frontend models.Frontend) (*models.Frontend, error) {
	url := c.base_url + "/services/haproxy/configuration/frontends/" + frontendName + "?transaction_id=" + transactionId
	bodyStr, _ := json.Marshal(frontend)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Frontend{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteFrontend(transactionId string, frontendName string) error {
	url := c.base_url + "/services/haproxy/configuration/frontends/" + frontendName + "?transaction_id=" + transactionId
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
