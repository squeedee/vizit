package main

import (
	"flag"
	"io/ioutil"
	"log"

	blueprint "github.com/bunniesandbeatings/vizit/blueprint"
	"github.com/bunniesandbeatings/vizit/il"
	"github.com/kr/pretty"
	"go.uber.org/zap"

	"sigs.k8s.io/yaml"
)

var source string

func init() {
	flag.StringVar(&source, "file", "source.yaml", "source to analyze")
}



func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	yamlFile, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalf("Could not open '%s': %v", source, err)
	}

	bp := &blueprint.Blueprint{}

	err = yaml.Unmarshal(yamlFile, bp)
	if err != nil {
		log.Fatalf("Could not unmarshal '%s': %v", source, err)
	}

	//_, _ = pretty.Println(bp)
	//
	//fmt.Printf("=====================================================")
	parsed := il.ParseBlueprint(*bp)
	_, _ = pretty.Println(parsed)
}
