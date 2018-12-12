package apidocmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/deckarep/golang-set"
)

const brTag = "<br>"

type apiData struct {
	Deprecated content `json:"deprecated"`

	Version string `json:"version"`

	Group      string `json:"group"`
	GroupTitle string `json:"groupTitle"`

	Title       string `json:"title"`
	Description string `json:"description"`
	Name        string `json:"name"`

	Type string `json:"type"`
	URL  string `json:"url"`

	Header    request `json:"header"`
	Parameter request `json:"parameter"`

	Success response `json:"success"`
	Error   response `json:"error"`

	Filename string `json:"filename"`
}

type apiDataList []apiData

func (l apiDataList) Grouping() groupDataList {
	gl := make(groupDataList, 0, len(l))
	gm := make(map[string]groupData)

	s := mapset.NewSet()
	for _, v := range l {
		if !s.Contains(v.Group) {
			g := groupData{
				ID:    v.Group,
				Title: v.GroupTitle,
				List:  make([]apiData, 0),
			}
			gl = append(gl, g)
			gm[v.Group] = g
			s.Add(v.Group)
		}

		g := gm[v.Group]
		g.List = append(g.List, v)
		gm[v.Group] = g
	}

	for i, v := range gl {
		sort.Slice(gm[v.ID].List, func(i, j int) bool {
			return gm[v.ID].List[i].Title < gm[v.ID].List[j].Title
		})
		gl[i].List = gm[v.ID].List
	}

	return gl
}

type groupData struct {
	ID    string
	Title string
	List  apiDataList
}

type groupDataList []groupData

type example struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type field struct {
	Group         string   `json:"group"`
	Type          string   `json:"type"`
	AllowedValues []string `json:"allowedValues"`
	Optional      bool     `json:"optional"`
	Field         string   `json:"field"`
	Description   string   `json:"description"`
	DefaultValue  string   `json:"defaultValue"`
}

type response struct {
	Fields   map[string][]field `json:"fields"`
	Examples []example          `json:"examples"`
}

type request struct {
	Fields   map[string][]field `json:"fields"`
	Examples []example          `json:"examples"`
}

type content struct {
	Content string `json:"content"`
}

// addRequestHeader write request header
func (d apiData) addRequestHeader(buf *buffer) error {
	for k, v := range d.Header.Fields {
		if _, err := buf.Writeln("### %s", k); err != nil {
			return err
		}

		tbl := `
| Field | Optional | Type | Description |
|---|---|---|---|
`
		for _, v2 := range v {
			// optional
			opt := "N"
			if v2.Optional {
				opt = "Y"
			}

			// description
			desc := strings.NewReplacer("<p>", "", "</p>", "").Replace(v2.Description)
			if v2.DefaultValue != "" {
				desc += brTag
				desc += "Default value: " + fmt.Sprintf("`%s`", v2.DefaultValue)
			}
			if len(v2.AllowedValues) > 0 {
				desc += brTag
				allows := v2.AllowedValues
				for i, j := range allows {
					allows[i] = fmt.Sprintf("`%s`", j)
				}
				desc += "Allowed values: " + strings.Join(allows, ",")
			}

			tbl += fmt.Sprintf("| %s | %s | %s | %s |\n", v2.Field, opt, v2.Type, desc)
		}
		if _, err := buf.Writeln(tbl); err != nil {
			return err
		}
	}

	return nil
}

// addRequestParameter write request parameter
func (d apiData) addRequestParameter(buf *buffer) error {
	for k, v := range d.Parameter.Fields {
		if _, err := buf.Writeln("### %s", k); err != nil {
			return err
		}

		tbl := `
| Field | Optional | Type | Description |
|---|---|---|---|
`
		sort.Slice(v, func(i, j int) bool {
			return v[i].Field < v[j].Field
		})
		for _, d := range v {
			keyPath := strings.Split(d.Field, ".")
			field := strings.Repeat("&nbsp;&nbsp;&nbsp;", len(keyPath)-1) + keyPath[len(keyPath)-1]

			// optional
			opt := "N"
			if d.Optional {
				opt = "Y"
			}

			// description
			desc := strings.NewReplacer("<p>", "", "</p>", "").Replace(d.Description)
			if d.DefaultValue != "" {
				desc += brTag
				desc += "Default value: " + fmt.Sprintf("`%s`", d.DefaultValue)
			}
			if len(d.AllowedValues) > 0 {
				desc += brTag
				allows := d.AllowedValues
				for i, j := range allows {
					allows[i] = fmt.Sprintf("`%s`", j)
				}
				desc += "Allowed values: " + strings.Join(allows, ",")
			}

			tbl += fmt.Sprintf("| %s | %s | %s | %s |\n", field, opt, d.Type, desc)
		}
		if _, err := buf.Writeln(tbl); err != nil {
			return err
		}
	}

	for _, v := range d.Parameter.Examples {
		if _, err := buf.Writeln("<details><summary><b>%s</b></summary><div>", v.Title); err != nil {
			return err
		}

		content := fmt.Sprintf("```%s\n", v.Type)
		content += fmt.Sprintf("%s\n", v.Content)
		content += fmt.Sprintf("```\n")
		if _, err := buf.Writeln(content); err != nil {
			return err
		}
		if _, err := buf.Writeln("</div></details>"); err != nil {
			return err
		}
	}

	return nil
}

// addSuccessResponse write success response
func (d apiData) addSuccessResponse(buf *buffer) error {
	for k, v := range d.Success.Fields {
		if _, err := buf.Writeln("### %s", k); err != nil {
			return err
		}

		tbl := `
| Field | Type | Description |
|---|---|---|
`
		sort.Slice(v, func(i, j int) bool {
			return v[i].Field < v[j].Field
		})
		for _, d := range v {
			keyPath := strings.Split(d.Field, ".")
			field := strings.Repeat("&nbsp;&nbsp;&nbsp;", len(keyPath)-1) + keyPath[len(keyPath)-1]
			desc := strings.NewReplacer("<p>", "", "</p>", "").Replace(d.Description)
			tbl += fmt.Sprintf("| %s | %s | %s |\n", field, d.Type, desc)
		}
		if _, err := buf.Writeln(tbl); err != nil {
			return err
		}
	}

	for _, v := range d.Success.Examples {
		if _, err := buf.Writeln("<details><summary><b>%s</b></summary><div>", v.Title); err != nil {
			return err
		}

		content := fmt.Sprintf("```%s\n", v.Type)
		content += fmt.Sprintf("%s\n", v.Content)
		content += fmt.Sprintf("```\n")
		if _, err := buf.Writeln(content); err != nil {
			return err
		}
		if _, err := buf.Writeln("</div></details>"); err != nil {
			return err
		}
	}

	return nil
}

// addErrorResponse write error response
func (d apiData) addErrorResponse(buf *buffer) error {
	for k, v := range d.Error.Fields {
		if _, err := buf.Writeln("### %s", k); err != nil {
			return err
		}

		tbl := `
| Name | Type | Description |
|---|---|---|
`
		sort.Slice(v, func(i, j int) bool {
			return v[i].Field < v[j].Field
		})
		for _, d := range v {
			keyPath := strings.Split(d.Field, ".")
			field := strings.Repeat("&nbsp;&nbsp;&nbsp;", len(keyPath)-1) + keyPath[len(keyPath)-1]
			desc := strings.NewReplacer("<p>", "", "</p>", "").Replace(d.Description)
			tbl += fmt.Sprintf("| %s | %s | %s |\n", field, d.Type, desc)
		}
		if _, err := buf.Writeln(tbl); err != nil {
			return err
		}
	}

	for _, v := range d.Error.Examples {
		if _, err := buf.Writeln("<details><summary><b>%s</b></summary><div>", v.Title); err != nil {
			return err
		}
		//buf.Writeln("#### %s", v.Title)
		content := fmt.Sprintf("```%s\n", v.Type)
		content += fmt.Sprintf("%s\n", v.Content)
		content += fmt.Sprintf("```\n")
		if _, err := buf.Writeln(content); err != nil {
			return err
		}
		if _, err := buf.Writeln("</div></details>"); err != nil {
			return err
		}
	}

	return nil
}

// addContents write contents
func (l apiDataList) addContents(buf *buffer) error {
	for _, v := range l {
		if _, err := buf.Writeln("<a name='%s'></a>", v.Name); err != nil {
			return err
		}
		if _, err := buf.Writeln("## %s", v.Title); err != nil {
			return err
		}
		if v.Deprecated.Content != "" {
			if _, err := buf.Writeln("<font color='red'>**DEPRECATED**</font> %s", v.Deprecated.Content); err != nil {
				return err
			}
		} else {
			if _, err := buf.Writeln("_v%s_", v.Version); err != nil {
				return err
			}
		}
		if _, err := buf.Writeln("[Back to top](#top)"); err != nil {
			return err
		}
		if _, err := buf.Writeln(""); err != nil {
			return err
		}

		if _, err := buf.Writeln(v.Description); err != nil {
			return err
		}
		if _, err := buf.Writeln(""); err != nil {
			return err
		}
		if _, err := buf.Writeln("\t%s %s", strings.ToUpper(v.Type), v.URL); err != nil {
			return err
		}
		if _, err := buf.Writeln(""); err != nil {
			return err
		}

		// header
		v.addRequestHeader(buf)

		// parameter
		v.addRequestParameter(buf)

		// success response
		v.addSuccessResponse(buf)

		// error response
		v.addErrorResponse(buf)
	}

	return nil
}

// addNavigation write navigation
func (gl groupDataList) addNavigation(buf *buffer) error {
	for _, v := range gl {
		if _, err := buf.Writeln("- [%s](#%s)", v.Title, v.ID); err != nil {
			return err
		}
		for _, v2 := range v.List {
			if _, err := buf.Writeln("\t- [%s](#%s)", v2.Title, v2.Name); err != nil {
				return err
			}
		}
	}

	return nil
}

// addContents write contents
func (gl groupDataList) addContents(buf *buffer) error {
	for _, v := range gl {
		if _, err := buf.Writeln("<a name='%s'></a>", v.ID); err != nil {
			return err
		}
		if _, err := buf.Writeln("# %s", v.Title); err != nil {
			return err
		}
		v.List.addContents(buf)
	}

	return nil
}
