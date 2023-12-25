package main

import (
	"bufio"
	"fmt"
	"github.com/aaronarduino/goqrsvg"
	"github.com/ajstarks/svgo"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image/png"
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

	output := promptUser("Enter output file name (default: " + ssid + "): ")
	if output == "" {
		output = ssid
	}

	format := promptUser("Enter output format (png, svg, or both, default: png): ")
	if format == "" {
		format = "png"
	} else {
		format = strings.ToLower(format)
	}
	validFormats := []string{"png", "svg", "both"}
	if !slices.Contains(validFormats, format) {
		fmt.Println("Invalid output format:", format)
		os.Exit(1)
	}

	var formats []string
	switch format {
	case "png":
		formats = []string{"png"}
	case "svg":
		formats = []string{"svg"}
	case "both":
		formats = []string{"png", "svg"}
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
	qrCode, err := qr.Encode(qrCodeContent, qr.M, qr.Auto)
	if err != nil {
		fmt.Println("Failed to generate QR code:", err)
		os.Exit(1)
	}

	// Save QR code
	// Save QR code in specified formats
	if slices.Contains(formats, "png") {
		outputFile := output + ".png"
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
		qrCode, err = barcode.Scale(qrCode, 256, 256)
		if err != nil {
			fmt.Println("Failed to scale QR code:", err)
			os.Exit(1)
		}
		err = png.Encode(file, qrCode)
		if err != nil {
			fmt.Println("Failed to save QR code:", err)
			os.Exit(1)
		}
	}
	if slices.Contains(formats, "svg") {
		outputFile := output + ".svg"
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
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(canvas)
		err = qs.WriteQrSVG(canvas)
		if err != nil {
			fmt.Println("Failed to save QR code:", err)
			os.Exit(1)
		}
		canvas.End()
	}

	if err != nil {
		fmt.Println("Failed to save QR code:", err)
		os.Exit(1)
	}

	fmt.Println("QR code saved to", output)
}
