package models

type GetBind struct {
	Version int  `json:"_version"`
	Data    Bind `json:"data"`
}

type Bind struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Port    int64 `json:"port"`
}

type GetBinds struct {
	Version int `json:"_version"`
	Data    []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Port    int64 `json:"port"`
	} `json:"data"`
}
