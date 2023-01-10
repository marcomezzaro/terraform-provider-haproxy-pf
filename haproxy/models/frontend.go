package models

type GetFrontend struct {
	Version int      `json:"_version"`
	Data    Frontend `json:"data"`
}

type Frontend struct {
	HTTPConnectionMode string `json:"http_connection_mode"`
	Maxconn            int64  `json:"maxconn"`
	Mode               string `json:"mode"`
	Name               string `json:"name"`
	DefaultBackend     string `json:"default_backend"`
}

type GetFrontends struct {
	Version int `json:"_version"`
	Data    []struct {
		HTTPConnectionMode string `json:"http_connection_mode"`
		Maxconn            int64  `json:"maxconn"`
		Mode               string `json:"mode"`
		Name               string `json:"name"`
		DefaultBackend     string `json:"default_backend"`
	} `json:"data"`
}
