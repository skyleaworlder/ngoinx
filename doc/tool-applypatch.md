# ApplyPatch

## Tool Description
Edits files using a strict patch language (structured diff-like format).  
Supports adding a new file or updating one existing file per call.

Critical behavior/constraints:
- Exactly one file operation per call.
- Must use patch envelope:
  - `*** Begin Patch`
  - file operation section
  - `*** End Patch`
- File operation headers:
  - `*** Add File: <absolute-path>`
  - `*** Update File: <absolute-path>`
- Update hunks use `@@` context markers and line prefixes:
  - space (` `) for context lines
  - `-` for removed lines
  - `+` for added lines
- All paths must be absolute.
- Read existing file before patching (unless creating a new file).
- Do not call this tool in parallel with other tools.

Accepted grammar (high level):
- `begin_patch -> file hunk(s) -> end_patch`
- add hunk requires only `+` lines
- update hunk supports context + additions/removals
- optional `*** End of File` marker in updates

## Tool Param
### patch_text
* description: Entire patch payload string in the required patch grammar (freeform input; not JSON).
* type: string (FREEFORM patch language)
* required: yes
* constraints:
  - must begin with `*** Begin Patch`
  - must end with `*** End Patch`
  - must include exactly one file operation (`Add File` or `Update File`)
  - must follow allowed hunk/line-prefix grammar
