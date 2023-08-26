package log

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func init() {

	zerolog.TimestampFunc = func() time.Time {
		loc := time.FixedZone("Asia/Jakarta", 7*60*60)
		return time.Now().In(loc)
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05.000"

}

func InitZerolog() {

	Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02/01/2006 15:04:05,999"}).With().Timestamp().Caller().Logger()

}

func Info() *zerolog.Event {
	return Logger.Info()
}

func Error() *zerolog.Event {
	return Logger.Error()
}

func Debug() *zerolog.Event {
	return Logger.Debug()
}

func Trace() *zerolog.Event {
	return Logger.Trace()
}

func Fatal() *zerolog.Event {
	return Logger.Fatal()
}
