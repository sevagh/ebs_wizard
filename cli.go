package main

import (
	"github.com/docopt/docopt-go"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	usage := `kraken - EC2/EBS utility

Usage:
  kraken <mountpath> [--raid-level=<level>]
  kraken -h | --help
  kraken --version

Options:
  --raid-level=<level>  0 or 1
  -h --help             Show this screen.
  --version             Show version.`
	arguments, _ := docopt.Parse(usage, nil, true, "Kraken 0.1", false)

	currTime := time.Now().UTC()
	logger := log.New(os.Stderr, "kraken: ", log.Lshortfile)
	logger.Printf("RUNNING KRAKEN: %s", currTime.Format(time.RFC850))
	mountPath := arguments["<mountpath>"].(string)
	raidLevel, err := strconv.Atoi(arguments["--raid-level"].(string))
	if err != nil {
		logger.Fatalf("Couldn't parse --raid-level as int")
		os.Exit(-1)
	}

	deviceNames, err := AttachEbsVolumes(logger)
	if err != nil {
		logger.Println(err)
		os.Exit(-1)
	}
	logger.Printf("Attached: %s\n", deviceNames)
	logger.Printf("Now mounting")
	if len(deviceNames) == 1 {
		if err := MountSingleDrive(deviceNames[0], mountPath, logger); err != nil {
			os.Exit(-1)
		}
	} else {
		if err := MountRaidDrives(deviceNames, mountPath, raidLevel, logger); err != nil {
			os.Exit(-1)
		}
	}
}
