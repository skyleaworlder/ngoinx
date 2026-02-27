# ReadLints

## Tool Description
Reads IDE/linter diagnostics (warnings/errors).  
Can return diagnostics for the entire workspace or for selected files/directories.

Behavior notes:
- Diagnostics can occasionally be outdated.
- Prefer limiting scope with `paths` when possible.
- Ignore issues that existed before current changes.

## Tool Param
### paths
* description: Optional list of file/directory paths to scope diagnostics.
* type: array<string>
* required: no
* default: all workspace diagnostics when omitted
