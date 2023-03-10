package middleware

import (
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

func (c *Client) GetConfiguration() (*models.Configuration, error) {
	url := c.base_url + "/services/haproxy/configuration/raw"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.Configuration{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	if res.Version == 0 {
		res.Version = 1
	}

	return &res, nil
}
