package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/andatoshiki/wakafetch/wakafetch/ui"
)

type Config struct {
	rangeFlag   *string
	apiKeyFlag  *string
	fullFlag    *bool
	daysFlag    *int
	dailyFlag   *bool
	heatmapFlag *bool
	jsonFlag    *bool
	helpFlag    *bool
	updateFlag  *bool
	timeoutFlag *int
}

type flagInfo struct {
	longName    string
	shortName   string
	defaultVal  any
	description string
	flagType    string
}

func parseFlags() Config {
	config := Config{}
	registeredFlags = nil

	config.rangeFlag = config.stringFlag("range", "r", "today", "Range of data to fetch (today/yesterday/7d/30d/6m/1y/all/year) (default: today)")
	config.daysFlag = config.intFlag("days", "d", 0, "Number of days to fetch data for (overrides --range)")
	config.fullFlag = config.boolFlag("full", "f", false, "Display full statistics")
	config.dailyFlag = config.boolFlag("daily", "D", false, "Display daily breakdown")
	config.heatmapFlag = config.boolFlag("heatmap", "H", false, "Display heatmap of daily activity")
	config.apiKeyFlag = config.stringFlag("api-key", "k", "", "Your WakaTime/Wakapi API key (overrides config)")
	config.jsonFlag = config.boolFlag("json", "j", false, "Output data in JSON format")
	config.helpFlag = config.boolFlag("help", "h", false, "Display help information")
	config.updateFlag = config.boolFlag("update", "u", false, "Check for updates and show install command if newer version exists")
	config.timeoutFlag = config.intFlag("timeout", "t", 10, "Request timeout in seconds")

	// Version flag
	versionFlag := config.boolFlag("version", "v", false, "Display version information")

	 flag.Usage = showCustomHelp
	 flag.Parse()

	 if !colorsShouldBeEnabled() {
	 	ui.DisableColors()
	 }

	 if *config.helpFlag {
	 	showCustomHelp()
	 }

	       if *versionFlag {
		       fmt.Println("wakafetch version: v" + Version)
		       os.Exit(0)
	       }

	 if *config.daysFlag < 0 {
	 	ui.Errorln("Invalid value for --days: must be a positive integer")
	 }

	 return config
}

var registeredFlags []flagInfo

func (c *Config) stringFlag(long, short, def, desc string) *string {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, "string"})
	val := flag.String(long, def, "")
	flag.StringVar(val, short, def, "")
	return val
}

func (c *Config) boolFlag(long, short string, def bool, desc string) *bool {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, ""})
	val := flag.Bool(long, def, "")
	flag.BoolVar(val, short, def, "")
	return val
}

func (c *Config) intFlag(long, short string, def int, desc string) *int {
	registeredFlags = append(registeredFlags, flagInfo{long, short, def, desc, "int"})
	val := flag.Int(long, def, "")
	flag.IntVar(val, short, def, "")
	return val
}

func showCustomHelp() {
	// ASCII title and version / credits
	fmt.Println(ui.Clr.Green + "▗▖ ▗▖ ▗▄▖ ▗▖ ▗▖ ▗▄▖ ▗▄▄▄▖▗▄▄▄▖▗▄▄▄▖▗▄▄▖▗▖ ▗▖" + ui.Clr.Reset)
	fmt.Println(ui.Clr.Green + "▐▌ ▐▌▐▌ ▐▌▐▌▗▞▘▐▌ ▐▌▐▌   ▐▌     █ ▐▌   ▐▌ ▐▌" + ui.Clr.Reset)
	fmt.Println(ui.Clr.Green + "▐▌ ▐▌▐▛▀▜▌▐▛▚▖ ▐▛▀▜▌▐▛▀▀▘▐▛▀▀▘  █ ▐▌   ▐▛▀▜▌" + ui.Clr.Reset)
	fmt.Println(ui.Clr.Green + "▐▙█▟▌▐▌ ▐▌▐▌ ▐▌▐▌ ▐▌▐▌   ▐▙▄▄▖  █ ▝▚▄▄▖▐▌ ▐▌" + ui.Clr.Reset)
	if Version != "" {
		fmt.Println(ui.Clr.Bold + ui.Clr.Yellow + " Wakafetch @" + Version + ui.Clr.Reset)
	}
	fmt.Println(ui.Clr.Blue + " A colorful WakaTime/Wakapi stats fetcher for your terminal (without needing to open/refresh the web dashboard every time)" + ui.Clr.Reset)
	fmt.Println(" Original author: " + ui.Clr.Green + "sahaj-b" + ui.Clr.Reset + "  |  Current maintainer: " + ui.Clr.Green + "andatoshiki" + ui.Clr.Reset + " (" + ui.Clr.Blue + "https://www.toshiki.dev" + ui.Clr.Reset + ")")
	fmt.Println()
	fmt.Println(ui.Clr.Bold + "Usage:" + ui.Clr.Reset + " wakafetch [options]")
	fmt.Println(ui.Clr.Bold + "Options:" + ui.Clr.Reset)

	maxWidth := 0
	for _, f := range registeredFlags {
		width := len("-" + f.shortName + ", --" + f.longName + " " + f.flagType)
		if width > maxWidth {
			maxWidth = width
		}
	}

	for _, f := range registeredFlags {
		flag := fmt.Sprintf("-%s, --%s", f.shortName, f.longName)
		flagLen := len(flag)
		if f.flagType != "" {
			flagLen = len(flag + " " + f.flagType)
			flag += " " + ui.Clr.Blue + f.flagType + ui.Clr.Reset
		}
		padding := strings.Repeat(" ", maxWidth-flagLen+2)
		fmt.Println("  " + ui.Clr.Green + flag + ui.Clr.Reset + padding + f.description)
	}
	os.Exit(0)
}
