# Shell

## Tool Description
Executes a given command in a stateful shell session for terminal operations (for example: git, npm, docker, python).  
The shell session persists across calls (current working directory and environment variables are preserved).  
This tool should **not** be used for file content operations (reading, writing, editing, searching files); dedicated tools should be used for those.

Operational requirements and constraints:
- Before commands that create files/directories, verify parent directory first (for example, run `ls`).
- Paths containing spaces must be quoted with double quotes.
- Do not run long-lived/hanging commands (watchers, forever dev servers, endless background tasks).
- Prefer `working_directory` instead of `cd ... && ...`.
- Prefer `rg`/`Glob` for searching and `ReadFile` for reading files.
- Avoid `find`, `grep`, `cat`, `head`, `tail` in normal workflows.
- When adding dependencies, use package managers (npm/pip/etc.) and avoid made-up versions.

## Tool Param
### command
* description: Shell command string to execute.
* type: string
* required: yes

### working_directory
* description: Absolute path where the command should run. If omitted, command runs in current shell directory.
* type: string
* required: no
* constraints: should be an absolute path

### timeout
* description: Command timeout in milliseconds.
* type: integer
* required: no
* default: 30000
* max: 600000

### description
* description: Clear concise summary of what the command does.
* type: string
* required: no
* constraints: recommended 5-10 words

### is_background
* description: Whether to run the command in background mode.
* type: boolean
* required: no
