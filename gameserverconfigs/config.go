package gameserverconfigs

import (
	"errors"
	"fmt"
	"strconv"

	"gopkg.in/gcfg.v1"
)

const (
	TypeLogin ServerType = iota
	TypeShard

	DBRead DBConfigType = iota
	DBWrite
	DBReadStatic
)

type ServerType int
type DBConfigType int

type Config struct {
	PathToConfig string

	Server struct {
		Type    ServerType
		ID      int
		Version string

		Port          string
		LogPath       string
		DebuggingMode int

		ServerSecretKey string

		Bugsnag string
	}

	LoginServer struct {
		RDBName string
		RDBUser string
		RDBHost string
		RDBPass string
		RDBPort string

		WDBName string
		WDBUser string
		WDBHost string
		WDBPass string
		WDBPort string

		RStaticDBName string
		RStaticDBUser string
		RStaticDBHost string
		RStaticDBPass string
		RStaticDBPort string
	}

	ShardServer struct {
		RDBName string
		RDBUser string
		RDBHost string
		RDBPass string
		RDBPort string

		WDBName string
		WDBUser string
		WDBHost string
		WDBPass string
		WDBPort string
	}
}

func New(serverType ServerType, pathToConfig string) (*Config, error) {
	c := new(Config)
	err := gcfg.ReadFileInto(c, pathToConfig)
	if err != nil {
		return nil, err
	}

	if serverType != TypeLogin && serverType != TypeShard {
		return nil, errors.New("unknown config type")
	}
	c.Server.Type = serverType

	c.PathToConfig = pathToConfig

	return c, nil
}

// DBPass()
// DBHost()
// DBPort()
// DBName()

func (c *Config) DBConfig(t DBConfigType) (
	*struct{ DBUser, DBPass, DBHost, DBPort, DBName string }, error) {
	if c.Server.Type == TypeLogin {
		switch t {
		default:
			return nil, fmt.Errorf("unknown login server db type: %s", string(t))

		case DBRead:
			return &struct{ DBUser, DBPass, DBHost, DBPort, DBName string }{
				DBName: c.LoginServer.RDBName,
				DBUser: c.LoginServer.RDBUser,
				DBHost: c.LoginServer.RDBHost,
				DBPass: c.LoginServer.RDBPass,
				DBPort: c.LoginServer.RDBPort}, nil

		case DBWrite:
			return &struct{ DBUser, DBPass, DBHost, DBPort, DBName string }{
				DBName: c.LoginServer.WDBName,
				DBUser: c.LoginServer.WDBUser,
				DBHost: c.LoginServer.WDBHost,
				DBPass: c.LoginServer.WDBPass,
				DBPort: c.LoginServer.WDBPort}, nil

		case DBReadStatic:
			return &struct{ DBUser, DBPass, DBHost, DBPort, DBName string }{
				DBName: c.LoginServer.RStaticDBName,
				DBUser: c.LoginServer.RStaticDBUser,
				DBHost: c.LoginServer.RStaticDBHost,
				DBPass: c.LoginServer.RStaticDBPass,
				DBPort: c.LoginServer.RStaticDBPort}, nil
		}

	} else if c.Server.Type == TypeShard {
		switch t {
		default:
			return nil, fmt.Errorf("unknown shard server db type: %s", string(t))

		case DBRead:
			return &struct{ DBUser, DBPass, DBHost, DBPort, DBName string }{
				DBName: c.ShardServer.RDBName,
				DBUser: c.ShardServer.RDBUser,
				DBHost: c.ShardServer.RDBHost,
				DBPass: c.ShardServer.RDBPass,
				DBPort: c.ShardServer.RDBPort}, nil

		case DBWrite:
			return &struct{ DBUser, DBPass, DBHost, DBPort, DBName string }{
				DBName: c.ShardServer.WDBName,
				DBUser: c.ShardServer.WDBUser,
				DBHost: c.ShardServer.WDBHost,
				DBPass: c.ShardServer.WDBPass,
				DBPort: c.ShardServer.WDBPort}, nil
		}
	}

	return nil, fmt.Errorf("unknown server type: %s", string(c.Server.Type))
}

func (c *Config) Type() ServerType {
	return c.Server.Type
}

func (c *Config) Version() (int64, error) {
	return strconv.ParseInt(c.Server.Version, 10, 0)
}
