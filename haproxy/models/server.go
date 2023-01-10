package models

type GetServer struct {
	Version int  `json:"_version"`
	Data    Server `json:"data"`
}

type Server struct {
	Address string `json:"address"`
	Name    string `json:"name"`
	Port    int64 `json:"port"`
	Check   string `json:"check"`
}

type GetServers struct {
	Version int `json:"_version"`
	Data    []struct {
		Address string `json:"address"`
		Name    string `json:"name"`
		Port    int64 `json:"port"`
		Check   string `json:"check"`
	} `json:"data"`
}
