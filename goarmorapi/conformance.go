package goarmorapi

type Conformance struct {
	Name   string             `json:"Name"`
	Server *ConformanceServer `json:"Server"`
	API    *ConformanceClient `json:"Api"`
}

type ConformanceServer struct {
	Type         string `json:"Type"`
	ID           uint64 `json:"Id"`
	Version      uint64 `json:"Version"`
	Architecture string `json:"Architecture"`
}

type ConformanceClient struct {
	Version string `json:"Version"`
}

func NewConformance(
	apiVersion string,
	infrastructureVersion uint64,
	serverTitle, serverName, serverArchitecture string,
	serverID uint64) *Conformance {
	return &Conformance{
		Name: serverTitle,
		API:  &ConformanceClient{Version: apiVersion},
		Server: &ConformanceServer{
			Type:         serverName,
			ID:           serverID,
			Version:      infrastructureVersion,
			Architecture: serverArchitecture}}
}
