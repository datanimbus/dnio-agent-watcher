package models

import (
	"time"
)

//LoginAPIRequest - agent login api request structure
type LoginAPIRequest struct {
	AgentID      string `json:"agentId"`
	Password     string `json:"password"`
	AgentVersion string `json:"agentVersion"`
}

//LoginAPIResponse - agent login api response structure
type LoginAPIResponse struct {
	Message string `json:"message"`
}

//AgentDataFromIM - agent information from IM
type AgentData struct {
	ID                     string `json:"_id"`
	Active                 bool   `json:"active"`
	AppName                string `json:"app"`
	AgentName              string `json:"name"`
	AgentVersion           int64  `json:"__v"`
	EncryptFile            bool   `json:"encryptFile"`
	RetainFileOnSuccess    bool   `json:"retainFileOnSuccess"`
	RetainFileOnError      bool   `json:"retainFileOnError"`
	Internal               bool   `json:"internal"`
	Token                  string `json:"token"`
	Secret                 string `json:"secret"`
	EncryptionKey          string `json:"encryptionKey"`
	UploadRetryCounter     string `json:"uploadRetryCounter"`
	DownloadRetryCounter   string `json:"downloadRetryCounter"`
	MaxConcurrentUploads   int    `json:"maxConcurrentUploads"`
	MaxConcurrentDownloads int    `json:"maxConcurrentDownloads"`
	SentinelMaxMissesCount string `json:"sentinelMaxMissesCount"`
}

//CentralHeartBeatRequest - agent central heartbeat request structure
type CentralHeartBeatRequest struct {
	MonitoringLedgerEntries []MonitoringLedgerEntry `json:"monitoringLedgerEntries"`
	TransferLedgerEntries   []TransferLedgerEntry   `json:"transferLedgerEntries"`
}

//CentralHeartBeatResponse - response structure to central heartbeat request
type CentralHeartBeatResponse struct {
	TransferLedgerEntries     []TransferLedgerEntry `json:"transferLedgerEntries"`
	Status                    string                `json:"status"`
	AgentMaxConcurrentUploads int                   `json:"agentMaxConcurrentUploads"`
}

//TransferLedgerEntry - TransferLedgerEntry structure of DB
type TransferLedgerEntry struct {
	ID             string    `storm:"id" bson:"id"`
	AgentID        string    `json:"agentID" bson:"AgentID"`
	AppName        string    `json:"appName" bson:"AppName"`
	AgentName      string    `json:"agentName" bson:"AgentName"`
	FlowName       string    `json:"flowName" bson:"FlowName"`
	FlowID         string    `json:"flowID" bson:"FlowID"`
	DeploymentName string    `json:"deploymentName" bson:"DeploymentName"`
	Action         string    `json:"action" bson:"Action"`
	MetaData       string    `json:"metaData" bson:"MetaData"`
	Timestamp      time.Time `json:"timestamp" bson:"Timestamp"`
	SentOrRead     bool      `json:"sentOrRead" bson:"SentOrRead"`
	Status         string    `json:"status" bson:"Status"`
}

//MonitoringLedgerEntry - entry structure for flowledger entry
type MonitoringLedgerEntry struct {
	AgentID            string         `storm:"id" json:"agentID" bson:"AgentID"`
	AppName            string         `json:"appName" bson:"AppName"`
	AgentName          string         `json:"agentName" bson:"AgentName"`
	HeartBeatFrequency string         `json:"heartBeatFrequency" bson:"HeartBeatFrequency"`
	MACAddress         string         `json:"macAddress" bson:"MACAddress"`
	IPAddress          string         `json:"ipAddress" bson:"IPAddress"`
	Status             string         `json:"status" bson:"Status"`
	Timestamp          time.Time      `json:"timestamp" bson:"Timestamp"`
	AbsolutePath       string         `json:"absolutePath" bson:"AbsolutePath"`
	PendingFilesCount  []PendingFiles `json:"pendingFiles" bson:"PendingFiles"`
	Release            string         `json:"release" bson:"Release"`
}

//PendingFiles - pending files struct
type PendingFiles struct {
	FlowID string `json:"flowID" bson:"FlowID"`
	Count  int    `json:"count" bson:"Count"`
}
