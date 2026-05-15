package storage

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	awsAlgorithm     = "AWS4-HMAC-SHA256"
	awsService       = "s3"
	unsignedPayload  = "UNSIGNED-PAYLOAD"
	maxPresignExpiry = 7 * 24 * time.Hour
)

type MinIOConfig struct {
	Endpoint       string
	PublicEndpoint string
	AccessKey      string
	SecretKey      string
	Bucket         string
	Region         string
	UseSSL         bool
	PresignTTL     time.Duration
}

type MinIO struct {
	endpoint       string
	publicEndpoint string
	accessKey      string
	secretKey      string
	bucket         string
	region         string
	useSSL         bool
	presignTTL     time.Duration
}

func NewMinIO(cfg MinIOConfig) (*MinIO, error) {
	if strings.TrimSpace(cfg.AccessKey) == "" {
		return nil, errors.New("minio access key is required")
	}
	if strings.TrimSpace(cfg.SecretKey) == "" {
		return nil, errors.New("minio secret key is required")
	}
	if strings.TrimSpace(cfg.Bucket) == "" {
		return nil, errors.New("minio bucket is required")
	}

	endpoint := normalizeEndpoint(cfg.Endpoint)
	publicEndpoint := normalizeEndpoint(cfg.PublicEndpoint)
	if publicEndpoint == "" {
		publicEndpoint = endpoint
	}
	if publicEndpoint == "" {
		return nil, errors.New("minio endpoint is required")
	}

	region := strings.TrimSpace(cfg.Region)
	if region == "" {
		region = "us-east-1"
	}

	presignTTL := cfg.PresignTTL
	if presignTTL <= 0 {
		presignTTL = 30 * time.Minute
	}
	if presignTTL > maxPresignExpiry {
		presignTTL = maxPresignExpiry
	}

	return &MinIO{
		endpoint:       endpoint,
		publicEndpoint: publicEndpoint,
		accessKey:      strings.TrimSpace(cfg.AccessKey),
		secretKey:      strings.TrimSpace(cfg.SecretKey),
		bucket:         strings.Trim(strings.TrimSpace(cfg.Bucket), "/"),
		region:         region,
		useSSL:         cfg.UseSSL,
		presignTTL:     presignTTL,
	}, nil
}

func (m *MinIO) PutObject(ctx context.Context, objectKey string, body io.ReadSeeker, size int64, contentType string) error {
	objectKey = strings.Trim(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return errors.New("minio object key is required")
	}
	if body == nil {
		return errors.New("minio object body is required")
	}
	if _, err := body.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind object body: %w", err)
	}

	payloadHash, err := sha256Hex(body)
	if err != nil {
		return fmt.Errorf("hash object body: %w", err)
	}
	if _, err = body.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("rewind object body after hash: %w", err)
	}

	endpointURL, err := m.internalURL()
	if err != nil {
		return err
	}

	canonicalURI := "/" + pathEncode(m.bucket) + "/" + pathEncode(objectKey)
	endpointURL.Path = "/" + m.bucket + "/" + objectKey
	endpointURL.RawPath = canonicalURI

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpointURL.String(), body)
	if err != nil {
		return fmt.Errorf("create minio put request: %w", err)
	}
	req.ContentLength = size
	if strings.TrimSpace(contentType) != "" {
		req.Header.Set("Content-Type", contentType)
	}

	m.signRequest(req, canonicalURI, "", payloadHash)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("put minio object: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("put minio object: unexpected status %s", resp.Status)
	}

	return nil
}

func (m *MinIO) PresignGetObject(objectKey string) (string, error) {
	objectKey = strings.Trim(strings.TrimSpace(objectKey), "/")
	if objectKey == "" {
		return "", errors.New("minio object key is required")
	}

	now := time.Now().UTC()
	date := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	scope := fmt.Sprintf("%s/%s/%s/aws4_request", date, m.region, awsService)

	endpointURL, err := m.publicURL()
	if err != nil {
		return "", err
	}

	canonicalURI := "/" + pathEncode(m.bucket) + "/" + pathEncode(objectKey)
	host := endpointURL.Host
	if host == "" {
		host = m.publicEndpoint
	}

	query := map[string]string{
		"X-Amz-Algorithm":     awsAlgorithm,
		"X-Amz-Credential":    m.accessKey + "/" + scope,
		"X-Amz-Date":          amzDate,
		"X-Amz-Expires":       fmt.Sprintf("%.0f", m.presignTTL.Seconds()),
		"X-Amz-SignedHeaders": "host",
	}

	canonicalQuery := canonicalQueryString(query)
	canonicalHeaders := "host:" + host + "\n"
	canonicalRequest := strings.Join([]string{
		"GET",
		canonicalURI,
		canonicalQuery,
		canonicalHeaders,
		"host",
		unsignedPayload,
	}, "\n")

	requestHash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		awsAlgorithm,
		amzDate,
		scope,
		hex.EncodeToString(requestHash[:]),
	}, "\n")

	signingKey := signingKey(m.secretKey, date, m.region, awsService)
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))

	query["X-Amz-Signature"] = signature
	endpointURL.Path = "/" + m.bucket + "/" + objectKey
	endpointURL.RawPath = canonicalURI
	endpointURL.RawQuery = canonicalQueryString(query)

	return endpointURL.String(), nil
}

func (m *MinIO) signRequest(req *http.Request, canonicalURI, query, payloadHash string) {
	now := time.Now().UTC()
	date := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	scope := fmt.Sprintf("%s/%s/%s/aws4_request", date, m.region, awsService)

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)

	host := req.Host
	if host == "" {
		host = req.URL.Host
	}
	headers := map[string]string{
		"host":                 host,
		"x-amz-content-sha256": payloadHash,
		"x-amz-date":           amzDate,
	}
	signedHeaders := []string{"host", "x-amz-content-sha256", "x-amz-date"}
	if contentType := strings.TrimSpace(req.Header.Get("Content-Type")); contentType != "" {
		headers["content-type"] = contentType
		signedHeaders = append([]string{"content-type"}, signedHeaders...)
	}

	canonicalHeaders := strings.Builder{}
	for _, header := range signedHeaders {
		canonicalHeaders.WriteString(header)
		canonicalHeaders.WriteString(":")
		canonicalHeaders.WriteString(headers[header])
		canonicalHeaders.WriteString("\n")
	}

	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		query,
		canonicalHeaders.String(),
		strings.Join(signedHeaders, ";"),
		payloadHash,
	}, "\n")

	requestHash := sha256.Sum256([]byte(canonicalRequest))
	stringToSign := strings.Join([]string{
		awsAlgorithm,
		amzDate,
		scope,
		hex.EncodeToString(requestHash[:]),
	}, "\n")

	signature := hex.EncodeToString(hmacSHA256(signingKey(m.secretKey, date, m.region, awsService), stringToSign))
	req.Header.Set("Authorization", fmt.Sprintf(
		"%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		awsAlgorithm,
		m.accessKey,
		scope,
		strings.Join(signedHeaders, ";"),
		signature,
	))
}

func (m *MinIO) internalURL() (*url.URL, error) {
	return m.endpointURL(m.endpoint, "minio endpoint")
}

func (m *MinIO) publicURL() (*url.URL, error) {
	return m.endpointURL(m.publicEndpoint, "minio public endpoint")
}

func (m *MinIO) endpointURL(endpoint, name string) (*url.URL, error) {
	scheme := "http"
	if m.useSSL {
		scheme = "https"
	}

	value := endpoint
	if !strings.Contains(value, "://") {
		value = scheme + "://" + value
	}

	parsed, err := url.Parse(value)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", name, err)
	}
	if parsed.Host == "" {
		return nil, fmt.Errorf("%s host is required", name)
	}

	return parsed, nil
}

func normalizeEndpoint(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}

	if strings.Contains(value, "://") {
		parsed, err := url.Parse(value)
		if err == nil && parsed.Host != "" {
			return parsed.Host
		}
	}

	host, port, err := net.SplitHostPort(value)
	if err == nil && host == "" {
		return net.JoinHostPort("localhost", port)
	}

	return strings.TrimRight(value, "/")
}

func canonicalQueryString(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(values[key]))
	}

	return strings.Join(parts, "&")
}

func pathEncode(value string) string {
	parts := strings.Split(value, "/")
	for i := range parts {
		parts[i] = strings.ReplaceAll(url.PathEscape(parts[i]), "+", "%20")
	}

	return strings.Join(parts, "/")
}

func signingKey(secretKey, date, region, service string) []byte {
	dateKey := hmacSHA256([]byte("AWS4"+secretKey), date)
	regionKey := hmacSHA256(dateKey, region)
	serviceKey := hmacSHA256(regionKey, service)
	return hmacSHA256(serviceKey, "aws4_request")
}

func hmacSHA256(key []byte, value string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return mac.Sum(nil)
}

func sha256Hex(r io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, r); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
