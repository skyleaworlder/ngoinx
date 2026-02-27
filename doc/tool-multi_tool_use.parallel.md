# multi_tool_use.parallel

## Tool Description
Wrapper tool to execute multiple tool calls concurrently when they are independent.  
Only tools in the `functions` namespace are allowed as recipients.

Behavior notes:
- Intended for parallelization and latency reduction.
- Each entry in `tool_uses` fully specifies one tool call.

## Tool Param
### tool_uses
* description: Array of tool invocation specifications to run in parallel.
* type: array<object>
* required: yes

#### tool_uses[].recipient_name
* description: Fully-qualified tool name to invoke (format: `functions.<tool_name>`).
* type: string
* required: yes

#### tool_uses[].parameters
* description: Parameter object passed to the selected tool; must satisfy that tool's schema.
* type: object
* required: yes
