package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("ERROR LOADING ENV VARS: " + err.Error())
		os.Exit(1)
	}

	sshCmdOut := bytes.Buffer{}

	sshCmd := exec.Command("ls", os.ExpandEnv("$HOME/.ssh"))
	sshCmd.Stdout = &sshCmdOut
	sshCmd.Stderr = os.Stderr

	err = sshCmd.Run()
	if err != nil {
		fmt.Println("ERROR CHECKING SSH DIRECTORY: " + err.Error())
		os.Exit(1)
	}

	sshDir := sshCmdOut.String()

	fmt.Println(sshDir)

	if !strings.Contains(sshDir, "bedrock_rsa") {
		keygenCmd := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", os.ExpandEnv("$HOME/.ssh/bedrock_rsa"), "-N", "\"\"")
		keygenCmd.Stderr = os.Stderr
		err := keygenCmd.Run()
		if err != nil {
			fmt.Println("ERROR CREATING SSH KEY: " + err.Error())
			os.Exit(1)
		}

	}

	file, err := os.Open("./ansible/hosts")
	if err != nil {
		fmt.Println("ERROR OPENING HOSTS FILE: " + err.Error())
		os.Exit(1)
	}

	hostsData, err := io.ReadAll(file)

	hostsStr := strings.Split(string(hostsData), "\n")

	ipAddresses := []string{}

	for _, str := range hostsStr {
		if strings.Contains(str, ".") {
			ipAddresses = append(ipAddresses, str)
			copyPubKeyCmd := exec.Command("sshpass", "-p", os.Getenv("HOST_PASS"), "ssh-copy-id", "-i", os.ExpandEnv("$HOME/.ssh/bedrock_rsa"), fmt.Sprintf("%v@%v", os.Getenv("HOST_USER"), str))
			fmt.Println(copyPubKeyCmd.String())
			copyPubKeyCmd.Stderr = os.Stderr
			copyPubKeyCmd.Run()
		}
	}

	// controllerIp := ipAddresses[0]

	// ipAddresses = ipAddresses[1:]

	// token := uuid.New().String()

	// setupMainCmd := exec.Command("ansible", "-i", "./ansible/hosts", "./ansible/main-node-playbook.yaml")

}
