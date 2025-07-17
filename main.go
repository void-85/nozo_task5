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

	"task5/paper_cell"
)

const (
	test_dir = "test"
)

/*
//-------------------------------------------------------------------------------------------------
                 (0;3)
   1   2   3   4
     x   x   x
  14	       5
     x   x   x
  13	       6
     x   x   x
  12           7
     x   x   x
  11  10   9   8
				 (4;3)

//-------------------------------------------------------------------------------------------------
*/

type Ray_point struct {
	y int
	x int
}

func calculate_fold_parameters(height, width, ray_from, ray_to int) (y1, x1, y2, x2, ray_direction int) {

	total_bounding_numbers := height*2 + width*2
	var ray_positions = make([]Ray_point, total_bounding_numbers)

	for i := range width {
		ray_positions[i] = Ray_point{0, i}
	}
	for i := range height {
		ray_positions[width+i] = Ray_point{i, width}
	}
	for i := range width {
		ray_positions[width+height+i] = Ray_point{height, width - i}
	}
	for i := range height {
		ray_positions[width*2+height+i] = Ray_point{height - i, 0}
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
	y1 = ray_positions[ray_from-1].y
	x1 = ray_positions[ray_from-1].x

	y2 = ray_positions[ray_to-1].y
	x2 = ray_positions[ray_to-1].x

	ray_direction = int(math.Atan2(float64(y1-y2), float64(x1-x2))*180/3.14) - 90
	if ray_direction <= 0 {
		ray_direction += 360
	}

	return y1, x1, y2, x2, ray_direction
	//return return_angle, return_fold_param
}

func main() {

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

	for dataset := range datasets {
		log.Printf("\033[32m----- TRYING %d DATASET -----\n", dataset+1)

		var n, m, k int
		fmt.Fscanln(inp, &n, &m, &k)

		span_top := n
		span_left := m

		log.Printf("\033[32mN==%d M==%d K==%d\n", n, m, k)

		table := make([][]paper_cell.Quarter, n*3)
		for i := range n * 3 {
			table[i] = make([]paper_cell.Quarter, m*3)
		}

		var chars string
		for y := range n {

			fmt.Fscanln(inp, &chars)
			for x := range m {
				if chars[x] == '#' {
					table[n+y][m+x] = 0b00001111

					/* } else if chars[x] == '.' {
					table[n+y][m+x] = 0b00000000*/

				} else {
					table[n+y][m+x] = 0b00000000
				}
			}

		}

		var s string = ""
		for y := range n * 3 {
			for x := range m * 3 {
				if table[y][x] == 0b00001111 {
					s += "# "
				} else if table[y][x] == 0b00000000 {
					s += ". "
				} else {
					s += "_ "
				}
			}
			s += "\n"
		}
		log.Printf("\033[32mtable loaded:\n%s\n", s)

		current_span_left := 0
		current_span_top := 0
		current_m := m
		current_n := n

		for range k {
			var ray_from, ray_to int
			fmt.Fscanln(inp, &ray_from, &ray_to)
			log.Printf("fold vector loaded: %d --> %d", ray_from, ray_to)

			y1, x1, y2, x2, ray_direction := calculate_fold_parameters(current_n, current_m, ray_from, ray_to)
			log.Printf("     (%d;%d) ---> (%d;%d) ang==%d\n", y1, x1, y2, x2, ray_direction)

			for y := range current_n {
				for x := range current_m {
					d := get_D(y1, x1, y2, x2, y, x)
					if d > 0 {
						//table[span_top+y][span_left+x] = 'R'

						new_y, new_x := get_mirror_point_pos(y1, x1, y2, x2, y, x)
						table[span_top+new_y][span_left+new_x] =
							table[span_top+new_y][span_left+new_x] |
								table[span_top+y][span_left+x]

						table[span_top+y][span_left+x] = 0b00000000
					} else if d < 0 {
						//table[span_top+y][span_left+x] = 'L'
					} else {
						//table[span_top+y][span_left+x] = '='
					}
				}
			}

			s = ""
			for y := range n * 3 {
				for x := range m * 3 {
					switch c := table[y][x]; c {
					case 0b00001111:
						s += "# "
					case 0b00000000:
						s += ". "
					case 0b00001000:
						s += "v "
					case 0b00000111:
						s += "v "
					case 0b00000001:
						s += "> "
					case 0b00001110:
						s += "> "
					case 0b00000100:
						s += "< "
					case 0b00001011:
						s += "< "
					case 0b00000010:
						s += "^ "
					case 0b00001101:
						s += "^ "
					case 0b00001100:
						s += "\\ "
					case 0b00000011:
						s += "\\ "
					case 0b00000110:
						s += "/ "
					case 0b00001001:
						s += "/ "
					case 0b00001010:
						s += "x "
					case 0b00000101:
						s += "x "
					default:
						s += "_ "
					}
				}
				s += "\n"
			}
			log.Printf("\033[32mjust folded:\n%s\n", s)

			//centering
			current_span_left = 0
			current_span_top = 0
			current_m = 0
			current_n = 0

			//-----------------------------------------------------
			for x := range m * 3 {
				column_sum := 0
				for y := range n * 3 {
					if table[y][x] == 0b00000000 {
						column_sum += 1
					}
				}
				if column_sum == n*3 {
					current_span_left += 1
				} else {
					break
				}
			}

			for x := m*3 - 1; m > 0; x-- { //current_span_left; x < m*3; x++ {
				column_sum := 0
				for y := range n * 3 {
					if table[y][x] == 0b00000000 {
						column_sum += 1
					}
				}
				if column_sum == n*3 {
					current_m += 1
				} else {
					break
				}
			}
			current_m = m*3 - current_span_left - current_m

			//-----------------------------------------------------

			//-----------------------------------------------------
			for y := range n * 3 {
				row_sum := 0
				for x := range m * 3 {
					if table[y][x] == 0b00000000 {
						row_sum += 1
					}
				}
				if row_sum == m*3 {

					current_span_top += 1
				} else {
					break
				}
			}

			//log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!! current_span_top ==%d", current_span_top)

			for y := n*3 - 1; y > 0; y-- { //current_span_top; y < n*3; y++ {
				row_sum := 0
				for x := range m * 3 {
					if table[y][x] == 0b00000000 {
						row_sum += 1
					}
				}
				if row_sum == m*3 {
					current_n += 1
				} else {
					break
				}
			}
			current_n = n*3 - current_span_top - current_n

			//-----------------------------------------------------

			log.Printf("empty left==%d   top==%d     new N==%d M==%d", current_span_left, current_span_top, current_n, current_m)

			temp_table := make([][]paper_cell.Quarter, current_n)
			for i := range current_n {
				temp_table[i] = make([]paper_cell.Quarter, current_m)
			}

			for y := range current_n {
				for x := range current_m {
					temp_table[y][x] = table[current_span_top+y][current_span_left+x]
					table[current_span_top+y][current_span_left+x] = 0b00000000
				}
			}

			for y := range current_n {
				for x := range current_m {
					table[n+y][m+x] = temp_table[y][x]
				}
			}

			s = ""
			for y := range n * 3 {
				for x := range m * 3 {
					switch c := table[y][x]; c {
					case 0b00001111:
						s += "# "
					case 0b00000000:
						s += ". "
					case 0b00001000:
						s += "v "
					case 0b00000111:
						s += "v "
					case 0b00000001:
						s += "> "
					case 0b00001110:
						s += "> "
					case 0b00000100:
						s += "< "
					case 0b00001011:
						s += "< "
					case 0b00000010:
						s += "^ "
					case 0b00001101:
						s += "^ "
					case 0b00001100:
						s += "\\ "
					case 0b00000011:
						s += "\\ "
					case 0b00000110:
						s += "/ "
					case 0b00001001:
						s += "/ "
					case 0b00001010:
						s += "x "
					case 0b00000101:
						s += "x "
					default:
						s += "_ "
					}
				}
				s += "\n"
			}
			log.Printf("\033[32mfolding passed:\n%s\n", s)

			s = ""
			for y := range current_n {
				for x := range current_m {
					switch c := table[n+y][m+x]; c {
					case 0b00001111:
						s += "#"
					case 0b00000000:
						s += "."
					case 0b00001000:
						s += "v"
					case 0b00000111:
						s += "v"
					case 0b00000001:
						s += ">"
					case 0b00001110:
						s += ">"
					case 0b00000100:
						s += "<"
					case 0b00001011:
						s += "<"
					case 0b00000010:
						s += "^"
					case 0b00001101:
						s += "^"
					case 0b00001100:
						s += "\\"
					case 0b00000011:
						s += "\\"
					case 0b00000110:
						s += "/"
					case 0b00001001:
						s += "/"
					case 0b00001010:
						s += "x"
					case 0b00000101:
						s += "x"
					default:
						s += "_"
					}
				}
				s += "\n"
			}
			fmt.Fprintf(out, "%s\n", s)

		}

	}

	//fmt.Fprint(out, "")
}

func get_D(y1, x1, y2, x2, y0, x0 int) float32 {

	return float32(x2-x1)*(float32(y0)+0.5-float32(y1)) - float32(y2-y1)*(float32(x0)+0.5-float32(x1))
}

func get_mirror_point_pos(y1, x1, y2, x2, y0, x0 int) (int, int) {
	switch {
	case x1 == x2: // Vertical
		return y0, 2*x1 - x0 - 1

	case y1 == y2: // Horizontal
		return 2*y1 - y0 - 1, x0

	case (y2 - y1) == (x2 - x1): // Diagonal
		dx := x1 - y1
		return y0 + dx, x0 - dx

	case (y2 - y1) == -(x2 - x1): // Diagonal
		dx := x1 + y1
		return -y0 + dx, -x0 + dx

	default:
		return -1, -1
	}
}
