package Interpolation

import (
	"errors"
	"math"
)

type (
	anyFloat interface {
		float32 | float64
	}
	anyInt interface {
		int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
	}
	anyUInt interface {
		uint | uint8 | uint16 | uint32 | uint64
	}
	anySInt interface {
		int | int8 | int16 | int32 | int64
	}
	anyNum interface {
		int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
	}
)

func HalfLengthValueSearcher[T anyNum](Function func(T) T, Xmin T, Xmax T, TargetY T, Accuracy T) (T, error) {

	x_lower := Xmin
	x_higher := Xmax
	y_lower := Function(x_lower)
	y_higher := Function(x_higher)
	if (TargetY > y_higher && TargetY > y_lower) || (TargetY < y_higher && TargetY < y_lower) {
		return 0, errors.New("out of range")
	}
	for math.Abs(float64(y_higher-y_lower)) > float64(Accuracy) {
		x_mid := (x_higher + x_lower) / 2
		y_mid := Function(x_mid)
		if y_mid > TargetY {
			x_higher = x_mid
			y_higher = y_mid
		} else {
			x_lower = x_mid
			y_lower = y_mid
		}
	}
	return (x_higher + x_lower) / 2, nil
}

func Round[T anyFloat](DecimalPlaces int, x ...*T) {
	k := T(math.Pow10(DecimalPlaces))
	for i := 0; i < len(x); i++ {
		*x[i] = T(math.Round(float64(*x[i]*k))) / k
	}
}

// Returns result of linear interpolation
// f( x ) = a * x + b, f( x1 ) = val1, f( x2 ) = val2, f( target_x ) = result
func Linear[T anyNum](target_x T, x1 T, x2 T, val1 T, val2 T) T {
	if x1 == x2 {
		return val1
	} else {
		var target_x_f, x1_f, x2_f, val1_f, val2_f float64
		target_x_f = float64(target_x)
		x1_f = float64(x1)
		x2_f = float64(x2)
		val1_f = float64(val1)
		val2_f = float64(val2)
		return T(-((x2_f*val1_f - x1_f*val2_f) / (x1_f - x2_f)) - ((-val1_f+val2_f)*target_x_f)/(x1_f-x2_f))
	}
}

func Linear2[T anyNum](targetX T, arrayOfX []T, ValueArray []T) (Result T, err error) {
	var (
		higherXid int
		lowerXid  int
	)
	if len(arrayOfX) != len(ValueArray) {
		return 0, errors.New("bad array sizes")
	}
	higherXid, lowerXid, err = SearchNearestId(targetX, &arrayOfX)
	if err != nil {
		return 0, err
	}
	if err != nil {
		return 0, err
	}
	return Linear(targetX, arrayOfX[lowerXid], arrayOfX[higherXid], ValueArray[lowerXid], ValueArray[higherXid]), nil
}

func BiLinear[T anyNum](targetX T, x1 T, x2 T, targetY T, y1 T, y2 T, valX1Y1 T, valX1Y2 T, valX2Y1 T, valX2Y2 T) T {
	if x1 == x2 {
		if y1 == y2 {
			return valX1Y1
		} else {
			return Linear(targetY, y1, y2, valX1Y1, valX1Y2)
		}
	} else if y1 == y2 {
		return Linear(targetX, x1, x2, valX1Y1, valX2Y1)
	} else {
		return Linear(targetX, x1, x2, Linear(targetY, y1, y2, valX1Y1, valX1Y2), Linear(targetY, y1, y2, valX2Y1, valX2Y2))
	}
}

func BiLinear2[T anyNum](targetX T, targetY T, arrayOfX []T, arrayOfY []T, ValueArray [][]T) (Result T, err error) {
	var (
		higherXid int
		lowerXid  int
		higherYid int
		lowerYid  int
	)
	if len(arrayOfX) != len(ValueArray) {
		return 0, errors.New("bad array sizes")
	}
	higherXid, lowerXid, err = SearchNearestId(targetX, &arrayOfX)
	if err != nil {
		return 0, err
	}
	if len(arrayOfY) != len(ValueArray[higherXid]) {
		return 0, errors.New("bad array sizes")
	}
	if len(arrayOfY) != len(ValueArray[lowerXid]) {
		return 0, errors.New("bad array sizes")
	}
	higherYid, lowerYid, err = SearchNearestId(targetY, &arrayOfY)
	if err != nil {
		return 0, err
	}
	return BiLinear(targetX, arrayOfX[lowerXid], arrayOfX[higherXid], targetY, arrayOfY[lowerYid], arrayOfY[higherYid], ValueArray[lowerXid][lowerYid], ValueArray[lowerXid][higherYid], ValueArray[higherXid][lowerYid], ValueArray[higherXid][higherYid]), nil
}

// This function search for the nearest higher and lower values in low to high sorted array & return it's indexes
func SearchNearestId[T anyNum](val T, arr *[]T) (int, int, error) {
	var (
		lower_id  int
		higher_id int
		err       error
	)
	err = errors.New("out of range")
	for i := 0; i < len(*arr)-1; i++ {
		if (*arr)[i] < val && (*arr)[i+1] > val {
			lower_id = i
			higher_id = i + 1
			err = nil
			break
		} else if (*arr)[i] == val {
			lower_id = i
			higher_id = i
			err = nil
			break
		}
	}
	if (*arr)[len(*arr)-1] == val {
		lower_id = len(*arr) - 1
		higher_id = len(*arr) - 1
		err = nil
	}
	return lower_id, higher_id, err
}

// Approximation by a 2nd order polynomial passing through zero
func ZeroParabolicApproximation[T anyNum](x1 T, x2 T, val1 T, val2 T, target_x T) (T, error) {
	if x1 == 0 || x2 == 0 || x1 == x2 {
		return 0, errors.New("divide dy zero")
	}

	a := -((-x2*val1 + x1*val2) / (x1 * (x1 - x2) * x2))
	b := -((x2*x2*val1 - x1*x1*val2) / (x1 * (x1 - x2) * x2))
	return a*target_x*target_x + b*target_x, nil
}
