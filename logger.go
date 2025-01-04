package pgxtx

import (
	"github.com/rs/zerolog"
)

var logger zerolog.Logger

func Init(in zerolog.Logger) {
	logger = in
}
