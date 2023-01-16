package service

type Kline struct {
	OpenTime  int64
	CloseTime int64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
}

type Klines struct {
	Length    int
	OpenTime  []int64
	CloseTime []int64
	Open      []float64
	High      []float64
	Low       []float64
	Close     []float64
	Volume    []float64
}

func (klines *Klines) Append(openTime, closeTime int64, open, high, low, close, volume float64) {
	klines.OpenTime = append(klines.OpenTime, openTime)
	klines.CloseTime = append(klines.CloseTime, closeTime)
	klines.Open = append(klines.Open, open)
	klines.High = append(klines.High, high)
	klines.Low = append(klines.Low, low)
	klines.Close = append(klines.Close, close)
	klines.Volume = append(klines.Volume, volume)
	klines.Length++
}

func (klines *Klines) Remove() {
	klines.OpenTime = klines.OpenTime[:klines.Length-1]
	klines.CloseTime = klines.CloseTime[:klines.Length-1]
	klines.Open = klines.Open[:klines.Length-1]
	klines.High = klines.High[:klines.Length-1]
	klines.Low = klines.Low[:klines.Length-1]
	klines.Close = klines.Close[:klines.Length-1]
	klines.Volume = klines.Volume[:klines.Length-1]
	klines.Length--
}
