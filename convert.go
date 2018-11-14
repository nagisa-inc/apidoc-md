package apidocmd

func Convert(inputDir, outputPath string) error {
	// read api project
	proj, err := readProject(inputDir)
	if err != nil {
		return err
	}

	// read api data
	list, err := readData(inputDir)
	if err != nil {
		return err
	}

	// buffer
	buf := NewBuffer(nil)

	// generate content
	if err := generateContent(buf, *proj, *list); err != nil {
		return err
	}

	// write buffer
	if err := buf.writeToFile(outputPath); err != nil {
		return err
	}

	return nil
}

// generateContent generate markdown contents
func generateContent(buf *buffer, p apiProject, gl groupDataList) error {
	// doc title
	p.addDocTitle(buf)

	// navigation
	p.addHeaderNavigation(buf)
	gl.addNavigation(buf)
	p.addFooterNavigation(buf)

	// header
	p.addHeaderContent(buf)

	// contents
	gl.addContents(buf)

	// footer
	p.addFooterContent(buf)

	// generator
	p.addGenerator(buf)

	return nil
}
