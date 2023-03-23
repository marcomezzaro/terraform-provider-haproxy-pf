package models

type GetResolvers struct {
	Version int `json:"_version"`
	Data    []struct {
		Name   string `json:"name"`
	} `json:"data"`
}
