package middleware

import (
	"fmt"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all resolvers
func (c *Client) GetResolvers() (*models.GetResolvers, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/resolvers", c.base_url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetResolvers{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
