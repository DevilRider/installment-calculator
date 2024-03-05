package beans

import (
	"github.com/shopspring/decimal"
	"installment-calculator/installment/installment_strategy"
	"time"
)

type Installment struct {

	//分期号
	Tenor int

	//还款日
	DueDate time.Time

	//分期还款 每期应还金额  本金+利息(不包含服务费)
	Amount decimal.Decimal

	//利息
	Interest decimal.Decimal

	//本金
	Principal decimal.Decimal

	//期初本金
	PrePrincipal decimal.Decimal

	//期末本金
	PostPrincipal decimal.Decimal
}

type InstallmentInfo struct {
	Strategy             installment_strategy.Strategy
	BorrowDate           time.Time
	DueDay               int
	Amount               float64
	Tenors               int
	AnnualPercentageRate float64
}
