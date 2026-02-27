# FetchMcpResource

## Tool Description
Reads a specific MCP resource by server name and resource URI.  
Optionally downloads that resource to a workspace-relative path instead of returning content inline.

## Tool Param
### server
* description: MCP server identifier.
* type: string
* required: yes

### uri
* description: Resource URI to read from the specified MCP server.
* type: string
* required: yes

### downloadPath
* description: Optional workspace-relative output path; when set, resource is saved to disk and not returned in response content.
* type: string
* required: no
