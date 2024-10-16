package test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"installment-calculator/installment"
	"installment-calculator/installment/installment_strategy"
	"installment-calculator/repayment"
	"installment-calculator/repayment/readjustment_strategy"
	"installment-calculator/repayment/repayment_strategy"
	"installment-calculator/utils"
	"testing"
	"time"
)

// 2y 3% - 5.735%

// 22.3 10y 7.2 - 12.126
func Test_Calculate_Installment_Annuity(t *testing.T) {
	borrowDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	generator := installment.NewInstance(
		installment.Strategy(installment_strategy.Annuity),
		installment.DueDay(28),
		installment.BorrowDate(borrowDate),
	)

	amount := float64(223000)
	tenors := 120
	annualPercentageRate := 0.12126
	is1, ta1 := generator.Generate(amount, tenors, annualPercentageRate)
	assert.Equal(t, tenors, len(is1))
	logrus.Infof("total: %v, interest: %v", ta1, ta1-amount)
	logrus.Infof("installment amount: %v", is1[0].Amount.RoundBank(2).InexactFloat64())

	utils.NewOutputor(amount, ta1, is1).Output()
}

func Test_Calculate_Partial_Repayment_Annuity_ShortenRepaymentTenors(t *testing.T) {
	borrowDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)
	installmentStrategy := installment_strategy.Annuity
	amount := float64(223000)
	tenors := 120
	annualPercentageRate := 0.108

	// time.Now()
	repaymentDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	_, amount2Pay, fee, reducedInterest := repayment.TrialCalculate(repayment_strategy.PartialRepayment,
		repayment.Installment(installmentStrategy, borrowDate, 28, amount, tenors, annualPercentageRate),
		repayment.Date(repaymentDate),
		repayment.Amount(50000),
		repayment.ReadjustmentStrategy(readjustment_strategy.ShortenRepaymentTenors),
	)
	logrus.Infof("[%s]Amount2Pay: %v, ReducedInterest: %v, Fee: %v", installmentStrategy, amount2Pay, reducedInterest, fee)
}
