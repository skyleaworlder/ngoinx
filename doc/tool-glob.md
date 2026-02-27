# Glob

## Tool Description
Searches for files matching a glob pattern, recursively and efficiently, and returns matching paths sorted by modification time.  
Useful for discovering files by name/path pattern across large repositories.

Behavior notes:
- If `glob_pattern` does not start with `**/`, the system automatically prepends `**/` for recursive matching.
- Typical use cases include finding file extensions, directories, or naming patterns.

## Tool Param
### target_directory
* description: Absolute directory path to search inside.
* type: string
* required: no
* default: workspace root

### glob_pattern
* description: Glob pattern used to match files.
* type: string
* required: yes
* behavior: patterns not starting with `**/` are automatically converted to recursive form
