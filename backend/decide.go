package main

import "time"

// WorkHours is the local working window; outside it counts as Night.
type WorkHours struct {
	Start int `json:"start"` // inclusive, 0-23
	End   int `json:"end"`   // exclusive, 0-23
}

// CountryConfig is everything decide() needs for one country. Loaded from the
// resource file; decide() never reads the file itself (stays pure).
type CountryConfig struct {
	Timezone    string    `json:"timezone"`     // IANA, e.g. "Asia/Jakarta"
	WeekendDays []int     `json:"weekend_days"` // time.Weekday ints (Sun=0), e.g. [6,0] or [5,6]
	WorkHours   WorkHours `json:"work_hours"`
	Holidays    []string  `json:"holidays"` // "YYYY-MM-DD"
}

// Reason codes — stable keys. Prose (reason + funny message) lives in the
// resource file (reasons.json), keyed by these codes. Add/remove a message =
// edit the file, never this code.
const (
	CodeOK              = "ok"
	CodeNight           = "night"
	CodeTodayWeekend    = "today_weekend"
	CodeTomorrowWeekend = "tomorrow_weekend"
	CodeTodayHoliday    = "today_holiday"
	CodeTomorrowHoliday = "tomorrow_holiday"
)

const dateFmt = "2006-01-02"

// decide is PURE: given a moment and a country's config, report whether it's
// safe to deploy plus the first failing reason code (CodeOK when safe).
// now is injected (never time.Now() inside) so tests can pin any clock.
//
// Safe to Deploy = NOT night AND today/tomorrow are neither weekend nor holiday.
func decide(now time.Time, cfg CountryConfig) (bool, string) {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		loc = time.UTC
	}
	local := now.In(loc)
	tomorrow := local.AddDate(0, 0, 1)

	h := local.Hour()
	if h < cfg.WorkHours.Start || h >= cfg.WorkHours.End {
		return false, CodeNight
	}
	if isWeekend(local.Weekday(), cfg.WeekendDays) {
		return false, CodeTodayWeekend
	}
	if isWeekend(tomorrow.Weekday(), cfg.WeekendDays) {
		return false, CodeTomorrowWeekend
	}
	if isHoliday(local, cfg.Holidays) {
		return false, CodeTodayHoliday
	}
	if isHoliday(tomorrow, cfg.Holidays) {
		return false, CodeTomorrowHoliday
	}
	return true, CodeOK
}

func isWeekend(d time.Weekday, weekendDays []int) bool {
	for _, wd := range weekendDays {
		if int(d) == wd {
			return true
		}
	}
	return false
}

func isHoliday(t time.Time, holidays []string) bool {
	ds := t.Format(dateFmt)
	for _, h := range holidays {
		if h == ds {
			return true
		}
	}
	return false
}
