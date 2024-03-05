package installment

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"installment-calculator/beans"
	"installment-calculator/constants"
	"installment-calculator/helper"
	"time"
)

// LinearGenerator 等额本金 按月计息
type LinearGenerator struct {
	// 借款日
	BorrowDate time.Time
	// 还款日
	DueDay int
}

// Generate
// 计算每月还款本金
// 计算每月利息 及还款金额
func (annuity LinearGenerator) Generate(amount float64, tenors int, annualPercentageRate float64) ([]beans.Installment, float64) {
	borrowAmount := decimal.NewFromFloat(amount)
	monthlyInterestRate := decimal.NewFromFloat(annualPercentageRate).Div(decimal.NewFromInt(constants.MonthsOfYear))
	firstDueDate := helper.GetFirstDueDate(annuity.DueDay, annuity.BorrowDate)

	principal := borrowAmount.Div(decimal.NewFromInt(int64(tenors))).RoundBank(constants.DefaultRoundPlaces)
	logrus.Infof("LinearGenerator Installment Principal : %v", principal)
	total := decimal.Zero
	var installments []beans.Installment
	for i := 0; i < tenors; i++ {
		dueDate := helper.CalculateDueDate(firstDueDate, i)
		interest := borrowAmount.Mul(monthlyInterestRate).RoundBank(constants.DefaultRoundPlaces)
		installmentAmount := principal.Add(interest)

		installment := beans.Installment{
			Tenor:         i + 1,
			DueDate:       dueDate,
			Amount:        installmentAmount,
			Principal:     principal,
			Interest:      interest,
			PrePrincipal:  borrowAmount,
			PostPrincipal: borrowAmount.Sub(principal),
		}

		//如果最后一期 还有剩余本金，则将多余的部分加到最后一期并重新计算
		if i == tenors-1 && installment.PostPrincipal != decimal.Zero {
			installment.Principal = installment.Principal.Add(installment.PostPrincipal)
			installment.Interest = installment.Principal.Mul(monthlyInterestRate).RoundBank(constants.DefaultRoundPlaces)
		}

		installments = append(installments, installment)

		borrowAmount = borrowAmount.Sub(principal)

		total = total.Add(installment.Amount)
	}

	return installments, total.RoundBank(constants.DefaultRoundPlaces).InexactFloat64()
}
