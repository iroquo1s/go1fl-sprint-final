package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("Ошибка: неверный формат данных")
	}

	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, errors.New("Ошибка: неверно введены шаги")
	}

	if steps <= 0 {
		return 0, "", 0, errors.New("Ошибка: число введенных шагов должно быть положительным")
	}

	activity := strings.TrimSpace(parts[1])
	if activity == "" {
		return 0, "", 0, errors.New("Ошибка: активность не может быть пустой")
	}

	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)

	if err != nil {
		return 0, "", 0, errors.New("Ошибка: неккоректно введена продолжительность")
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("Ошибка: продолжительность должна быть положительным числом")
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {

	stepLength := height * stepLengthCoefficient
	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)
	hours := duration.Hours()
	speed := dist / hours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", err
	}

	if weight <= 0 || height <= 0 {
		return "", errors.New("Ошибка: вес и рост должны быть положительными")
	}

	var distanceKm, speed, calories float64
	var errCalc error

	switch strings.ToLower(activity) {
	case "бег":
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, errCalc = RunningSpentCalories(steps, weight, height, duration)
	case "ходьба":
		distanceKm = distance(steps, height)
		speed = meanSpeed(steps, height, duration)
		calories, errCalc = WalkingSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("Ошибка: неизвестный тип тренировки")
	}

	if errCalc != nil {
		return "", errCalc
	}

	durationHours := fmt.Sprintf("%.2f", duration.Hours())

	info := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %s ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		activity,
		durationHours,
		distanceKm,
		speed,
		calories,
	)

	return info, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Ошибка: число шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("Ошибка: вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("Ошибка: рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("Ошибка: продолжительность должна быть положительным числом")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Ошибка: шагов должно быть положительным")
	}
	if weight <= 0 {
		return 0, errors.New("Ошибка: вес должен быть положительным")
	}
	if height <= 0 {
		return 0, errors.New("Ошибка: рост должен быть положительным")
	}
	if duration <= 0 {
		return 0, errors.New("Ошибка: продолжительность должна быть положительным числом")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	baseCalories := (weight * speed * durationMinutes) / minInH
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}
