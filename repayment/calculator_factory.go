package repayment

import (
	"installment-calculator/beans"
	"installment-calculator/installment/installment_strategy"
	"installment-calculator/repayment/readjustment_strategy"
	"installment-calculator/repayment/repayment_strategy"
	"time"
)

type EarlyRepaymentCalculator interface {
	// Calculate 返回新的还款计划, 还款总金额, 提前还款费用(本期利息+服务费), 节省利息
	Calculate() (*[]beans.Installment, float64, float64, float64)
}

type Config struct {
	Installment          beans.InstallmentInfo
	ReadjustmentStrategy readjustment_strategy.Strategy // enums.ReduceMonthlyRepayment, enums.ShortenRepaymentTenors
	Amount               int64
	Date                 time.Time
}

type Option func(*Config)

func Installment(strategy installment_strategy.Strategy, borrowDate time.Time, dueDay int, amount float64, tenors int, annualPercentageRate float64) Option {
	return func(s *Config) {
		s.Installment = beans.InstallmentInfo{
			Strategy:             strategy,
			BorrowDate:           borrowDate,
			DueDay:               dueDay,
			Amount:               amount,
			Tenors:               tenors,
			AnnualPercentageRate: annualPercentageRate,
		}
	}
}

func ReadjustmentStrategy(strategy readjustment_strategy.Strategy) Option {
	return func(s *Config) {
		s.ReadjustmentStrategy = strategy
	}
}

func Amount(amount int64) Option {
	return func(s *Config) {
		s.Amount = amount
	}
}

func Date(date time.Time) Option {
	return func(s *Config) {
		s.Date = date
	}
}

func TrialCalculate(strategy repayment_strategy.Strategy, opts ...Option) (*[]beans.Installment, float64, float64, float64) {
	config := &Config{}
	for _, option := range opts {
		option(config)
	}

	if config.Date.IsZero() {
		config.Date = time.Now()
	}

	var calc EarlyRepaymentCalculator
	switch strategy {
	case repayment_strategy.Payoff:
		calc = PayoffCalculator{
			installmentInfo: config.Installment,
			repaymentDate:   config.Date,
		}
	case repayment_strategy.PartialRepayment:
		if config.Amount <= 0 {
			panic("invalid repayment amount")
		}

		// 月供不变 缩短期限 仅对等额本息生效
		if (config.Installment.Strategy == installment_strategy.Outright || config.Installment.Strategy == installment_strategy.Linear) &&
			config.ReadjustmentStrategy == readjustment_strategy.ShortenRepaymentTenors {
			panic("Outright & Linear not supported ShortenRepaymentTenors")
		}

		if config.ReadjustmentStrategy == "" {
			config.ReadjustmentStrategy = readjustment_strategy.ReduceMonthlyRepayment
		}

		calc = PartialRepaymentCalculator{
			installmentInfo:      config.Installment,
			repaymentDate:        config.Date,
			repaymentAmount:      config.Amount,
			readjustmentStrategy: config.ReadjustmentStrategy,
		}
	default:
		panic("invalid repayment strategy")
	}
	return calc.Calculate()
}
