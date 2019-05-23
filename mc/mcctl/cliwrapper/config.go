package cliwrapper

import "github.com/mobiledgex/edge-cloud-infra/mc/ormapi"

func (s *Client) UpdateConfig(uri, token string, config map[string]interface{}) (int, error) {
	args := []string{"config", "update"}
	return s.runObjs(uri, token, args, config, nil)
}

func (s *Client) ShowConfig(uri, token string) (*ormapi.Config, int, error) {
	args := []string{"config", "show"}
	config := ormapi.Config{}
	st, err := s.runObjs(uri, token, args, nil, &config)
	return &config, st, err
}