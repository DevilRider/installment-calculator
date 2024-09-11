package repayment

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"installment-calculator/beans"
	"installment-calculator/constants"
	"installment-calculator/installment"
	"installment-calculator/repayment/readjustment_strategy"
	"installment-calculator/utils"
	"time"
)

type PartialRepaymentCalculator struct {
	installmentInfo      beans.InstallmentInfo
	repaymentDate        time.Time
	repaymentAmount      int64
	readjustmentStrategy readjustment_strategy.Strategy // ReduceMonthlyRepayment, ShortenRepaymentTenors
}

// Calculate 返回新的还款计划, 还款总金额, 提前还款费用(本期利息+服务费), 节省利息
// 服务费: 1% * 提前还款金额 服务费不加在还款总额里，需额外计算
// 采用期初本金进行利息计算，
//
//	interest = pre principal * daily interest rate * interest days(repayment date - pre due date)
//
// 使用本期还款后金额进行新的分期计划生成
//
//	 还款期限不变，减少月供
//		borrow amount = pre principal - repayment amount
//		borrow date = repayment date
//		tenors = tenors - current tenor
func (calculator PartialRepaymentCalculator) Calculate() (*[]beans.Installment, float64, float64, float64) {
	installments := backtrackInstallments(calculator.installmentInfo)
	currentTenor, preDueDate := locateCurrentInstallment(installments, calculator.repaymentDate)
	if currentTenor == 0 {
		panic("repayment date is invalid")
	}
	currentInstallment := installments[currentTenor-1]
	dailyInterestRate := decimal.NewFromFloat(calculator.installmentInfo.AnnualPercentageRate).Div(decimal.NewFromInt(constants.DaysOfYear))
	interestDays := preDueDate.Sub(calculator.repaymentDate).Hours() / 24
	if interestDays < 0 {
		interestDays = 0
	}
	interest := currentInstallment.PrePrincipal.Mul(dailyInterestRate).Mul(decimal.NewFromFloat(interestDays))
	fineFee := decimal.NewFromFloat(constants.RepaymentFineRate).Mul(decimal.NewFromInt(calculator.repaymentAmount))

	fee := interest.Add(fineFee).RoundBank(2).InexactFloat64()
	originalInstallmentInterest := calculateReduceInterest4Payoff(currentTenor, installments)
	p := installment.NewInstance(
		installment.Strategy(calculator.installmentInfo.Strategy),
		installment.DueDay(calculator.installmentInfo.DueDay),
		installment.BorrowDate(calculator.repaymentDate),
	)
	amount := currentInstallment.PrePrincipal.Sub(decimal.NewFromInt(calculator.repaymentAmount))
	tenors := calculator.installmentInfo.Tenors - currentTenor

	switch calculator.readjustmentStrategy {
	case readjustment_strategy.ReduceMonthlyRepayment:
		newInstallments, total := p.Generate(amount.InexactFloat64(), tenors, calculator.installmentInfo.AnnualPercentageRate)
		finalInstallmentInterest := decimal.NewFromFloat(total).Sub(amount)
		logrus.Infof("[%s]ReduceMonthlyRepayment: original installments: %v(%d), new installments: %v(%d)", calculator.installmentInfo.Strategy, currentInstallment.Amount.RoundBank(constants.DefaultRoundPlaces), calculator.installmentInfo.Tenors, newInstallments[0].Amount.RoundBank(constants.DefaultRoundPlaces), tenors)
		utils.NewOutputor(amount.InexactFloat64(), total, newInstallments).Output()
		return &newInstallments, total, fee, originalInstallmentInterest.Sub(finalInstallmentInterest).RoundDown(2).InexactFloat64()
	case readjustment_strategy.ShortenRepaymentTenors:
		newInstallments, total := calculator.rematchInstallments(tenors, amount, currentInstallment.Amount, calculator.installmentInfo.AnnualPercentageRate)
		finalInstallmentInterest := decimal.NewFromFloat(total).Sub(amount)
		logrus.Infof("[%s]ShortenRepaymentTenors: original installments: %v(%d), new installments: %v(%d)", calculator.installmentInfo.Strategy, currentInstallment.Amount.RoundBank(constants.DefaultRoundPlaces), calculator.installmentInfo.Tenors, newInstallments[0].Amount.RoundBank(constants.DefaultRoundPlaces), len(newInstallments))
		utils.NewOutputor(amount.InexactFloat64(), total, newInstallments).Output()
		return &newInstallments, total, fee, originalInstallmentInterest.Sub(finalInstallmentInterest).RoundDown(2).InexactFloat64()
	default:
		panic("unsupported readjustment_strategy strategy")
	}
}

// 月供不变仅对等额本息生效
// amount 为还款后剩余本金； 还款日，借款日不变
// 循环计算 直到分期金额与当前分期金额相差最小
func (calculator PartialRepaymentCalculator) rematchInstallments(tenors int, amount decimal.Decimal, target decimal.Decimal, AnnualPercentageRate float64) ([]beans.Installment, float64) {
	p := installment.NewInstance(
		installment.DueDay(calculator.installmentInfo.DueDay),
		installment.BorrowDate(calculator.installmentInfo.BorrowDate),
	)
	threshold := decimal.NewFromInt(constants.One)
	trs := tenors
	for {
		is, t := p.Generate(amount.InexactFloat64(), trs, AnnualPercentageRate)
		logrus.Debugf("rematch installments: is[0]: %v, target: %v", is[0].Amount.RoundBank(constants.DefaultRoundPlaces), target.RoundBank(constants.DefaultRoundPlaces))
		if target.Sub(is[0].Amount).Abs().LessThan(threshold) {
			logrus.Debugf("[AnnuityGenerator]ShortenRepaymentTenors: original installments: %v(%d), new installments: %v(%d)", target.RoundBank(constants.DefaultRoundPlaces), tenors, is[0].Amount.RoundBank(constants.DefaultRoundPlaces), trs)
			return is, t
		} else {
			trs--
		}

		if trs <= 0 {
			// upper threshold and retry up by one
			trs = tenors
			logrus.Debugf("upper threshold: %v, tenors: %d, trs: %d", threshold, tenors, trs)
			threshold = threshold.Add(decimal.NewFromInt(constants.One))
		}
	}
}
