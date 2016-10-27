package main

import (
	"bufio"
	"flag"
	"math"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gizak/termui"
)

var (
	colorForeground = termui.ColorWhite
)

// TODO:
// - Different chart types
// - Check sizing

func drawTermChart(label []string, data []int) {
	// Determine width of chart
	dataPoints := len(data)
	barWidth := 4
	barGap := 1
	width := (barWidth * dataPoints) + (barGap * (dataPoints - 2))
	termWidth := termui.TermWidth()
	termHeight := termui.TermHeight()

	// If we can fit it all in one screen
	if width < termWidth {
		bc := termui.NewBarChart()
		bc.Data = data
		bc.Height = termHeight
		bc.Width = termWidth
		bc.BarWidth = barWidth
		bc.BarGap = barGap
		bc.DataLabels = label
		bc.TextColor = colorForeground
		bc.BarColor = termui.ColorBlue
		bc.NumColor = colorForeground
		termui.Body.AddRows(
			termui.NewRow(
				termui.NewCol(12, 0, bc),
			))
	} else {
		requiredCharts := int(math.Ceil(float64(width) / 200.0))
		var datas = make([][]int, requiredCharts)
		var labels = make([][]string, requiredCharts)
		// Split the data
		for i := 0; i < dataPoints; i++ {
			chartNo := (i % requiredCharts)
			datas[chartNo] = append(datas[chartNo], data[i])
			labels[chartNo] = append(labels[chartNo], label[i])
		}
		// Create the charts
		barCharts := make([]termui.Bufferer, requiredCharts)
		for z := 0; z < requiredCharts; z++ {
			bc := termui.NewBarChart()
			bclabels := labels[z]
			bc.Data = datas[z]
			bc.Height = termHeight / requiredCharts
			bc.BarWidth = int(float64(barWidth) * 1.5)
			bc.BarGap = barGap
			bc.DataLabels = bclabels
			bc.TextColor = colorForeground
			bc.BarColor = termui.ColorBlue
			bc.NumColor = colorForeground
			bc.SetY(z * 30)
			barCharts[z] = bc
			termui.Body.AddRows(
				termui.NewRow(
					termui.NewCol(12, 0, bc),
				),
			)
		}
	}
	termui.Body.Align()
	termui.Render(termui.Body)

}

func tryDetectDelimeter(l string) string {
	var bd string
	md := 0
	ds := []string{",", "|", "\t", "."}
	for _, d := range ds {
		c := strings.Split(l, d)
		if len(c) > md {
			md = len(c)
			bd = d
		}
	}
	return bd
}

func tryDetectColTypes(l string, d string) []string {
	cols := strings.Split(l, d)
	colTypes := make([]string, len(cols))
	for i, col := range cols {
		_, err := strconv.ParseInt(col, 10, 64)
		if err == nil {
			colTypes[i] = "int"
			continue
		}
		_, err = strconv.ParseFloat(col, 10)
		if err == nil {
			colTypes[i] = "float"
			continue
		}
		_, err = strconv.ParseBool(col)
		if err == nil {
			colTypes[i] = "bool"
			continue
		}
		colTypes[i] = "string"
	}
	return colTypes
}

func mustParseFlags() {
	var theme string
	flag.StringVar(&theme, "theme", "dark", "color theme to use; one of: light, dark")
	flag.Parse()
	if theme != "light" && theme != "dark" {
		log.WithFields(log.Fields{"value": theme}).Fatal("unsupported theme name")
	}

	if theme == "light" {
		colorForeground = termui.ColorBlack
	}
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	var (
		data   []int
		labels []string
		delim  string
		rows   int
	)
	mustParseFlags()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if rows == 0 {
			delim = tryDetectDelimeter(line)
		}
		values := strings.Split(line, delim)
		intVal, err := strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			log.WithField("value", values[1]).Fatal("Failed to convert input data to int")
		}
		labels = append(labels, values[0])
		data = append(data, int(intVal))
		rows++
	}
	if err := termui.Init(); err != nil {
		log.WithField("value", err).Fatal("Failed to start termui")
	}
	defer termui.Close()

	// Handlers
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/C-d", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/C-x", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/wnd/resize", func(e termui.Event) {
		termui.Body.Width = termui.TermWidth()
		termui.Body.Align()
		termui.Clear()
		termui.Render(termui.Body)
	})

	drawTermChart(labels, data)
	termui.Loop()

}
