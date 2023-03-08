package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
	"math/rand"
	"os"
	"statBot/utils"
	"strings"
	"time"
)

var (
	chromeContext context.Context
	cancel        context.CancelFunc
)

func init() {
	rand.Seed(time.Now().Unix())
	chromeContext, cancel = chromedp.NewContext(context.Background())
	//defer cancel()
	fmt.Println(chromeContext.Done())

}

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

	var buf []byte
	ctx, _cancel := chromedp.NewContext(chromeContext)
	defer _cancel()
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

var (
	subtitlePatterns = []string{
		"Разведка доложила, что по этим данным, %{Name}s — самый опасный флудераст на весь чат. С этим надо что-то делать.",
		"Эй-эй-эй, это что у нас... Ага, Больше всех нафлудил %{Name}s за данный срок. Настоящий флудераст.",
		"Woop-woop! That's the sound of da police. Выходи с поднятыми руками, %{Name}s.",
		"Ты пидор, %{Name}s.",
		"Zeit zur Arbeit zu gehen, %{Name}s.",
	}
)

func Tprintf(format string, params map[string]interface{}) string {
	for key, val := range params {
		format = strings.Replace(format, "%{"+key+"}s", fmt.Sprintf("%s", val), -1)
	}
	return fmt.Sprintf(format)
}

func RenderActiveUsers(elements []utils.SomePlaceholder, fileName string, limit int, fromTimeText string) {
	pie := charts.NewPie()
	flooderName := fmt.Sprintf("%s %s [%s]", elements[0].User.FirstName, elements[0].User.LastName, elements[0].User.UserName)
	textStyle := opts.TextStyle{
		FontFamily: "Noto Sans",
		FontSize:   24,
	}
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(
			opts.Initialization{
				Theme:  types.ThemeChalk,
				Width:  "2000px",
				Height: "1200px",
			}),
		charts.WithTitleOpts(
			opts.Title{
				Title:         "Флудильщики за " + fromTimeText,
				Subtitle:      Tprintf(subtitlePatterns[rand.Int()%len(subtitlePatterns)], map[string]interface{}{"Name": flooderName}),
				SubtitleStyle: &textStyle,
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
