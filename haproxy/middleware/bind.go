package middleware

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"terraform-provider-haproxy-pf/haproxy/models"
)

// return all binds
func (c *Client) GetBinds() (*models.GetBinds, error) {
	url := c.base_url + "/services/haproxy/configuration/binds/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBinds{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// return single bind
func (c *Client) GetBind(bindName string, parentName string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?parent_type=frontend&parent_name=%s&frontend=%s", c.base_url, bindName, parentName, parentName)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res := models.GetBind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res.Data, nil
}

func (c *Client) CreateBind(transactionId string, bind models.Bind, parentName string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/?parent_type=frontend&parent_name=%s&transaction_id=%s&frontend=%s", c.base_url, parentName, transactionId, parentName)
	bodyStr, _ := json.Marshal(bind)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Bind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) UpdateBind(transactionId string, bind models.Bind, parentName string) (*models.Bind, error) {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?transaction_id=%s&parent_type=frontend&parent_name=%s&frontend=%s", c.base_url, bind.Name, transactionId, parentName, parentName)
	bodyStr, _ := json.Marshal(bind)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyStr))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res := models.Bind{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) DeleteBind(transactionId string, bindName string, parentName string) error {
	url := fmt.Sprintf("%s/services/haproxy/configuration/binds/%s?transaction_id=%s&parent_type=frontend&parent_name=%s&frontend=%s", c.base_url, bindName, transactionId, parentName, parentName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	if err := c.sendRequest(req, nil); err != nil {
		return err
	}

	return nil
}
