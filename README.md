# Installment Calculator
* 提供两大功能，分期计算器和提前还款计算器
  * golang 分期计算器， 支持等额本息，等额本金，到期还本付息等分期方式的计算。
  * 提前还款计算器，支持提前还款后的还款计划调整。

## Installment Generator
* Installment strategy
  * Annuity
  * Linear
  * Outright

### Demo
```golang
  generator := installment.NewInstance(
        installment.Strategy(installment_strategy.Annuity),
        installment.DueDay(9),
        installment.BorrowDate(time.Date(2024, 1, 9, 0, 0, 0, 0, time.UTC)), 
    )
  
  installments, total := generator.Generate(100000, 12, 0.0345) // borrow amount, tenors, interest rate
```

## Early Repayment Calculator
* Early repayment strategy
  * Payoff
  * Partial Repayment

* Readjustment strategy
  * Reduce Monthly Repayment
  * Shorten Repayment Tenors