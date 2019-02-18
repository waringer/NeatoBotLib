package NeatoBotLib

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Robot holds the basic infos about an robot
type Robot struct {
	Serial                            string   `json:"serial"`
	Prefix                            string   `json:"prefix"`
	Name                              string   `json:"name"`
	Model                             string   `json:"model"`
	Timezone                          string   `json:"timezone"`
	SecretKey                         string   `json:"secret_key"`
	PurchasedAt                       string   `json:"purchased_at"`
	LinkedAt                          string   `json:"linked_at"`
	NucleoURL                         string   `json:"nucleo_url"`
	Traits                            []string `json:"traits"`
	ProofOfPurchaseURL                string   `json:"proof_of_purchase_url"`
	ProofOfPurchaseURLValidForSeconds string   `json:"proof_of_purchase_url_valid_for_seconds"`
	ProofOfPurchaseGeneratedAt        string   `json:"proof_of_purchase_generated_at"`
	MacAddress                        string   `json:"mac_address"`
	CreatedAt                         string   `json:"created_at"`
	LatestExplorationMapID            string   `json:"latest_exploration_map_id"`
	PersistentMaps                    []string `json:"persistent_maps"`
}

type robotCommand struct {
	ID      string `json:"reqId"`
	Command string `json:"cmd"`
}

// RobotState holds the status of an robot
type RobotState struct {
	Version int    `json:"version"`
	ReqID   string `json:"reqId"`
	Result  string `json:"result"`
	Error   string `json:"error"`
	Data    string `json:"data"`
	State   int    `json:"state"`
	Action  int    `json:"action"`

	Cleaning struct {
		Category   int `json:"category"`
		Mode       int `json:"mode"`
		Modifier   int `json:"modifier"`
		SpotWidth  int `json:"spotWidth"`
		SpotHeight int `json:"spotHeight"`
	} `json:"cleaning"`

	Details struct {
		IsCharging        bool `json:"isCharging"`
		IsDocked          bool `json:"isDocked"`
		IsScheduleEnabled bool `json:"isScheduleEnabled"`
		DockHasBeenSeen   bool `json:"dockHasBeenSeen"`
		Charge            int  `json:"charge"`
	} `json:"details"`

	AvailableCommands struct {
		Start    bool `json:"start"`
		Stop     bool `json:"stop"`
		Pause    bool `json:"pause"`
		Resume   bool `json:"resume"`
		GoToBase bool `json:"goToBase"`
	} `json:"availableCommands"`

	AvailableServices struct {
		HouseCleaning  string `json:"houseCleaning"`
		SpotCleaning   string `json:"spotCleaning"`
		ManualCleaning string `json:"manualCleaning"`
		EasyConnect    string `json:"easyConnect"`
		Schedule       string `json:"schedule"`
	} `json:"availableServices"`

	Meta struct {
		ModelName string `json:"modelName"`
		Firmware  string `json:"firmware"`
	} `json:"meta"`
}

// AuthResponse holds the session infos
type AuthResponse struct {
	AccessToken string `json:"access_token"`
	CurrentTime string `json:"current_time"`
}

// DashResponse holds all account infos
type DashResponse struct {
	ID              string  `json:"id"`
	EMail           string  `json:"email"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Locale          string  `json:"locale"`
	CountryCode     string  `json:"country_code"`
	Developer       string  `json:"developer"`
	Newsletter      string  `json:"newsletter"`
	CreatedAt       string  `json:"created_at"`
	VerifiedAt      string  `json:"verified_at"`
	Robots          []Robot `json:"robots"`
	RecentFirmwares string  `json:"recent_firmwares"`
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func makeAuth(rob Robot, commandData []byte) (string, string, string) {
	vendor := "neato"

	switch rob.Model {
	case "VR200":
		vendor = "vorwerk"
	}

	utcDate := time.Now().UTC().Format("Mon, 2 Jan 2006 15:04:05 GMT")

	commandMessage := strings.Join([]string{strings.ToLower(rob.Serial), utcDate, string(commandData)}, "\n")

	authHMAC := hmac.New(sha256.New, []byte(rob.SecretKey))
	authHMAC.Write([]byte(commandMessage))
	authString := "NEATOAPP " + hex.EncodeToString(authHMAC.Sum(nil))

	CommandURL := rob.NucleoURL + "/vendors/" + vendor + "/robots/" + rob.Serial + "/messages"

	return utcDate, CommandURL, authString
}

// Auth authenticate a user and start a session
func Auth(URL string, Username string, Password string) (retValue AuthResponse) {
	token, _ := randomHex(32)
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	authResp, err := client.PostForm(URL+"sessions", url.Values{"platform": {"ios"}, "email": {Username}, "token": {token}, "password": {Password}})

	if err == nil {
		defer authResp.Body.Close()
		json.NewDecoder(authResp.Body).Decode(&retValue)
	} else {
		retValue.AccessToken = ""
		retValue.CurrentTime = ""
	}

	return
}

// GetDashboard retrieves infos about an account
func GetDashboard(URL string, Auth AuthResponse) (retValue DashResponse) {
	dashReq, _ := http.NewRequest("GET", URL+"dashboard", nil)
	dashReq.Header.Add("Authorization", "Token token="+Auth.AccessToken)

	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	dashResp, err := client.Do(dashReq)

	if err == nil {
		defer dashResp.Body.Close()
		json.NewDecoder(dashResp.Body).Decode(&retValue)
	}

	return
}

// GetRobotState returns the state of an robot
func GetRobotState(Auth AuthResponse, rob Robot) (retValue RobotState) {
	commandData, _ := json.Marshal(robotCommand{ID: "1", Command: "getRobotState"})

	utcDate, CommandURL, authString := makeAuth(rob, commandData)

	CommandReq, _ := http.NewRequest("POST", CommandURL, bytes.NewBuffer(commandData))
	CommandReq.Header.Add("Date", utcDate)
	CommandReq.Header.Add("Authorization", authString)
	CommandReq.Header.Add("Accept", "application/vnd.neato.nucleo.v1")

	client := &http.Client{}
	CommandResp, err := client.Do(CommandReq)

	if err == nil {
		defer CommandResp.Body.Close()
		json.NewDecoder(CommandResp.Body).Decode(&retValue)
	}

	return
}
