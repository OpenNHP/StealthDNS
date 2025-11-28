package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/OpenNHP/StealthDNS/dns"
	"github.com/urfave/cli/v2"
)

// Run the program as an administrator to automatically set up the DNS proxy.
// On Linux/MacOS, the program needs to be started with the root account or
// sudo command to ensure it can properly listen on port 53.
func main() {
	if len(os.Args) == 1 {
		// On Windows, launch the application by double-clicking its executable.
		runApp()
		return
	}

	app := cli.NewApp()
	app.Name = "StealthDNS"
	app.Usage = "local DNS proxy service"
	app.Version = "0.1"

	runCmd := &cli.Command{
		Name:  "run",
		Usage: "create and run local DNS proxy service",
		Action: func(c *cli.Context) error {
			return runApp()
		},
	}

	app.Commands = []*cli.Command{
		runCmd,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(os.Stderr, err)
	}
}

func runApp() error {
	fmt.Println("Stealth DNS starting")
	exeFilePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDirPath := filepath.Dir(exeFilePath)

	p := &dns.ProxyService{}
	err = p.Start(exeDirPath, 4)
	if err != nil {
		return err
	}
	//react to terminate signals
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM, os.Interrupt, syscall.SIGABRT)

	// block until terminated
	<-termCh
	p.Stop()
	return nil
}
