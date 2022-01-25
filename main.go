package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"compress/zlib"

	blueprint "github.com/bunniesandbeatings/vizit/blueprint"
	"github.com/bunniesandbeatings/vizit/il"
	"go.uber.org/zap"
	"sigs.k8s.io/yaml"
)

var source string

func init() {
	flag.StringVar(&source, "file", "source.yaml", "source to analyze")
}

func mermaid(parsed il.Blueprint) string {
	lines := []string{}
	edges := []string{}

	lines = append(lines, "flowchart RL")

	for name, resource := range parsed.Resources {
		lines = append(lines, fmt.Sprintf("  subgraph res_%s[%s]", name, name))
		lines = append(lines, fmt.Sprintf("  direction RL"))
		for i, opt := range resource.Options {
			lines = append(lines, fmt.Sprintf("    opt_%d_%s[%s]", i, opt.TemplateRef.Name, opt.TemplateRef.Name))
			for _, input := range opt.Inputs {
				edges = append(edges, fmt.Sprintf("  opt_%d_%s --> res_%s", i, opt.TemplateRef.Name, input))
			}
		}
		lines = append(lines, fmt.Sprintf("  end"))
	}

	s1 := strings.Join(lines, "\n")
	s2 := strings.Join(edges, "\n")
	return s1 + "\n" + s2
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
	//_, _ = pretty.Println(parsed)
	//_, _ = pretty.Println(parsed.Entrypoints())

	mermaidString := mermaid(parsed)
	siteString := fmt.Sprintf("{\"code\":\"%s\",\"mermaid\":\"{\\n  \\\"theme\\\": \\\"dark\\\"\\n}\",\"updateEditor\":false,\"autoSync\":true,\"updateDiagram\":false}", strings.Replace(mermaidString, "\n", "\\n",-1))

	var b bytes.Buffer

	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write([]byte(siteString))
	w.Close()

	sEnc := base64.URLEncoding.EncodeToString(b.Bytes())
	fmt.Println("https://mermaid.live/edit/#pako:" + sEnc)
}
