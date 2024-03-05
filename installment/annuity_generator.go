package installment

import (
	"github.com/shopspring/decimal"
	"installment-calculator/beans"
	"installment-calculator/constants"
	"installment-calculator/helper"
	"time"
)

// AnnuityGenerator 等额本息 按月计息
type AnnuityGenerator struct {
	// 借款日
	BorrowDate time.Time
	// 还款日
	DueDay int
}

func (annuity AnnuityGenerator) Generate(amount float64, tenors int, annualPercentageRate float64) ([]beans.Installment, float64) {
	borrowAmount := decimal.NewFromFloat(amount)
	monthlyInterestRate := decimal.NewFromFloat(annualPercentageRate).Div(decimal.NewFromInt(constants.MonthsOfYear))
	firstDueDate := helper.GetFirstDueDate(annuity.DueDay, annuity.BorrowDate)

	// 每月还款额 = 每月利息 + 每月本金
	installmentAmount := annuity.calculateInstallmentAmount(borrowAmount, tenors, monthlyInterestRate)

	total := decimal.Zero
	var installments []beans.Installment
	for i := 0; i < tenors; i++ {
		dueDate := helper.CalculateDueDate(firstDueDate, i)
		interest := borrowAmount.Mul(monthlyInterestRate) // 每月利息 = 待还本金 * 月利率
		principal := installmentAmount.Sub(interest)      // 每月本金 = 每月还款额 - 每月利息

		installment := beans.Installment{
			Tenor:         i + 1,
			DueDate:       dueDate,
			Amount:        installmentAmount,
			Interest:      interest,
			Principal:     principal,
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
		total = total.Add(installmentAmount)
	}

	return installments, total.RoundBank(constants.DefaultRoundPlaces).InexactFloat64()
}

// CalculateInstallmentAmount 等额本息每月还款额计算
// 每月还款额 = 贷款金额 * 月利率 * (1+月利率)^还款月数 / （(1+月利率)^还款月数 - 1）
func (annuity AnnuityGenerator) calculateInstallmentAmount(amount decimal.Decimal, tenors int, monthlyInterestRate decimal.Decimal) decimal.Decimal {
	one := decimal.NewFromInt(constants.One)
	installmentAmount := amount.Mul(monthlyInterestRate).Mul(one.Add(monthlyInterestRate).Pow(decimal.NewFromInt(int64(tenors)))).
		Div(one.Add(monthlyInterestRate).Pow(decimal.NewFromInt(int64(tenors))).Sub(one))
	return installmentAmount
}
