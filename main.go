// Copyright 2015 Daniel Theophanes.
// Use of this source code is governed by a zlib-style
// license that can be found in the LICENSE file.

// Simple service that only works by printing a log message every few seconds.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"ds-agent-watcher/models"
	"ds-agent-watcher/sentinel"
	"ds-agent-watcher/utils"

	"github.com/appveen/go-log/logger"
	"github.com/howeyc/gopass"
	"github.com/kardianos/service"
)

var Logger logger.Logger
var svcFlag = flag.String("service", "", "Control the system service.")
var password = flag.String("p", "", "Password")
var confFilePath = flag.String("c", "./conf/agent.conf", "Conf File Path")
var BMResponse = models.LoginAPIResponse{}
var AgentDataFromIM = models.AgentData{}

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	if service.Interactive() {
		Logger.Info("Running sentinel in terminal/interactive mode not allowed. Exiting.")
		verifyAgentPassword(*password)
		os.Exit(0)
	} else {
		Logger.Info("Running sentinel under service manager.")
		verifyAgentPassword(*password)
	}
	p.exit = make(chan struct{})
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() error {
	sentinel.RunSentinel(AgentDataFromIM)
	return nil
}
func (p *program) Stop(s service.Service) error {
	close(p.exit)
	return nil
}

func main() {
	dir, _, err := utils.GetExecutablePathAndName()
	confFilePath := filepath.Join(dir, "..", "conf", "agent.conf")
	confData, err := utils.ReadSentinelConfFile(confFilePath)
	serviceName := "DATASTACKB2BAgentSentinel"
	serviceName = serviceName + confData["sentinel-port-number"]

	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: "Sentinel to monitor b2b agent health",
		Arguments:   []string{"-p", *password},
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		Logger.Fatal(err)
	}
	errs := make(chan error, 5)
	go func() {
		for {
			err := <-errs
			if err != nil {
				Logger.Info(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			if *svcFlag == "stop" && (err.Error() != "Failed to stop "+serviceName+": The specified service does not exist as an installed service." &&
				err.Error() != "Failed to stop "+serviceName+": The service has not been started.") {
				Logger.Info("Action performed : %q\n", *svcFlag)
				Logger.Fatal(err)
			} else if *svcFlag == "uninstall" && err.Error() != "Failed to uninstall "+serviceName+": service "+serviceName+" is not installed" {
				Logger.Info("Action performed : %q\n", *svcFlag)
				Logger.Fatal(err)
			}
		}
		if *svcFlag == "start" {
			Logger.Info(serviceName + " started succcessfully")
		} else if *svcFlag == "install" {
			Logger.Info(serviceName + " installed succcessfully")
		}
		return
	}
	err = s.Run()
	if err != nil {
		Logger.Error(err)
	}
}

func verifyAgentPassword(password string) string {
	pass := ""
	if password == "" {
		fmt.Print("Enter Password : ")
		passBytes, _ := gopass.GetPasswdMasked()
		pass = string(passBytes)
	} else {
		pass = password
	}

	confData, err := utils.ReadSentinelConfFile(*confFilePath)
	payload := models.LoginAPIRequest{
		AgentID:      confData["agent-id"],
		Password:     pass,
		AgentVersion: confData["agent-version"],
	}
	data, err := json.Marshal(payload)
	if err != nil {
		data = nil
		Logger.Error(err)
	}

	URL := "https://{BaseURL}/b2b/bm/auth/login"
	URL = strings.Replace(URL, "{BaseURL}", confData["base-url"], -1)
	Logger.Info("Connecting to integration manager - " + URL)
	client := utils.GetNewHTTPClient(nil)
	req, err := http.NewRequest("POST", URL, bytes.NewReader(data))
	if err != nil {
		data = nil
		Logger.Error("Error from Integration Manager - " + err.Error())
		os.Exit(0)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Close = true
	res, err := client.Do(req)
	if err != nil {
		data = nil
		Logger.Error("Error from Integration Manager - " + err.Error())
		os.Exit(0)
	}
	if res.StatusCode != 200 {
		data = nil
		if res.Body != nil {
			responseData, _ := ioutil.ReadAll(res.Body)
			err = json.Unmarshal(responseData, &BMResponse)
			if err != nil {
				responseData = nil
				Logger.Error("Error unmarshalling response from IM - " + err.Error())
				os.Exit(0)
			}
			Logger.Error("Error from Integration Manager - " + BMResponse.Message)
			os.Exit(0)
		} else {
			Logger.Error("Error from Integration Manager - " + http.StatusText(res.StatusCode))
			os.Exit(0)
		}
	} else {
		bytesData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			data = nil
			bytesData = nil
			Logger.Error("Error reading response body from IM - " + err.Error())
			os.Exit(0)
		}
		if res.Body != nil {
			res.Body.Close()
		}
		err = json.Unmarshal([]byte(bytesData), &AgentDataFromIM)
		if err != nil {
			data = nil
			bytesData = nil
			Logger.Error("Error unmarshalling agent data from IM - " + err.Error())
			os.Exit(0)
		}
		Logger.Info("Agent Successfuly Logged In")
		Logger.Debug("Agent details fetched -  %v ", AgentDataFromIM)
	}
	return string(pass)
}
