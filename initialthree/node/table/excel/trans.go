package excel

import (
	"fmt"
	"github.com/hqpko/hutils"
	"initialthree/node/common/enumType"
	"strconv"
)

func ReadInt32(str string) int32 {
	if str == "" {
		return 0
	} else {
		return int32(ReadFloat(str))
	}
}

func ReadInt64(str string) int64 {
	if str == "" {
		return 0
	} else {
		return int64(ReadFloat(str))
	}
}

func ReadFloat(str string) float64 {
	if str == "" {
		return 0
	} else {
		vv, err := strconv.ParseFloat(str, 64)
		if nil != err {
			panic(err)
		}
		return vv
	}
}

func ReadBool(str string) bool {
	if str == "" {
		return false
	} else {
		b, err := strconv.ParseBool(str)
		if nil != err {
			panic(err)
		}
		return b
	}
}

func ReadStr(str string) string {
	return fmt.Sprintf("%s", str)
}

func ReadEnum(str string, def ...string) int32 {
	if str == "" && len(def) > 0 {
		if def[0] == "" {
			return 0
		}

		str = def[0]
	}

	return hutils.Must(enumType.GetEnumType(str)).(int32)
}

func ToStr(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
