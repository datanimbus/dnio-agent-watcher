package utils

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//GetNewHTTPClient - get new http client with/without TLS
func GetNewHTTPClient(transport *http.Transport) *http.Client {
	if transport != nil {
		return &http.Client{Transport: transport}
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}
}

//TLSConfigWithActualFileData - tls config with actual certs
func TLSConfigWithActualFileData(certFile []byte, keyFile []byte, serverCertFile []byte) *http.Transport {
	cert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error1 : ", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(serverCertFile)
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()
	tlsConfig.InsecureSkipVerify = true
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return transport
}

//WriteLinesToFile - write lines to a file from string slice
func WriteLinesToFile(lines []string, path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	file.Truncate(0)
	file.Seek(0, 0)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	return w.Flush()
}

//ReadLinesFromFile - Read lines from a file in string slice
func ReadLinesFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

//CreateOrUpdateSentinelConfFile - create/update sentinel conf file
func CreateOrUpdateSentinelConfFile(data map[string]string, filePath string) error {
	items := []string{}
	for k, v := range data {
		items = append(items, k+"="+v)
	}
	err := WriteLinesToFile(items, filePath)
	if err != nil {
		return err
	}
	return nil
}

//ReadSentinelConfFile - read sentinel conf file into map
func ReadSentinelConfFile(filePath string) (map[string]string, error) {
	data, err := ReadLinesFromFile(filePath)
	if err != nil {
		return nil, err
	}
	mappedValues := make(map[string]string)
	for _, item := range data {
		values := strings.Split(item, "=")
		if len(values) == 2 {
			mappedValues[values[0]] = values[1]
		} else {
			mappedValues[values[0]] = ""
		}
	}
	return mappedValues, nil
}

//TimeDifference - get difference between two time.Time objects in all forms
func TimeDifference(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

//GetExecutablePathAndName - get executable path and name
func GetExecutablePathAndName() (string, string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", "", err
	}
	d, f := filepath.Split(executablePath)
	return d, f, nil
}

//MakeJSONRequest - utility function to make JSON request
func MakeJSONRequest(client *http.Client, url string, payload interface{}, headers map[string]string, response interface{}) error {

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, er := http.NewRequest("POST", url, bytes.NewReader(data))
	if er != nil {
		return er
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")
	res, errr := client.Do(req)
	if errr != nil {

		return errr
	}
	if res.StatusCode != 200 {
		return errors.New("request failed")
	}
	bytesData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytesData, &response)
	if err != nil {
		return err
	}
	return nil
}

//PrepareTLSConfigWithEncodedTrustStore - returns only tls config, not transport and it is not retrievable from transport
func PrepareTLSConfigWithEncodedTrustStore(certFile []byte, keyFile []byte, trustCerts []string) *tls.Config {
	cert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error1 : ", err)
	}
	caCertPool := x509.NewCertPool()
	for i := 0; i < len(trustCerts); i++ {
		currentEncodedCert := trustCerts[i]
		decodedCert, err := base64.StdEncoding.DecodeString(currentEncodedCert)
		if err != nil {
			log.Fatal("Error1 : ", err)
		}
		caCertPool.AppendCertsFromPEM(decodedCert)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.BuildNameToCertificate()
	return tlsConfig
}

//PrepareTLSTransportConfigWithEncodedTrustStore - preparing tls config using encoded trust store
func PrepareTLSTransportConfigWithEncodedTrustStore(certFile []byte, keyFile []byte, trustCerts []string) *http.Transport {
	cert, err := tls.X509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal("Error1 : ", err)
	}
	caCertPool := x509.NewCertPool()
	for i := 0; i < len(trustCerts); i++ {
		currentEncodedCert := trustCerts[i]
		decodedCert, err := base64.StdEncoding.DecodeString(currentEncodedCert)
		if err != nil {
			log.Fatal("Error1 : ", err)
		}
		caCertPool.AppendCertsFromPEM(decodedCert)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.InsecureSkipVerify = true
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	return transport
}

func UpdateValuesInStopServicesFile(path string, agentPortNumber string, sentinelPortNumber string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	input1, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return err
	}

	if !(bytes.Contains(input1, []byte("DATASTACKB2BAgent"+agentPortNumber)) || bytes.Contains(input1, []byte("DATASTACKB2BAgentSentinel"+sentinelPortNumber))) {
		input2 := bytes.Replace(input1, []byte("DATASTACKB2BAgent"), []byte("DATASTACKB2BAgent"+agentPortNumber), -1)
		output := bytes.Replace(input2, []byte("DATASTACKB2BAgent"+agentPortNumber+"Sentinel"), []byte("DATASTACKB2BAgentSentinel"+sentinelPortNumber), -1)

		if err := os.Truncate(path, 0); err != nil {
			log.Printf("Failed to truncate: %v", err)
			return err
		}

		if err = ioutil.WriteFile(path, output, 0666); err != nil {
			return err
		}
	}
	return nil
}
