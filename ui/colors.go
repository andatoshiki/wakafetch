package ui

type Colors struct {
	MidGray  string
	Red      string
	Yellow   string
	BoldBlue string
	Blue     string
	Green    string
	Gray     string
	Bold     string
	Reset    string
}

var Clr Colors = Colors{
	MidGray:  "\x1b[38;2;128;128;128m",
	Red:      "\x1b[31m",
	Yellow:   "\x1b[33m",
	BoldBlue: "\x1b[1;34m",
	Blue:     "\x1b[34m",
	Green:    "\x1b[32m",
	Gray:     "\x1b[90m",
	Bold:     "\x1b[1m",
	Reset:    "\x1b[0m",
}

func DisableColors() {
	Clr = Colors{
		MidGray:  "",
		Red:      "",
		Yellow:   "",
		BoldBlue: "",
		Blue:     "",
		Green:    "",
		Gray:     "",
		Bold:     "",
		Reset:    "",
	}
}
