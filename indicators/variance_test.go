package indicators_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/thetruetrade/gotrade"
	. "github.com/thetruetrade/gotrade/indicators"
	"time"
)

var _ = Describe("when calculating a variance(var)", func() {
	var (
		period     int = 3
		variance   *Variance
		indicator  Indicator
		sourceData = []gotrade.DOHLCV{gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 5.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 6.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 7.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 8.0, 0.0),
			gotrade.NewDOHLCVDataItem(time.Now(), 0.0, 0.0, 0.0, 9.0, 0.0)}
	)

	BeforeEach(func() {
		variance, _ = NewVariance(period, gotrade.UseClosePrice)
		indicator = variance
	})

	Context("and the variance has received less ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i < period-1; i++ {
				variance.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the variance should have no result data", func() {
			Expect(len(variance.Data)).To(Equal(0))
		})

		It("the indicator stream length should be zero", func() {
			Expect(indicator.Length()).To(Equal(0))
		})

		It("the indicator stream should have no valid bars", func() {
			Expect(indicator.ValidFromBar()).To(Equal(-1))
		})
	})

	Context("and the variance has received ticks equal to the lookback period", func() {

		BeforeEach(func() {
			for i := 0; i <= period-1; i++ {
				variance.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the variance should have result data with a single entry", func() {
			Expect(len(variance.Data)).To(Equal(1))
		})

		It("the indicator stream length should be one", func() {
			Expect(indicator.Length()).To(Equal(1))
		})

		It("the indicator stream min and max should be equal", func() {
			Expect(indicator.MaxValue()).To(Equal(indicator.MinValue()))
		})

		It("the indicator stream should have valid bars from the lookback period", func() {
			Expect(indicator.ValidFromBar()).To(Equal(period))
		})
	})

	Context("and the variance has received more ticks than the lookback period", func() {

		BeforeEach(func() {
			for i := range sourceData {
				variance.ReceiveDOHLCVTick(sourceData[i], i+1)
			}
		})

		It("the variance should have result data with entries equal to the number of ticks less the (lookback period - 1)", func() {
			Expect(len(variance.Data)).To(Equal(len(sourceData) - (period - 1)))
		})

		It("the indicator stream min should equal the result data minimum", func() {
			Expect(indicator.MinValue()).To(Equal(GetDataMin(variance.Data)))
		})

		It("the indicator stream max should equal the result data maximum", func() {
			Expect(indicator.MaxValue()).To(Equal(GetDataMax(variance.Data)))
		})
	})
})
