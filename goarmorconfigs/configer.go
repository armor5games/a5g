package goarmorconfigs

type Configer interface {
	ServerDebuggingLevel() int
}

func (c *Config) ServerDebuggingLevel() int {
	return c.Server.DebuggingLevel
}
