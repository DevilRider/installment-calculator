package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"installment-calculator/beans"
)

const outputPath = "../app-output/"
const splitString = "-----------------------------------------------"

type Outputor struct {
	total        decimal.Decimal
	principal    decimal.Decimal
	installments []beans.Installment

	// calculate from above
	tenors            int
	years             int
	installmentAmount decimal.Decimal
	interest          decimal.Decimal
}

func NewOutputor(principal float64, total float64, installments []beans.Installment) *Outputor {
	o := &Outputor{
		total:        decimal.NewFromFloat(total),
		principal:    decimal.NewFromFloat(principal),
		installments: installments}

	o.tenors = len(installments)
	o.years = o.tenors / 12
	o.interest = o.total.Sub(o.principal).RoundBank(2)
	o.installmentAmount = o.installments[0].Amount
	return o
}

func (o Outputor) Output() {
	filename := fmt.Sprintf("%s %sw~%dYears[%s].txt", outputPath, o.principal.Div(decimal.NewFromInt(10000)).RoundBank(2).String(), o.years, o.installmentAmount.RoundBank(2).String())
	WriteCsv(filename, o.rows())
	logrus.Infof("filename: %s", filename)
}

func (o Outputor) rows() [][]string {
	var rows [][]string
	rows = append(rows, []string{splitString})
	for _, val := range o.header() {
		rows = append(rows, val)
	}
	rows = append(rows, []string{splitString})
	rows = append(rows, []string{"Tenor", "Principal", "Interest", "Left Principal"})
	for _, val := range o.results() {
		rows = append(rows, val)
	}
	return rows
}

func (o Outputor) header() [][]string {
	var rows [][]string
	rows = append(rows, []string{fmt.Sprintf("Borrow %s for %d Years.", o.principal.String(), o.years)})
	rows = append(rows, []string{fmt.Sprintf("Installment Amount: %s", o.installmentAmount.RoundBank(2).String())})
	rows = append(rows, []string{fmt.Sprintf("Total Repayment Amount: %s", o.total.String())})
	rows = append(rows, []string{fmt.Sprintf("Repayment Interest: %s", o.interest.RoundBank(2).String())})

	return rows
}

func (o Outputor) results() [][]string {
	var rows [][]string
	for _, i := range o.installments {
		rows = append(rows, []string{
			fmt.Sprintf("%d", i.Tenor),
			i.Principal.RoundBank(2).String(),
			i.Interest.RoundBank(2).String(),
			i.PostPrincipal.RoundBank(2).String()})
	}
	return rows
}
