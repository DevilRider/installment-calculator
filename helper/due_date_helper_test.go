package helper

import (
	"github.com/stretchr/testify/assert"
	"installment-calculator/utils"
	"testing"
	"time"
)

func Test_GetFirstDueDate(t *testing.T) {
	normalDay := time.Date(2019, 2, 27, 0, 0, 0, 0, time.UTC)
	day29 := time.Date(2019, 2, 29, 0, 0, 0, 0, time.UTC)
	day30 := time.Date(2019, 2, 30, 0, 0, 0, 0, time.UTC)
	day31 := time.Date(2019, 2, 31, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, 27, GetFirstDueDate(0, normalDay).Day())
	assert.Equal(t, 01, GetFirstDueDate(0, day29).Day())
	assert.Equal(t, 01, GetFirstDueDate(1, day29).Day())
	assert.Equal(t, 02, GetFirstDueDate(0, day30).Day())
	assert.Equal(t, 01, GetFirstDueDate(1, day30).Day())
	assert.Equal(t, 03, GetFirstDueDate(0, day31).Day())
	assert.Equal(t, 01, GetFirstDueDate(1, day31).Day())
}

func Test_CalculateDueDate(t *testing.T) {
	firstDueDate := time.Date(2019, 2, 27, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "2019-04-27", CalculateDueDate(firstDueDate, 2).Format(utils.DefaultDayLayout))
	assert.Equal(t, "2019-03-27", CalculateDueDate(firstDueDate, 1).Format(utils.DefaultDayLayout))
	assert.Equal(t, "2020-01-27", CalculateDueDate(firstDueDate, 11).Format(utils.DefaultDayLayout))
}
