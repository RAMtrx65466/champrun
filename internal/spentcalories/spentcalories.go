package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчёта длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчёта калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается три элемента")
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка преобразования количества шагов: %v", err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("количество шагов должно быть положительным")
	}

	activityType := parts[1]

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка парсинга продолжительности: %v", err)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("продолжительность должна быть больше нуля")
	}

	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	totalDistanceMeters := float64(steps) * stepLength
	return totalDistanceMeters / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	distKm := distance(steps, height)
	durationHours := duration.Hours()
	return distKm / durationHours
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	calories := (weight * speed * durationMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, fmt.Errorf("некорректные входные параметры")
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()
	baseCalories := (weight * speed * durationMinutes) / minInH
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	distanceKm := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch activityType {
	case "Ходьба":
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Бег":
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activityType)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType, duration.Hours(), distanceKm, speed, calories,
	), nil
}
