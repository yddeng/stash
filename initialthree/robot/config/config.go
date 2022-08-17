package config

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/BurntSushi/toml"
)

type FloatingDuration struct {
	// 单位(ms)
	Unit uint `toml:"Unit"`

	// min
	Min uint `toml:"Min"`

	// max
	Max uint `toml:"max"`
}

func (d *FloatingDuration) GenerateDuration() time.Duration {
	if d.Unit == 0 {
		return 0
	}

	var interval uint
	if d.Min == d.Max {
		interval = d.Min
	} else {
		interval = d.Min + uint(rand.Intn(int(d.Max-d.Min)+1))
	}

	return time.Duration(interval) * time.Duration(d.Unit) * time.Millisecond
}

type Config struct {
	// address of login service
	Service string `toml:"Service"`

	// server id
	ServerID uint `toml:"ServerID"`

	// tick cycle, ms.
	TickCycle uint `toml:"TickCycle"`

	Robot struct {
		// robot count
		Count uint `toml:"Count"`

		// the prefix of robot user-id
		UserIDPrefix string `toml:"UserIDPrefix"`

		// the initial number of robot user-id
		InitialNo int `toml:"InitialNo"`
	} `toml:"Robot"`

	Resource struct {
		// directory of excel files
		ExcelPath string `toml:"ExcelPath"`

		// directory of quest files
		QuestPath string `toml:"QuestPath"`

		// behavior config file
		BehaviorConfigPath string `toml:"BehaviorConfigPath"`

		// behavior tree name
		BehaviorTree string `toml:"BehaviorTree"`
	} `toml:"Resource"`

	DisconnectTest struct {
		// whether enbale disconnection test.
		Enable bool `toml:"Enable"`

		// interval(s) to disconnect robot
		Interval uint `toml:"Interval"`

		// robot count to disconnect
		Count uint `toml:"Count"`
	} `toml:"DisconnectTest"`

	Statistics struct {
		// interval(s) to log statistics
		OutputInterval uint `toml:"OutputInterval"`

		// statistics output file
		OutputFile string `toml:"OutputFile"`

		// the port of statistics metrics api
		MetricsPort int `toml:"MetricsPort"`
	} `toml:"Statistics"`

	PProf struct {
		// whether to start pprof
		Enable bool `toml:"Enable"`

		// pprof port
		Port int `toml:"Port"`
	} `toml:"PProf"`

	Log struct {
		// directory of log files
		Dir string `toml:"Dir"`

		// log level
		Level string `toml:"Level"`

		// whether enable std output.
		EnableStdOut bool `toml:"EnableStdOut"`

		// max file size.
		MaxSize int `toml:"MaxSize"`

		// max age.
		MaxAge int `toml:"MaxAge"`

		// log max backups
		MaxBackups uint `toml:"MaxBackups"`
	}
}

func (c *Config) Check() error {
	// if c.IntervalOfBehaviorTreePerform.Max < c.IntervalOfBehaviorTreePerform.Min {
	// 	return errors.New("IntervalOfBehaviorTreePeform: invalid range")
	// }

	if c.TickCycle == 0 {
		return errors.New("TickCycle zero")
	}

	return nil
}

func (c *Config) MakeRobotUserID(no uint) string {
	return fmt.Sprintf("%s%d", c.Robot.UserIDPrefix, int64(no)+int64(c.Robot.InitialNo))
}

func (c *Config) GetStatisticsOutputInterval() time.Duration {
	return time.Second * time.Duration(c.Statistics.OutputInterval)
}

func LoadConfig(file string) (*Config, error) {
	config := new(Config)
	if _, err := toml.DecodeFile(file, config); err != nil {
		return nil, err
	}

	return config, config.Check()
}

var cfg *Config

func Init(file string) (*Config, error) {
	var err error
	cfg, err = LoadConfig(file)
	return cfg, err
}

func GetConfig() *Config { return cfg }
