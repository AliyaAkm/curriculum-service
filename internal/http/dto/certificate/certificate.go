package certificate

import "time"

type IssueCertificateResponse struct {
	CertificateNumber string    `json:"certificate_number"`
	IssuedAt          time.Time `json:"issued_at"`
	CompletedAt       time.Time `json:"completed_at"`
	DownloadURL       string    `json:"download_url"`
}
