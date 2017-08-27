package goarmorconfigs

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"gopkg.in/gcfg.v1"
)

const (
	TypeLogin ServerType = iota
	TypeShard

	DBRead DBConfigType = iota
	DBWrite
	DBReadStatic

	MongoDBClientMetrics MongoDBConfigType = iota
)

type (
	ServerType        int
	DBConfigType      int
	MongoDBConfigType int
)

type Config struct {
	PathToConfig string

	Server struct {
		Type    ServerType
		ID      uint64
		Version uint64
		URL     string

		ListenAddress  string
		LogPath        string
		DebuggingLevel int

		ServerSecretKey string
		ClientSecretKey string

		Bugsnag string

		APITimeoutSeconds uint64
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

		MongoDBHost string
		MongoDBName string
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

	Google struct {
		// AppLicenseKey
		// <https://developer.android.com/google/play/billing/billing_integrate.html#billing-security>.
		AppLicenseKey string
	}

	Static struct {
		Directory string
	}
}

func New(
	serverType ServerType, serverVersion uint64, pathToConfig string) (*Config, error) {
	c := new(Config)
	err := gcfg.ReadFileInto(c, pathToConfig)
	if err != nil {
		return nil, err
	}

	if serverType != TypeLogin && serverType != TypeShard {
		return nil, errors.New("unknown config type")
	}
	c.Server.Type = serverType

	if serverVersion == 0 {
		return nil, errors.New("server version undefined")
	}
	c.Server.Version = serverVersion

	if c.Server.URL == "" {
		return nil, errors.New("server url undefined")
	}

	if _, err = url.Parse(c.Server.URL); err != nil {
		return nil, fmt.Errorf("url.Parse fn error: %s", err.Error())
	}

	c.PathToConfig = pathToConfig

	if c.Server.APITimeoutSeconds == 0 {
		return nil, errors.New("undefined api timeout")
	}

	return c, nil
}

// FIXME: Remove me!
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

func (c *Config) MongoDBConfig(t MongoDBConfigType) (
	*struct{ MongoDBHost, MongoDBName string }, error) {
	if c.Server.Type == TypeLogin {
		switch t {
		default:
			return nil, errors.Errorf("unknown login server mongodb config type: %s",
				string(t))

		case MongoDBClientMetrics:
			return &struct{ MongoDBHost, MongoDBName string }{
				MongoDBHost: c.LoginServer.MongoDBHost,
				MongoDBName: c.LoginServer.MongoDBName}, nil
		}
	}

	return nil, fmt.Errorf("unknown server type: %s", string(c.Server.Type))
}

func (c *Config) SEConfig() (*struct{ UserSecure string }, error) {
	return &struct{ UserSecure string }{
		UserSecure: c.Server.ClientSecretKey,
	}, nil
}

func (c *Config) StaticsDir() (string, error) {
	return c.Static.Directory, nil
}
