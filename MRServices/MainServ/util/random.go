package util

import (
	"math/rand"
)

func IsPick(weight, totalWeight uint32) bool {
	return totalWeight <= weight || uint32(rand.Int31n(int32(totalWeight))+1) <= weight
}

func RandInRange(min, max uint32) uint32 {
	if min > max {
		min, max = max, min
	}
	return uint32(rand.Int31n(int32(max-min+1))) + min
}

func RandFromPool[T any](pool []T) (ret T, err error) {
	l := len(pool)
	if l == 0 {
		return ret, errRandomPoolEmpty
	}
	return pool[rand.Int31n(int32(l))], nil
}

func RandSlice[T any](pool []T, needNum uint32) (ret []T, err error) {
	if needNum == 0 {
		return
	}
	l := len(pool)
	if l == 0 {
		return nil, errRandomPoolEmpty
	}
	if l < int(needNum) {
		return nil, errRandomPoolTooSmall
	}
	newPool := make([]T, l)
	copy(newPool, pool)
	l--
	for i := needNum; i > 0; i-- {
		j := int(rand.Int31n(int32(l + 1)))
		newPool[l], newPool[j] = newPool[j], newPool[l]
		l--
	}
	return newPool[len(newPool)-int(needNum):], nil
}

func RandByWeight[T any](pool []T, weight func(i int) uint32) (ret T, err error) {
	if len(pool) == 0 {
		return ret, errRandomPoolEmpty
	}
	if weight == nil {
		return ret, errRandomWeightFuncNil
	}
	var tWeight uint32
	for i := range pool {
		tWeight += weight(i)
	}
	if tWeight == 0 {
		return ret, errRandomTotalWeightZero
	}
	tmp := uint32(rand.Int31n(int32(tWeight))) + 1
	for i := range pool {
		w := weight(i)
		if tmp > w {
			tmp -= w
			continue
		}
		return pool[i], nil
	}
	return ret, errRandomTotalWeightZero
}

func RandSliceByWeight[T any](pool []T, needNum uint32, weight func(i int) uint32) (ret []T, err error) {
	if len(pool) == 0 {
		return ret, errRandomPoolEmpty
	}
	if weight == nil {
		return ret, errRandomWeightFuncNil
	}
	if needNum == 0 {
		return nil, err
	}
	var tWeight uint32
	for i := range pool {
		tWeight += weight(i)
	}
	if tWeight == 0 {
		return ret, errRandomTotalWeightZero
	}

	l := len(pool)
	tmpDst := make([]T, needNum)
	outed := make([]bool, l)
	for n := needNum; n > 0; n-- {
		tmp := uint32(rand.Int31n(int32(tWeight)) + 1)
		for i := 0; i < l; i++ {
			if outed[i] {
				continue
			}
			w := weight(i)
			if tmp > w {
				tmp -= w
				continue
			}
			tmpDst[needNum-n] = pool[i]
			outed[i] = true
			tWeight -= w
			if tWeight == 0 {
				return tmpDst[:needNum-n], nil
			}
			break
		}
	}
	return tmpDst, nil
}
