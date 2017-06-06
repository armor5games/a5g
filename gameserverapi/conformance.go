package gameserverapi

type Conformance struct {
	Name   string
	Server *ConformanceServer
	API    *ConformanceClient
}

type ConformanceServer struct {
	Type         string
	ID           int
	Version      int
	Architecture string
}

type ConformanceClient struct {
	Version string
}

func NewConformance(
	apiVersion string,
	infrastructureVersion int,
	serverTitle, serverName, serverArchitecture string,
	serverID int) *Conformance {
	return &Conformance{
		Name: serverTitle,
		API:  &ConformanceClient{Version: apiVersion},
		Server: &ConformanceServer{
			Type:         serverName,
			ID:           serverID,
			Version:      infrastructureVersion,
			Architecture: serverArchitecture}}
}
