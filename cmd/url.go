// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/tls"
	"fmt"
	"github.com/parnurzeal/gorequest"
	chrm "go-screenshot/chrome"
	"go-screenshot/storage"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// urlCmd represents the url command
var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		urls := strings.Split(screenshotURL, ";")
		for _, urlString := range urls {
			do(urlString)
		}
	},
}

func do(urlString string) {
	u, err := url.ParseRequestURI(urlString)
	if err != nil {
		fmt.Println("Invalid URL specified")
	}
	// Process this URL
	ProcessURL(u, &chrome, &db, waitTimeout)
}

func init() {
	rootCmd.AddCommand(urlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// urlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// urlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	urlCmd.Flags().StringVarP(&screenshotURL, "values", "v", "", "The URL to screenshot")

}

const (
	// HTTP is the prefix for http:// urls
	HTTP string = "http://"
	// HTTPS is the prefox for https:// urls
	HTTPS string = "https://"
)

// ProcessURL processes a URL
func ProcessURL(url *url.URL, chrome *chrm.Chrome, db *storage.Storage, timeout int) {

	// prepare some storage for this URL
	HTTPResponseStorage := storage.HTTResponse{URL: url.String()}

	// prepare a storage instance for this URL
	//log.WithField("url", url).Debug("Processing URL")

	request := gorequest.New().Timeout(time.Duration(timeout)*time.Second).
		TLSClientConfig(&tls.Config{InsecureSkipVerify: true}).
		Set("User-Agent", chrome.UserAgent)

	resp, _, errs := request.Get(url.String()).End()
	if errs != nil {
		//log.WithFields(log.Fields{"url": url, "error": errs}).Error("Failed to query url")

		return
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
		storageHeader := storage.HTTPHeader{Key: k, Value: headerValue}
		HTTPResponseStorage.Headers = append(HTTPResponseStorage.Headers, storageHeader)

		//log.WithFields(log.Fields{"url": url, k: headerValue}).Info("Response header")
	}

	// Parse any TLS information
	if resp.TLS != nil {

		// storage for the TLS information
		SSLCertificate := storage.SSLCertificate{}

		for _, c := range resp.TLS.PeerCertificates {

			SSLCertificateAttributes := storage.SSLCertificateAttributes{
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
	chrome.ScreenshotURL(finalURL, dst)

	// Update the database with this entry
	db.SetHTTPData(&HTTPResponseStorage)
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
