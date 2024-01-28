package ftracker

import (
	"fmt"
	"math"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// action int - количество совершенных действий (число шагов при ходьбе и беге, либо гребков при плавании).
func distance(action int) float64 {
	return float64(action) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки км/duration.
//
// Параметры:
//
// action int - количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 - длительность тренировки в часах.
func meanSpeed(action int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	distance := distance(action)
	return distance / duration
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// action int - количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// trainingType string - вид тренировки(Бег, Ходьба, Плавание).
// duration float64 - длительность тренировки в часах.
// weight float64 - вес пользователя в кг.
// height float64 - рост пользователя в см.
// lengthPool int - длинна басейна в метрах.
// countPool int - сколько раз пользователь переплыл бассейн.
func ShowTrainingInfo(action int, trainingType string, duration, weight, height float64, lengthPool, countPool int) string {
	const resultOutput string = `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`
	var trainingDistance, trainingSpeed, trainingCalories float64

	switch {
	case trainingType == "Бег":
		trainingDistance = distance(action)
		trainingSpeed = meanSpeed(action, duration)
		trainingCalories = RunningSpentCalories(action, weight, duration)
	case trainingType == "Ходьба":
		trainingDistance = distance(action)
		trainingSpeed = meanSpeed(action, duration)
		trainingCalories = WalkingSpentCalories(action, duration, weight, height)
	case trainingType == "Плавание":
		trainingDistance = distance(action)
		trainingSpeed = swimmingMeanSpeed(lengthPool, countPool, duration)
		trainingCalories = SwimmingSpentCalories(lengthPool, countPool, duration, weight)

	default:
		return "неизвестный тип тренировки"
	}
	return fmt.Sprintf(resultOutput, trainingType, duration, trainingDistance, trainingSpeed, trainingCalories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18   // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 1.79 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// action int - количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// weight float64 - вес пользователя в кг.
// duration float64 - длительность тренировки в часах.
// ((18 * СредняяСкоростьВКм/ч * 1.79) * ВесСпортсменаВКг / mInKM * ВремяТренировкиВЧасах * minInH)
func RunningSpentCalories(action int, weight, duration float64) float64 {
	return runningCaloriesMeanSpeedMultiplier * meanSpeed(action, duration) * runningCaloriesMeanSpeedShift * weight / mInKm * duration * minInH
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// action int - количество совершенных действий(число шагов при ходьбе и беге, либо гребков при плавании).
// duration float64 - длительность тренировки в часах.
// weight float64 - вес пользователя в кг.
// height float64 - рост пользователя в см.
// ((0.035 * ВесСпортсменаВКг + (СредняяСкоростьВМетрахВСекунду**2 / РостВМетрах) * 0.029 * ВесСпортсменаВКг) * ВремяТренировкиВЧасах * minInH).
func WalkingSpentCalories(action int, duration, weight, height float64) float64 {
	meanSpeedInMs := meanSpeed(action, duration) * kmhInMsec
	heightInM := height / cmInM
	return (walkingCaloriesWeightMultiplier*weight + math.Pow(meanSpeedInMs, 2)/heightInM*walkingSpeedHeightMultiplier*weight) * duration * minInH
}

// Константы для расчета калорий, расходуемых при плавании.
const (
	swimmingCaloriesMeanSpeedShift   = 1.1 // среднее количество сжигаемых колорий при плавании относительно скорости.
	swimmingCaloriesWeightMultiplier = 2   // множитель веса при плавании.
)

// swimmingMeanSpeed возвращает среднюю скорость при плавании в км/ч.
//
// Параметры:
//
// lengthPool int - длина бассейна в метрах.
// countPool int - сколько раз пользователь переплыл бассейн.
// duration float64 - длительность тренировки в часах.
func swimmingMeanSpeed(lengthPool, countPool int, duration float64) float64 {
	if duration == 0 {
		return 0
	}
	return float64(lengthPool) * float64(countPool) / mInKm / duration
}

// SwimmingSpentCalories возвращает количество потраченных калорий при плавании.
//
// Параметры:
//
// lengthPool int - длина бассейна в метрах.
// countPool int - сколько раз пользователь переплыл бассейн.
// duration float64 - длительность тренировки в часах.
// weight float64 - вес пользователя в кг.
// (СредняяСкоростьВКм/ч + 1.1) * 2 * ВесСпортсменаВКг * ВремяТренеровкиВЧасах.
func SwimmingSpentCalories(lengthPool, countPool int, duration, weight float64) float64 {
	return (swimmingMeanSpeed(lengthPool, countPool, duration) + swimmingCaloriesMeanSpeedShift) * swimmingCaloriesWeightMultiplier * weight * duration
}
