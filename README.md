# threeBythreeApp

## Description
`threeBythreeApp` is a Go program that processes CSV files with the following features:
- Adds an empty row every Nth row
- Adds a column with empty values every Nth column
- Rotates any 3x3 block counterclockwise around a nonempty center cell, only if all 8 neighbors are nonempty

## Usage
```
./threeBythreeApp <input_csv> <n_row> <n_col> <output_csv>
```
- `<input_csv>`: Path to the input CSV file
- `<n_row>`: Add an empty row every Nth row (integer)
- `<n_col>`: Add an empty column every Nth column (integer)
- `<output_csv>`: Path to the output CSV file

## Example
```
./threeBythreeApp test1.csv 3 3 newTest1.csv
```
This command will:
- Read `test1.csv`
- Add an empty row every 3 rows
- Add an empty column every 3 columns
- Rotate all eligible 3x3 blocks
- Write the result to `newTest1.csv`

## Requirements
- Go 1.18 or newer (for building from source)

## Build
To build the program from source:
```
go build -o threeBythreeApp main.go
```

## Notes
- The program expects a well-formed CSV file with consistent row lengths.
- Rotation only occurs for nonempty cells with 8 nonempty neighbors.
