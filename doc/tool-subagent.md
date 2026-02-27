# Subagent

## Tool Description
Launches a specialized autonomous subagent to handle complex multi-step tasks.  
Subagents can be used for deeper investigations, broad code search, or longer scoped workflows.

Behavior notes:
- Requires selecting a `subagent_type`.
- Supports resuming prior subagent sessions via `resume`.
- Can be launched in readonly mode (`readonly: true`) for ask-only behavior.
- Up to 4 concurrent subagents is the recommended maximum.
- Available `subagent_type` in this environment: `generalPurpose`.
- Available model option documented in schema: `fast`.

## Tool Param
### description
* description: Short 3-5 word summary of what the subagent will do.
* type: string
* required: yes

### prompt
* description: Detailed task instructions for the subagent, including required output expectations.
* type: string
* required: yes

### model
* description: Optional model choice for the subagent.
* type: string
* required: no
* enum: ["fast"]

### resume
* description: Existing subagent ID to continue prior execution context.
* type: string
* required: no

### readonly
* description: If true, run subagent in readonly/ask mode with restricted write operations.
* type: boolean
* required: no

### subagent_type
* description: Subagent capability type to use.
* type: string
* required: no
* enum: ["generalPurpose"]

### attachments
* description: Optional list of image/video file paths to pass into subagent context.
* type: array<string>
* required: no
* supported_formats:
  - images: png, jpg, gif, webp
  - videos: mp4, webm
