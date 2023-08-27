package service

import (
	"math"
	"sync"
	"time"

	"scalper/config"
	"scalper/log"
	"scalper/ta"
)

type Scalper struct {
	mutex *sync.RWMutex
}

func (t *Scalper) Call(klines *Klines, ticker *Kline, isFinal bool) {
	if isFinal {
		t.run(klines, ticker)
	}

	if position.Hold != "NONE" {
		t.trail(klines, ticker)
	}
}

func (t *Scalper) run(klines *Klines, ticker *Kline) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	start := time.Now()

	ssl := t.ssl(klines, config.Param.SSL.Length)

	n := klines.Length - 1
	bull := klines.Close[n] > ssl[n]
	bear := klines.Close[n] < ssl[n]

	elapsed := time.Since(start)

	if bull {
		if position.Hold != "LONG" {
			position.Open("LONG", ticker.Close)
			return
		}
	} else if bear {
		if position.Hold != "SHORT" {
			position.Open("SHORT", ticker.Close)
			return
		}
	}
	log.Infof("Scalper elapsed: %s", elapsed)
}

func (t *Scalper) trail(klines *Klines, ticker *Kline) {
	if !t.mutex.TryLock() {
		return
	}
	defer t.mutex.Unlock()

	if position.Peak < 0 {
		position.Peak = position.Entry
	}

	sign := func() float64 {
		if position.Hold == "SHORT" {
			return -1
		}

		return 1
	}()

	roe := sign * (ticker.Close - position.Entry) / position.Entry * 100

	if sign < 0 {
		position.Peak = func() float64 {
			if ticker.OpenTime > 0 {
				return math.Min(position.Peak, ticker.Low)
			}
			return math.Min(position.Peak, ticker.Close)
		}()
	} else {
		position.Peak = func() float64 {
			if ticker.OpenTime > 0 {
				return math.Max(position.Peak, ticker.High)
			}
			return math.Max(position.Peak, ticker.Close)
		}()
	}

	if roe < -config.Param.TSL.StopLoss {
		position.Close(position.Hold, ticker.Close)
		return
	}

	if roe >= config.Param.TSL.TrailProfit {
		position.Reach = true
	}

	if position.Reach {
		offset := sign * (ticker.Close - position.Peak) / position.Peak * 100
		if offset < -config.Param.TSL.TrailOffset {
			position.Close(position.Hold, ticker.Close)
		}
	}
}

func (t *Scalper) ssl(klines *Klines, length int) []float64 {
	hh := ta.Hma(klines.High, length)
	ll := ta.Hma(klines.Low, length)
	hlv := make([]int, klines.Length)
	for i := 0; i < klines.Length; i++ {
		hlv[i] = func() int {
			if i == 0 {
				return 0
			}
			if klines.Close[i] > hh[i] {
				return 1
			}
			if klines.Close[i] < ll[i] {
				return -1
			}
			return hlv[i-1]
		}()
	}
	ssl := make([]float64, klines.Length)
	for i := 0; i < klines.Length; i++ {
		ssl[i] = func() float64 {
			if hlv[i] < 0 {
				return hh[i]
			}
			return ll[i]
		}()
	}

	return ssl
}
