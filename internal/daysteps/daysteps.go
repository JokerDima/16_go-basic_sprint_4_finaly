package daysteps

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

const StepLength = 0.65 // длина шага в метрах

var (
	errData = errors.New("error data. No data for conversion [ds]")
	errConv = errors.New("error conversion. The data has not been converted correctly [ds]")
)

func parsePackage(data string) (int, time.Duration, error) {
	if len(data) == 0 {
		return 0, 0, errData
	}

	dataParse := strings.Split(data, ",")
	if len(dataParse) != 2 {
		return 0, 0, errConv
	}

	isStringSteps := reflect.TypeOf(dataParse[0])
	isStringDuration := reflect.TypeOf(dataParse[1])
	isString := reflect.TypeOf("String")

	if isStringSteps != isString || isStringDuration != isString {
		return 0, 0, errConv
	}

	steps, err := strconv.Atoi(dataParse[0])
	if err != nil {
		return 0, 0, err
	} else if steps <= 0 {
		return 0, 0, errConv
	}

	duration, err := time.ParseDuration(dataParse[1])
	if err != nil {
		return 0, 0, err
	} else if duration <= 0 {
		return 0, 0, errConv
	}

	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	} else if (steps == 0) || (duration == 0) {
		return ""
	}

	distance := float64(steps) * StepLength
	distanceKm := distance * 1000

	spentCalories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	title := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, spentCalories)

	return title
}
