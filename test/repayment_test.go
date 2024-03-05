package test

import (
	"github.com/sirupsen/logrus"
	"installment-calculator/installment/installment_strategy"
	"installment-calculator/repayment"
	"installment-calculator/repayment/readjustment_strategy"
	"installment-calculator/repayment/repayment_strategy"
	"testing"
	"time"
)

func Test_PartialRepayment_Annuity_ShortenRepaymentTenors(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Annuity
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.PartialRepayment,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
		repayment.Amount(20000),
		repayment.ReadjustmentStrategy(readjustment_strategy.ShortenRepaymentTenors),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_PartialRepayment_Annuity_ReduceMonthlyRepayment(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Annuity
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.PartialRepayment,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
		repayment.Amount(20000),
		repayment.ReadjustmentStrategy(readjustment_strategy.ReduceMonthlyRepayment),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_PartialRepayment_Linear_ReduceMonthlyRepayment(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Linear
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.PartialRepayment,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
		repayment.Amount(20000),
		repayment.ReadjustmentStrategy(readjustment_strategy.ReduceMonthlyRepayment),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_PartialRepayment_Outright_ReduceMonthlyRepayment(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Outright
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.PartialRepayment,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
		repayment.Amount(20000),
		repayment.ReadjustmentStrategy(readjustment_strategy.ReduceMonthlyRepayment),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_Payoff_Annuity(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Annuity
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.Payoff,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_Payoff_Linear(t *testing.T) {
	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Linear
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.Payoff,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}

func Test_Payoff_Outright(t *testing.T) {

	borrowDate := time.Date(2023, 1, 10, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Outright
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.Payoff,
		repayment.Installment(installmentStrategy, borrowDate, 20, 100000, 36, 0.032),
		repayment.Date(time.Now()),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}
