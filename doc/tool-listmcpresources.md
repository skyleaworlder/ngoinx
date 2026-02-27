# ListMcpResources

## Tool Description
Lists available resources from configured MCP servers.  
Each resource includes standard MCP fields plus a `server` field indicating source server.

Behavior notes:
- MCP resources are not the same thing as MCP tools.
- Can query all servers or a specific server.

## Tool Param
### server
* description: Optional MCP server identifier to filter results.
* type: string
* required: no
* default: resources from all configured servers
