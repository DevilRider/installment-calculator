package repayment

import (
	"github.com/shopspring/decimal"
	"installment-calculator/beans"
	"installment-calculator/constants"
	"time"
)

type PayoffCalculator struct {
	installmentInfo beans.InstallmentInfo
	repaymentDate   time.Time
}

// Calculate 返回新的还款计划, 还款总金额(本金+提前还款费用), 提前还款费用(本期利息+服务费), 节省利息
// 返回所需还款金额, 根据repayment date 定位还到第几期，当期利息按天计息，额外收取还款金额（待还本金）1%为服务费，剩余本金全部还清
// 还款金额为 期初剩余本金 + 本期日息（上个还款日到提前还款时的天数 * 期初待还本金 * 日息）+ 服务费（期初本金 * 1%）
func (c PayoffCalculator) Calculate() (*[]beans.Installment, float64, float64, float64) {
	installments := backtrackInstallments(c.installmentInfo)
	currentTenor, preDueDate := locateCurrentInstallment(installments, c.repaymentDate)
	currentInstallment := installments[currentTenor-1]

	dailyInterestRate := decimal.NewFromFloat(c.installmentInfo.AnnualPercentageRate).Div(decimal.NewFromInt(constants.DaysOfYear))
	interestDays := preDueDate.Sub(c.repaymentDate).Hours() / constants.HoursOfDay
	interest := currentInstallment.PrePrincipal.Mul(dailyInterestRate).Mul(decimal.NewFromFloat(interestDays))
	serviceFee := currentInstallment.PrePrincipal.Mul(decimal.NewFromFloat(constants.RepaymentFineRate))
	amountToRepay := currentInstallment.PrePrincipal.Add(interest).Add(serviceFee)
	fee := interest.Add(serviceFee)

	return nil, amountToRepay.RoundBank(constants.DefaultRoundPlaces).InexactFloat64(),
		fee.RoundBank(constants.DefaultRoundPlaces).InexactFloat64(),
		calculateReduceInterest4Payoff(currentTenor, installments).RoundDown(2).InexactFloat64()
}
