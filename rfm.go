package main

import (
	"math"
	"reflect"

	"gonum.org/v1/gonum/stat"
)

type Class int

const (
	H Class = iota
	G
	F
	E
	D
	C
	B
	A
)

type RFM[T Customer] interface {
	//取得因子標準差
	getStandardDeviationByFactor(factor FactorProperty) float64
	//取得因子平均數
	getAverageByFactor(factor FactorProperty) float64
	//取得客戶某因子的得分
	getCustomerScoreByFactor(c T, factor FactorProperty) float64
	//取得所有客戶得分平均值
	getScoreAvg() float64
	//取得所有客戶得分標準差
	getScoreStandardDeviation() float64
	//取得群體，透過class
	GetClass(c Class) []T
	//取得顧客總得分
	getCustomerScore(c T) float64
}

type rfm[T Customer] struct {
	data        []T
	factorAvg   map[string]float64
	factorStdev map[string]float64
	scoreAvg    float64
	scoreStdev  float64
}

func NewRFM[T Customer](data []T) RFM[T] {
	r := &rfm[T]{
		data:        data,
		factorAvg:   make(map[string]float64),
		factorStdev: make(map[string]float64),
	}
	r.initFactor()
	return r
}
func (r *rfm[T]) initFactor() {
	f := r.data[0].GetFactorName()

	for i := 0; i < len(f); i++ {
		factor := f[i]
		r.factorAvg[factor.Name] = r.getAverageByFactor(factor)
		r.factorStdev[factor.Name] = r.getStandardDeviationByFactor(factor)
	}
	r.scoreAvg = r.getScoreAvg()
	r.scoreStdev = r.getScoreStandardDeviation()
}
func (r *rfm[T]) GetClass(c Class) []T {
	var result map[Class][]T = make(map[Class][]T)
	for _, v := range r.data {
		s := r.getCustomerScore(v)
		key := Class(int(s))
		if _, ok := result[key]; !ok {
			result[key] = make([]T, 0)
		}
		result[key] = append(result[key], v)
	}
	return result[c]
}
func (r *rfm[T]) getCustomerScore(c T) float64 {
	var score float64
	f := c.GetFactorName()
	for i := 0; i < len(f); i++ {
		score += r.getCustomerScoreByFactor(c, f[i])
	}
	var result float64
	scoreAvg := r.scoreAvg
	scoreStdev := r.scoreStdev

	if (score - scoreAvg) > 0 {
		result = math.Ceil((score - scoreAvg) / scoreStdev)
	} else {
		result = math.Floor((score - scoreAvg) / scoreStdev)
	}
	if result > 4 {
		result = 4
	}
	if result < -4 {
		result = -4
	}
	if result > 0 {
		result += 3
	} else {
		result += 4
	}
	return result
}
func (r *rfm[T]) getScoreStandardDeviation() float64 {
	var total []float64
	f := r.data[0].GetFactorName()
	for _, v := range r.data {
		var score float64
		for i := 0; i < len(f); i++ {
			score += r.getCustomerScoreByFactor(v, f[i])
		}
		total = append(total, score)
	}
	result := stat.StdDev(total, nil)
	return result
}
func (r *rfm[T]) getScoreAvg() float64 {
	var score float64
	f := r.data[0].GetFactorName()
	for _, v := range r.data {
		for i := 0; i < len(f); i++ {
			score += r.getCustomerScoreByFactor(v, f[i])
		}
	}
	return score / float64(len(r.data))
}
func (r *rfm[T]) getCustomerScoreByFactor(c T, factor FactorProperty) float64 {
	value := reflect.ValueOf(c)
	fieldValue := value.FieldByName(factor.Name).Float()
	var avg, stdev, result float64

	avg = r.factorAvg[factor.Name]
	stdev = r.factorStdev[factor.Name]

	if (fieldValue - avg) >= 0 {
		result = math.Ceil((fieldValue - avg) / stdev)
	} else {
		result = math.Floor((fieldValue - avg) / stdev)
	}
	if result > 4 {
		result = 4
	}
	if result < -4 {
		result = -4
	}
	if factor.ReversWeight {
		result = result * -1
	}

	if result > 0 {
		result += 3
	} else {
		result += 4
	}
	return result
}
func (r *rfm[T]) getAverageByFactor(factor FactorProperty) float64 {
	var total float64
	for _, v := range r.data {
		value := reflect.ValueOf(v)
		fieldValue := value.FieldByName(factor.Name).Float()
		total += fieldValue
	}
	result := total / float64(len(r.data))
	return result
}
func (r *rfm[T]) getStandardDeviationByFactor(factor FactorProperty) float64 {
	var total []float64
	for _, v := range r.data {
		value := reflect.ValueOf(v)
		fieldValue := value.FieldByName(factor.Name).Float()
		total = append(total, fieldValue)
	}
	result := stat.StdDev(total, nil)
	return result
}

func (c Class) GetClass() string {
	switch c {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	case E:
		return "E"
	case F:
		return "F"
	case G:
		return "G"
	case H:
		return "H"
	default:
		return "Z"
	}
}

func StringToClass(s string) Class {
	switch s {
	case "A":
		return A
	case "B":
		return B
	case "C":
		return C
	case "D":
		return D
	case "E":
		return E
	case "F":
		return F
	case "G":
		return G
	case "H":
		return H
	default:
		return A
	}
}
