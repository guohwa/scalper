package service

import (
	"context"

	"scalper/config"

	"github.com/uncle-gua/gobinance/futures"
	"github.com/uncle-gua/log"
)

// func init() {
// 	if err := service.Start(); err != nil {
// 		log.Fatal(err)
// 	}
// }

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
	Ticker:   &Kline{},
	Status:   "Stopped",
	policies: make([]Policy, 0),
}

type klineService struct {
	Klines   *Klines
	Ticker   *Kline
	Status   string
	policies []Policy
	done     chan struct{}
	stop     chan struct{}
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

	for _, p := range serv.policies {
		go func(policy Policy) {
			defer func() {
				if err := recover(); err != nil {
					log.Error(err)
				}
			}()
			policy.Call(serv.Klines, serv.Ticker, event.Kline.IsFinal)
		}(p)
	}
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

	serv.done, serv.stop, err = futures.WsKlineServe(config.Param.Symbol.Name, config.Param.Symbol.Period, serv.wsKlineHandler, serv.errHandler)
	if err != nil {
		return err
	}

	serv.Status = "Started"
	return nil
}

func (serv *klineService) Register(p Policy) {
	serv.policies = append(serv.policies, p)
}

func (serv *klineService) Stop() {
	defer close(serv.done)
	serv.Status = "Stopped"
	serv.stop <- struct{}{}
}
