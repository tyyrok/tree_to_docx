## Utility to export files/directories tree in docx with file content

So it works in a very simple way - put your binary executable in directory where you want to get results and execute it.
Result could be found in file `structure.docx`


### Repository Includes:
- Go utility code
- structure.docx - an output example
- folder_to_test - folder with some code to test the utility

### Setup Instructions and Usage Instructions
1. Use executable (for Windows /bin/tree_to_docx.exe, for Linux you could build it by `GOOS=linux GOARCH=amd64 go build -o bin/tree_to_docx main.go`)

### Some Warnings
1. Code ignores files starting with `.`, but some configs might be exported to .docx too, so be careful
