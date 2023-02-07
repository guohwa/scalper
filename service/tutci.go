package service

import (
	"math"
	"sync"
	"time"

	"scalper/config"
	"scalper/log"
	"scalper/ta"
)

func NewTuTCI() Policy {
	return &TuTCI{
		mutex: &sync.RWMutex{},
	}
}

type TuTCI struct {
	mutex *sync.RWMutex
}

func (t *TuTCI) Call(klines *Klines, ticker *Kline, isFinal bool) {
	if isFinal {
		t.run(klines, ticker)
	}

	if position.Hold != "NONE" {
		t.trail(klines, ticker)
	}
}

func (t *TuTCI) run(klines *Klines, ticker *Kline) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	start := time.Now()

	trend := t.supertrend(klines, config.Param.SuperTrend.DemaLength, config.Param.SuperTrend.AtrLength, config.Param.SuperTrend.AtrMult)
	upper, lower := t.tutci(klines, config.Param.TuTCI.Entry)
	pv := t.pv(klines)
	ssl := t.ssl(klines, config.Param.SSL.Length)

	n := klines.Length - 1
	bull := trend[n] > 0 &&
		klines.High[n] > upper[n-1] &&
		(!config.Param.PV.Enable || pv[n] < config.Param.PV.Threshold) &&
		(!config.Param.SSL.Enable || klines.Close[n] > ssl[n])
	bear := trend[n] < 0 &&
		klines.Low[n] < lower[n-1] &&
		(!config.Param.PV.Enable || pv[n] < config.Param.PV.Threshold) &&
		(!config.Param.SSL.Enable || klines.Close[n] < ssl[n])

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

func (t *TuTCI) trail(klines *Klines, ticker *Kline) {
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

func (t *TuTCI) supertrend(klines *Klines, length1, length2 int, mult float64) []int {
	ema1 := ta.Ema(klines.Close, length1)
	ema2 := ta.Ema(ema1, length1)
	ema3 := ta.Ema(ema2, length1)
	s1 := make([]float64, klines.Length)
	ta.Fill(s1, 3.0)
	basis := ta.Add(ta.Mult(ta.Sub(ema1, ema2), s1), ema3)

	atr := ta.Atr(klines.High, klines.Low, klines.Close, length2)

	s2 := make([]float64, klines.Length)
	ta.Fill(s2, mult)
	up := ta.Sub(basis, ta.Mult(atr, s2))
	dn := ta.Add(basis, ta.Mult(atr, s2))

	tup := make([]float64, klines.Length)
	tdn := make([]float64, klines.Length)
	for i := 0; i < klines.Length; i++ {
		tup[i] = up[i]
		tdn[i] = dn[i]

		if i == 0 {
			continue
		}

		if klines.Close[i-1] > tup[i-1] {
			tup[i] = math.Max(up[i], tup[i-1])
		}
		if klines.Close[i-1] < tdn[i-1] {
			tdn[i] = math.Min(dn[i], tdn[i-1])
		}
	}

	trend := make([]int, klines.Length)
	for i := 0; i < klines.Length; i++ {
		if i == 0 {
			trend[i] = 1
			continue
		}

		if klines.Close[i] > tdn[i-1] {
			trend[i] = 1
		} else if klines.Close[i] < tup[i-1] {
			trend[i] = -1
		} else {
			trend[i] = trend[i-1]
		}
	}

	tsl := make([]float64, klines.Length)
	for i := 0; i < klines.Length; i++ {
		if trend[i] > 0 {
			tsl[i] = tup[i]
		} else {
			tsl[i] = tdn[i]
		}
	}

	return trend
}

func (t *TuTCI) tutci(klines *Klines, entry int) ([]float64, []float64) {
	upper := ta.Max(klines.High, entry)
	lower := ta.Min(klines.Low, entry)
	return upper, lower
}

func (t *TuTCI) pv(klines *Klines) []float64 {
	pv := make([]float64, klines.Length)
	for i := 0; i < klines.Length; i++ {
		pv[i] = (klines.High[i] - klines.Low[i]) / math.Abs(klines.Close[i]-klines.Open[i])
	}
	return pv
}

func (t *TuTCI) ssl(klines *Klines, length int) []float64 {
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
