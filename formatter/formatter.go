package formatter

import (
	"fmt"

	"github.com/exchangedataset/streamcommons/formatter/jsonf"
)

var formatters = make(map[string]map[string]Formatter)

// Formatter formats raw line into desired format.
type Formatter interface {
	FormatStart(urlStr string) ([][]byte, error)
	FormatMessage(channel string, line []byte) ([][]byte, error)
	IsSupported(channel string) bool
}

// GetFormatter returns the right formatter for given parameters.
func GetFormatter(exchange string, channels []string, format string) (Formatter, error) {
	formats, ok := formatters[exchange]
	if !ok {
		return nil, fmt.Errorf("exchange '%s' is not supported", exchange)
	}
	formatter, ok := formats[format]
	if !ok {
		return nil, fmt.Errorf("format '%s' is not supported for exchange '%s'", format, exchange)
	}
	for _, ch := range channels {
		if !formatter.IsSupported(ch) {
			return nil, fmt.Errorf("channel '%s' of exchange '%s' is not supported for format '%s'", ch, exchange, format)
		}
	}
	return formatter, nil
}

func init() {
	formatters["bitflyer"] = make(map[string]Formatter)
	formatters["bitfinex"] = make(map[string]Formatter)
	formatters["bitmex"] = make(map[string]Formatter)
	formatters["binance"] = make(map[string]Formatter)

	formatters["bitflyer"]["json"] = new(jsonf.BitflyerFormatter)
	formatters["bitfinex"]["json"] = &jsonf.BitfinexFormatter{}
	formatters["bitmex"]["json"] = &jsonf.BitmexFormatter{}
	formatters["binance"]["json"] = &jsonf.BinanceFormatter{}
}
