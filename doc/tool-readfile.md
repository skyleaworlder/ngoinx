# ReadFile

## Tool Description
Reads a local filesystem file directly.  
Supports text files, images (jpeg/jpg/png/gif/webp), and PDF files (converted to text).  
For long files, supports partial reads via line offset and line limit.

Behavior notes:
- If a provided file path is invalid or missing, the tool returns an error.
- Line-oriented output is numbered in format: `LINE_NUMBER|LINE_CONTENT`.
- Empty files return `File is empty.`

## Tool Param
### path
* description: Absolute path of the file to read.
* type: string
* required: yes

### offset
* description: Starting line number for partial reading.
* type: integer
* required: no

### limit
* description: Maximum number of lines to read.
* type: integer
* required: no
