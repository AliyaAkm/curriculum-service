package certificate

import (
	"bytes"
	"curriculum-service/internal/domain/certificate"
	"fmt"
	"math"
	"strconv"
	"strings"
)

const (
	pageW = 842.0
	pageH = 595.0
)

func renderCertificatePDF(data certificate.IssueData) []byte {
	c := newPDFCanvas()

	c.rectFill(0, 0, pageW, pageH, "0.98 0.98 0.97")
	c.rectFill(18, 18, pageW-36, pageH-36, "1 1 1")
	c.rectStroke(18, 18, pageW-36, pageH-36, "0.48 0.50 0.52", 1.2)
	c.rectStroke(26, 26, pageW-52, pageH-52, "0.72 0.74 0.76", 0.7)

	c.drawCornerMarks()
	c.drawWavePattern()
	c.drawRibbon()
	c.drawSeal()

	c.text("ZERDE", 86, 470, 46, "F2", "0.03 0.40 0.70")
	c.text("STUDY", 260, 470, 46, "F2", "0.03 0.40 0.70")

	c.text(strings.ToUpper(data.IssuedAt.Format("January 02, 2006")), 86, 356, 8, "F3", "0.23 0.23 0.26")
	c.text(fitPDFText(data.UserName, 34), 86, 314, 24, "F1", "0.12 0.14 0.18")
	c.text("has successfully completed", 86, 272, 10, "F3", "0.28 0.28 0.30")
	c.text(fitPDFText(data.CourseTitle, 54), 86, 238, 17, "F3", "0.08 0.08 0.10")
	c.text("an online course offered through Zerde Study", 86, 210, 9, "F1", "0.28 0.28 0.30")

	c.text("Completed: "+data.CompletedAt.Format("2006-01-02"), 86, 162, 9, "F1", "0.20 0.20 0.22")
	c.text("Certificate ID: "+data.CertificateNumber, 86, 142, 9, "F1", "0.20 0.20 0.22")

	c.centeredText("COURSE", 670, 430, 15, "F3", "0 0 0")
	c.centeredText("CERTIFICATE", 670, 406, 15, "F3", "0 0 0")
	c.centeredText("Verify this certificate by its number", 670, 92, 8, "F3", "0.20 0.20 0.22")
	c.centeredText(data.CertificateNumber, 670, 76, 8, "F3", "0.20 0.20 0.22")

	return c.pdf()
}

type pdfCanvas struct {
	content bytes.Buffer
}

func newPDFCanvas() *pdfCanvas {
	return &pdfCanvas{}
}

func (c *pdfCanvas) text(value string, x, y, size float64, font, color string) {
	c.content.WriteString("BT\n")
	c.content.WriteString(color + " rg\n")
	c.content.WriteString(fmt.Sprintf("/%s %.1f Tf\n", font, size))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f Td\n", x, y))
	c.content.WriteString("(" + pdfEscape(toPDFText(value)) + ") Tj\n")
	c.content.WriteString("ET\n")
}

func (c *pdfCanvas) centeredText(value string, centerX, y, size float64, font, color string) {
	c.text(value, centerX-textWidth(value, size)/2, y, size, font, color)
}

func (c *pdfCanvas) rectFill(x, y, w, h float64, color string) {
	c.content.WriteString(color + " rg\n")
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f re f\n", x, y, w, h))
}

func (c *pdfCanvas) rectStroke(x, y, w, h float64, color string, width float64) {
	c.content.WriteString(color + " RG\n")
	c.content.WriteString(fmt.Sprintf("%.1f w\n", width))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f re S\n", x, y, w, h))
}

func (c *pdfCanvas) line(x1, y1, x2, y2 float64, color string, width float64) {
	c.content.WriteString(color + " RG\n")
	c.content.WriteString(fmt.Sprintf("%.2f w\n", width))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f m %.1f %.1f l S\n", x1, y1, x2, y2))
}

func (c *pdfCanvas) circle(cx, cy, r float64, color string, width float64) {
	k := 0.5522847498 * r
	c.content.WriteString(color + " RG\n")
	c.content.WriteString(fmt.Sprintf("%.2f w\n", width))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f m\n", cx+r, cy))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f c\n", cx+r, cy+k, cx+k, cy+r, cx, cy+r))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f c\n", cx-k, cy+r, cx-r, cy+k, cx-r, cy))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f c\n", cx-r, cy-k, cx-k, cy-r, cx, cy-r))
	c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f c S\n", cx+k, cy-r, cx+r, cy-k, cx+r, cy))
}

func (c *pdfCanvas) drawRibbon() {
	c.rectFill(574, 18, 192, 456, "0.88 0.90 0.93")
	c.content.WriteString("0.88 0.90 0.93 rg\n")
	c.content.WriteString("574 474 m 766 474 l 766 104 l 670 42 l 574 104 l h f\n")
	c.rectStroke(574, 18, 192, 456, "0.66 0.68 0.70", 0.8)
	c.content.WriteString("0.66 0.68 0.70 RG 0.8 w\n")
	c.content.WriteString("574 474 m 766 474 l 766 104 l 670 42 l 574 104 l h S\n")
	for x := 585.0; x < 758; x += 8 {
		c.line(x, 110, x, 556, "0.82 0.85 0.88", 0.15)
	}
}

func (c *pdfCanvas) drawSeal() {
	c.circle(670, 250, 70, "0.42 0.43 0.44", 2.0)
	c.circle(670, 250, 58, "0.42 0.43 0.44", 0.8)
	c.circle(670, 250, 44, "0.42 0.43 0.44", 0.8)
	for i := 0; i < 56; i++ {
		angle := float64(i) * 2 * math.Pi / 56
		x := 670 + math.Cos(angle)*64
		y := 250 + math.Sin(angle)*64
		c.circle(x, y, 0.8, "0.42 0.43 0.44", 0.5)
	}
	c.centeredText("ZERDE", 670, 260, 18, "F2", "0.36 0.36 0.37")
	c.centeredText("STUDY", 670, 238, 18, "F2", "0.36 0.36 0.37")
	c.centeredText("COURSE CERTIFICATE", 670, 204, 7.5, "F3", "0.36 0.36 0.37")
}

func (c *pdfCanvas) drawWavePattern() {
	for i := 0; i < 18; i++ {
		y := 112.0 + float64(i)*17
		c.content.WriteString("0.96 0.74 0.52 RG\n")
		c.content.WriteString("0.22 w\n")
		c.content.WriteString(fmt.Sprintf("30 %.1f m\n", y))
		for x := 70.0; x < 812; x += 80 {
			y1 := y + math.Sin((x+float64(i)*17)/38)*26
			y2 := y + math.Cos((x+float64(i)*19)/41)*26
			y3 := y + math.Sin((x+float64(i)*23)/43)*20
			c.content.WriteString(fmt.Sprintf("%.1f %.1f %.1f %.1f %.1f %.1f c\n", x-35, y1, x-10, y2, x+40, y3))
		}
		c.content.WriteString("S\n")
	}
}

func (c *pdfCanvas) drawCornerMarks() {
	c.rectStroke(16, 16, 6, 6, "0.62 0.64 0.66", 0.8)
	c.rectStroke(820, 16, 6, 6, "0.62 0.64 0.66", 0.8)
	c.rectStroke(16, 573, 6, 6, "0.62 0.64 0.66", 0.8)
	c.rectStroke(820, 573, 6, 6, "0.62 0.64 0.66", 0.8)
}

func (c *pdfCanvas) pdf() []byte {
	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"<< /Type /Pages /Kids [3 0 R] /Count 1 >>",
		"<< /Type /Page /Parent 2 0 R /MediaBox [0 0 842 595] /Resources << /Font << /F1 4 0 R /F2 5 0 R /F3 6 0 R >> >> /Contents 7 0 R >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica-Bold >>",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Times-Roman >>",
		"<< /Length " + strconv.Itoa(c.content.Len()) + " >>\nstream\n" + c.content.String() + "endstream",
	}

	pdf := bytes.Buffer{}
	pdf.WriteString("%PDF-1.4\n")

	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for i, object := range objects {
		offsets = append(offsets, pdf.Len())
		pdf.WriteString(strconv.Itoa(i+1) + " 0 obj\n")
		pdf.WriteString(object)
		pdf.WriteString("\nendobj\n")
	}

	xrefOffset := pdf.Len()
	pdf.WriteString("xref\n")
	pdf.WriteString(fmt.Sprintf("0 %d\n", len(offsets)))
	pdf.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		pdf.WriteString(fmt.Sprintf("%010d 00000 n \n", offset))
	}
	pdf.WriteString("trailer\n")
	pdf.WriteString(fmt.Sprintf("<< /Size %d /Root 1 0 R >>\n", len(offsets)))
	pdf.WriteString("startxref\n")
	pdf.WriteString(strconv.Itoa(xrefOffset))
	pdf.WriteString("\n%%EOF\n")

	return pdf.Bytes()
}

func fitPDFText(value string, maxRunes int) string {
	value = toPDFText(value)
	runes := []rune(value)
	if len(runes) <= maxRunes {
		return value
	}

	return string(runes[:maxRunes-3]) + "..."
}

func toPDFText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	var b strings.Builder
	for _, r := range value {
		if r >= 32 && r <= 126 {
			b.WriteRune(r)
			continue
		}
		if repl, ok := cyrillicTranslit[r]; ok {
			b.WriteString(repl)
			continue
		}
		b.WriteRune(' ')
	}

	return strings.Join(strings.Fields(b.String()), " ")
}

func textWidth(value string, size float64) float64 {
	return float64(len([]rune(toPDFText(value)))) * size * 0.52
}

func pdfEscape(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"(", "\\(",
		")", "\\)",
		"\r", " ",
		"\n", " ",
	)

	return replacer.Replace(value)
}

var cyrillicTranslit = map[rune]string{
	'\u0410': "A", '\u0430': "a", '\u04d8': "A", '\u04d9': "a",
	'\u0411': "B", '\u0431': "b",
	'\u0412': "V", '\u0432': "v",
	'\u0413': "G", '\u0433': "g", '\u0492': "G", '\u0493': "g",
	'\u0414': "D", '\u0434': "d",
	'\u0415': "E", '\u0435': "e",
	'\u0401': "Yo", '\u0451': "yo",
	'\u0416': "Zh", '\u0436': "zh",
	'\u0417': "Z", '\u0437': "z",
	'\u0418': "I", '\u0438': "i", '\u0406': "I", '\u0456': "i",
	'\u0419': "Y", '\u0439': "y",
	'\u041a': "K", '\u043a': "k", '\u049a': "K", '\u049b': "k",
	'\u041b': "L", '\u043b': "l",
	'\u041c': "M", '\u043c': "m",
	'\u041d': "N", '\u043d': "n", '\u04a2': "N", '\u04a3': "n",
	'\u041e': "O", '\u043e': "o", '\u04e8': "O", '\u04e9': "o",
	'\u041f': "P", '\u043f': "p",
	'\u0420': "R", '\u0440': "r",
	'\u0421': "S", '\u0441': "s",
	'\u0422': "T", '\u0442': "t",
	'\u0423': "U", '\u0443': "u", '\u04b0': "U", '\u04b1': "u", '\u04ae': "U", '\u04af': "u",
	'\u0424': "F", '\u0444': "f",
	'\u0425': "Kh", '\u0445': "kh", '\u04ba': "H", '\u04bb': "h",
	'\u0426': "Ts", '\u0446': "ts",
	'\u0427': "Ch", '\u0447': "ch",
	'\u0428': "Sh", '\u0448': "sh",
	'\u0429': "Shch", '\u0449': "shch",
	'\u042a': "", '\u044a': "",
	'\u042b': "Y", '\u044b': "y",
	'\u042c': "", '\u044c': "",
	'\u042d': "E", '\u044d': "e",
	'\u042e': "Yu", '\u044e': "yu",
	'\u042f': "Ya", '\u044f': "ya",
}
