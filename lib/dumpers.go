package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

func DumpLauncher(ranks Ranks, location *string) error {
	if err := DumpJson(ranks); err != nil {
		return fmt.Errorf("JSON Dump failed: %s", err)
	}
	if err := DumpMarkdown(ranks, location); err != nil {
		return fmt.Errorf("MD Dump failed: %s", err)
	}
	return nil
}

func DumpJson(data interface{}) error {
	Trace.Printf("Writing JSON to: %s\n", rspJSONpath)

	jsonenc, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	f, err := os.Create(rspJSONpath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString(string(jsonenc))
	if err != nil {
		return err
	}
	Trace.Printf("Wrote %d bytes to %s\n", n4, rspJSONpath)

	w.Flush()

	return nil
}

func DumpMarkdown(ranks Ranks, location *string) error {
	Trace.Printf("Writing MD to: %s\n", rspMDpath)

	head := `# soranks

[Stackoverflow](http://stackoverflow.com/) rankings by **location**.

### Area%s


Rank|Name|Rep|Top Tags|Location|Web
----|----|---|--------|--------|---
`
	var fmtLocation string

	if *location == "." {
		fmtLocation = ": WorldWide"
	} else {
		fmtLocation = fmt.Sprintf(" *pattern*: %s", *location)
	}

	userfmt := "{{.Rank}}|[{{.DisplayName}}]({{.Link}})|{{.Reputation}}|<ul>{{.TopTags}}</ul>|{{.Location}}|[![Web]({{.ProfileImage}})]({{.WebsiteURL}})\n"

	f, err := os.Create(rspMDpath)
	if err != nil {
		return err
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	n4, err := w.WriteString(fmt.Sprintf(head, fmtLocation))
	if err != nil {
		return err
	}
	w.Flush()

	tmpl, err := template.New("Ranking").Parse(userfmt)
	if err != nil {
		return err
	}

	for _, userRank := range ranks {
		_ = tmpl.Execute(f, userRank)
	}
	Trace.Printf("Wrote %d bytes to %s\n", n4, rspMDpath)
	w.Flush()

	return nil
}
