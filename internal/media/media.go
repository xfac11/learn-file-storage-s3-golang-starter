package media

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
)

type ffproberesult struct {
	Streams []ffprobestream `json:"streams"`
}
type ffprobestream struct {
	Index  int `json:"index"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type AspectRatio struct {
	WidthRatio  int
	HeightRatio int
}

func NewAspectRatio(width, height int) AspectRatio {
	gcd := greatestCommonDivider(width, height)
	return AspectRatio{
		WidthRatio:  width / gcd,
		HeightRatio: height / gcd,
	}
}

func (ratio *AspectRatio) AsString() string {
	return fmt.Sprintf("%d:%d", ratio.WidthRatio, ratio.HeightRatio)
}

// Makes the actual ratios inside a AspectRatio into common ratios.
// If not close to any supported ratios the ratios will stay the same.
//
// Support for 16:9 and 9:16.
func (ratio *AspectRatio) MakeCommonRatio() {
	var decimal float32 = float32(ratio.WidthRatio) / float32(ratio.HeightRatio)
	var close169 float32 = (16.0 / 9.0) - decimal
	var close916 float32 = (9.0 / 16.0) - decimal

	if math.Abs(float64(close169)) <= 0.1 {
		ratio.WidthRatio = 16
		ratio.HeightRatio = 9
	} else if math.Abs(float64(close916)) <= 0.1 {
		ratio.WidthRatio = 9
		ratio.HeightRatio = 16
	}
}

func greatestCommonDivider(a, b int) int {
	if b == 0 {
		return a
	}
	return greatestCommonDivider(b, a%b)
}

// Returns the video aspect ratio in 16:9, 9:16 or other
func GetVideoAspectRatio(filePath string) (string, error) {
	cmdName := "ffprobe"
	args := []string{"-v", "error", "-print_format", "json", "-show_streams", filePath}
	cmd := exec.Command(cmdName, args...)
	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Could not run the command: %s, with arguments: %s.\n Error: %s", cmdName, args, err)
	}

	var result ffproberesult
	json.Unmarshal(buffer.Bytes(), &result)

	ratio := NewAspectRatio(result.Streams[0].Width, result.Streams[0].Height)
	ratio.MakeCommonRatio()
	if ratio.AsString() == "16:9" || ratio.AsString() == "9:16" {
		return ratio.AsString(), nil
	}
	return "other", nil
}

func ProcessVideoForFastStart(filePath string) (string, error) {
	fastStartPath := filePath + ".processing"
	cmdName := "ffmpeg"
	args := []string{"-i", filePath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", fastStartPath}
	cmd := exec.Command(cmdName, args...)
	buffer := bytes.Buffer{}
	cmd.Stdout = &buffer
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Could not run the command: %s, with arguments: %s.\n Error: %s", cmdName, args, err)
	}
	return fastStartPath, nil
}
