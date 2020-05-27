package csvlike

// HeaderOrderBook is the header row of standard orderbook snapshot channel
var HeaderOrderBook = "Pair,Price,Size"

// HeaderTrade is the header row of standard trade channel
var HeaderTrade = "Pair,Price,Size"

// HeaderBitfinexOrderBook is the header for book from bitfinex exchange
var HeaderBitfinexOrderBook = "Pair,Price,Count,Size"

// HeaderBitfinexFundingTrade is the header for funding trade from bitfinex exchange
var HeaderBitfinexFundingTrade = "Pair,Rate,Period,Amount"

// HeaderBitmexOrderBookL2 is the header for orderBookL2 channel in bitmex exchange
var HeaderBitmexOrderBookL2 = "Pair,Price,ID,Size"

// HeaderBitmexTrade is the header for trade channel in bitmex exchange
var HeaderBitmexTrade = "Pair,Price,Size,Timestamp,TrdMatchID,TickDirection,GrossValue,HomeNotional,ForeignNotional"
