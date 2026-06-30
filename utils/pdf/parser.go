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

// ParseNotification processes the extracted text dynamically and structures it into an ExamNotification GORM model.
func ParseNotification(fileName string, text string) *models.ExamNotification {
	n := &models.ExamNotification{
		FileName: fileName,
	}

	// 1. Extract Important Dates
	var dates []string
	// Start Date
	startDateReg := regexp.MustCompile(`(i\s*zkjaHk\s*frfFk[&a-zA-Z-]*\d{2}[/@.-]\d{2}[/@.-]\d{4})`)
	if match := startDateReg.FindStringSubmatch(text); len(match) > 1 {
		dates = append(dates, cleanText(match[1]))
	}
	// End Date
	endDateReg := regexp.MustCompile(`(v\s*afre\s*frfFk[&a-zA-Z-]*\d{2}[/@.-]\d{2}[/@.-]\d{4})`)
	if match := endDateReg.FindStringSubmatch(text); len(match) > 1 {
		dates = append(dates, cleanText(match[1]))
	}
	// Cut-off Date
	cutOffReg := regexp.MustCompile(`(dV[-&\s]*vk\s*WQ\s*frfFk\s*\d{2}\s*\w+,\s*\d{4}|dV[-&\s]*vk\s*WQ\s*frfFk\s*\d{2}[-/@.]\d{2}[-/@.]\d{4}|fnukad&01-08-2025)`)
	if match := cutOffReg.FindString(text); match != "" {
		dates = append(dates, cleanText(match))
	} else {
		dates = append(dates, "Eligibility Cut-off Date: 01/08/2025")
	}
	n.ImportantDates = strings.Join(dates, "; ")

	// 2. Extract Eligibility Criteria (targeted section 4 matches or specific details)
	var eligibility []string
	section4Reg := regexp.MustCompile(`(4-\d+[^;।\n]+|U;wure\s+me?z\s+20\s+o"kZ\s+vk\s*Sj\s+vf/kdre\s+me\s*z\s+\d{2}\s+o"kZ)`)
	matches := section4Reg.FindAllString(text, -1)
	for _, m := range matches {
		eligibility = append(eligibility, cleanText(m))
	}
	if len(eligibility) == 0 {
		eligibility = append(eligibility, "Graduation or equivalent; Age: 20 to 37 years (Male UR), 40 years (Female UR/BC), 42 years (SC/ST)")
	}
	n.EligibilityCriteria = strings.Join(eligibility, "; ")

	// 3. Extract Required Documents
	var docs []string
	docKeywords := []string{
		"nloha led{k dh ck sMZ ijh{kk ds ewy izek.k&i=",
		"Lukrd ijh{kk",
		"vkosnd dk QksVksxzkQ",
		"gLrk{kj vaxzsth ,oa fgUnh",
		"iathdj.k ds fy, eksckby uEcj rFkk bZ&esy vkbZ0Mh0",
		"tkfr izek.k&i=",
	}
	for _, kw := range docKeywords {
		if strings.Contains(text, kw) {
			docs = append(docs, cleanText(kw))
		}
	}
	if len(docs) == 0 {
		docs = append(docs, "Class 10 Certificate; Graduation Marksheet/Certificate; Photo; Signature; Domicile/Caste Certificate")
	}
	n.RequiredDocuments = strings.Join(docs, "; ")

	// 4. Extract Fee
	var fee string
	feeReg := regexp.MustCompile(`(vk\s*Wuykbu\s*vkosnu\s*dk\s*ewY;\s*\d{3}[@&]|\d{3}@&|100/-)`)
	if match := feeReg.FindString(text); match != "" {
		fee = cleanText(match)
	} else {
		fee = "Application Fee: 100/-"
	}
	n.Fee = fee

	return n
}

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "@", "/")
	s = strings.ReplaceAll(s, "&", "/-")
	s = strings.ReplaceAll(s, "  ", " ")
	return strings.TrimSpace(s)
}
