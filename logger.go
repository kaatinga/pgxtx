package pgxtx

import (
	"github.com/kaatinga/dummylogger"
)

var l = dummylogger.Get()

func Init(logger dummylogger.I) {
	dummylogger.Set(logger)
}
