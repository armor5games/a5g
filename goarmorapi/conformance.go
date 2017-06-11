package goarmorapi

type Conformance struct {
	Name   string
	Server *ConformanceServer
	API    *ConformanceClient
}

type ConformanceServer struct {
	Type         string
	ID           uint64
	Version      uint64
	Architecture string
}

type ConformanceClient struct {
	Version string
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
