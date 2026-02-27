# EditNotebook

## Tool Description
Edits Jupyter notebook cells and can also create new cells.  
This tool is the required mechanism for notebook content changes.

Behavior rules:
- Use `is_new_cell = false` to edit an existing cell (replace one occurrence of `old_string` with `new_string`).
- Use `is_new_cell = true` to create a new cell at `cell_idx`.
- Cell indices are 0-based.
- `old_string` must uniquely identify exactly one location in the target cell; include substantial context.
- Only one replacement occurrence is changed per call.
- For multiple edits, use multiple calls.
- Always provide arguments in this order:
  1. `target_notebook`
  2. `cell_idx`
  3. `is_new_cell`
  4. `cell_language`
  5. `old_string`
  6. `new_string`
- Always provide all required arguments, including both `old_string` and `new_string`.
- If creating a new notebook, use `is_new_cell = true` and `cell_idx = 0`.

## Tool Param
### target_notebook
* description: Path (relative or absolute) of the notebook file to modify.
* type: string
* required: yes

### cell_idx
* description: 0-based index of the target cell.
* type: integer
* required: yes

### is_new_cell
* description: If true, create a new cell at `cell_idx`; if false, edit existing cell at `cell_idx`.
* type: boolean
* required: yes

### cell_language
* description: Language/type of the target cell.
* type: string
* required: yes
* enum: ["python", "markdown", "javascript", "typescript", "r", "sql", "shell", "raw", "other"]
* constraints: must be one of the listed values

### old_string
* description: Exact text to replace in the target cell. Must be unique within that cell.
* type: string
* required: yes

### new_string
* description: Replacement text, or full content for a newly created cell.
* type: string
* required: yes
