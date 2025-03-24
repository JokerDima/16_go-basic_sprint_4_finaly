package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

const (
	StepLength = 0.65 // длина шага в метрах
	mInKm      = 1000 // количество метров в километре.
)

func parsePackage(data string) (int, time.Duration, error) {

	dataParse := strings.Split(data, ",")
	if len(dataParse) != 2 {
		return 0, 0, errors.New("error conversion. The data has not been converted correctly")
	}

	steps, err := strconv.Atoi(dataParse[0])
	if (err != nil) || (steps <= 0) {
		return 0, 0, errors.New("error data. Step count error")
	}

	duration, err := time.ParseDuration(dataParse[1])
	if (err != nil) || (duration <= 0) {
		return 0, 0, errors.New("error data. Duration error")
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
	if (len(data) == 0) || (weight == 0) || (height == 0) {
		return ""
	}

	steps, duration, err := parsePackage(data)
	if err != nil {
		return err.Error()
	}

	distance := float64(steps) * StepLength
	distanceKm := distance / mInKm

	spentCalories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	title := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceKm, spentCalories)

	return title
}
