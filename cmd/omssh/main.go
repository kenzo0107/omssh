package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli"

	"github.com/kenzo0107/omssh"

	latest "github.com/tcnksm/go-latest"
)

const version = "0.0.3"

var (
	buildDate string
)

func main() {
	var (
		showVersion bool
	)

	flag.BoolVar(&showVersion, "v", false, "show version")
	flag.BoolVar(&showVersion, "version", false, "show version")

	if showVersion {
		fmt.Println("version:", version)
		fmt.Println("build:", buildDate)
		checkLatest(version)
		return
	}

	app := cli.NewApp()

	app.Name = "Oreno mssh"
	app.Version = version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "region, r",
			Value: "ap-northeast-1",
			Usage: "aws region",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "22",
			Usage: "ssh port",
		},
		cli.BoolFlag{
			Name:  "user, u",
			Usage: "select ssh user",
		},
	}

	app.Action = omssh.Pre

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func checkLatest(version string) {
	version = fixVersionStr(version)
	githubTag := &latest.GithubTag{
		Owner:             "kenzo0107",
		Repository:        "omssh",
		FixVersionStrFunc: fixVersionStr,
	}
	res, err := latest.Check(githubTag, version)
	if err != nil {
		fmt.Println(err)
		return
	}
	if res.Outdated {
		fmt.Printf("%s is not latest, you should upgrade to %s\n", version, res.Current)
	}
}

func fixVersionStr(v string) string {
	v = strings.TrimPrefix(v, "v")
	vs := strings.Split(v, "-")
	return vs[0]
}
