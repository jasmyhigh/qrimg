package main

import (
	"flag"
	"fmt"
	"image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
)

func main() {
	size := flag.Int("x", 1024, "qr image pixel size")
	output := flag.String("o", "", "qr image path (optional, defaults to qrimg-{timestamp})")

	flag.Parse()

	var text string
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		text = string(stdin)
	} else if flag.NArg() > 0 {
		text = strings.Join(flag.Args(), " ")
	} else {
		fmt.Println("Please provide text to encode, either via stdin or command-line arguments")
		return
	}

	tmpDir := os.TempDir()
	// Generate filename using timestamp unless output filename is specified
	filename := fmt.Sprintf("qrimg-%d.png", time.Now().UnixMilli())
	if *output != "" {
		filename = *output
	}
	filePath := filepath.Join(tmpDir, filename)

	qr, err := qrcode.New(text, qrcode.Medium)
	if err != nil {
		fmt.Println("Error generating QR code:", err)
		return
	}

	// Save QR code as PNG image file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	err = png.Encode(file, qr.Image(*size))
	if err != nil {
		fmt.Println("Error encoding PNG image:", err)
		return
	}

	fmt.Println("QR code saved to:", filePath)

	// Open the image with the default system application
	err = openFile(filePath)
	if err != nil {
		fmt.Println("Error opening image:", err)
	}
}

// openFile open file with default system application
func openFile(filePath string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("xdg-open", filePath)
	case "windows":
		cmd = exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", filePath)
	case "darwin":
		cmd = exec.Command("open", filePath)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Run()
}
