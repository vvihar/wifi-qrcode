package main

import (
	"bufio"
	"fmt"
	"github.com/ajstarks/svgo"
	"github.com/skip2/go-qrcode"
	"os"
	"slices"
	"strings"
)

// promptUser prompts the user for input and returns the input
func promptUser(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func main() {
	// Prompt user for input
	ssid := promptUser("Enter SSID of the WiFi network: ")

	password := promptUser("Enter password of the WiFi network: ")

	output := promptUser("Enter output file name (default: output): ")
	if output == "" {
		output = "output"
	}

	format := promptUser("Enter output format (png or svg, default: png): ")
	if format == "" {
		format = "png"
	} else {
		format = strings.ToLower(format)
	}
	validFormats := []string{"png", "svg"}
	if !slices.Contains(validFormats, format) {
		fmt.Println("Invalid output format:", format)
		os.Exit(1)
	}

	auth := promptUser("Enter auth type (WPA or WEP, default: WPA): ")
	if auth == "" {
		auth = "WPA"
	} else {
		auth = strings.ToUpper(auth)
	}
	validAuthTypes := []string{"WPA", "WEP"}
	if !slices.Contains(validAuthTypes, auth) {
		fmt.Println("Invalid auth type:", auth)
		os.Exit(1)
	}

	// Generate WiFi QR code
	qrCodeContent := fmt.Sprintf("WIFI:S:%s;T:%s;P:%s;;", ssid, auth, password)
	qrCode, err := qrcode.New(qrCodeContent, qrcode.Medium)
	if err != nil {
		fmt.Println("Failed to generate QR code:", err)
		os.Exit(1)
	}

	// Save QR code
	outputFile := fmt.Sprintf("%s.%s", output, format)
	switch strings.ToLower(format) {
	case "png":
		err = qrCode.WriteFile(256, outputFile)
	case "svg":
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Println("Failed to create output file:", err)
			os.Exit(1)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				fmt.Println("Failed to close output file:", err)
				os.Exit(1)
			}
		}(file)

		canvas := svg.New(file)
		canvas.Start(256, 256)
		canvas.Path(qrCode.ToSmallString(false), "fill:black;stroke:none")
		canvas.End()
	default:
		fmt.Println("Invalid output format:", format)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Failed to save QR code:", err)
		os.Exit(1)
	}

	fmt.Println("QR code saved to", outputFile)
}
