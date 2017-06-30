package goarmorconfigs

type Configer interface {
	ServerDebuggingLevel() uint64
}

func (c *Config) ServerDebuggingLevel() uint64 {
	return c.Server.DebuggingLevel
}
