# rg

## Tool Description
Searches workspace contents using ripgrep-compatible regex queries, respecting ignore rules (`.gitignore`, `.cursorignore`).  
Supports content output, file-path output, and count output, with optional context lines and pagination-like controls.

Behavior notes:
- Prefer narrowing by `type` or `path` for performance.
- Use multiline mode only when needed (it can degrade performance).
- If output says results were truncated ("at least ..."), tighten query or increase `head_limit`.

## Tool Param
### pattern
* description: Regular expression pattern to search for in file contents.
* type: string
* required: yes

### path
* description: File or directory path scope for the search (`rg pattern -- PATH` semantics).
* type: string
* required: no
* default: workspace root

### glob
* description: Glob filter for files (mapped to `rg --glob`).
* type: string
* required: no

### output_mode
* description: Controls output format.
* type: string
* required: no
* default: "content"
* enum: ["content", "files_with_matches", "count"]

### -B
* description: Number of lines shown before each match (`rg -B`).
* type: integer
* required: no
* constraints: only effective when `output_mode` is `"content"`

### -A
* description: Number of lines shown after each match (`rg -A`).
* type: integer
* required: no
* constraints: only effective when `output_mode` is `"content"`

### -C
* description: Number of context lines before and after each match (`rg -C`).
* type: integer
* required: no
* constraints: only effective when `output_mode` is `"content"`

### -i
* description: Case-insensitive search toggle.
* type: boolean
* required: no
* default: false

### type
* description: File type filter (`rg --type`), e.g. js, py, go, java.
* type: string
* required: no

### head_limit
* description: Maximum returned entries (matches for content mode, files for other modes).
* type: number
* required: no
* minimum: 0

### offset
* description: Number of initial entries to skip (supports pagination with `head_limit`).
* type: number
* required: no
* minimum: 0

### multiline
* description: Enables multiline mode (`rg -U --multiline-dotall`) so `.` can span lines.
* type: boolean
* required: no
* default: false
