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

	"github.com/bunniesandbeatings/vizit/blueprint"
	"github.com/bunniesandbeatings/vizit/il"
	"go.uber.org/zap"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

var source string

func init() {
	flag.StringVar(&source, "file", "source.yaml", "source to analyze")
}

func selectorString(selector *v1.LabelSelector) string {
	var selectors []string
	for name, value := range selector.MatchLabels {
		selectors = append(selectors, fmt.Sprintf(
			"#8226; %s=%s",
			name,
			value))
	}
	for _, expr := range selector.MatchExpressions {
		selectors = append(selectors, fmt.Sprintf(
			"#8226; %s %s %s",
			expr.Key,
			expr.Operator,
			strings.Join(expr.Values,","),
			))
	}
	return strings.Join(selectors,"\\\n")
}

func mermaid(parsed il.Blueprint) string {
	var lines []string
	var edges []string

	lines = append(lines, "flowchart RL")
	lines = append(lines, "  classDef sourceNode stroke:#66f,stroke-width:2px;")

	for name, resource := range parsed.Resources {
		resourceNodeName := fmt.Sprintf(
			"res_%s",
			name)
		lines = append(lines, fmt.Sprintf(
			"  subgraph %s[%s]",
			resourceNodeName,
			name))
		lines = append(lines, fmt.Sprintf(
			"    direction RL"))
		for i, opt := range resource.Options {
			optNodeName := fmt.Sprintf("opt_%d_%s", i, opt.TemplateRef.Name)
			lines = append(lines, fmt.Sprintf(
				"    %s[\\\"%s\\\\n%s\\\"]",
				optNodeName,
				opt.TemplateRef.Name,
				selectorString(opt.Selector),
			))
			for _, input := range opt.Inputs {
				edges = append(edges, fmt.Sprintf(
					"  %s --> res_%s",
					optNodeName,
					input))
			}
			if len(opt.Inputs) < 1 {
				lines = append(lines, fmt.Sprintf(
					"  class %s sourceNode;",
					optNodeName))
			}

		}
		lines = append(lines, fmt.Sprintf(
			"  end"))
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
	siteString := fmt.Sprintf("{\"code\":\"%s\",\"mermaid\":\"{\\n  \\\"theme\\\": \\\"dark\\\"\\n}\",\"updateEditor\":false,\"autoSync\":false,\"updateDiagram\":false}", strings.Replace(mermaidString, "\n", "\\n", -1))

	var b bytes.Buffer

	w, _ := zlib.NewWriterLevel(&b, zlib.BestCompression)
	w.Write([]byte(siteString))
	w.Close()

	sEnc := base64.URLEncoding.EncodeToString(b.Bytes())
	fmt.Println("https://mermaid.live/edit/#pako:" + sEnc)
}
