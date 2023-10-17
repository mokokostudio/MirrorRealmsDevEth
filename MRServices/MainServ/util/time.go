package util

import (
	"time"

	"github.com/oldjon/gutil/timeutil"
)

// relative time
type TimeRel float32

func (tr TimeRel) ToMiliseconds() TimeRelMS {
	return TimeRelMS(float32(tr)*1000 + 0.5)
}
func (tr TimeRel) ToFloat32() float32 {
	return float32(tr)
}

type TimeRelMS int32

func (trm TimeRelMS) ToSeconds() TimeRel {
	return TimeRel(float32(trm) / 1000)
}
func (trm TimeRelMS) ToInt() int32 {
	return int32(trm)
}

// absolute time
type TimeAbs float32

func (ta TimeAbs) ToMiliseconds() TimeAbsMS {
	return TimeAbsMS(float32(ta)*1000 + 0.5)
}
func (ta TimeAbs) ToFloat32() float32 {
	return float32(ta)
}
func (ta TimeAbs) Forward(t TimeRel) TimeAbs {
	return TimeAbs(ta.ToFloat32() + t.ToFloat32())
}
func (ta TimeAbs) Backward(t TimeRel) TimeAbs {
	return TimeAbs(ta.ToFloat32() - t.ToFloat32())
}
func (ta TimeAbs) Interval(t TimeAbs) TimeRel {
	return TimeRel(ta - t)
}

type TimeAbsMS int32

func (tam TimeAbsMS) ToSeconds() TimeAbs {
	return TimeAbs(float32(tam) / 1000)
}
func (tam TimeAbsMS) ToInt() int32 {
	return int32(tam)
}
func (tam TimeAbsMS) Forward(t TimeRelMS) TimeAbsMS {
	return TimeAbsMS(tam.ToInt() + t.ToInt())
}
func (tam TimeAbsMS) Backward(t TimeRelMS) TimeAbsMS {
	return TimeAbsMS(tam.ToInt() - t.ToInt())
}
func (tam TimeAbsMS) Interval(t TimeAbsMS) TimeRelMS {
	return TimeRelMS(tam - t)
}

func NowInMilliSeconds() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}
func NowInSeconds() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

type _TimerHandler func()

type Timer struct {
	isRepeated  bool
	duration    TimeRelMS
	expiredTime TimeAbsMS
	action      _TimerHandler
}

func (t *Timer) Create(currentTime TimeAbsMS, after TimeRelMS, action _TimerHandler, isRepeated bool) *Timer {
	t.SetExpiredTime(currentTime.Forward(after))
	t.duration = after
	t.isRepeated = isRepeated
	t.action = action
	return t
}

func (t *Timer) SetExpiredTime(expiredTime TimeAbsMS) {
	t.expiredTime = expiredTime
}

func (t *Timer) IsExpired(currentTime TimeAbsMS) bool {
	return t.expiredTime.ToInt() >= 0 && currentTime.Interval(t.expiredTime) >= 0
}

func (t *Timer) Update(currentTime TimeAbsMS) bool {
	var actionToCall _TimerHandler
	if t.IsExpired(currentTime) && t.action != nil {
		actionToCall = t.action
		isFinished := true
		if t.isRepeated {
			exceededTimeSlice := currentTime.Interval(t.expiredTime)
			t.SetExpiredTime(currentTime.Forward(t.duration - exceededTimeSlice))
			isFinished = false
		} else {
			t.action = nil
		}
		if actionToCall != nil {
			actionToCall()
		}
		return isFinished
	}
	return false
}

func LocalDateInUint32(t time.Time, region string) uint32 {
	loc, err := timeutil.RegionTimeZoneLocation(region)
	if err != nil {
		return 0
	}
	t = t.In(loc)
	return uint32(t.Year()*10000) + uint32(t.Month()*100) + uint32(t.Day())
}

func ResetTime(now time.Time, hour, minute, second uint32) time.Time {
	return time.Date(now.Year(), now.Month(), now.Day(), int(hour), int(minute), int(second), 0, now.Location())
}

func GetNowTimeMilliStr() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

type TimeUnixRange interface {
	GetStartTime() int64
	GetEndTime() int64
}

func InTimeUnixRange(r TimeUnixRange, now int64) bool {
	return now >= r.GetStartTime() && (r.GetEndTime() == 0 || now <= r.GetEndTime())
}
