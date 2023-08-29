package service

import (
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"time"
)

func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
func shuffleUsers(userIds []int, percentRaw int) ([]int, error) {
	if len(userIds) == 0 {
		return nil, errors.New("null users")
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(userIds), func(i, j int) { userIds[i], userIds[j] = userIds[j], userIds[i] })
	var percent float64 = float64(percentRaw) / 100
	randomCount := math.Ceil(float64(len(userIds)) * float64(percent))
	shuffledUserIds := userIds[0:int(randomCount)]
	return shuffledUserIds, nil
}
