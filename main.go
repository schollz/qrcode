package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/urfave/cli"
	"rsc.io/qr"
)

var version string

func main() {
	app := cli.NewApp()
	app.Name = "qrcode"
	app.Usage = "generate a QR code from a text file"
	app.Version = version
	app.Compiled = time.Now()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "text file to convert to QR code",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "qr.png",
			Usage: "file to output",
		},
	}

	app.Action = func(c *cli.Context) error {
		err := qrcode(c.GlobalString("file"), c.GlobalString("output"))
		if err != nil {
			fmt.Print(err)
		} else {
			fmt.Printf("QR code from '%s' written to '%s'\n", c.GlobalString("file"), c.GlobalString("output"))
		}
		return err
	}

	app.Run(os.Args)
}

func qrcode(readFile, outputFile string) (err error) {
	b, err := ioutil.ReadFile(readFile)
	if err != nil {
		return
	}
	q, err := qr.Encode(string(b), qr.L)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(outputFile, q.PNG(), 0644)
	return
}
