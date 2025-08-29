package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
)

// rotateSquareMatrixCounterClockwise rotates every square submatrix of nonempty cells counterclockwise 90 degrees
func rotateSquareMatrixCounterClockwise(records [][]string) [][]string {
    rows := len(records)
    if rows == 0 {
        return records
    }
    cols := len(records[0])
    minDim := rows
    if cols < rows {
        minDim = cols
    }
    // For every possible size (2x2 up to minDim x minDim)
    for size := 2; size <= minDim; size++ {
        for i := 0; i <= rows-size; i++ {
            for j := 0; j <= cols-size; j++ {
                // Check if all cells in the submatrix are nonempty
                allNonEmpty := true
                for r := i; r < i+size; r++ {
                    for c := j; c < j+size; c++ {
                        if records[r][c] == "" {
                            allNonEmpty = false
                            break
                        }
                    }
                    if !allNonEmpty {
                        break
                    }
                }
                if !allNonEmpty {
                    continue
                }
                // Check if the submatrix is surrounded by empty cells or the border
                surrounded := true
                // Top border
                if i > 0 {
                    for c := j; c < j+size; c++ {
                        if records[i-1][c] != "" {
                            surrounded = false
                            break
                        }
                    }
                }
                // Bottom border
                if surrounded && i+size < rows {
                    for c := j; c < j+size; c++ {
                        if records[i+size][c] != "" {
                            surrounded = false
                            break
                        }
                    }
                }
                // Left border
                if surrounded && j > 0 {
                    for r := i; r < i+size; r++ {
                        if records[r][j-1] != "" {
                            surrounded = false
                            break
                        }
                    }
                }
                // Right border
                if surrounded && j+size < cols {
                    for r := i; r < i+size; r++ {
                        if records[r][j+size] != "" {
                            surrounded = false
                            break
                        }
                    }
                }
                if surrounded {
                    // Rotate the submatrix counterclockwise
                    rotateSubMatrixCCW(records, i, j, size)
                }
            }
        }
    }
    return records
}

// rotateSubMatrixCCW rotates a square submatrix in place
func rotateSubMatrixCCW(records [][]string, startRow, startCol, size int) {
	// Layer by layer rotation
	for layer := 0; layer < size/2; layer++ {
		first := layer
		last := size - 1 - layer
		for k := first; k < last; k++ {
			offset := k - first
			// Save top
			top := records[startRow+first][startCol+k]
			// right -> top
			records[startRow+first][startCol+k] = records[startRow+k][startCol+last]
			// bottom -> right
			records[startRow+k][startCol+last] = records[startRow+last][startCol+last-offset]
			// left -> bottom
			records[startRow+last][startCol+last-offset] = records[startRow+last-offset][startCol+first]
			// top -> left
			records[startRow+last-offset][startCol+first] = top
		}
	}
}

// rotate3x3CounterClockwise rotates each 3x3 matrix counterclockwise around the center cell if it is fully surrounded
func rotate3x3CounterClockwise(records [][]string) [][]string {
    rows := len(records)
    if rows < 3 {
        return records
    }
    cols := len(records[0])
    if cols < 3 {
        return records
    }
    // Make a deep copy to avoid in-place modification
    result := make([][]string, rows)
    for i := range records {
        result[i] = make([]string, len(records[i]))
        copy(result[i], records[i])
    }
    for i := 1; i < rows-1; i++ {
        for j := 1; j < cols-1; j++ {
            center := records[i][j]
            if center == "" {
                continue
            }
            // Check all 8 neighbors are nonempty
            if records[i-1][j-1] == "" || records[i-1][j] == "" || records[i-1][j+1] == "" ||
               records[i][j-1] == ""   || records[i][j+1] == ""   ||
               records[i+1][j-1] == "" || records[i+1][j] == "" || records[i+1][j+1] == "" {
                continue
            }
            // Perform counterclockwise rotation
            result[i-1][j-1] = records[i-1][j+1] // top left <- top right
            result[i-1][j]   = records[i][j+1]   // top center <- middle right
            result[i-1][j+1] = records[i+1][j+1] // top right <- bottom right
            result[i][j+1]   = records[i+1][j]   // middle right <- bottom center
            result[i+1][j+1] = records[i+1][j-1] // bottom right <- bottom left
            result[i+1][j]   = records[i][j-1]   // bottom center <- middle left
            result[i+1][j-1] = records[i-1][j-1] // bottom left <- top left
            result[i][j-1]   = records[i-1][j]   // middle left <- top center
            // center cell remains unchanged
        }
    }
    return result
}

// addEmptyRowEveryN adds an empty row every n rows
func addEmptyRowEveryN(records [][]string, n int) [][]string {
    var result [][]string
    for i, row := range records {
        result = append(result, row)
        if (i+1)%n == 0 {
            result = append(result, make([]string, len(row)))
        }
    }
    return result
}

// addEmptyColumnEveryN adds a column with empty values every n columns
func addEmptyColumnEveryN(records [][]string, n int) [][]string {
    for i := range records {
        row := records[i]
        newRow := []string{}
        count := 0
        for j := 0; j < len(row); j++ {
            newRow = append(newRow, row[j])
            count++
            if n > 0 && count%n == 0 && j != len(row)-1 {
                newRow = append(newRow, "")
            }
        }
        records[i] = newRow
    }
    return records
}

func main() {
    if len(os.Args) < 5 {
        fmt.Println("Usage: go run main.go <csv_file> <n_row> <n_col> <output_file>")
        return
    }
    filePath := os.Args[1]
    nRow, err := strconv.Atoi(os.Args[2])
    if err != nil {
        fmt.Println("Invalid n_row value")
        return
    }
    nCol, err := strconv.Atoi(os.Args[3])
    if err != nil {
        fmt.Println("Invalid n_col value")
        return
    }
    outFile := os.Args[4]

    f, err := os.Open(filePath)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer f.Close()

    reader := csv.NewReader(f)
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error reading CSV:", err)
        return
    }

		records = addEmptyColumnEveryN(records, nCol)
    records = addEmptyRowEveryN(records, nRow)
    records = rotate3x3CounterClockwise(records)

    out, err := os.Create(outFile)
    if err != nil {
        fmt.Println("Error creating output file:", err)
        return
    }
    defer out.Close()

    writer := csv.NewWriter(out)
    err = writer.WriteAll(records)
    if err != nil {
        fmt.Println("Error writing CSV:", err)
        return
    }
    writer.Flush()
    fmt.Println("Processed CSV written to", outFile)
}
