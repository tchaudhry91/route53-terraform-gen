package main

import "encoding/json"
import "text/template"
import "flag"
import "fmt"
import "io/ioutil"
import "os"

type Record struct {
	Name   string `json:"name"`
	Record string `json:"record"`
	ZoneID string `json:"zone_id"`
	TTL    string `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type Records struct {
	Records []Record `json:"records"`
}

func getTempl() (*template.Template, error) {
	baseTemplate :=
		` resource "aws_route53_record" "{{.Name}}" {
          zone_id = "{{.ZoneID}}"
          name    = "{{.Record}}"
          type    = "{{.Type}}"
          ttl     = "{{.TTL}}"
          records = ["{{.Value}}"]
        }`
	return template.New("r53").Parse(baseTemplate)
}

func main() {
	var recordsFile = flag.String("recordsFile", "", "Please specify file to read records from")
	flag.Parse()
	var records Records
	fmt.Println(*recordsFile)
	recs, err := ioutil.ReadFile(*recordsFile)
	if err != nil {
		fmt.Errorf("Unable to read records file:%s", *recordsFile)
		panic(err)
	}

	err = json.Unmarshal(recs, &records)
	if err != nil {
		fmt.Errorf("Unable to unmarshal json")
		panic(err)
	}

	tmpl, err := getTempl()
	if err != nil {
		fmt.Errorf("Error reading template")
		panic(err)
	}

	for _, record := range records.Records {
		filename := fmt.Sprintf("%s.tf", record.Record)
		f, err := os.Create(filename)
		if err != nil {
			fmt.Errorf("Unable to create file:%s", filename)
			panic(err)
		}
		defer f.Close()
		err = tmpl.Execute(f, record)
		if err != nil {
			fmt.Errorf("Unable to execute template")
			panic(err)
		}
		f.Sync()
	}
}
