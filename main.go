package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
	"xmasx/xmasx"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	var parser = argparse.NewParser("nwCraft", "Provides shopping list for desired craft")
	file := parser.String("f", "file", &argparse.Options{Required: false, Help: "File of craftData to load", Default: "sampleData.csv"})
	debugLevel := parser.Selector("d", "debug-level", []string{"INFO", "DEBUG"}, &argparse.Options{Required: false, Help: "Logging debug level"})
	pretty := parser.Flag("p", "pretty", &argparse.Options{Required: false, Help: "Pretty output"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(parser.Usage(err))
		return
	}

	switch *debugLevel {
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	if *pretty {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("| %s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		log.Logger = zerolog.New(output).With().Timestamp().Logger()
	}
	xmasx.Run(*file)
}
