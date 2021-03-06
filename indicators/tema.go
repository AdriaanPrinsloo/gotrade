// Triple Exponential Moving Average (TEMA)
package indicators

// TEMA(X) = (2 * EMA(X, CLOSE)) - (EMA(X, EMA(X, CLOSE)))

import (
	"github.com/thetruetrade/gotrade"
)

type baseTEMA struct {
	*baseIndicatorWithLookback

	// private variables
	valueAvailableAction ValueAvailableAction
	ema1                 *EMA
	ema2                 *EMA
	ema3                 *EMA
	currentEMA           float64
	currentEMA2          float64
}

func newBaseTEMA(lookbackPeriod int) *baseTEMA {
	newTEMA := baseTEMA{baseIndicatorWithLookback: newBaseIndicatorWithLookback(lookbackPeriod)}
	return &newTEMA
}

// A Double Exponential Moving Average Indicator
type TEMA struct {
	*baseTEMA

	// public variables
	Data []float64
}

// NewTEMA returns a new Double Exponential Moving Average (TEMA) configured with the
// specified lookbackPeriod. The TEMA results are stored in the DATA field.
func NewTEMA(lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *TEMA, err error) {
	newTEMA := TEMA{baseTEMA: newBaseTEMA(lookbackPeriod)}
	newTEMA.ema1, _ = NewEMA(lookbackPeriod, selectData)

	newTEMA.ema1.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.currentEMA = dataItem
		newTEMA.ema2.ReceiveTick(dataItem, streamBarIndex)
	}

	newTEMA.ema2, _ = NewEMA(lookbackPeriod, selectData)
	newTEMA.ema2.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.currentEMA2 = dataItem
		newTEMA.ema3.ReceiveTick(dataItem, streamBarIndex)
	}

	newTEMA.ema3, _ = NewEMA(lookbackPeriod, selectData)

	newTEMA.ema3.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.dataLength += 1
		if newTEMA.validFromBar == -1 {
			newTEMA.validFromBar = streamBarIndex
		}

		//T-EMA = (3*EMA – 3*EMA(EMA)) + EMA(EMA(EMA))
		tema := (3*newTEMA.currentEMA - 3*newTEMA.currentEMA2) + dataItem

		if tema > newTEMA.maxValue {
			newTEMA.maxValue = tema
		}

		if tema < newTEMA.minValue {
			newTEMA.minValue = tema
		}

		newTEMA.valueAvailableAction(tema, streamBarIndex)
	}

	newTEMA.selectData = selectData
	newTEMA.valueAvailableAction = func(dataItem float64, streamBarIndex int) {
		newTEMA.Data = append(newTEMA.Data, dataItem)
	}
	return &newTEMA, nil
}

func NewTEMAForStream(priceStream *gotrade.DOHLCVStream, lookbackPeriod int, selectData gotrade.DataSelectionFunc) (indicator *TEMA, err error) {
	newTEMA, err := NewTEMA(lookbackPeriod, selectData)
	priceStream.AddTickSubscription(newTEMA)
	return newTEMA, err
}

func (tema *baseTEMA) ReceiveDOHLCVTick(tickData gotrade.DOHLCV, streamBarIndex int) {
	var selectedData = tema.selectData(tickData)
	tema.ReceiveTick(selectedData, streamBarIndex)
}

func (tema *baseTEMA) ReceiveTick(tickData float64, streamBarIndex int) {
	tema.ema1.ReceiveTick(tickData, streamBarIndex)
}
