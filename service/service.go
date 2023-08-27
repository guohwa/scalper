package service

import (
	"context"
	"strconv"
	"sync"

	"scalper/config"
	"scalper/log"

	"github.com/uncle-gua/gobinance/futures"
)

func init() {
	if err := service.Start(); err != nil {
		log.Fatal(err)
	}
	if err := position.Load(); err != nil {
		log.Fatal(err)
	}
}

func Start() error {
	return service.Start()
}

func Stop() {
	service.Stop()
}

func Status() string {
	return service.Status
}

var service = &klineService{
	Klines: &Klines{
		Length:    0,
		OpenTime:  make([]int64, 0),
		CloseTime: make([]int64, 0),
		Open:      make([]float64, 0),
		High:      make([]float64, 0),
		Low:       make([]float64, 0),
		Close:     make([]float64, 0),
		Volume:    make([]float64, 0),
	},
	Ticker: &Kline{},
	Status: "Stopped",
	Policy: &Scalper{
		mutex: &sync.RWMutex{},
	},
}

type klineService struct {
	Klines        *Klines
	Ticker        *Kline
	Status        string
	Policy        Policy
	ch_kline_done chan struct{}
	ch_trade_done chan struct{}
}

func (serv *klineService) errHandler(err error) {
	log.Error(err)
}

func (serv *klineService) wsKlineHandler(event *futures.WsKlineEvent) {
	if serv.Klines.Length == 0 {
		return
	}

	kline := event.Kline

	serv.Ticker.OpenTime = kline.StartTime
	serv.Ticker.CloseTime = kline.EndTime
	serv.Ticker.Open = kline.Open
	serv.Ticker.High = kline.High
	serv.Ticker.Low = kline.Low
	serv.Ticker.Close = kline.Close
	serv.Ticker.Volume = kline.Volume

	if serv.Klines.OpenTime[serv.Klines.Length-1] == serv.Ticker.OpenTime {
		serv.Klines.Remove()
	}
	if event.Kline.IsFinal {
		serv.Klines.Append(
			serv.Ticker.OpenTime,
			serv.Ticker.CloseTime,
			serv.Ticker.Open,
			serv.Ticker.High,
			serv.Ticker.Low,
			serv.Ticker.Close,
			serv.Ticker.Volume,
		)
	}

	go func(policy Policy) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		policy.Call(serv.Klines, serv.Ticker, event.Kline.IsFinal)
	}(serv.Policy)
}

func (serv *klineService) wsAggTradeHandler(event *futures.WsAggTradeEvent) {
	price, err := strconv.ParseFloat(event.Price, 64)
	if err != nil {
		log.Error(err)
		return
	}

	kline := &Kline{
		OpenTime: 0,
		Close:    price,
	}

	go func(policy Policy) {
		defer func() {
			if err := recover(); err != nil {
				log.Error(err)
			}
		}()
		policy.Call(serv.Klines, kline, false)
	}(serv.Policy)
}

func (serv *klineService) Start() error {
	client := futures.NewClient("", "")

	klines, err := client.NewKlinesService().
		Symbol(config.Param.Symbol.Name).
		Interval(config.Param.Symbol.Period).
		Limit(config.Param.Symbol.Limit).Do(context.Background())
	if err != nil {
		return err
	}

	for i := 0; i < len(klines)-1; i++ {
		kline := klines[i]
		serv.Klines.Append(kline.OpenTime, kline.CloseTime, kline.Open, kline.High, kline.Low, kline.Close, kline.Volume)
	}

	serv.ch_kline_done, err = futures.WsKlineServe(config.Param.Symbol.Name, config.Param.Symbol.Period, serv.wsKlineHandler, serv.errHandler)
	if err != nil {
		return err
	}

	serv.ch_trade_done, err = futures.WsAggTradeServe(config.Param.Symbol.Name, serv.wsAggTradeHandler, serv.errHandler)
	if err != nil {
		return err
	}

	serv.Status = "Started"
	return nil
}

func (serv *klineService) Stop() {
	serv.Status = "Stopped"
	serv.ch_kline_done <- struct{}{}
	serv.ch_kline_done <- struct{}{}
}
