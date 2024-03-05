package installment

import (
	"installment-calculator/beans"
	"installment-calculator/constants"
	"installment-calculator/installment/installment_strategy"
	"time"
)

type Generator interface {
	// Generate 分期还款计划，还款总额
	Generate(amount float64, tenors int, annualPercentageRate float64) ([]beans.Installment, float64)
}

type Config struct {
	Strategy   installment_strategy.Strategy
	DueDay     int
	BorrowDate time.Time
}

type Option func(*Config)

func Strategy(strategy installment_strategy.Strategy) Option {
	return func(s *Config) {
		s.Strategy = strategy
	}
}

func DueDay(dueDay int) Option {
	return func(s *Config) {
		s.DueDay = dueDay
	}
}

func BorrowDate(borrowDate time.Time) Option {
	return func(s *Config) {
		s.BorrowDate = borrowDate
	}
}

func NewInstance(options ...Option) Generator {
	config := &Config{
		Strategy: installment_strategy.Annuity,
		DueDay:   constants.DefaultDueDay,
	}

	for _, option := range options {
		option(config)
	}

	return newGenerator(config.Strategy, config.DueDay, config.BorrowDate)
}

func newGenerator(strategy installment_strategy.Strategy, dueDay int, borrowDate time.Time) Generator {
	switch strategy {
	case installment_strategy.Linear:
		return LinearGenerator{
			BorrowDate: borrowDate,
			DueDay:     dueDay,
		}
	case installment_strategy.Annuity:
		return AnnuityGenerator{
			BorrowDate: borrowDate,
			DueDay:     dueDay,
		}
	case installment_strategy.Outright:
		return OutrightGenerator{
			BorrowDate: borrowDate,
			DueDay:     dueDay,
		}
	default:
		panic("unsupported installment_calculator strategy")

	}
}
