package pdf

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/dslipak/pdf"
	"github.com/examsync/pdf-parser/internal/models"
)

// ExtractText extracts the plain text from the PDF file bytes.
func ExtractText(fileBytes []byte) (string, error) {
	reader, err := pdf.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return "", err
	}

	b, err := reader.GetPlainText()
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(b); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ParseNotification processes the extracted text dynamically and structures it into an ExamNotification GORM model.
func ParseNotification(fileName string, text string) *models.ExamNotification {
	n := &models.ExamNotification{
		FileName: fileName,
	}

	lines := strings.Split(text, "\n")

	var dateLines []string
	var eligibilityLines []string
	var documentLines []string
	var feeLines []string

	// Compile useful regexes
	dateRegex := regexp.MustCompile(`\d{2}[/@.-]\d{2}[/@.-]\d{4}`)

	// Category Keywords (Low-case for easier matching)
	dateKeywords := []string{"frfFk", "date", "fnukad", "cutoff", "cut-off", "last", "start", "time"}
	eligibilityKeywords := []string{"eligibility", "qualification", "criteria", "Lukrd", "graduate", "Lukr", "mez", "lhek", "height", "chest", "weight", "ÅWapkbZ", "lhuk", "ot+u", "physical", "vgZrk", "ik=rk", "age"}
	documentKeywords := []string{"document", "certificate", "marksheet", "photo", "signature", "izek.k&i=", "passport", "aadhar", "pan card", "email", "mobile", "iathdj.k", "gLrk{kj", "QksVks"}
	feeKeywords := []string{"fee", "शुल्क", "charges", "'kqYd", "ewY;", "pktZ", "/-", "@&", "amount", "payment", "transaction"}

	seenDates := make(map[string]bool)
	seenEligibilities := make(map[string]bool)
	seenDocs := make(map[string]bool)
	seenFees := make(map[string]bool)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lowerLine := strings.ToLower(trimmed)

		// 1. Scan for Dates
		isDateLine := dateRegex.MatchString(trimmed)
		if !isDateLine {
			for _, kw := range dateKeywords {
				if strings.Contains(lowerLine, kw) {
					isDateLine = true
					break
				}
			}
		}
		if isDateLine && len(trimmed) < 150 {
			if !seenDates[trimmed] {
				dateLines = append(dateLines, trimmed)
				seenDates[trimmed] = true
			}
		}

		// 2. Scan for Eligibility
		isEligibility := false
		for _, kw := range eligibilityKeywords {
			if strings.Contains(lowerLine, kw) {
				isEligibility = true
				break
			}
		}
		if isEligibility && len(trimmed) < 200 {
			if !seenEligibilities[trimmed] {
				eligibilityLines = append(eligibilityLines, trimmed)
				seenEligibilities[trimmed] = true
			}
		}

		// 3. Scan for Documents
		isDoc := false
		for _, kw := range documentKeywords {
			if strings.Contains(lowerLine, kw) {
				isDoc = true
				break
			}
		}
		if isDoc && len(trimmed) < 200 {
			if !seenDocs[trimmed] {
				documentLines = append(documentLines, trimmed)
				seenDocs[trimmed] = true
			}
		}

		// 4. Scan for Fee
		isFee := false
		for _, kw := range feeKeywords {
			if strings.Contains(lowerLine, kw) {
				isFee = true
				break
			}
		}
		if isFee && len(trimmed) < 200 {
			if !seenFees[trimmed] {
				feeLines = append(feeLines, trimmed)
				seenFees[trimmed] = true
			}
		}
	}

	n.ImportantDates = strings.Join(dateLines, "; ")
	n.EligibilityCriteria = strings.Join(eligibilityLines, "; ")
	n.RequiredDocuments = strings.Join(documentLines, "; ")
	n.Fee = strings.Join(feeLines, "; ")

	return n
}
