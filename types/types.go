package types

type StatItem struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	// Seconds      float64 `json:"seconds"`
	// Percent      float64 `json:"percent"`
	// Digital      string  `json:"digital"`
	// Hours        int     `json:"hours"`
	// Minutes      int     `json:"minutes"`
	// Text         string  `json:"text"`
}

type DayData struct {
	Entities         []StatItem `json:"entities"`
	Branches         []StatItem `json:"branches"`
	Categories       []StatItem `json:"categories"`
	Dependencies     []StatItem `json:"dependencies"`
	Editors          []StatItem `json:"editors"`
	Languages        []StatItem `json:"languages"`
	Machines         []StatItem `json:"machines"`
	OperatingSystems []StatItem `json:"operating_systems"`
	Projects         []StatItem `json:"projects"`
	GrandTotal       struct {
		Digital      string  `json:"digital"`
		Hours        int     `json:"hours"`
		Minutes      int     `json:"minutes"`
		Text         string  `json:"text"`
		TotalSeconds float64 `json:"total_seconds"`
	} `json:"grand_total"`
	Range struct {
		Date     string `json:"date"`
		End      string `json:"end"`
		Start    string `json:"start"`
		Text     string `json:"text"`
		Timezone string `json:"timezone"`
	} `json:"range"`
}

// /summaries
type SummaryResponse struct {
	Data            []DayData `json:"data"`
	CumulativeTotal struct {
		Digital string  `json:"digital"`
		Seconds float64 `json:"seconds"`
		Text    string  `json:"text"`
	} `json:"cumulative_total"`
	DailyAverage struct {
		DaysIncludingHolidays int     `json:"days_including_holidays"`
		DaysMinusHolidays     int     `json:"days_minus_holidays"`
		Holidays              int     `json:"holidays"`
		Seconds               float64 `json:"seconds"`
		Text                  string  `json:"text"`
	} `json:"daily_average"`
	End   string `json:"end"`
	Start string `json:"start"`
}

// /stats
type StatsResponse struct {
	Data struct {
		Branches                  []StatItem `json:"branches"`
		Categories                []StatItem `json:"categories"`
		Editors                   []StatItem `json:"editors"`
		Languages                 []StatItem `json:"languages"`
		Machines                  []StatItem `json:"machines"`
		OperatingSystems          []StatItem `json:"operating_systems"`
		Projects                  []StatItem `json:"projects"`
		Range                     string     `json:"range"`
		Status                    string     `json:"status"`
		TotalSeconds              float64    `json:"total_seconds"`
		UserID                    string     `json:"user_id"`
		Username                  string     `json:"username"`
		DailyAverage              float64    `json:"daily_average"`
		DaysIncludingHolidays     int        `json:"days_including_holidays"`
		Start                     string     `json:"start"`
		End                       string     `json:"end"`
		HumanReadableDailyAverage string     `json:"human_readable_daily_average"`
		HumanReadableRange        string     `json:"human_readable_range"`
		HumanReadableTotal        string     `json:"human_readable_total"`
		IsCodingActivityVisible   bool       `json:"is_coding_activity_visible"`
		IsOtherUsageVisible       bool       `json:"is_other_usage_visible"`
	} `json:"data"`
}
