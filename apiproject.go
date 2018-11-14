package apidocmd

type apiProject struct {
	Name           string       `json:"name"`
	Version        string       `json:"version"`
	Description    string       `json:"description"`
	URL            string       `json:"url"`
	Title          string       `json:"title"`
	Order          []string     `json:"order"`
	Header         titleContent `json:"header"`
	Footer         titleContent `json:"footer"`
	SampleURL      interface{}  `json:"sampleurl"`
	DefaultVersion string       `json:"defaultVersion"`
	Apidoc         string       `json:"apidoc"`
	Generator      generator    `json:"generator"`
}

type titleContent struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type generator struct {
	Name    string `json:"name"`
	Time    string `json:"time"`
	URL     string `json:"url"`
	Version string `json:"version"`
}

// addDocTitle write document title
func (p apiProject) addDocTitle(buf *buffer) error {
	if _, err := buf.Writeln("<a name='top'></a>"); err != nil {
		return err
	}
	if _, err := buf.Writeln("# %s v%s", p.Name, p.Version); err != nil {
		return err
	}
	if _, err := buf.Writeln("%s", p.Description); err != nil {
		return err
	}
	return nil
}

// addHeaderNavigation write header title
func (p apiProject) addHeaderNavigation(buf *buffer) error {
	if _, err := buf.Writeln("- [%s](#%s)", p.Header.Title, p.Header.Title); err != nil {
		return err
	}
	return nil
}

// addFooterNavigation write footer title
func (p apiProject) addFooterNavigation(buf *buffer) error {
	if _, err := buf.Writeln("- [%s](#%s)", p.Footer.Title, p.Footer.Title); err != nil {
		return err
	}
	return nil
}

// addHeaderContent write header content
func (p apiProject) addHeaderContent(buf *buffer) error {
	if p.Header.Content != "" {
		if _, err := buf.Writeln("<a name='%s'></a>", p.Header.Title); err != nil {
			return err
		}
		if _, err := buf.Writeln(p.Header.Content); err != nil {
			return err
		}
	}
	return nil
}

// addFooterContent write footer content
func (p apiProject) addFooterContent(buf *buffer) error {
	if p.Footer.Content != "" {
		if _, err := buf.Writeln("<a name='%s'></a>", p.Footer.Title); err != nil {
			return err
		}
		if _, err := buf.Writeln(p.Footer.Content); err != nil {
			return err
		}
	}
	return nil
}

// addGenerator write generator
func (p apiProject) addGenerator(buf *buffer) error {
	if _, err := buf.Writeln("---"); err != nil {
		return err
	}
	if _, err := buf.Writeln("Generated with [%s](%s) %s - %s", p.Generator.Name, p.Generator.URL, p.Generator.Version, p.Generator.Time); err != nil {
		return err
	}
	return nil
}
