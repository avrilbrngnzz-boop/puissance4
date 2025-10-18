package routes

func Seq(start, end int) []int {
	s := make([]int, end-start+1)
	for i := range s {
		s[i] = start + i
	}
	return s
}

func CellClass(cell string) string {
	switch cell {
	case "X":
		return "X"
	case "O":
		return "O"
	default:
		return "empty"
	}
}

func CheckWin(grid [][]string, symbol string) bool {
	rows := len(grid)
	cols := len(grid[0])
	for i := 0; i < rows; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i][j+1] == symbol &&
				grid[i][j+2] == symbol && grid[i][j+3] == symbol {
				return true
			}
		}
	}
	for i := 0; i <= rows-4; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == symbol && grid[i+1][j] == symbol &&
				grid[i+2][j] == symbol && grid[i+3][j] == symbol {
				return true
			}
		}
	}
	for i := 0; i <= rows-4; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i+1][j+1] == symbol &&
				grid[i+2][j+2] == symbol && grid[i+3][j+3] == symbol {
				return true
			}
		}
	}
	for i := 3; i < rows; i++ {
		for j := 0; j <= cols-4; j++ {
			if grid[i][j] == symbol && grid[i-1][j+1] == symbol &&
				grid[i-2][j+2] == symbol && grid[i-3][j+3] == symbol {
				return true
			}
		}
	}
	return false
}

func IsDraw(grid [][]string) bool {
	for i := range grid {
		for j := range grid[i] {
			if grid[i][j] == "" {
				return false
			}
		}
	}
	return true
}
