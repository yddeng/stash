package main

import (
	"errors"
	"flag"
	"fmt"
	"initialthree/robot/config"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type int8Value int8

func newInt8Value(v int8, p *int8) *int8Value {
	*p = v
	return (*int8Value)(p)
}

func (p *int8Value) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

func (p *int8Value) Set(s string) error {
	if i, err := strconv.ParseInt(s, 10, 8); err != nil {
		return err
	} else {
		*p = int8Value(i)
		return nil
	}
}

type uint8Value uint8

func newUint8Value(v uint8, p *uint8) *uint8Value {
	*p = v
	return (*uint8Value)(p)
}

func (p *uint8Value) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

func (p *uint8Value) Set(s string) error {
	if i, err := strconv.ParseInt(s, 10, 8); err != nil {
		return err
	} else {
		*p = uint8Value(i)
		return nil
	}
}

type int16Value int16

func newInt16Value(v int16, p *int16) *int16Value {
	*p = v
	return (*int16Value)(p)
}

func (p *int16Value) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

func (p *int16Value) Set(s string) error {
	if i, err := strconv.ParseInt(s, 10, 16); err != nil {
		return err
	} else {
		*p = int16Value(i)
		return nil
	}
}

type uint16Value uint16

func newUint16Value(v uint16, p *uint16) *uint16Value {
	*p = v
	return (*uint16Value)(p)
}

func (p *uint16Value) String() string {
	return strconv.FormatInt(int64(*p), 10)
}

func (p *uint16Value) Set(s string) error {
	if i, err := strconv.ParseInt(s, 10, 16); err != nil {
		return err
	} else {
		*p = uint16Value(i)
		return nil
	}
}

type floatingDurationValue config.FloatingDuration

func newFloatingDurationValue(v config.FloatingDuration, p *config.FloatingDuration) *floatingDurationValue {
	*p = v
	return (*floatingDurationValue)(p)
}

func (p *floatingDurationValue) String() string {
	return fmt.Sprintf("%d:%d:%d", p.Unit, p.Min, p.Max)
}

func (p *floatingDurationValue) Set(s string) error {
	if n, err := fmt.Sscanf(s, "%d:%d:%d", &p.Unit, &p.Min, &p.Max); err != nil {
		return err
	} else if n != 3 {
		return errors.New("invalid format")
	} else {
		return nil
	}
}

type flagBinding struct {
	name  string
	usage string
	field string
	value interface{}
}

func (fb *flagBinding) bind(c *config.Config, fs *flag.FlagSet) {
	v := reflect.ValueOf(c)
	v = v.Elem()

	var fv reflect.Value
	fields := strings.Split(fb.field, ".")
	if len(fields) == 0 {
		panic(fmt.Errorf("field invalid of flag %s", fb.name))
	}

	for _, field := range fields {
		if !fv.IsValid() {
			fv = v.FieldByName(field)
		} else {
			fv = fv.FieldByName(field)
		}

		if !fv.IsValid() {
			panic(fmt.Errorf("config do not have field %s", fb.field))
		}
	}

	if reflect.TypeOf(fb.value) != fv.Type() {
		panic(fmt.Errorf("field and value types do not match of flag %s", fb.name))
	}

	switch fv.Interface().(type) {
	case string:
		fs.StringVar((*string)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(string), fb.usage)

	case int8:
		fs.Var(newInt8Value(fb.value.(int8), (*int8)(unsafe.Pointer(fv.UnsafeAddr()))), fb.name, fb.usage)

	case uint8:
		fs.Var(newUint8Value(fb.value.(uint8), (*uint8)(unsafe.Pointer(fv.UnsafeAddr()))), fb.name, fb.usage)

	case int16:
		fs.Var(newInt16Value(fb.value.(int16), (*int16)(unsafe.Pointer(fv.UnsafeAddr()))), fb.name, fb.usage)

	case uint16:
		fs.Var(newUint16Value(fb.value.(uint16), (*uint16)(unsafe.Pointer(fv.UnsafeAddr()))), fb.name, fb.usage)

	case uint, uint32:
		fs.UintVar((*uint)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(uint), fb.usage)

	case int, int32:
		fs.IntVar((*int)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(int), fb.usage)

	case uint64:
		fs.Uint64Var((*uint64)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(uint64), fb.usage)

	case int64:
		fs.Int64Var((*int64)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(int64), fb.usage)

	case float64:
		fs.Float64Var((*float64)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(float64), fb.usage)

	case bool:
		fs.BoolVar((*bool)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(bool), fb.usage)

	case time.Duration:
		fs.DurationVar((*time.Duration)(unsafe.Pointer(fv.UnsafeAddr())), fb.name, fb.value.(time.Duration), fb.usage)

	case config.FloatingDuration:
		fs.Var(newFloatingDurationValue(fb.value.(config.FloatingDuration), (*config.FloatingDuration)(unsafe.Pointer(fv.UnsafeAddr()))), fb.name, fb.usage)

	default:
		panic(fmt.Errorf("field type of flag %s not support", fb.name))
	}

}

func (fb *flagBinding) copy(dst, src *config.Config) {
	vDst, vSrc := reflect.ValueOf(dst).Elem(), reflect.ValueOf(src).Elem()

	var fieldDst, fieldSrc reflect.Value
	fields := strings.Split(fb.field, ".")
	for _, field := range fields {
		if !fieldDst.IsValid() {
			fieldDst = vDst.FieldByName(field)
		} else {
			fieldDst = fieldDst.FieldByName(field)
		}

		if !fieldSrc.IsValid() {
			fieldSrc = vSrc.FieldByName(field)
		} else {
			fieldSrc = fieldSrc.FieldByName(field)
		}
	}

	fieldDst.Set(fieldSrc)
}

var flagBindings = map[string]*flagBinding{
	"service": {
		name:  "service",
		usage: "address of login service",
		field: "Service",
		value: "127.0.0.1:9201",
	},
	"server-id": {
		name:  "server-id",
		usage: "server id",
		field: "ServerID",
		value: uint(1),
	},
	"robot-id-prefix": {
		name:  "robot-id-prefix",
		usage: "robot user-id prefix",
		field: "Robot.UserIDPrefix",
		value: "robot_",
	},
	"robot-count": {
		name:  "robot-count",
		usage: "robot count",
		field: "Robot.Count",
		value: uint(100),
	},
	"robot-initail-no": {
		name:  "robot-initail-no",
		usage: "initail number of robot user-id",
		field: "Robot.InitialNo",
		value: int(1),
	},
	"excel-path": {
		name:  "excel-path",
		usage: "directory of excel files",
		field: "Resource.ExcelPath",
		value: "./config/excel",
	},
	"quest-path": {
		name:  "quest-path",
		usage: "directory of quest files",
		field: "Resource.QuestPath",
		value: "./config/quest",
	},
	"behavior-config-path": {
		name:  "behavior-config-path",
		usage: "behavior config xml file",
		field: "Resource.BehaviorConfigPath",
		value: "./config/behavior/config.xml",
	},
	"behavior-tree": {
		name:  "behavior-tree",
		usage: "behavior tree name",
		field: "Resource.BehaviorTree",
		value: "tree",
	},
	"tick-cycle": {
		name:  "tick-cycle",
		usage: "tick cycle, ms",
		field: "TickCycle",
		value: uint(1000),
	},
	"disconnect-enable": {
		name:  "diconnect-enable",
		usage: "whether enable disconnection test",
		field: "DisconnectTest.Enable",
		value: true,
	},
	"disconnect-interval": {
		name:  "disconnect-interval",
		usage: "interval(s) to disconnect robot",
		field: "DisconnectTest.Interval",
		value: uint(300),
	},
	"disconnect-count": {
		name:  "disconnect-count",
		usage: "robot count to disconnect",
		field: "DisconnectTest.Count",
		value: uint(10),
	},
	"statistics-output-interval": {
		name:  "statistics-output-interval",
		usage: "interval(s) to output statistics",
		field: "Statistics.OutputInterval",
		value: uint(60),
	},
	"statistics-output-file": {
		name:  "statistics-output-file",
		usage: "statistics output file",
		field: "Statistics.OutputFile",
		value: "./statistics.txt",
	},
	"statistics-metrics-port": {
		name:  "statistics-metrics-port",
		usage: "the port of statistics metrics api",
		field: "Statistics.MetricsPort",
		value: int(2112),
	},
	"pprof-enable": {
		name:  "pprof-enable",
		usage: "whether to enable pprof",
		field: "PProf.Enable",
		value: false,
	},
	"pprof-port": {
		name:  "pprof-port",
		usage: "pprof port",
		field: "PProf.Port",
		value: int(2222),
	},
	"log-dir": {
		name:  "log-dir",
		usage: "directory of log files",
		field: "Log.Dir",
		value: "./log",
	},
	"log-level": {
		name:  "log-level",
		usage: "log level",
		field: "Log.Level",
		value: "debug",
	},
	"log-enable-std-out": {
		name:  "log-enable-std-out",
		usage: "whether enable std output",
		field: "Log.EnableStdOut",
		value: false,
	},
	"log-max-size": {
		name:  "log-max-size",
		usage: "log file max size",
		field: "Log.MaxSize",
		value: 100,
	},
	"log-max-age": {
		name:  "log-max-age",
		usage: "log file max age",
		field: "Log.MaxAge",
		value: 14,
	},
	"log-max-backups": {
		name:  "log-max-backups",
		usage: "log max backups",
		field: "Log.MaxBackups",
		value: uint(0),
	},
}

func bindFlags(c *config.Config, fs *flag.FlagSet) {
	for _, v := range flagBindings {
		v.bind(c, fs)
	}
}
