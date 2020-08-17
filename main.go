package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func readConfig(fileName string) (*viper.Viper, error) {
	conf := viper.New()
	conf.SetConfigFile(fileName)
	conf.SetConfigType("env")

	err := conf.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func combine(s, t *viper.Viper, forceMerge, quite bool) *viper.Viper {

	conf := t.AllSettings()
	for k, v := range conf {
		isSet := s.IsSet(k)
		if !isSet {
			vv := v

			// interactive
			if !quite {
				reader := bufio.NewReader(os.Stdin)
				fmt.Printf("%s (default: %s): ", k, v)
				readLine, _ := reader.ReadString('\n')
				text := strings.TrimSuffix(readLine, "\n")

				if text != "" {
					vv = text
				}
			}

			s.Set(k, vv)
		}
	}

	return s
}

// try using it to prevent further errors.
// https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {

	var mainEnvFilename, distEnvFilename string

	flagSet := flag.NewFlagSet("env", flag.ContinueOnError)
	quiteFlag := flagSet.Bool("quite", false, "--quite true")

	args := os.Args
	if len(args) < 3 {
		log.Fatal("not enough arguments")
	}

	mainEnvFilename = args[1]
	distEnvFilename = args[2]
	err := flagSet.Parse(os.Args[3:])
	if err != nil {
		log.Fatal(err)
	}

	// create source file
	if !fileExists(mainEnvFilename) {
		_, err = os.Create(mainEnvFilename)
		if err != nil {
			log.Fatal(err)
		}
	}

	// main env
	mainEnv, err := readConfig(mainEnvFilename)
	if err != nil {
		log.Fatal(err)
	}

	// dist
	distViper, err := readConfig(distEnvFilename)
	if err != nil {
		log.Fatal(err)
	}

	mainEnv = combine(mainEnv, distViper, false, *quiteFlag)

	err = mainEnv.WriteConfig()
	if err != nil {
		log.Fatal(err)
	}
}
