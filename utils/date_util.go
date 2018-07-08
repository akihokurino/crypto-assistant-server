package utils

import "time"

type DateUtil interface {
	CurrentTime() time.Time
}

type dateUtil struct {
}

func NewDateUtil() DateUtil {
	return &dateUtil{}
}

func (u *dateUtil) CurrentTime() time.Time {
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return nowUTC.In(jst)
}
