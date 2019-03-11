package web

import (
	"crypto/tls"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/parnurzeal/gorequest"
)

// ProcessURL processes a URL
func (chrome *Chrome) ProcessURL(url *url.URL, timeout int) *HTTResponse {

	// prepare some storage for this URL
	HTTPResponseStorage := HTTResponse{URL: url.String()}

	// prepare a storage instance for this URL
	//log.WithField("url", url).Debug("Processing URL")

	request := gorequest.New().Timeout(time.Duration(timeout)*time.Second).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("User-Agent", chrome.UserAgent)

	resp, _, errs := request.Get(url.String()).End()
	if errs != nil {
		//log.WithFields(log.Fields{"url": url, "error": errs}).Error("Failed to query url")

		return nil
	}

	// update the response code
	HTTPResponseStorage.ResponseCode = resp.StatusCode
	HTTPResponseStorage.ResponseCodeString = resp.Status
	//log.WithFields(log.Fields{"url": url, "status": resp.Status}).Info("Response code")

	finalURL := resp.Request.URL
	HTTPResponseStorage.FinalURL = resp.Request.URL.String()
	//log.WithFields(log.Fields{"url": url, "final-url": finalURL}).Info("Final URL after redirects")

	// process response headers
	for k, v := range resp.Header {
		headerValue := strings.Join(v, ", ")
		storageHeader := HTTPHeader{Key: k, Value: headerValue}
		HTTPResponseStorage.Headers = append(HTTPResponseStorage.Headers, storageHeader)

		//log.WithFields(log.Fields{"url": url, k: headerValue}).Info("Response header")
	}

	// Parse any TLS information
	if resp.TLS != nil {

		// storage for the TLS information
		SSLCertificate := SSLCertificate{}

		for _, c := range resp.TLS.PeerCertificates {

			SSLCertificateAttributes := SSLCertificateAttributes{
				SubjectCommonName:  c.Subject.CommonName,
				IssuerCommonName:   c.Issuer.CommonName,
				SignatureAlgorithm: c.SignatureAlgorithm.String(),
			}

			//log.WithFields(log.Fields{"url": url, "common_name": c.Subject.CommonName}).Info("Certificate chain common name")
			//log.WithFields(log.Fields{"url": url, "signature-alg": c.SignatureAlgorithm}).Info("Signature algorithm")
			//log.WithFields(log.Fields{"url": url, "pubkey-alg": c.PublicKeyAlgorithm}).Info("Public key algorithm")
			//log.WithFields(log.Fields{"url": url, "issuer": c.Issuer.CommonName}).Info("Issuer")

			for _, d := range c.DNSNames {

				SSLCertificateAttributes.DNSNames = append(SSLCertificateAttributes.DNSNames, d)
				//log.WithFields(log.Fields{"url": url, "dns-names": d}).Info("DNS Name")
			}

			SSLCertificate.PeerCertificates = append(SSLCertificate.PeerCertificates, SSLCertificateAttributes)
		}

		SSLCertificate.CipherSuite = resp.TLS.CipherSuite
		HTTPResponseStorage.SSL = SSLCertificate
		//log.WithFields(log.Fields{"url": url, "cipher-suite": resp.TLS.CipherSuite}).Info("Cipher suite in use")
	}

	// Generate a safe filename to use
	fname := SafeFileName(url.String()) + ".png"

	// Get the tull path where we will be saving the screenshot to
	dst := filepath.Join(chrome.ScreenshotPath, fname)

	HTTPResponseStorage.ScreenshotFile = dst
	//log.WithFields(log.Fields{"url": url, "file-name": fname, "destination": dst}).
	//	Debug("Generated filename for screenshot")

	// Screenshot the URL
	chrome.screenshotURL(finalURL, dst)

	return &HTTPResponseStorage
}

// SafeFileName return a safe string that can be used in file names
func SafeFileName(str string) string {

	name := strings.ToLower(str)
	name = strings.Trim(name, " ")

	separators, err := regexp.Compile(`[ &_=+:]`)
	if err == nil {
		name = separators.ReplaceAllString(name, "-")
	}

	legal, err := regexp.Compile(`[^[:alnum:]-.]`)
	if err == nil {
		name = legal.ReplaceAllString(name, "")
	}

	for strings.Contains(name, "--") {
		name = strings.Replace(name, "--", "-", -1)
	}

	return name
}
