package formatter

import (
	"fmt"
	"strings"

	"github.com/exchangedataset/streamcommons"
)

// Formatter formats raw line into desired format.
type Formatter interface {
	FormatStart(urlStr string) ([]Result, error)
	FormatMessage(channel string, line []byte) ([]Result, error)
	IsSupported(channel string) bool
}

// GetFormatter returns the right formatter for given parameters.
func GetFormatter(exchange string, channels []string, format string) (Formatter, error) {
	var f Formatter
	switch format {
	case "json":
		switch exchange {
		case streamcommons.ExchangeBinance:
			f = newBinanceFormatter()
		case streamcommons.ExchangeBitflyer:
			f = newBitflyerFormatter()
		case streamcommons.ExchangeBitmex:
			// Converts raw channels (user specified) to formatter channels.
			set := make(map[string]bool)
			for _, ch := range channels {
				if ri := strings.IndexRune(ch, '_'); ri != -1 {
					set[ch] = true
				} else {
					set[ch] = true
				}
				// Ignore channels that does not have a symbol
			}
			targets := make([]string, len(set))
			i := 0
			for ch := range set {
				targets[i] = ch
				i++
			}
			f = newBitmexFormatter(targets)
		case streamcommons.ExchangeLiquid:
			f = newLiquidFormatter()
		case streamcommons.ExchangeBitfinex:
			f = newBitfinexFormatter()
		default:
			return nil, fmt.Errorf("format '%s' is not supported for exchange '%s'", format, exchange)
		}
	default:
		return nil, fmt.Errorf("format '%s' is not supported", format)
	}
	for _, ch := range channels {
		if !f.IsSupported(ch) {
			return nil, fmt.Errorf("channel '%s' of exchange '%s' is not supported for format '%s'", ch, exchange, format)
		}
	}
	return f, nil
}
