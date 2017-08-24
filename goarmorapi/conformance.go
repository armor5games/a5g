package goarmorapi

type Conformance struct {
	Name   string             `json:"name"`
	Server *ConformanceServer `json:"server"`
	API    *ConformanceClient `json:"api"`
}

type ConformanceServer struct {
	Type         string `json:"type"`
	ID           uint64 `json:"id"`
	Version      uint64 `json:"version"`
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
	startedAt string) *Conformance {
	return &Conformance{
		Name: serverTitle,
		API:  &ConformanceClient{Version: apiVersion},
		Server: &ConformanceServer{
			Type:         serverName,
			ID:           serverID,
			Version:      infrastructureVersion,
			Architecture: serverArchitecture,
			StartedAt:    startedAt}}
}
