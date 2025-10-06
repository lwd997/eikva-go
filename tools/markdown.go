package tools

import (
	"strings"

	"github.com/xuri/excelize/v2"
)

func XlsxToMD(xlsx *excelize.File) string {
	var md strings.Builder

	for _, sheetName := range xlsx.GetSheetList() {
		rows, err := xlsx.GetRows(sheetName)
		if err != nil {
			continue
		}

		for i, row := range rows {
			if i == 0 {
				md.WriteString("| " + strings.Join(row, " | ") + " |\n")
				md.WriteString("|" + strings.Repeat(" --- |", len(row)) + "\n")
			} else {
				md.WriteString("| " + strings.Join(row, " | ") + " |\n")
			}
		}

		md.WriteString("\n")
	}

	return md.String()
}

/*func DocxToMD(docx *document.Document) string {
	var md strings.Builder

	for _, p := range docx.Paragraphs() {
		pRun := p.AddRun()
		text := strings.TrimSpace(pRun.Text())

		if text == "" {
			md.WriteString("\n")
			continue
		}
		style := p.Style()
		if style != "" {
			if strings.Contains(p.Style(), "heading") {
				hLevel := 2
				for i := 1; i <= 6; i++ {
					if strings.Contains(style, fmt.Sprintf("%d", i)) {
						hLevel = i
						break
					}
				}
				md.WriteString(fmt.Sprintf("%s %s\n\n", strings.Repeat("#", hLevel), text))
				continue
			}
		}

		md.WriteString(text + "\n\n")
	}

	return md.String()
}*/
