package service

type Policy interface {
	Call(klines *Klines, ticker *Kline, isFinal bool)
}
