package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"os"
	"statBot/utils"
)
func generatePieItems(elements []utils.SomePlaceholder, limit int) []opts.PieData {

	items := make([]opts.PieData, 0)
	for _, user := range elements[:limit] {
		items = append(items, opts.PieData{
			Name:  fmt.Sprintf("%s %s [%s]", user.User.FirstName, user.User.LastName, user.User.UserName),
			Value: user.Messages})
	}

	return items
}
func makeScreenshot(fileName string) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var buf []byte
	cwd, err := os.Getwd()
	utils.PanicErr(err)
	task := chromedp.Tasks{
		// use file:// for a local file
		chromedp.Navigate(fmt.Sprintf("file://%s/%s.html", cwd, fileName)),
		// set resolution for screenshot
		chromedp.EmulateViewport(2000, 1200),
		// take screenshot
		chromedp.FullScreenshot(&buf, 125),
	}

	// run tasks
	if err := chromedp.Run(ctx, task); err != nil {
		utils.PanicErr(err)
	}
	// write png
	finalFile, err := os.Create(fileName)
	_, err = finalFile.Write(buf)
}
func RenderActiveUsers(elements []utils.SomePlaceholder, fileName string, limit int, fromTimeText string) {
	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(
			opts.Initialization{
				Theme:  types.ThemeChalk,
				Width:  "2000px",
				Height: "1200px",
			}),
		charts.WithTitleOpts(
			opts.Title{
				Title:    "Флудильщики за " + fromTimeText,
				Subtitle: "",
			},
		),
	)
	pie.SetSeriesOptions()
	pie.AddSeries("", generatePieItems(elements, limit))
	//pie
	pie.SetSeriesOptions(
		charts.WithPieChartOpts(
			opts.PieChart{
				Center: []string{"50%", "50%"},
				Radius: nil,
			}),
		charts.WithLabelOpts(
			opts.Label{
				FontFamily:    "Noto Sans",
				FontSize:      22,
				VerticalAlign: "auto",
				FontWeight:    "bold",
				Show:          true,
				Formatter:     "{b}: {c}",
			},
		),
	)
	f, err := os.Create(fmt.Sprintf("%s.html", fileName))

	utils.PanicErr(err)
	err = pie.Render(f)
	utils.PanicErr(err)

	// screenshot buffer
	makeScreenshot(fileName)
	utils.PanicErr(err)

}
