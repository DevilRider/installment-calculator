package test

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"installment-calculator/installment"
	"installment-calculator/installment/installment_strategy"
	"installment-calculator/utils"
	"testing"
	"time"
)

// https://fin.paas.cmbchina.com/fininfo/calloanper

func Test_Installment_Annuity(t *testing.T) {
	generator := installment.NewInstance(
		installment.Strategy(installment_strategy.Annuity),
		installment.DueDay(9),
		installment.BorrowDate(time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC)),
	)

	is1, ta1 := generator.Generate(100000, 12, 0.0345)
	assert.Equal(t, 12, len(is1))
	assert.Equal(t, 101878.59, ta1)
	assert.Equal(t, 8489.88, is1[0].Amount.RoundBank(2).InexactFloat64())

	is2, ta2 := generator.Generate(350000, 360, 0.07)
	assert.Equal(t, 360, len(is2))
	assert.Equal(t, 838281.14, ta2)
	assert.Equal(t, 2328.56, is2[0].Amount.RoundBank(2).InexactFloat64())
}

func Test_Installment_Linear(t *testing.T) {

	generator := installment.NewInstance(
		installment.Strategy(installment_strategy.Linear),
		installment.DueDay(9),
		installment.BorrowDate(time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC)),
	)

	is1, ta1 := generator.Generate(100000, 12, 0.0345)
	assert.Equal(t, 12, len(is1))
	assert.Equal(t, 101868.72, ta1) // 8333.33
	assert.Equal(t, 8620.83, is1[0].Amount.RoundBank(2).InexactFloat64())

	is2, ta2 := generator.Generate(350000, 360, 0.07)
	assert.Equal(t, 360, len(is2))
	assert.Equal(t, 718520.88, ta2)
	assert.Equal(t, 3013.89, is2[0].Amount.RoundBank(2).InexactFloat64())
	assert.Equal(t, 2951.50, is2[11].Amount.RoundBank(2).InexactFloat64())
}

func Test_Installment_Outright(t *testing.T) {
	generator := installment.NewInstance(
		installment.Strategy(installment_strategy.Outright),
		installment.DueDay(10),
		installment.BorrowDate(time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)),
	)

	is1, ta1 := generator.Generate(100000, 12, 0.0345)
	assert.Equal(t, 12, len(is1))
	assert.Equal(t, 103507.48, ta1) // 103507.5
	assert.Equal(t, 297.08, is1[0].Amount.RoundBank(2).InexactFloat64())
	assert.Equal(t, 100297.08, is1[11].Amount.RoundBank(2).InexactFloat64())

	is2, ta2 := generator.Generate(350000, 60, 0.07)
	assert.Equal(t, 60, len(is2))
	assert.Equal(t, 474337.50, ta2)
	assert.Equal(t, 2109.72, is2[0].Amount.RoundBank(2).InexactFloat64())
	assert.Equal(t, 2041.67, is2[3].Amount.RoundBank(2).InexactFloat64())
	assert.Equal(t, 1905.56, is2[13].Amount.RoundBank(2).InexactFloat64())
}

func Test_Installment_Annuity_r(t *testing.T) {
	generator := installment.NewInstance(
		installment.Strategy(installment_strategy.Annuity),
		installment.DueDay(28),
		installment.BorrowDate(time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC)),
	)

	amount := float64(225000)
	tenors := 120
	annualPercentageRate := 0.11
	is1, ta1 := generator.Generate(amount, tenors, annualPercentageRate)
	assert.Equal(t, tenors, len(is1))
	logrus.Infof("total: %v, interest: %v", ta1, ta1-amount)
	logrus.Infof("installment amount: %v", is1[0].Amount.RoundBank(2).InexactFloat64())

	utils.NewOutputor(amount, ta1, is1).Output()
}
