package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	assets "hola/Blacklisted"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/corpix/uarand"
)

type Foreground struct {
	Red        string
	Green      string
	Yellow     string
	Cyan       string
	LightWhite string
}

var Fore = Foreground{
	Red:        "\033[31m",
	Green:      "\033[32m",
	Yellow:     "\033[33m",
	Cyan:       "\033[36m",
	LightWhite: "\033[97m",
}

type Styling struct {
	Bold      string
	Italic    string
	Blink     string
	Reset_all string
}

var Style = Styling{
	Bold:      "\033[1m",
	Italic:    "\033[3m",
	Blink:     "\033[5m",
	Reset_all: "\033[0m",
}

const (
	url               = "https://login.microsoftonline.com/common/GetCredentialType?mkt=en-US"
	numWorkersForFile = 250
	ValidEmailRegex   = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

func ExtractFromCombo(email string) string {
	isValidEmail, _ := regexp.MatchString(ValidEmailRegex, email)
	if !isValidEmail {
		return ""
	}
	domainParts := strings.Split(email, "@")
	if len(domainParts) > 1 {
		return strings.ToLower(domainParts[1])
	}
	return ""
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
func progressBar(duration time.Duration) {
	for i := 0; i <= 80; i++ {
		bar := ""
		for j := 0; j < i; j++ {
			bar += fmt.Sprintf("%v█", Style.Bold)
		}
		for j := i; j < 80; j++ {
			bar += " "
		}
		fmt.Printf("\r              |%s|", bar)
		time.Sleep(duration / 80)
	}
	fmt.Println()
}

func fileWorker(id int, jobs <-chan string, results chan<- string, blacklistCount *int, wg *sync.WaitGroup) {
	defer wg.Done()

	for line := range jobs {
		domain := ExtractFromCombo(line)
		if contains(assets.BlackListed, domain) {
			*blacklistCount++
			continue
		}
		results <- line
	}
}

func ProcessFileWithWorkers(filename string) (int, int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()
	var totalLines int
	var blacklistCount int
	jobs := make(chan string, numWorkersForFile)
	results := make(chan string, numWorkersForFile)
	var wg sync.WaitGroup
	for i := 0; i < numWorkersForFile; i++ {
		wg.Add(1)
		go fileWorker(i, jobs, results, &blacklistCount, &wg)
	}
	go func() {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			jobs <- line
			totalLines++
		}
		close(jobs)
		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			close(results)
			return
		}
	}()
	go func() {
		wg.Wait()
		close(results)
	}()
	var filteredLines []string
	for result := range results {
		filteredLines = append(filteredLines, result)
	}
	err = os.WriteFile(filename, []byte(strings.Join(filteredLines, "\n")), 0644)
	if err != nil {
		return 0, 0, err
	}

	return blacklistCount, len(filteredLines), nil
}

type Payload struct {
	IsOtherIdpSupported            bool   `json:"isOtherIdpSupported"`
	CheckPhones                    bool   `json:"checkPhones"`
	IsRemoteNGCSupported           bool   `json:"isRemoteNGCSupported"`
	IsCookieBannerShown            bool   `json:"isCookieBannerShown"`
	IsFidoSupported                bool   `json:"isFidoSupported"`
	Country                        string `json:"country"`
	Forceotclogin                  bool   `json:"forceotclogin"`
	IsExternalFederationDisallowed bool   `json:"isExternalFederationDisallowed"`
	IsRemoteConnectSupported       bool   `json:"isRemoteConnectSupported"`
	FederationFlags                int    `json:"federationFlags"`
	IsSignup                       bool   `json:"isSignup"`
	FlowToken                      string `json:"flowToken"`
	IsAccessPassSupported          bool   `json:"isAccessPassSupported"`
	Username                       string `json:"username,omitempty"`
}

type Response struct {
	IfExistsResult int `json:"IfExistsResult"`
}

var (
	client    = &http.Client{}
	wg        sync.WaitGroup
	semaphore = make(chan struct{}, 20)
)

func processEmail(email string, resultFile *os.File) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()

	payload := Payload{
		IsOtherIdpSupported:            true,
		CheckPhones:                    false,
		IsRemoteNGCSupported:           true,
		IsCookieBannerShown:            false,
		IsFidoSupported:                true,
		Country:                        "NG",
		Forceotclogin:                  false,
		IsExternalFederationDisallowed: false,
		IsRemoteConnectSupported:       false,
		FederationFlags:                0,
		IsSignup:                       false,
		FlowToken:                      "AQABAAEAAAAtyolDObpQQ5VtlI4uGjEPmVTvB5eZTaL_xvRdNX8zoF_M9oCPfpR1_3-Wz9ETrDbl5ca9Avq0LYJkoyoMgY5LIhrw_zFYKZPKDynsKoHPZfgeKmWiIAs1DXbLOrj1FwddvGzTm1ABWqIWhpNkryjIGv9-pljgbUnhPWj9pTn9ffvUpp8V2MtaAhrj46pyDne0WQmgpo5yxrOcie6NRDmX-vIRN1MIuXjLJ27VP51D0OM2hEp1OD47P6dtU6fk3-n2oCqUh1nP1tJCP1Pr47Uw2d3Hx3uCPYHHQJ8S3DkYwNqi4ieYGWQoRIaGrswGKuHiQRsyIuf8jtXEVXyOGqJhVIrV13orhsMe8QFdNAQE95yTcr7oSV6cXL7EWJdelszMsPUosCWSNdpwVI3lFGrKkYHetSaT2PrQem5AJIKBpKpvdLzk4q_P1P5_HA5hrOLCjH451cW4GJ2aVLL8wejgiEdppAzICHiHJOAthyGUP1R7w0z62wD6Ml9QOrRuqGS1KRxOCycJSUhLQcXDX5yIL1PCokaNJIAca5y1wkJb4zMbwhsGoVaUnWZK8XjTWYovsLqEn1dvUW_GrQxdQzwyIAA",
		IsAccessPassSupported:          true,
		Username:                       email,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("[ERROR] creating request:", err)
		return
	}

	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "close")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", uarand.GetRandom())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[ERROR] sending request:", err)
		time.Sleep(15 * time.Second)
		go processEmail(email, resultFile)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var response Response
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("[ERROR] decoding response:", err)
			return
		}

		result := fmt.Sprintf("%s", email)
		var Res_str string
		if len(result) > 50 {
			Res_str = result[:47] + "..."
		} else {
			Res_str = result
		}
		if response.IfExistsResult == 5 || response.IfExistsResult == 0 {
			fmt.Printf("%v[+]%v  |%v   VALID   %v| %-30s | %vCRUX-CORE%v \n", Style.Bold+Fore.Green, Style.Reset_all+Style.Bold, Style.Bold+Fore.Green, Style.Reset_all+Style.Bold, Res_str, Style.Bold+Style.Italic+Fore.Cyan, Style.Reset_all)
			if _, err := resultFile.WriteString(result + "\n"); err != nil {
				fmt.Println("[ERROR] writing to result file:", err)
			}
		} else {
			fmt.Printf("%v[-]%v  |%v  INVALID  %v| %-30s | %vCRUX-CORE%v \n", Style.Bold+Fore.Red, Style.Reset_all+Style.Bold, Fore.Red, Style.Reset_all+Style.Bold, Res_str, Style.Bold+Style.Italic+Fore.Cyan, Style.Reset_all)
		}
	} else {
		fmt.Printf("API request failed for email %s with status code: %d\n", email, resp.StatusCode)
		time.Sleep(15 * time.Second)
		go processEmail(email, resultFile)
	}
}
func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
func logo() {
	brightRed := Fore.Red + Style.Blink
	fmt.Printf("                       %v ▄▄░%v  ▄▌    %v ▄▒░%v  ▄▌\n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold, brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                   ▀▄▄ %v▀▀%v▄▄░░█ ▀▄▄ %v▀▀%v▄▄░░█\n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold, brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                     ▓▒░░█▀▀▀█▌  ▓▒░░▀▀▀▀█▌         █▀▀█  █▀▀█  █  █ ▀▄ ▄▀   █▀▀█  █▀▀▀█  █▀▀█  █▀▀▀ \n")
	fmt.Printf("                  %v▄▀%v▐▀▀  ▄▄   ▀ ▐▀▀ ▄▄▄   ▀         █     █▄▄▀  █  █   █     █     █   █  █▄▄▀  █▀▀▀ \n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                  %v▀▄%v    ▐▓▒█▄ %v▄▄  ▀%v▐▒░▀▀  ▄         █▄▄█  █  █  ▀▄▄▀ ▄▀ ▀▄   █▄▄█  █▄▄▄█  █  █  █▄▄▄ \n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold, brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                  	  █▀▀▀░▄ %v▀▀%v  ▓▒▄▄▄▄▄▌%v▄%v\n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold, brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                  	 ▀ %v▀▒%v  ▀▄   ▐▀░ ▀▀░█▀%v ▒%v     O  F  F  I  C  E  -  V  A  L  I  D  A  T  O  R\n", brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold, brightRed, Style.Reset_all+Fore.LightWhite+Style.Bold)
	fmt.Printf("                  	   %v░%v\n\n\n", brightRed, Style.Reset_all)

}

func main() {
	clearConsole()
	logo()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("                   %v[ENTER]%v Path of the email file: ", Style.Bold+Fore.Red, Style.Reset_all+Style.Bold)
	filePath, _ := reader.ReadString('\n')
	filePath = strings.TrimSpace(filePath)

	fmt.Printf("                   %v[ENTER]%v Path of the output file: ", Style.Bold+Fore.Red, Style.Reset_all+Style.Bold)
	resultFilePath, _ := reader.ReadString('\n')
	resultFilePath = strings.TrimSpace(resultFilePath)

	emailFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening email file:", err)
		return
	}
	defer emailFile.Close()

	resultFile, err := os.Create(resultFilePath)
	if err != nil {
		fmt.Println("Error creating result file:", err)
		return
	}
	defer resultFile.Close()

	fmt.Printf("\n\n")
	fmt.Printf("%v                 ⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯\n\n", Style.Bold+Fore.LightWhite)
	fmt.Printf("                               %vP  R  O  C  E  S  S  I  N  G  -  L  I  N  E  S%v\n\n", Fore.Yellow, Style.Reset_all)
	Blacklist, lines, err := ProcessFileWithWorkers(filePath)
	if err != nil {
		fmt.Printf(" %v[ERROR]%v UNABLE TO PROCESS LINEs", Fore.Red, Style.Reset_all)
	} else {
		fmt.Printf("                                  %v[BLACKLISTED LINES] :%v %v\n", Fore.Cyan+Style.Italic, Style.Reset_all, Blacklist)
		fmt.Printf("                                  %v[GOOD  LINES]       :%v %v\n\n", Fore.Cyan+Style.Italic, Style.Reset_all, lines)
	}
	fmt.Printf("                               %vS  T  A  R  T  I  N  G  -  C  H  E  C  K  E  R%v\n\n", Fore.Yellow, Style.Reset_all)
	fmt.Printf("%v                 ⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯⎯\n", Style.Bold+Fore.LightWhite)
	progressBar(4 * time.Second)
	fmt.Println("")
	fmt.Println("")

	scanner := bufio.NewScanner(emailFile)
	for scanner.Scan() {
		line := scanner.Text()

		email := line
		wg.Add(1)
		go processEmail(email, resultFile)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading email file:", err)
	}

	wg.Wait()
	fmt.Println("Processing completed.")
}
