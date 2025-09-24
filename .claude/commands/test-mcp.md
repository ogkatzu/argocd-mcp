# Testing MCP Server with Claude Code

After restarting Claude Code, you should have access to your MCP server tools. Here's how to test them:

## 🔧 Available Tools

Your MCP server provides these tools that should now be available in Claude Code:

### 1. **echo** - Text Echo Tool
Test with: "Use the echo tool to repeat 'Hello from MCP!'"

### 2. **calculate** - Math Calculator
Test with: "Use the calculate tool to add 123 and 456"
Or: "Calculate 15.5 divided by 3.2"

### 3. **system_info** - System Information
Test with: "Use the system_info tool to show my system details"

### 4. **read_file** - File Reader
Test with: "Use the read_file tool to read the contents of README.md"

## 📚 Available Resources

Your server also provides resources (though these are accessed differently):

- **config://server** - Server configuration
- **status://server** - Server runtime status
- **help://tools** - Tool documentation

## 🧪 Test Commands to Try

Copy and paste these into Claude Code to test each tool:

```
1. Test Echo:
"Can you use the echo tool to repeat the message 'MCP is working great!'"

2. Test Calculator:
"Please use the calculate tool to multiply 42 by 17"

3. Test System Info:
"Show me my system information using the system_info tool"

4. Test File Reading:
"Use the read_file tool to read the sample.txt file"
```

## ✅ What You Should See

- Claude Code should automatically recognize and use your MCP tools
- Tool calls will show up with proper formatting
- You'll see the actual tool results integrated in the conversation
- Error handling will work properly (try dividing by zero with calculate!)

## 🔍 Troubleshooting

If the tools don't appear:

1. **Check Configuration**: Verify the path in `~/.claude-code/mcp_servers.json` is correct
2. **Restart Claude Code**: Make sure you've restarted after adding the configuration
3. **Check Permissions**: Ensure the `mcp-server` binary is executable (`chmod +x mcp-server`)
4. **Test Manually**: Run `./mcp-server` in the terminal to verify it works standalone

## 🎯 Success Indicators

When everything is working correctly:
- ✅ Tools appear in Claude Code's available tools
- ✅ Tool calls execute without errors
- ✅ Results are properly formatted and displayed
- ✅ Error cases are handled gracefully
- ✅ Server logs show requests being processed

Enjoy testing your MCP server integration! 🚀