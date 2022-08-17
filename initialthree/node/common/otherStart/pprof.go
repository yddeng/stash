package otherStart

import (
	"initialthree/pkg/util"
	"net/http"
	_ "net/http/pprof"
)

// pprof
func StartPProf(address string) {
	address = util.ParseAddress(address)
	go func() {
		err := http.ListenAndServe(address, nil)
		if err != nil {
			panic(err)
		}
	}()
}

func PProfHasAndRun(args []string) {
	DisposeArgs(args)
	if address, ok := Has(PPROF); ok {
		StartPProf(address)
	}
}
