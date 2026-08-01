package pdf

import (
	"bytes"
	"strings"

	"github.com/dslipak/pdf"
	"github.com/examsync/pdf-parser/internal/models"
)

// ExtractText extracts plain text from PDF file bytes.
func ExtractText(fileBytes []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return "", err
	}

	var textBuilder strings.Builder
	totalPage := reader.NumPage()
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := reader.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		texts := p.Content().Text
		for _, txt := range texts {
			textBuilder.WriteString(txt.S)
		}
		textBuilder.WriteByte('\n')
	}

	return textBuilder.String(), nil
}

// ParseNotification maps extracted text to the ExamNotification model.
func ParseNotification(fileName string, text string) *models.ExamNotification {
	return &models.ExamNotification{
		FileName: fileName,
		RawText:  text,
	}
}

