package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	test_dir = "test"
)

func main() {

	fmt.Printf("\n\033[33m[ STARTED ]\n")

	angle, fold_param := calculate_fold_parameters(4, 3, 13, 9)
	fmt.Printf("ang == %d param == %d", angle, fold_param)

}

type Ray_coords struct {
	y int
	x int
}

func calculate_fold_parameters(height, width, ray_from, ray_to int) (ray_direction int, fold_start_shift int) {

	//var return_angle int
	//var return_fold_param int

	total_bounding_numbers := height*2 + width*2
	var ray_positions = make([]Ray_coords, total_bounding_numbers)

	for i := range width {
		ray_positions[i] = Ray_coords{0, i}
	}

	for i := range height {
		ray_positions[width+i] = Ray_coords{i, width}
	}

	for i := range width {
		ray_positions[width+height+i] = Ray_coords{height, width - i}
	}

	for i := range height {
		ray_positions[width*2+height+i] = Ray_coords{height - i, 0}
	}

	//----------------------------------------------------------------------------
	// DEBUG PRINT
	//----------------------------------------------------------------------------
	/*for i := range total_bounding_numbers {

		if ray_positions[i].x != 0 || ray_positions[i].y != 0 {
			fmt.Printf("\033[35m")
		} else {
			fmt.Printf("\033[30m")
		}

		fmt.Printf("ray_pos %d == (%d ; %d)\n",
			i+1,
			ray_positions[i].y,
			ray_positions[i].x,
		)
	} */
	//----------------------------------------------------------------------------

	// (y1,x1) -----> (y2,x2)
	y1 := ray_positions[ray_from-1].y
	x1 := ray_positions[ray_from-1].x

	y2 := ray_positions[ray_to-1].y
	x2 := ray_positions[ray_to-1].x

	/*
		//check for H-ray
		if y1 == y2 {
			if x1 < x2 {
				ray_direction = 90
				fold_start_shift = y1
			} else {
				ray_direction = 270
				fold_start_shift = y1
			}
		} else

		//check for V-ray
		if x1 == x2 {
			if y1 < y2 {
				ray_direction = 180
				fold_start_shift = x1
			} else {
				ray_direction = 360
				fold_start_shift = x1
			}
		} else

		//check for correct diagonal ray
		if math.Abs(float64(x1)-float64(x2)) == math.Abs(float64(y1)-float64(y2)) {
			ray_direction = 77
		}
	*/

	ray_direction = int(math.Atan2(float64(y1-y2), float64(x1-x2))*180/3.14) - 90
	if ray_direction <= 0 {
		ray_direction += 360
	}

	return ray_direction, fold_start_shift
	//return return_angle, return_fold_param
}

func main2() {

	fmt.Printf("\n\033[33m[ STARTED ]\n")

	test_input_files, err := filepath.Glob(filepath.Join(test_dir, "*"))
	if err != nil {
		fmt.Printf("\033[31mfailed to list input files: %v", err)
	}

	for _, in_file := range test_input_files {

		file_info, err := os.Stat(in_file)
		if err != nil || !file_info.Mode().IsRegular() || strings.Contains(file_info.Name(), ".a") {
			continue
		}

		fmt.Printf("\033[37mtesting file \"%s\"\n", in_file)
		out_file := in_file + ".a"

		input, err := os.ReadFile(in_file)
		if err != nil {
			fmt.Printf("\033[31mfailed to read input file %s: %v\n", in_file, err)
			continue
		}

		expectedOutput, err := os.ReadFile(out_file)
		if err != nil {
			fmt.Printf("\033[31mfailed to read output file %s: %v\n", out_file, err)
			continue
		}

		origStdin := os.Stdin
		origStdout := os.Stdout
		rIn, wIn, _ := os.Pipe()
		wIn.Write(input)
		wIn.Close()
		os.Stdin = rIn
		rOut, wOut, _ := os.Pipe()
		os.Stdout = wOut
		start := time.Now()
		test_func()
		wOut.Close()
		duration := time.Since(start)
		var buf bytes.Buffer
		io.Copy(&buf, rOut)
		os.Stdin = origStdin
		os.Stdout = origStdout

		actualOutput := strings.TrimSpace(buf.String())
		expected := strings.TrimSpace(string(expectedOutput))

		if actualOutput != expected {
			fmt.Printf("\n\033[31mFAILED %s (worked %s)\nExpected:\n%s\nGot:\n%s\n", in_file, duration, expected, actualOutput)
		} else {
			fmt.Printf("\033[32mPASSED %s 	(worked %s)\n", in_file, duration)
		}

	}

	fmt.Printf("\033[33m[ FINISHED ]\n\n")
}

func test_func() {
	inp := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var datasets int
	fmt.Fscanln(inp, &datasets)

	log.Printf("\033[35m TRYING %d DATASETS\n", datasets)
	for range datasets {

		var n, m, k int
		fmt.Fscanln(inp, &n, &m, &k)

		/* var arr [n][m]byte

		for i = range n {
			for j = range m {
				fmt.Fscanf(inp, &arr[i][j])
			}
		} */

	}

	fmt.Fprint(out, "")
}
