package service

import (
	"html/template"
	"os"
	"strings"
)

type StaticPath string

const (
	HOME_PATH   StaticPath = "home.html"
	ORDERS_PATH StaticPath = "orders.html"
	ORDER_PATH  StaticPath = "order.html"
)

func ParseStaticPath(
	path StaticPath,
	parseFiles func(filenames ...string) (*template.Template, error),
	config *Config,
) (*template.Template, error) {
	return parseFiles(strings.Join([]string{config.Staticpath, string(os.PathSeparator), string(path)}, ""))
}
