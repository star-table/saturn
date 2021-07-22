package welink

import "github.com/galaxy-book/saturn/proxy"

type welinkProxy struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	Ticket       proxy.Ticket
}

func NewWelinkProxy(ClientID, ClientSecret string, ticket proxy.Ticket) *welinkProxy {
	return &welinkProxy{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Ticket:       ticket,
	}
}
