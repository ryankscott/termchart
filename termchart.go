package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gizak/termui"
)

// TODO:
// - Different chart types
// - Check sizing

func drawTermChart(label []string, data []int) {
	if err := termui.Init(); err != nil {
		log.WithFields(log.Fields{
			"message": "Failed to start termui",
			"value":   err,
		}).Panic("Failed to start termui")
	}
	defer termui.Close()

	// Determine width of chart
	dataPoints := len(data)
	barWidth := 4
	barGap := 1
	width := (barWidth * dataPoints) + (barGap * (dataPoints - 2))

	// If we can fit it all in one screen
	if width < 200 {
		bc := termui.NewBarChart()
		bclabels := label
		bc.Data = data
		bc.Height = 30
		bc.BarWidth = barWidth
		bc.BarGap = barGap
		bc.DataLabels = bclabels
		bc.TextColor = termui.ColorWhite
		bc.BarColor = termui.ColorBlue
		bc.NumColor = termui.ColorWhite
		termui.Body.AddRows(
			termui.NewRow(
				termui.NewCol(12, 0, bc),
			))
		termui.Body.Align()
		termui.Render(termui.Body)
	} else {
		requiredCharts := int(math.Ceil(float64(width) / 200.0))
		fmt.Println("Number of charts required:", requiredCharts)
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
			bc.Height = 30
			bc.BarWidth = int(float64(barWidth) * 1.5)
			bc.BarGap = barGap
			bc.DataLabels = bclabels
			bc.TextColor = termui.ColorWhite
			bc.BarColor = termui.ColorBlue
			bc.NumColor = termui.ColorWhite
			bc.SetY(z * 30)
			barCharts[z] = bc
			termui.Body.AddRows(
				termui.NewRow(
					termui.NewCol(12, 0, bc),
				),
			)
		}
		termui.Body.Align()
		termui.Render(termui.Body)
	}

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Loop()

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

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	var data []int
	var labels []string
	var delim string
	var rows int
	x := bufio.NewScanner(os.Stdin)
	for {
		z := x.Scan()
		if z == false {
			break
		} else {

			line := x.Text()
			if rows == 0 {
				delim = tryDetectDelimeter(line)
			}
			values := strings.Split(line, delim)
			intVal, err := strconv.ParseInt(values[1], 10, 64)
			if err != nil {
				log.WithFields(log.Fields{
					"message": "Failed to convert value to int",
					"value":   values[1],
				}).Panic("Failed to convert input data to int")
			}

			labels = append(labels, values[0])
			data = append(data, int(intVal))
		}
		rows++
	}
	drawTermChart(labels, data)
}
