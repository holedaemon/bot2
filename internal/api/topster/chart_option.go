package topster

type chartOptions struct {
	User            string  `json:"user"`
	Period          string  `json:"period"`
	Title           string  `json:"title"`
	BackgroundColor string  `json:"background_color"`
	TextColor       string  `json:"text_color"`
	Gap             float64 `json:"gap"`
	ShowNumbers     bool    `json:"show_numbers"`
	ShowTitles      bool    `json:"show_titles"`
}

// ChartOption configures a chart.
type ChartOption func(*chartOptions)

// Period sets a chart's period.
func Period(period string) ChartOption {
	return func(co *chartOptions) {
		co.Period = period
	}
}

// Title sets a chart's title.
func Title(title string) ChartOption {
	return func(co *chartOptions) {
		co.Title = title
	}
}

// BackgroundColor sets a chart's background color.
func BackgroundColor(bgc string) ChartOption {
	return func(co *chartOptions) {
		co.BackgroundColor = bgc
	}
}

// TextColor sets a chart's text color.
func TextColor(tc string) ChartOption {
	return func(co *chartOptions) {
		co.TextColor = tc
	}
}

// Gap sets a chart's gap.
func Gap(gap float64) ChartOption {
	return func(co *chartOptions) {
		co.Gap = gap
	}
}

// ShowNumbers enables showing numbers on a chart.
func ShowNumbers() ChartOption {
	return func(co *chartOptions) {
		co.ShowNumbers = true
	}
}

// ShowTitles enables showing titles on a chart.
func ShowTitles() ChartOption {
	return func(co *chartOptions) {
		co.ShowTitles = true
	}
}
