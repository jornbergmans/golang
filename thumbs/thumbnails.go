package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func getfps() (float64, error) {
	cmdString, err := probe("r_frame_rate")
	if err != nil {
		return 0, err
	}

	stringParts := strings.Split(cmdString, "/")

	if len(stringParts) == 1 {
		strconv.ParseFloat(stringParts[0], 64)
	} else if len(stringParts) == 2 {
		part1, part2 := stringParts[0], stringParts[1]
		if part1 != "" && part2 != "" {
			part1Float, err := strconv.ParseFloat(part1, 64)
			if err != nil {
				return 0, err
			}
			part2Float, err := strconv.ParseFloat(part2, 64)
			if err != nil {
				return 0, err
			}
			return part1Float / part2Float, nil
		} else if part1 != "" && part2 == "" {
			part1Float, err := strconv.ParseFloat(part1, 64)
			if err != nil {
				return 0, err
			}
			return part1Float, nil
		} else if part1 == "" && part2 != "" {
			part2Float, err := strconv.ParseFloat(part2, 64)
			if err != nil {
				return 0, err
			}
			return part2Float, nil
		}
	}

	return 0, fmt.Errorf("Unknown framerate format:", cmdString)
}

func getframes() (uint64, error) {

	cmdString, err := probe("nb_frames")
	if err != nil {
		return 0, err
	}
	cmdInt, err := strconv.ParseUint(cmdString, 10, 64)
	if err != nil {
		return 0, err
	}

	return cmdInt, err
}

func probe(streamToExtract string) (string, error) {
	entryToExtract := fmt.Sprintf("stream=%s", streamToExtract)
	cmdName := "ffprobe"
	cmdArguments := []string{
		"-hide_banner", "-loglevel", "panic",
		"-select_streams", "v:0", "-show_entries", entryToExtract,
		"-of", "default=noprint_wrappers=1:nokey=1", "-i", "/Users/vtr/Downloads/Coco_ManCream_ONLINE_1859.mov",
	}

	cmd := exec.Command(cmdName, cmdArguments...)

	cmdBytes, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	cmdString := strings.TrimSpace(string(cmdBytes))

	return cmdString, err
}

func thumbgen() {

	fps, _ := getfps()
	frames, _ := getframes()

	outframes := float64(frames) / fps

	infile := bufio.NewReader(os.Stdin)

	cmdName := "/usr/bin/env"
	// cmdArgument1 := []string{
	// 	"ffmpeg",
	// 	"-hide_banner", "-loglevel", "panic",
	// 	"-y", "-i",
	// }
	cmdInput := fmt.Sprintln(infile)
	// cmdArgument2 := []string{
	// 	"-frames", "1", "-q", "5",
	// 	"-vf", "\"mpdecimate,select=not(mod(n\\,2)),scale='480:-1',tile=",
	// }
	cmdFrameCount := fmt.Sprintln(outframes)
	// cmdArgument3 := []string{
	// 	"x1",
	// 	"____OUTPUT____.png",
	// }

	cmdArguments := append([]string{
		"ffmpeg",
		"-hide_banner", "-loglevel", "panic",
		"-y", "-i"},
		cmdInput,
		"-frames", "1", "-q", "5",
		"-vf", "\"mpdecimate,select=not(mod(n\\,2)),scale='480:-1',tile=",
		cmdFrameCount,
		"x1",
		"____OUTPUT____.png")

	// cmdArguments := []string{"ffmpeg", "-version"}
	cmd := exec.Command(cmdName, cmdArguments...)

	cmdBytes, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error:", err.Error())
	}
	fmt.Printf("command output: %q\n", cmdBytes)
}

func main() {
	fps, _ := getfps()
	fmt.Println(fps)
	frames, _ := getframes()
	fmt.Println(frames)
}
