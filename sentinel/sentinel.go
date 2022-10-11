package sentinel

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"ds-agent-watcher/log"
	"ds-agent-watcher/models"
	"ds-agent-watcher/utils"

	"github.com/gorilla/mux"
	"github.com/ian-kent/go-log/logger"
)

//Logger - Central Logger
var Logger logger.Logger
var confFilePath = ""
var monitoringLedger *MonitoringLedger
var agentDetails *AgentDetails
var passwordString string
var executableDir string

//RunSentinel - run the sentinel
func RunSentinel(agentData models.AgentData) {
	dir, _, err := utils.GetExecutablePathAndName()
	if err != nil {
		Logger.Info("%s", err)
		os.Exit(0)
	}
	executableDir = dir
	confFilePath = filepath.Join(dir, "..", "conf", "agent.conf")
	dbFilePath := filepath.Join(dir, "..", "conf", "sentinel.db")
	data, err := utils.ReadSentinelConfFile(confFilePath)
	agentDetails = &AgentDetails{
		AgentName:              data["agent-name"],
		AgentID:                data["agent-id"],
		AgentVersion:           agentData.AgentVersion,
		AppName:                agentData.AppName,
		UploadRetryCounter:     agentData.UploadRetryCounter,
		DownloadRetryCounter:   agentData.DownloadRetryCounter,
		BaseURL:                data["base-url"],
		HeartBeatFrequency:     data["heartbeat-frequency"],
		LogLevel:               data["log-level"],
		SentinelPortNumber:     data["sentinel-port-number"],
		SentinelMaxMissesCount: agentData.SentinelMaxMissesCount,
	}

	agentDetails.BaseURL = "https://" + string(agentDetails.BaseURL)
	Logger = log.GetLogger(data["log-level"], "SENTINEL", "001")
	stopServiceFileName := ""
	if runtime.GOOS == "windows" {
		stopServiceFileName = "stop-services.bat"
	} else {
		stopServiceFileName = "stop-services.sh"
	}
	stopServicesFilePath := filepath.Join(dir, "..", stopServiceFileName)
	err = utils.UpdateValuesInStopServicesFile(stopServicesFilePath, data["agent-port-number"], data["sentinel-port-number"])
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}

	getRunningORPendingFlowFromPartnerManagerAfterRestart()
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	fl, err := InitMonitoringLedger(dbFilePath)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	monitoringLedger = fl
	go runServer()
	go updateAgentStatus(confFilePath)
}

func healAgent(data map[string]string) error {
	Logger.Info("Trying to start agent...")
	serviceName := "DATASTACKB2BAgent"
	serviceName = serviceName + data["agent-port-number"]
	cmd := exec.Command("service", serviceName, "start")
	var out bytes.Buffer
	cmd.Stderr = &out
	if runtime.GOOS == "windows" {
		cmd = exec.Command("sc", "query", serviceName)
		output, err := cmd.Output()
		if err != nil {
			Logger.Error(err)
			Logger.Error("Output = " + out.String())
			os.Exit(0)
		}
		if !strings.Contains(string(output), "START_PENDING") {
			Logger.Info("running windows command = " + "datastack-agent.exe" + " -p " + "*******" + " -c " + executableDir + string(os.PathSeparator) + ".." + string(os.PathSeparator) + "conf" + string(os.PathSeparator) + "agent.conf" + " -service" + " start")
			cmd = exec.Command(executableDir+"datastack-agent.exe", "-p", passwordString, "-c", executableDir+".."+string(os.PathSeparator)+"conf"+string(os.PathSeparator)+"agent.conf", "-service", "start")
		} else {
			Logger.Info("Agent already Running")
			return nil
		}
	} else if runtime.GOOS == "darwin" {
		Logger.Info("running mac command = ./datastack-agent -p <password> -c <conf-file-path> -service start")
		cmd.Dir = executableDir
		cmd = exec.Command(executableDir+"datastack-agent", "-p", passwordString, "-c", executableDir+string(os.PathSeparator)+".."+string(os.PathSeparator)+"conf"+string(os.PathSeparator)+"agent.conf", "-service", "start")
	} else {
		Logger.Info("running linux command = service DATASTACKB2BAgent start")
		Logger.Info("service " + serviceName + " start")
	}
	output, err := cmd.Output()
	if err != nil {
		Logger.Error(err)
		Logger.Error("Output = " + out.String())
		os.Exit(0)
	} else {
		Logger.Info("Agent start command response = " + string(output))
	}
	return nil
}

func runServer() {
	var tempServer *http.Server
	r := mux.NewRouter()
	tempServer = &http.Server{
		Addr:    "0.0.0.0:" + agentDetails.SentinelPortNumber,
		Handler: r,
	}
	err := tempServer.ListenAndServe()
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	Logger.Info("Server started for sentinel")
	Logger.Info("Begin listening port: %s", agentDetails.SentinelPortNumber)
	for {
		Logger.Info("Request Received....")
		if err != nil {
			Logger.Error(fmt.Sprintf("%s", err))
			continue
		}
	}
}

func stopAgent(data map[string]string) error {
	Logger.Info("Trying to stop agent...")
	serviceName := "DATASTACKB2BAgent"
	serviceName = serviceName + data["agent-port-number"]
	cmd := exec.Command("service", serviceName, "stop")
	if runtime.GOOS == "windows" {
		Logger.Info("Running windows command = net stop DATASTACKB2BAgent")
		cmd = exec.Command(`FOR /F "tokens=3" %%A IN ('sc queryex DATASTACKB2BAgent ^| findstr PID') DO (SET pid=%%A)
		IF "%pid%" NEQ "0" (
		 taskkill /F /PID %pid%
	   )`)
	} else if runtime.GOOS == "darwin" {
		Logger.Info("Running mac command = ./datastack-agent -p <password> -c <conf-file-path> -service stop")
		cmd.Dir = executableDir
		cmd = exec.Command(executableDir+"datastack-agent", "-p", passwordString, "-c", executableDir+string(os.PathSeparator)+".."+string(os.PathSeparator)+"conf"+string(os.PathSeparator)+"agent.conf", "-service", "stop")
	} else {
		Logger.Info("Running linux command = launchctl stop DATASTACKB2BAgent")
		Logger.Info("Service DATASTACKB2BAgent stop")
	}
	var out bytes.Buffer
	cmd.Stderr = &out
	output, err := cmd.Output()
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		Logger.Error("Output = ", out.String())
		return err
	} else {
		Logger.Info("Agent stop command response = " + string(output))
	}
	return nil
}

func updateAgentStatus(configFilePath string) {
	Logger.Info("Starting Sentinel...")
	data, err := utils.ReadSentinelConfFile(confFilePath)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	maxMissesCount, err := strconv.Atoi(agentDetails.SentinelMaxMissesCount)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	heartBeatFrequency, err := strconv.Atoi(data["heartbeat-frequency"])
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	for {
		Logger.Info("Fetching Ledger State...")
		entries, err := monitoringLedger.GetAllEntries()
		if err != nil {
			Logger.Error(fmt.Sprintf("%s", err))
			os.Exit(0)
		}
		if len(entries) == 0 {
			time.Sleep(time.Duration(heartBeatFrequency) * time.Second)
			continue
		}
		Logger.Info("Number of entries: %s", len(entries))
		for _, entry := range entries {
			Logger.Info("Status of the entry: %S", entry.Status)
			if entry.Status == "STOPPED" {
				Logger.Info("Agent received stop/delete request from server. Stopping agent and sentinel..")
				entry.Status = "RUNNING"
				err = monitoringLedger.AddOrUpdateEntry(&entry)
				if err != nil {
					Logger.Error(fmt.Sprintf("%s", err))
					continue
				}
				stopAgent(data)
				os.Exit(0)
			}
			lastTimeStamp := entry.Timestamp
			frequency, err := strconv.Atoi(entry.HeartBeatFrequency)
			if err != nil {
				Logger.Error(fmt.Sprintf("%s", err))
				os.Exit(0)
			}
			latestTimeStamp := time.Now()
			_, _, _, _, _, seconds := utils.TimeDifference(latestTimeStamp, lastTimeStamp)
			Logger.Info("Seconds: %s", seconds)
			temp := maxMissesCount * frequency
			Logger.Info("Temp: %s", temp)
			if seconds > maxMissesCount*frequency*3 {
				err := monitoringLedger.AddOrUpdateEntry(&entry)
				if err != nil {
					Logger.Error(fmt.Sprintf("%s", err))
					os.Exit(0)
				}
				go healAgent(data)
			}
		}
		time.Sleep(time.Duration(heartBeatFrequency) * time.Second)
	}
}

func getRunningORPendingFlowFromPartnerManagerAfterRestart() {
	Logger.Info("Fetching Sentinel/Agent Status")
	var client = utils.GetNewHTTPClient(nil)
	data := CentralHeartBeatResponse{}
	headers := make(map[string]string)
	headers["AgentID"] = agentDetails.AgentID
	headers["AgentName"] = agentDetails.AgentName
	headers["AgentType"] = "Agent"
	request := CentralHeartBeatRequest{}
	err := utils.MakeJSONRequest(client, agentDetails.BaseURL+"/getrunningorpendingflowsfrompartnermanager", request, headers, &data)
	if err != nil {
		Logger.Error(fmt.Sprintf("%s", err))
		os.Exit(0)
	}
	if data.Status == "DISABLED" {
		Logger.Error("Sentinel cannot be started.The agent has been stopped/disabled from the server.")
		os.Exit(0)
	}
}
