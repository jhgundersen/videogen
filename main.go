package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	apiBase      = "https://api.defapi.org"
	pollInterval = 5 * time.Second
)

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type taskData struct {
	TaskID       string `json:"task_id"`
	Status       string `json:"status"`
	Result       struct {
		Video string `json:"video"`
	} `json:"result"`
	StatusReason struct {
		Message *string `json:"message"`
	} `json:"status_reason"`
}

func apiKey() string {
	key := os.Getenv("DEFAPI_API_KEY")
	if key == "" {
		fmt.Fprintln(os.Stderr, "error: DEFAPI_API_KEY environment variable not set")
		os.Exit(1)
	}
	return key
}

func post(endpoint string, body map[string]any, key string) json.RawMessage {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", apiBase+endpoint, strings.NewReader(string(b)))
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return readResponse(resp)
}

func get(endpoint string, key string) json.RawMessage {
	req, _ := http.NewRequest("GET", apiBase+endpoint, nil)
	req.Header.Set("Authorization", "Bearer "+key)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "request error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	return readResponse(resp)
}

func readResponse(resp *http.Response) json.RawMessage {
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		fmt.Fprintf(os.Stderr, "HTTP %d: %s\n", resp.StatusCode, string(raw))
		os.Exit(1)
	}
	var ar apiResponse
	if err := json.Unmarshal(raw, &ar); err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}
	if ar.Code != 0 {
		fmt.Fprintf(os.Stderr, "API error %d: %s\n", ar.Code, ar.Message)
		os.Exit(1)
	}
	return ar.Data
}

func extractTaskID(data json.RawMessage) string {
	var d struct {
		TaskID string `json:"task_id"`
	}
	json.Unmarshal(data, &d)
	return d.TaskID
}

func poll(taskID, key string) string {
	fmt.Printf("Task submitted: %s\nPolling", taskID)
	for {
		time.Sleep(pollInterval)
		data := get("/api/task/query?task_id="+taskID, key)
		var td taskData
		json.Unmarshal(data, &td)

		switch td.Status {
		case "success":
			fmt.Println(" done.")
			if td.Result.Video == "" {
				fmt.Fprintln(os.Stderr, "error: no video URL in response")
				os.Exit(1)
			}
			return td.Result.Video
		case "failed":
			msg := "unknown reason"
			if td.StatusReason.Message != nil {
				msg = *td.StatusReason.Message
			}
			fmt.Fprintf(os.Stderr, "\ngeneration failed: %s\n", msg)
			os.Exit(1)
		default:
			fmt.Print(".")
		}
	}
}

func download(videoURL, taskID string) string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "Downloads")
	os.MkdirAll(dir, 0755)
	dest := filepath.Join(dir, "videogen_"+taskID+".mp4")

	fmt.Println("Downloading...")
	resp, err := http.Get(videoURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "download error: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	io.Copy(f, resp.Body)
	return dest
}

// --- model subcommands ---

func cmdSeedance(args []string) {
	fs := flag.NewFlagSet("seedance", flag.ExitOnError)
	duration := fs.Int("duration", 5, "Duration in seconds: 5, 10, 15")
	ratio := fs.String("ratio", "16:9", "Aspect ratio: 16:9 9:16 1:1 4:3 3:4 21:9")
	image := fs.String("image", "", "Reference image URL for image-to-video (up to 9 supported)")
	fs.Usage = func() {
		fmt.Println("Usage: videogen seedance [flags] <prompt>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}
	validDuration := map[int]bool{5: true, 10: true, 15: true}
	if !validDuration[*duration] {
		fmt.Fprintln(os.Stderr, "error: --duration must be 5, 10, or 15")
		os.Exit(1)
	}

	key := apiKey()
	fmt.Printf("Model: seedance | Duration: %ds | Ratio: %s\nPrompt: %s\n\n", *duration, *ratio, prompt)

	content := []map[string]any{{"type": "text", "text": prompt}}
	if *image != "" {
		content = append(content, map[string]any{
			"type":      "image_url",
			"image_url": map[string]any{"url": *image},
		})
	}

	data := post("/api/video/seedance/gen", map[string]any{
		"model":    "seedance-2.0",
		"content":  content,
		"duration": *duration,
		"ratio":    *ratio,
	}, key)

	taskID := extractTaskID(data)
	videoURL := poll(taskID, key)
	dest := download(videoURL, taskID)
	fmt.Printf("\nSaved to: \033]8;;file://%s\033\\%s\033]8;;\033\\\n", dest, dest)
}

func cmdGrok(args []string) {
	fs := flag.NewFlagSet("grok", flag.ExitOnError)
	duration := fs.Int("duration", 10, "Duration in seconds: 10, 15")
	ratio := fs.String("ratio", "16:9", "Aspect ratio: 16:9 9:16 1:1 2:3 3:2")
	image := fs.String("image", "", "Reference image URL for image-to-video")
	fs.Usage = func() {
		fmt.Println("Usage: videogen grok [flags] <prompt>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}
	validDuration := map[int]bool{10: true, 15: true}
	if !validDuration[*duration] {
		fmt.Fprintln(os.Stderr, "error: --duration must be 10 or 15")
		os.Exit(1)
	}

	key := apiKey()
	fmt.Printf("Model: grok | Duration: %ds | Ratio: %s\nPrompt: %s\n\n", *duration, *ratio, prompt)

	body := map[string]any{
		"prompt":       prompt,
		"model":        "grok-imagine-video",
		"duration":     fmt.Sprintf("%d", *duration),
		"aspect_ratio": *ratio,
	}
	if *image != "" {
		body["images"] = []string{*image}
	}

	data := post("/api/grok-imagine-video/gen", body, key)

	taskID := extractTaskID(data)
	videoURL := poll(taskID, key)
	dest := download(videoURL, taskID)
	fmt.Printf("\nSaved to: \033]8;;file://%s\033\\%s\033]8;;\033\\\n", dest, dest)
}

func cmdSora(args []string) {
	fs := flag.NewFlagSet("sora", flag.ExitOnError)
	duration := fs.Int("duration", 10, "Duration in seconds: 10, 15, 25 (25 requires --variant sora-2-pro)")
	ratio := fs.String("ratio", "16:9", "Aspect ratio: 16:9 9:16")
	variant := fs.String("variant", "sora-2", "Model variant: sora-2, sora-2-hd, sora-2-pro")
	image := fs.String("image", "", "Reference image URL for image-to-video")
	fs.Usage = func() {
		fmt.Println("Usage: videogen sora [flags] <prompt>")
		fs.PrintDefaults()
	}
	fs.Parse(args)

	prompt := strings.Join(fs.Args(), " ")
	if prompt == "" {
		fs.Usage()
		os.Exit(1)
	}
	validDuration := map[int]bool{10: true, 15: true, 25: true}
	if !validDuration[*duration] {
		fmt.Fprintln(os.Stderr, "error: --duration must be 10, 15, or 25")
		os.Exit(1)
	}
	if *duration == 25 && *variant != "sora-2-pro" {
		fmt.Fprintln(os.Stderr, "warning: 25s duration requires sora-2-pro, switching variant")
		*variant = "sora-2-pro"
	}

	key := apiKey()
	fmt.Printf("Model: sora (%s) | Duration: %ds | Ratio: %s\nPrompt: %s\n\n", *variant, *duration, *ratio, prompt)

	body := map[string]any{
		"prompt":       prompt,
		"model":        *variant,
		"duration":     fmt.Sprintf("%d", *duration),
		"aspect_ratio": *ratio,
	}
	if *image != "" {
		body["images"] = []string{*image}
	}

	data := post("/api/sora2/gen", body, key)

	taskID := extractTaskID(data)
	videoURL := poll(taskID, key)
	dest := download(videoURL, taskID)
	fmt.Printf("\nSaved to: \033]8;;file://%s\033\\%s\033]8;;\033\\\n", dest, dest)
}

func usage() {
	fmt.Println(`Usage: videogen <model> [flags] <prompt>

Models:
  seedance   ByteDance Seedance 2.0
  grok       xAI Grok Imagine Video
  sora       OpenAI Sora 2 Stable

Run 'videogen <model> --help' for model-specific flags.`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "seedance":
		cmdSeedance(os.Args[2:])
	case "grok":
		cmdGrok(os.Args[2:])
	case "sora":
		cmdSora(os.Args[2:])
	case "-h", "--help", "help":
		usage()
	default:
		fmt.Fprintf(os.Stderr, "unknown model: %s\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}
}
