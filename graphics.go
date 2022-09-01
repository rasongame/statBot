package main

import (
	"fmt"
	"github.com/wcharczuk/go-chart"
	"os"
)

func RenderActiveUsers(elements []SomePlaceholder, fileName string, limit int, fromTimeText string) {
	var activeStat []chart.Value
	for _, v := range elements[:limit] {
		activeStat = append(activeStat, chart.Value{
			Style: chart.Style{},
			Label: fmt.Sprintf("%s %s %d", v.User.FirstName, v.User.LastName, v.Messages),
			Value: float64(v.Messages),
		})
	}
	finalChart := chart.PieChart{
		Title:  fmt.Sprintln("Активные флудильщики за ", fromTimeText),
		Values: activeStat,
		Width:  4096,
		Height: 4096,
	}

	f, _ := os.Create(fileName)
	defer f.Close()
	finalChart.Render(chart.PNG, f)

}
