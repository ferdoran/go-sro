package utils

import "time"

//BIT-MASKS
//YEAR   : 0000 0000 0000 0000 0000 0000 0011 1111
//MONTH  : 0000 0000 0000 0000 0000 0011 1100 0000
//DAY    : 0000 0000 0000 0000 0111 1100 0000 0000
//HOUR   : 0000 0000 0000 1111 1000 0000 0000 0000
//MINUTE : 0000 0011 1111 0000 0000 0000 0000 0000
//SECOND : 1111 1100 0000 0000 0000 0000 0000 0000

const (
	YearMask   = ^uint32(0) >> (32 - 6)
	MonthMask  = ^uint32(0) >> (32 - 4)
	DayMask    = ^uint32(0) >> (32 - 5)
	HourMask   = DayMask
	MinuteMask = YearMask
	SecondMask = YearMask
)

func ParseSilkroadTime(timestamp uint32) time.Time {
	year := int(timestamp&YearMask) + 2000
	month := int(timestamp >> 6 & MonthMask)
	day := int(timestamp >> 10 & DayMask)
	hour := int(timestamp >> 15 & HourMask)
	minute := int(timestamp >> 20 & MinuteMask)
	second := int(timestamp >> 26 & SecondMask)

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
}

func ToSilkroadTime(t time.Time) uint32 {
	var sroTime uint32 = 0
	sroTime += uint32(t.Year()-2000) & YearMask
	sroTime += uint32(t.Month()) & MonthMask << 6
	sroTime += uint32(t.Day()) & DayMask << 10
	sroTime += uint32(t.Hour()) & HourMask << 15
	sroTime += uint32(t.Minute()) & MinuteMask << 20
	sroTime += uint32(t.Second()) & SecondMask << 26

	return sroTime
}
