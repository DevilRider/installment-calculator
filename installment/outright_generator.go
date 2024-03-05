package installment

import (
	"github.com/shopspring/decimal"
	"installment-calculator/beans"
	"installment-calculator/constants"
	"installment-calculator/helper"
	"time"
)

// OutrightGenerator 到期还本付息 按日计息 按月付息
type OutrightGenerator struct {
	// 借款日
	BorrowDate time.Time
	// 还款日
	DueDay int
}

// Generate Plan 利息= 本金 * 日利率 * 计息天数
func (outright OutrightGenerator) Generate(amount float64, tenors int, annualPercentageRate float64) ([]beans.Installment, float64) {
	borrowAmount := decimal.NewFromFloat(amount)
	dailyInterestRate := decimal.NewFromFloat(annualPercentageRate).Div(decimal.NewFromInt(constants.DaysOfYear))
	firstDueDate := helper.GetFirstDueDate(outright.DueDay, outright.BorrowDate)

	preDueDate := outright.BorrowDate
	total := decimal.Zero
	var installments []beans.Installment
	for i := 0; i < tenors; i++ {
		dueDate := helper.CalculateDueDate(firstDueDate, i)
		interestDays := decimal.NewFromInt(int64(dueDate.Sub(preDueDate).Hours() / constants.HoursOfDay))
		if tenors != 0 {
			preDueDate = dueDate
		}
		interest := borrowAmount.Mul(dailyInterestRate).Mul(interestDays).RoundBank(constants.DefaultRoundPlaces) // 每月利息 = 待还本金 * 日利率 * 计息天数

		installmentAmount := interest
		if i == tenors-1 {
			installmentAmount = installmentAmount.Add(borrowAmount)
		}

		installment := beans.Installment{
			Tenor:         i + 1,
			DueDate:       dueDate,
			Amount:        installmentAmount,
			Interest:      interest,
			Principal:     borrowAmount,
			PrePrincipal:  borrowAmount,
			PostPrincipal: borrowAmount,
		}

		installments = append(installments, installment)
		total = total.Add(installmentAmount)
	}

	return installments, total.InexactFloat64()
}
