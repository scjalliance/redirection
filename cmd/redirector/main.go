package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/scjalliance/redirection"
)

var config string

func main() {
	flag.StringVar(&config, "config", "redirect.json", "configuration file")
	flag.Parse()

	element, err := load(config)
	if err != nil {
		log.Fatal(err)
	}
	handler := redirection.NewHandler(element)

	errchan := make(chan error)
	go func(errchan chan error) {
		for {
			err := http.ListenAndServe("0.0.0.0:80", handler)
			errchan <- fmt.Errorf("redirector exited: %s", err)
			time.Sleep(time.Second * 1)
		}
	}(errchan)

	for {
		fmt.Println(<-errchan)
	}
}

func load(path string) (mapper redirection.Mapper, err error) {
	contents, fileErr := ioutil.ReadFile(path)
	if fileErr != nil {
		return nil, fmt.Errorf("unable to read configuration file \"%s\": %v", path, fileErr)
	}

	// TODO: Use json.Decoder and stream the file instead of slurping?

	var element redirection.Element
	dataErr := json.Unmarshal(contents, &element)
	if dataErr != nil {
		return nil, fmt.Errorf("decoding error while parsing configuration file \"%s\": %v", path, dataErr)
	}

	return &element, nil
}
