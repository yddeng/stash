package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hqpko/hutils"
	"initialthree/cluster"
	"initialthree/cluster/addr"
	"path"
	"strconv"
)

type Config struct {
	Common     *Common
	Log        *Log
	Dir        []*Dir
	Login      []*Login
	Game       []*Game
	Map        []*Map
	World      []*World
	Team       []*Team
	Gate       []*Gate
	WebService []*WebService
	Rank       []*Rank
}

func LoadConfig(path string) (*Config, error) {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)

	if err != nil {
		return nil, err
	}

	if len(conf.Common.ServerGroups) == 0 {
		conf.Common.ServerGroups = []int32{1}
	}

	return conf, err
}

func GetConfig(conf *Config, cxt, index string) (ret interface{}, err error) {

	var idx int
	idx, err = strconv.Atoi(index)
	if err != nil {
		return
	}
	switch cxt {
	case "Gate":
		ret = conf.Gate[idx]
	case "Game":
		ret = conf.Game[idx]
	case "Login":
		ret = conf.Login[idx]
	case "Dir":
		ret = conf.Dir[idx]
	case "Map":
		ret = conf.Map[idx]
	case "World":
		ret = conf.World[idx]
	case "Team":
		ret = conf.Team[idx]
	case "WebService":
		ret = conf.WebService[idx]
	case "Rank":
		ret = conf.Rank[idx]
	default:
		err = fmt.Errorf("invalid context:%s", cxt)
		return
	}
	return
}

func MustStartCluster(conf *Config, serverAddr addr.Addr, srvType uint32, uniLocker cluster.UniLocker, export ...bool) {
	if serverAddr.Logic.Type() != srvType {
		panic(fmt.Errorf("invalid server type:%s should be:%d", serverAddr.Logic.String(), srvType))
		return
	}
	hutils.Must(nil, cluster.Start(conf.Common.CenterAddr, serverAddr, uniLocker, export...))
}

type Addr struct {
	LogicAddr   string `toml:"logicAddr"`
	ClusterAddr string `toml:"clusterAddr"`
}

func (this *Addr) MakeAddr() (addr.Addr, error) {
	return addr.MakeAddr(this.LogicAddr, this.ClusterAddr)
}

type Common struct {
	CenterAddr      []string `toml:"centerAddr"`
	DbConfig        []string `toml:"dbConfig"`
	ServerGroups    []int32  `toml:"serverGroups"`
	CfgPathRoot     string   `toml:"cfgPathRoot"`
	ExcelPath       string   `toml:"excelPath"`
	WordsFilterPath string   `toml:"wordsFilterPath"`
}

func (c *Common) GetCfgFilePath(filepath string) string {
	if path.IsAbs(filepath) {
		return filepath
	}
	return path.Join(c.CfgPathRoot, filepath)
}

func (c *Common) GetExcelPath() string {
	return c.GetCfgFilePath(c.ExcelPath)
}

func (c *Common) GetWordsFilterPath() string {
	return c.GetCfgFilePath(c.WordsFilterPath)
}

type Log struct {
	Path            string `toml:"path"`
	Level           string `toml:"level"`
	MaxSize         int    `toml:"maxSize"`
	MaxAge          int    `toml:"maxAge"`
	MaxBackups      int    `toml:"maxBackups"`
	EnableLogStdout bool   `toml:"enableLogStdout"`
}

type Dir struct {
	*Addr
	BBQueryMax   int32  `toml:"dbQueryMax"`
	ExternalAddr string `toml:"externalAddr"`
	Servers      []struct {
		ServerId   int32  `toml:"serverID"`
		ServerName string `toml:"serverName"`
	} `toml:"Servers"`
}

type Login struct {
	*Addr
	ExternalAddr string `toml:"externalAddr"`
}

type Game struct {
	*Addr
}

type World struct {
	*Addr
}

type Map struct {
	*Addr
	WorldAddr string `toml:"worldAddr"`
}

type Team struct {
	*Addr
}

type Gate struct {
	*Addr
	ExternalAddr string `toml:"externalAddr"`
}

type WebService struct {
	*Addr
	WebAddress  string `toml:"webAddress"`
	AccessToken string `toml:"accessToken"`
}

type Rank struct {
	*Addr

	SqlType    string `toml:"SqlType"`
	DbHost     string `toml:"DbHost"`
	DbPort     int    `toml:"DbPort"`
	DbUser     string `toml:"DbUser"`
	DbPassword string `toml:"DbPassword"`
	DbDataBase string `toml:"DbDataBase"`
}
