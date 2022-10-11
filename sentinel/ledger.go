package sentinel

import (
	"fmt"
	"time"

	"github.com/asdine/storm"
)

//RUNNING - status when "agent/edge is successfully running"
const RUNNING = "RUNNING"

//STOPPED - status when "agent/edge is stopped"
const STOPPED = "STOPPED"

//AgentDetails - Required details for agent
type AgentDetails struct {
	AgentID                string
	Mode                   string
	AgentName              string
	AgentVersion           int64
	AppName                string
	UploadRetryCounter     string
	DownloadRetryCounter   string
	BaseURL                string
	HeartBeatFrequency     string
	LogLevel               string
	SentinelPortNumber     string
	SentinelMaxMissesCount string
	EncryptFile            string
	RetainFileOnSuccess    string
	RetainFileOnError      string
}

//TransferLedgerEntry - TransferLedgerEntry structure of DB
type TransferLedgerEntry struct {
	ID             string    `storm:"id"`
	AgentID        string    `json:"agentID"`
	AppName        string    `json:"appName"`
	AgentName      string    `json:"agentName"`
	Namespace      string    `json:"namespace"`
	PartnerName    string    `json:"partnerName"`
	PartnerID      string    `json:"partnerID"`
	FlowName       string    `json:"flowName"`
	FlowID         string    `json:"flowID"`
	DeploymentName string    `json:"deploymentName"`
	Action         string    `json:"action"`
	MetaData       string    `json:"metaData"`
	Timestamp      time.Time `json:"timestamp"`
	EntryType      string    `json:"entryType"`
	SentOrRead     bool      `json:"sentOrRead"`
}

//CentralHeartBeatRequest - agent central heartbeat request structure
type CentralHeartBeatRequest struct {
	MonitoringLedgerEntries []MonitoringLedgerEntry `json:"monitoringLedgerEntries"`
	TransferLedgerEntries   []TransferLedgerEntry   `json:"transferLedgerEntries"`
}

//CentralHeartBeatResponse - central heartbeat response
type CentralHeartBeatResponse struct {
	TransferLedgerEntries []TransferLedgerEntry `json:"transferLedgerEntries"`
	Status                string                `json:"status"`
}

//MonitoringLedger - ledger to main flows creation/deletion
type MonitoringLedger struct {
	DB *storm.DB
}

//MonitoringLedgerEntry - entry structure for flowledger entry
type MonitoringLedgerEntry struct {
	AgentID            string    `storm:"id" json:"agentID"`
	AppName            string    `json:"appName"`
	AgentName          string    `json:"agentName"`
	AgentType          string    `json:"agentType"`
	ResponseAgentID    string    `json:"responseAgentID"`
	HeartBeatFrequency string    `json:"heartBeatFrequency"`
	MACAddress         string    `json:"macAddress"`
	IPAddress          string    `json:"ipAddress"`
	Status             string    `json:"status"`
	Timestamp          time.Time `json:"timestamp"`
	AbsolutePath       string    `json:"absolutePath"`
}

//InitMonitoringLedger - intialize ledger ledger
func InitMonitoringLedger(filePath string) (*MonitoringLedger, error) {
	db, err := storm.Open(filePath)
	if err != nil {
		return nil, err
	}
	monitoringLedger := MonitoringLedger{
		DB: db,
	}
	return &monitoringLedger, nil
}

//AddOrUpdateEntry - add new entry to monitoring ledger
func (db *MonitoringLedger) AddOrUpdateEntry(entry *MonitoringLedgerEntry) error {
	var newEntry MonitoringLedgerEntry
	err := db.DB.One("AgentID", entry.AgentID, &newEntry)
	if fmt.Sprintf("%s", err) == "not found" {
		err = db.DB.Save(entry)
	} else {
		if err != nil {
			return err
		}
		err = db.DB.Update(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetAllEntries - get all entries of monitoring table
func (db *MonitoringLedger) GetAllEntries() ([]MonitoringLedgerEntry, error) {
	query := db.DB.Select()
	monitoringLedgerEntries := []MonitoringLedgerEntry{}
	err := query.Find(&monitoringLedgerEntries)
	if fmt.Sprintf("%s", err) == "not found" {
		return monitoringLedgerEntries, nil
	}
	if err != nil {
		return nil, err
	}

	return monitoringLedgerEntries, nil
}
