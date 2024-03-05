package helper

import (
	"installment-calculator/constants"
	"time"
)

func GetFirstDueDate(dueDay int, borrowDate time.Time) time.Time {
	if dueDay != 0 {
		return calculateFirstDueDate(dueDay, borrowDate)
	} else {
		return addMonth(borrowDate, 1)
	}
}

func CalculateDueDate(firstDueDate time.Time, installmentNumber int) time.Time {
	return addMonth(firstDueDate, installmentNumber)
}

// calculateFirstDueDate
// 因为还款日只有日期，所以还款日存在小于借款日的情况，
// 如果还款日与借款日的差值 小于 {@LeastInstallmentDays} 设定的最小账单周期，则还款日后延一个月
func calculateFirstDueDate(dueDay int, borrowDate time.Time) time.Time {
	tmpDate := time.Date(borrowDate.Year(), borrowDate.Month(), dueDay, 0, 0, 0, 0, time.UTC)
	for tmpDate.Day() != dueDay {
		tmpDate.AddDate(0, -2, 0)
		tmpDate = time.Date(tmpDate.Year(), tmpDate.Month(), dueDay, 0, 0, 0, 0, time.UTC)
	}

	for tmpDate.Sub(borrowDate).Hours()/constants.HoursOfDay <= constants.LeastInstallmentDays {
		tmpDate = tmpDate.AddDate(0, constants.DefaultDateIncreaseStep, 0)
	}
	return tmpDate
}

// addMonth 还款日 月累加
// 29->1, 30->2, 31-> 3,如果大于28号，需要延迟到下下个月还款，所以需要修改日期，并且多加一个月
func addMonth(firstDueDate time.Time, months int) time.Time {
	tempDate := firstDueDate
	year := tempDate.Year()
	month := tempDate.Month()
	day := tempDate.Day()

	if day > 28 {
		tempDate = time.Date(year, month, day-28, 0, 0, 0, 0, time.UTC)
		tempDate.AddDate(0, months, 0)
	}

	return tempDate.AddDate(0, months, 0)
}
