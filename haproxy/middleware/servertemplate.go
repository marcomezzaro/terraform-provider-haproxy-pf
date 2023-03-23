package middleware

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all server_templates
func (c *Client) GetServerTemplates(parentName string) (*models.GetServerTemplates, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/server_templates/?backend=%s", c.base_url, parentName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetServerTemplates{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// return single server_templates
func (c *Client) GetServerTemplate(serverTemplateName string, parentName string) (*models.ServerTemplate, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/server_templates/%s?backend=%s", c.base_url, serverTemplateName, parentName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetServerTemplate{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateServerTemplate(transactionId string, serverTemplate models.ServerTemplate, parentName string) (*models.ServerTemplate, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/server_templates/?transaction_id=%s&backend=%s", c.base_url, transactionId, parentName)
	bodyStr, _ := json.Marshal(serverTemplate)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.ServerTemplate{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateServerTemplate(transactionId string, serverTemplate models.ServerTemplate, parentName string) (*models.ServerTemplate, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/server_templates/%s?transaction_id=%s&backend=%s", c.base_url, serverTemplate.Prefix, transactionId, parentName)
	bodyStr, _ := json.Marshal(serverTemplate)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.ServerTemplate{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteServerTemplate(transactionId string, serverTemplateName string, parentName string) error {
	url := fmt.Sprintf("%s/services/haproxy/configuration/server_templates/%s?transaction_id=%s&parent_type=backend&parent_name=%s&backend=%s", c.base_url, serverTemplateName, transactionId, parentName, parentName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
