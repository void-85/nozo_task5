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

		table := make([][]uint8, n*3)
		for i := range n * 3 {
			table[i] = make([]uint8, m*3)
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

			s = ""
			for y := range current_n {
				for x := range current_m {

					if table[span_top+y][span_left+x] != 0b00000000 {

						d := get_D(y1, x1, y2, x2, y, x, ray_direction)
						if d > 0 {
							//table[span_top+y][span_left+x] = 'R'

							s += "R "
							new_y, new_x := get_mirror_point_pos(y1, x1, y2, x2, y, x)
							table[span_top+new_y][span_left+new_x] =
								table[span_top+new_y][span_left+new_x] |
									Invert_at_angle(
										table[span_top+y][span_left+x],
										ray_direction-90,
										false,
									)

							table[span_top+y][span_left+x] = 0b00000000

						} else if d == 0 {

							s += "M "
							//table[span_top+y][span_left+x] = '='
							table[span_top+y][span_left+x] =
								Invert_at_angle(
									table[span_top+y][span_left+x],
									ray_direction-90,
									true,
								)

						} else if d < 0 {
							//table[span_top+y][span_left+x] = 'L'
							s += "L "
						}

					} else {
						s += "_ "
					}
				}

				s += "\n"
			}
			log.Printf("current working window LINE division:\n%s", s)

			/* s = ""
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
			log.Printf("\033[32mjust folded:\n%s\n", s) */

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

			temp_table := make([][]uint8, current_n)
			for i := range current_n {
				temp_table[i] = make([]uint8, current_m)
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
					case 0b11110101:
						s += "R "
					case 0b11110110:
						s += "M "
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

func get_mirror_point_pos_OLD(y1, x1, y2, x2, y0, x0 int) (int, int) {

	switch {
	case x1 == x2: // Vertical
		log.Printf("p0(%d;%d) mirrored to p'0(%d;%d)", y0, x0, y0, 2*x1-x0-1)
		return y0, 2*x1 - x0 - 1

	case y1 == y2: // Horizontal
		log.Printf("p0(%d;%d) mirrored to p'0(%d;%d)", y0, x0, 2*y1-y0-1, x0)
		return 2*y1 - y0 - 1, x0

	case (y2 - y1) == (x2 - x1): // Diagonal
		dx := x1 - y1
		log.Printf("p0(%d;%d) mirrored to p'0(%d;%d)", y0, x0, y0+dx, x0-dx)
		return y0 + dx, x0 - dx

	case (y2 - y1) == -(x2 - x1): // Diagonal
		dx := x1 + y1
		log.Printf("p0(%d;%d) mirrored to p'0(%d;%d)", y0, x0, -y0+dx, -x0+dx)
		return -y0 + dx, -x0 + dx

	default:
		return -1, -1
	}
}

func get_mirror_point_pos(y1, x1, y2, x2, y0, x0 int) (int, int) {
	// Convert inputs to float64 for calculations
	dy := float64(y2 - y1)
	dx := float64(x2 - x1)

	// Check line type based on direction (must be 0 or 45 degrees)
	if dy == 0 && dx != 0 {
		// Horizontal line: y = y1
		// Reflection: (y0, x0) -> (2*y1 - y0, x0)
		yPrime := 2*y1 - y0
		xPrime := x0
		return yPrime, xPrime
	} else if dx == 0 && dy != 0 {
		// Vertical line: x = x1
		// Reflection: (y0, x0) -> (y0, 2*x1 - x0)
		yPrime := y0
		xPrime := 2*x1 - x0
		return yPrime, xPrime
	} else if math.Abs(dy) == math.Abs(dx) && dy != 0 {
		// Diagonal line (slope = ±1, i.e., 45 degrees)
		// For a 45-degree line, we can use a coordinate transformation or direct reflection.
		// Method: Reflect by swapping relative coordinates based on line slope.

		// Slope is +1 (y = x + c) or -1 (y = -x + c)
		slope := dy / dx // Either +1 or -1
		var yPrime, xPrime float64

		if slope == 1 {
			// Line y = x + c, where c = y1 - x1
			c := float64(y1 - x1)
			// Reflection over y = x + c: Swap x and y relative to the line
			// Midpoint of (x0, y0) and (x', y') lies on the line
			yPrime = float64(x0) + c
			xPrime = float64(y0) - c
		} else if slope == -1 {
			// Line y = -x + c, where c = y1 + x1
			c := float64(y1 + x1)
			// Reflection over y = -x + c
			yPrime = c - float64(x0)
			xPrime = c - float64(y0)
		} else {
			// Should not happen given constraints
			log.Printf("#\n#\n#\n#\n##########\nget_mirror_point_pos ERROR\n##########\n#\n#\n#\n#\n")
			return -333, -333
		}

		// Round to nearest integer
		return int(math.Round(yPrime)), int(math.Round(xPrime))
	} else {
		// Invalid line (not horizontal, vertical, or diagonal)
		//panic("Line must be horizontal, vertical, or 45-degree diagonal")
		log.Printf("#\n#\n#\n#\n##########\nget_mirror_point_pos ERROR\n##########\n#\n#\n#\n#\n")
		return -333, -333
	}
}

func Invert_at_angle(bits uint8, angle int, is_folding_edge bool) (res uint8) {

	if angle < 0 {
		angle += 360
	}

	switch a := angle; a {

	case 360, 180:
		res = (bits & 0b00000101) |
			((bits & 0b00001000) >> 2) |
			((bits & 0b00000010) << 2)

	case 90, 270:
		res = (bits & 0b00001010) |
			((bits & 0b00000100) >> 2) |
			((bits & 0b00000001) << 2)

	case 45, 225:
		res = ((bits & 0b00001000) >> 3) |
			((bits & 0b00000100) >> 1) |
			((bits & 0b00000010) << 1) |
			((bits & 0b00000001) << 3)

	case 135, 315:
		res = ((bits & 0b00001000) >> 1) |
			((bits & 0b00000100) << 1) |
			((bits & 0b00000010) >> 1) |
			((bits & 0b00000001) << 1)

	default:
		res = bits
	}

	if is_folding_edge {
		switch a := angle; a {
		case 45:
			res = (bits | res) & 0b00001100
		case 135:
			res = (bits | res) & 0b00000110
		case 225:
			res = (bits | res) & 0b00000011
		case 315:
			res = (bits | res) & 0b00001001
		}
	}

	return res
}

func OLD_get_D(y1, x1, y2, x2, y0, x0 int) float32 {

	return float32(x2-x1)*(float32(y0)+0.5-float32(y1)) - float32(y2-y1)*(float32(x0)+0.5-float32(x1))
	//return float32((x2*2-x1*2)*(y0*2+0-y1*2) - (y2*2-y1*2)*(x0*2+0-x1*2))
}

func get_D(y1, x1, y2, x2, y0, x0, angle int) int {

	dx := float32(x2) - float32(x1)
	dy := float32(y2) - float32(y1)
	px := float32(x0) - float32(x1)
	py := float32(y0) - float32(y1)

	//if angle%90 == 0 {
	px += 0.5
	py += 0.5
	//}

	D := int(dx*py - dy*px)

	//log.Printf("line: (%d;%d)--->(%d;%d)   point:(%d;%d)    D==%d", y1, x1, y2, x2, y0, x0, D)

	return D
}
