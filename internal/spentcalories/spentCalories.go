package spentcalories

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                            = 0.65  // средняя длина шага.
	mInKm                              = 1000  // количество метров в километре.
	minInH                             = 60    // количество минут в часе.
	runningCaloriesMeanSpeedMultiplier = 18.0  // множитель средней скорости при беге.
	runningCaloriesMeanSpeedShift      = 20.0  // среднее количество сжигаемых калорий при беге.
	walkingCaloriesWeightMultiplier    = 0.035 // множитель массы тела при ходьбе.
	walkingSpeedHeightMultiplier       = 0.029 // множитель роста при ходьбе.
)

var (
	errData = errors.New("error data. No data for conversion [sc]")
	errConv = errors.New("error conversion. The data has not been converted correctly [sc]")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	if len(data) == 0 {
		return 0, "", 0, errData
	}

	dataParse := strings.Split(data, ",")
	if len(dataParse) != 3 {
		return 0, "", 0, errConv
	}

	isStringSteps := reflect.TypeOf(dataParse[0])
	isStringActivity := reflect.TypeOf(dataParse[1])
	isStringDuration := reflect.TypeOf(dataParse[2])
	isString := reflect.TypeOf("String")

	if (isStringSteps != isString) || (isStringActivity != isString) || (isStringDuration != isString) {
		return 0, "", 0, errConv
	}

	steps, err := strconv.Atoi(dataParse[0])
	if err != nil {
		return 0, "", 0, err
	} else if steps <= 0 {
		return 0, "", 0, errConv
	}

	activity := dataParse[1]

	duration, err := time.ParseDuration(dataParse[2])
	if err != nil {
		return 0, "", 0, err
	} else if duration <= 0 {
		return 0, "", 0, errConv
	}

	return steps, activity, duration, nil

}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	if steps <= 0 {
		return 0
	}

	distance := (float64(steps) * lenStep) / float64(mInKm)

	return distance

}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if (steps <= 0) || (duration <= 0) {
		return 0
	}

	distance := distance(steps)
	durationInHours := duration.Hours()

	averageSpeed := distance / durationInHours

	return averageSpeed

}

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	if (steps <= 0) || (weight <= 0) || (duration <= 0) {
		return 0
	}

	meanSpeed := meanSpeed(steps, duration)

	spentCalories := ((runningCaloriesMeanSpeedMultiplier * meanSpeed) - runningCaloriesMeanSpeedShift) * weight

	return spentCalories
}

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	if (steps <= 0) || (weight <= 0) || (height <= 0) || (duration <= 0) {
		return 0
	}

	meanSpeed := meanSpeed(steps, duration)
	durationInHours := duration.Hours()

	spentCalories := ((walkingCaloriesWeightMultiplier * weight) + (meanSpeed*meanSpeed/height)*walkingSpeedHeightMultiplier) * durationInHours * float64(minInH)

	return spentCalories
}

// TrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	if (len(data) == 0) || (weight <= 0) || (height <= 0) {
		return ""
	}

	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return ""
	} else if (steps == 0) || (len(activity) == 0) || (duration == 0) {
		return ""
	}

	durationInHours := duration.Hours()
	distance := distance(steps)
	averageSpeed := meanSpeed(steps, duration)

	var spentCalories float64 = 0

	switch activity {
	case "Ходьба":
		spentCalories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		spentCalories = RunningSpentCalories(steps, weight, duration)
	default:
		return "Неизвестный тип тренировки"
	}

	title := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f", activity, durationInHours, distance, averageSpeed, spentCalories)

	return title
}
