package main

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

// Indonesia: weekend Sat/Sun, holiday on a Wednesday (2026-08-19) so tomorrow-
// holiday can be tested with a workday today.
var idCfg = CountryConfig{
	Timezone:    "Asia/Jakarta", // UTC+7
	WeekendDays: []int{6, 0},
	WorkHours:   WorkHours{Start: 9, End: 17},
	Holidays:    []string{"2026-08-19"},
}

// Saudi: weekend Fri/Sat, no holidays.
var saCfg = CountryConfig{
	Timezone:    "Asia/Riyadh", // UTC+3
	WeekendDays: []int{5, 6},
	WorkHours:   WorkHours{Start: 9, End: 17},
	Holidays:    []string{},
}

func TestDecide(t *testing.T) {
	cases := []struct {
		name     string
		now      string
		cfg      CountryConfig
		wantSafe bool
		wantCode string
	}{
		// 2026-08-12 = Wed, 2026-08-13 = Thu → clean workday
		{"weekday daytime is safe", "2026-08-12T10:00:00+07:00", idCfg, true, CodeOK},
		{"night blocks", "2026-08-12T02:00:00+07:00", idCfg, false, CodeNight},
		{"17:00 counts as night", "2026-08-12T17:00:00+07:00", idCfg, false, CodeNight},
		// 2026-08-15 = Sat
		{"today weekend blocks", "2026-08-15T10:00:00+07:00", idCfg, false, CodeTodayWeekend},
		// 2026-08-14 = Fri, tomorrow Sat
		{"tomorrow weekend blocks", "2026-08-14T10:00:00+07:00", idCfg, false, CodeTomorrowWeekend},
		// 2026-08-19 = Wed holiday
		{"today holiday blocks", "2026-08-19T10:00:00+07:00", idCfg, false, CodeTodayHoliday},
		// 2026-08-18 = Tue, tomorrow Wed holiday
		{"tomorrow holiday blocks", "2026-08-18T10:00:00+07:00", idCfg, false, CodeTomorrowHoliday},
		// Saudi: 2026-08-14 = Fri → weekend for SA
		{"SA Friday is weekend", "2026-08-14T10:00:00+03:00", saCfg, false, CodeTodayWeekend},
		// Saudi: 2026-08-16 = Sun → workday for SA, tomorrow Mon
		{"SA Sunday is a workday", "2026-08-16T10:00:00+03:00", saCfg, true, CodeOK},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			safe, code := decide(mustTime(c.now), c.cfg)
			if safe != c.wantSafe || code != c.wantCode {
				t.Errorf("decide() = (%v, %q); want (%v, %q)", safe, code, c.wantSafe, c.wantCode)
			}
		})
	}
}
