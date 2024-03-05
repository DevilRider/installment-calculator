package repayment

import (
	"github.com/shopspring/decimal"
	"installment-calculator/beans"
	"installment-calculator/installment"
	"time"
)

func calculateReduceInterest4Payoff(currentTenor int, installments []beans.Installment) decimal.Decimal {
	reduceInterest := decimal.Zero
	for i := currentTenor; i < len(installments); i++ {
		reduceInterest = reduceInterest.Add(installments[i].Interest)
	}
	return reduceInterest
}

// return current installment index, pre installment due date
func locateCurrentInstallment(installments []beans.Installment, repaymentDate time.Time) (int, *time.Time) {
	preDueDate := installments[0].DueDate
	for _, installment := range installments {
		if preDueDate.Before(repaymentDate) && installment.DueDate.After(repaymentDate) {
			return installment.Tenor, &preDueDate
		}
		preDueDate = installment.DueDate
	}

	panic("repayment date is invalid")
}

func backtrackInstallments(installmentInfo beans.InstallmentInfo) []beans.Installment {
	installmentPlanner := installment.NewInstance(
		installment.Strategy(installmentInfo.Strategy),
		installment.BorrowDate(installmentInfo.BorrowDate),
		installment.DueDay(installmentInfo.DueDay),
	)
	installments, _ := installmentPlanner.Generate(installmentInfo.Amount, installmentInfo.Tenors, installmentInfo.AnnualPercentageRate)
	return installments
}
