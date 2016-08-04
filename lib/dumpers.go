package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

func DumpJson(path *string, ranks *Ranks) {
	//Trace.Printf("Writing JSON to: %s\n", *path)
	jsonenc, _ := json.MarshalIndent(*ranks, "", " ")
	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(string(jsonenc))
	if err != nil {
		panic(err)
	}
	//Trace.Printf("Wrote %d bytes to %s\n", n4, *path)

	w.Flush()
}

func DumpMarkdown(path *string, ranks Ranks, location *string) {
	//Trace.Printf("Writing MD to: %s\n", *path)

	head := `# soranks

[Stackoverflow](http://stackoverflow.com/) rankings by **location**.

### Area%s


Rank|Name|Rep|Location|Web|Avatar
----|----|---|--------|---|------
`
	var fmtLocation string

	if *location == "." {
		fmtLocation = ": WorldWide"
	} else {
		fmtLocation = fmt.Sprintf(" *pattern*: %s", location)
	}

	userfmt := "{{.Rank}}|[{{.DisplayName}}]({{.Link}})|{{.Reputation}}|{{.Location}}|{{.WebsiteURL}}|![Avatar]({{.ProfileImage}})\n"

	f, err := os.Create(*path)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	_, err = w.WriteString(fmt.Sprintf(head, fmtLocation))
	if err != nil {
		panic(err)
	}
	w.Flush()

	tmpl, _ := template.New("Ranking").Parse(userfmt)
	for _, userRank := range ranks {
		_ = tmpl.Execute(f, userRank)
	}
	//Trace.Printf("Wrote %d bytes to %s\n", n4, *path)
	w.Flush()
}
