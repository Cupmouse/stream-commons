package formatter

import (
	"fmt"
)

var formatters = make(map[string]map[string]Formatter)

// Formatter formats raw line into desired format.
type Formatter interface {
	FormatStart(urlStr string) ([]StartReturn, error)
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

	formatters["bitflyer"]["json"] = new(bitflyerFormatter)
	formatters["bitfinex"]["json"] = new(bitfinexFormatter)
	formatters["bitmex"]["json"] = new(bitmexFormatter)
	formatters["binance"]["json"] = new(binanceFormatter)
}
