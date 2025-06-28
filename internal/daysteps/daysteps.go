package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Ошибка: некорректный формат данных: ожидается две части")
	}

	stepsStr := parts[0]
	durationStr := parts[1]

	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка: ошибка преобразования количества шагов: %w", err)
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка: ошибка преобразования продолжительности: %w", err)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Ошибка:%s", err)
		return ""
	}

	if steps <= 0 {
		err := errors.New("parsePackage: the number of steps is zero")
		log.Println(err) // Логируем ошибку
		return ""
	}

	distanceM := float64(steps) * stepLength
	distanceKm := distanceM / mInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Printf("Ошибка: ошибка расчета калорий:%s", err)
		return ""
	}

	info := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps,
		distanceKm,
		calories,
	)

	return info
}
