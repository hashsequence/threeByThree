package main

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
)

var tpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html>
<head>
	<title>threeBythree CSV Processor</title>
</head>
<body>
	<h1>threeBythree CSV Processor</h1>
	<form method="POST" enctype="multipart/form-data">
		<label>Input CSV: <input type="file" name="inputcsv" required></label><br><br>
		<label>Empty Row Every N: <input type="number" name="nrow" min="1" required></label><br><br>
		<label>Empty Col Every N: <input type="number" name="ncol" min="1" required></label><br><br>
		<label>Output CSV Name: <input type="text" name="outputcsv" required></label><br><br>
		<input type="submit" value="Process">
	</form>
	{{if .Output}}
		<h2>Download Output:</h2>
		<a href="/download?file={{.Output}}">{{.Output}}</a>
	{{end}}
</body>
</html>
`))

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/download", downloadHandler)
	fmt.Println("Serving on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	var output string
	if r.Method == http.MethodPost {
		file, header, err := r.FormFile("inputcsv")
		if err != nil {
			http.Error(w, "Error reading input file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		outputcsv := r.FormValue("outputcsv")
		nrow, _ := strconv.Atoi(r.FormValue("nrow"))
		ncol, _ := strconv.Atoi(r.FormValue("ncol"))

		// Save uploaded file temporarily
		infile := "tmp_" + header.Filename
		f, err := os.Create(infile)
		if err != nil {
			http.Error(w, "Error saving input file", http.StatusInternalServerError)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		// Process CSV
		processCSV(infile, nrow, ncol, outputcsv)
		output = outputcsv
	}
	tpl.Execute(w, map[string]string{"Output": output})
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	file := r.URL.Query().Get("file")
	w.Header().Set("Content-Disposition", "attachment; filename="+file)
	w.Header().Set("Content-Type", "text/csv")
	http.ServeFile(w, r, file)
}

// processCSV uses the same logic as main.go
func processCSV(input string, nRow, nCol int, output string) {
	f, err := os.Open(input)
	if err != nil {
		return
	}
	defer f.Close()
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return
	}
	records = addEmptyRowEveryN(records, nRow)
	records = addEmptyColumnEveryN(records, nCol)
	records = rotateSquareMatrixCounterClockwise(records)
	out, err := os.Create(output)
	if err != nil {
		return
	}
	defer out.Close()
	writer := csv.NewWriter(out)
	writer.WriteAll(records)
	writer.Flush()
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

