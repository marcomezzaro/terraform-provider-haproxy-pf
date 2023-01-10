package models

type GetBackend struct {
	Version int     `json:"_version"`
	Data    Backend `json:"data"`
}

type Balance struct {
	Algorithm string `json:"algorithm"`
}

type Backend struct {
	Balance Balance `json:"balance"`
	Mode    string  `json:"mode"`
	Name    string  `json:"name"`
}

type GetBackends struct {
	Version int `json:"_version"`
	Data    []struct {
		Balance struct {
			Algorithm string `json:"algorithm"`
		} `json:"balance"`
		Mode string `json:"mode"`
		Name string `json:"name"`
	} `json:"data"`
}
