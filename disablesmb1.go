// SMB1Disabler.exe
// Disabled SMB1 on Windows Machines
// Author: Ron Egli - ron.egli@tvrms.com
// Version 0.3

package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "time"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func main(){
	checkIfAdmin()
	fmt.Println("SMBv1 Status: Loading....")
	out, err := exec.Command("PowerShell", "-Command", "Get-WindowsOptionalFeature -Online -FeatureName smb1protocol | Select State").Output()
    if err != nil {
        log.Fatal(err)
    } 
    if strings.Contains(string(out), "DisablePending") {
    	fmt.Println("SMBv1 Status: Disable Pending")
    	fmt.Println("Please reboot to complete the patching process.")
    }
    if strings.Contains(string(out), "Disabled") {
    	fmt.Println("SMBv1 Status: Disabled")
    	fmt.Println("All Clear!")
    }
    if strings.Contains(string(out), "Enabled") || strings.Contains(string(out), "EnablePending") {
    	fmt.Println("SMBv1 Status: Enabled - System Vulnerable")
    	fmt.Println("Running Fix")
    	runFix()
    	fmt.Println("Please reboot to complete the patching process.")
    }
	time.Sleep(5 * time.Second)
}

func runFix(){
	out, err := exec.Command("PowerShell", "-Command", "Disable-WindowsOptionalFeature -Online -FeatureName smb1protocol -n").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(out)
}

func checkIfAdmin(){
	Block{
		Try: func() {
			fo, err := os.Create("c:\\test.txt")
			if err != nil {
				Throw("Needs to be run with Administrative Priviledges")
			}
			defer fo.Close()
		},
		Catch: func(e Exception) {
			fmt.Printf("%v\n", e)
			os.Exit(1)
		},
		Finally: func() {
		},
	}.Do()
}

func (tcf Block) Do() {
	if tcf.Finally != nil {
		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()
}

func in_array(val string, array []string) (exists bool, index int) {
    exists = false
    index = -1;

    for i, v := range array {
        if val == v {
            index = i
            exists = true
            return
        }   
    }

    return
}
