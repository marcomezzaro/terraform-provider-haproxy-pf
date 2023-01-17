package middleware

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all servers
func (c *Client) GetServers() (*models.GetServers, error) {
	url := c.base_url + "/services/haproxy/configuration/servers/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetServers{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// return single server
func (c *Client) GetServer(serverName string, parentName string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?parent_type=backend&parent_name=%s&backend=%s", c.base_url, serverName, parentName, parentName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetServer{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateServer(transactionId string, server models.Server, parentName string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/?parent_type=backend&parent_name=%s&transaction_id=%s&backend=%s", c.base_url, parentName, transactionId, parentName)
	bodyStr, _ := json.Marshal(server)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Server{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateServer(transactionId string, server models.Server, parentName string) (*models.Server, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?transaction_id=%s&parent_type=backend&parent_name=%s&backend=%s", c.base_url, server.Name, transactionId, parentName, parentName)
	bodyStr, _ := json.Marshal(server)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Server{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteServer(transactionId string, serverName string, parentName string) error {
	url := fmt.Sprintf("%s/services/haproxy/configuration/servers/%s?transaction_id=%s&parent_type=backend&parent_name=%s&backend=%s", c.base_url, serverName, transactionId, parentName, parentName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
