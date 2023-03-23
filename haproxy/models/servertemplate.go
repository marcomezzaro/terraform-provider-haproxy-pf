package models

type GetServerTemplate struct {
	Version int  `json:"_version"`
	Data    ServerTemplate `json:"data"`
}

type ServerTemplate struct {
	Fqdn    string `json:"fqdn"`
	Num_or_range    string `json:"num_or_range"`
	Port   int64 `json:"port"`
	Prefix string `json:"prefix"`
	Check string `json:"check"`
	Resolvers string `json:"resolvers"`
}

type GetServerTemplates struct {
	Version int `json:"_version"`
	Data    []struct {
		Fqdn    string `json:"fqdn"`
		Num_or_range    string `json:"num_or_range"`
		Port   int64 `json:"port"`
		Prefix string `json:"prefix"`
		Check string `json:"check"`
		Resolvers string `json:"resolvers"`
	} `json:"data"`
}
