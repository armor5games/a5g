package goarmorapi

import "github.com/pkg/errors"

type Conformance struct {
	Name   string             `json:"name"`
	Server *ConformanceServer `json:"server"`
	API    *ConformanceClient `json:"api"`
}

type ConformanceServer struct {
	Type         string `json:"type"`
	ID           uint64 `json:"id"`
	Version      uint64 `json:"version"`
	Sharded      bool   `json:"sharded,omitempty"`
	ShardsCount  uint64 `json:"shardsCount,omitempty"`
	StartedAt    string `json:"startedAt,omitempty"`
	Architecture string `json:"architecture"`
}

type ConformanceClient struct {
	Version string `json:"version"`
}

func NewConformance(
	apiVersion string,
	infrastructureVersion uint64,
	serverTitle, serverName, serverArchitecture string,
	serverID uint64,
	shardsCount int,
	startedAt string) (*Conformance, error) {
	if shardsCount < 0 {
		return nil, errors.New("unexpected shard servers count")
	}

	v := &Conformance{
		Name: serverTitle,
		API:  &ConformanceClient{Version: apiVersion},
		Server: &ConformanceServer{
			Type:         serverName,
			ID:           serverID,
			Version:      infrastructureVersion,
			Sharded:      shardsCount > 0,
			ShardsCount:  uint64(shardsCount),
			Architecture: serverArchitecture,
			StartedAt:    startedAt}}

	return v, nil
}
