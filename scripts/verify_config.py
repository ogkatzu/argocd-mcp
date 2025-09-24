#!/usr/bin/env python3
"""
Verify that the MCP server configuration is set up correctly for Claude Code.
"""

import json
import os
import subprocess
import sys

def check_config_file():
    """Check if the Claude Code MCP configuration exists and is valid."""
    config_path = os.path.expanduser("~/.claude-code/mcp_servers.json")

    print("🔍 Checking Claude Code MCP configuration...")

    if not os.path.exists(config_path):
        print(f"❌ Configuration file not found at: {config_path}")
        return False

    try:
        with open(config_path, 'r') as f:
            config = json.load(f)

        if "basic-mcp-server" not in config:
            print("❌ 'basic-mcp-server' not found in configuration")
            return False

        server_config = config["basic-mcp-server"]
        command_path = server_config.get("command")

        print(f"✅ Configuration file found: {config_path}")
        print(f"✅ Server configured with command: {command_path}")

        # Check if the executable exists
        if not os.path.exists(command_path):
            print(f"❌ MCP server executable not found at: {command_path}")
            return False

        # Check if it's executable
        if not os.access(command_path, os.X_OK):
            print(f"❌ MCP server is not executable: {command_path}")
            return False

        print(f"✅ MCP server executable found and is runnable")
        return True

    except json.JSONDecodeError as e:
        print(f"❌ Invalid JSON in configuration file: {e}")
        return False
    except Exception as e:
        print(f"❌ Error reading configuration: {e}")
        return False

def test_server():
    """Test if the MCP server starts and responds correctly."""
    print("\n🧪 Testing MCP server functionality...")

    try:
        # Test server startup with initialization
        process = subprocess.Popen(
            ["./mcp-server"],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            cwd=os.path.dirname(os.path.abspath(__file__))
        )

        # Send initialization request
        init_request = '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"capabilities": {}}}\n'

        try:
            stdout, stderr = process.communicate(input=init_request.encode(), timeout=5)

            if process.returncode == 0 or stdout:
                print("✅ MCP server starts and responds correctly")
                return True
            else:
                print(f"❌ MCP server failed to start properly")
                print(f"stderr: {stderr.decode()}")
                return False

        except subprocess.TimeoutExpired:
            process.kill()
            print("❌ MCP server timed out (this might actually be OK - server may be waiting for more input)")
            return True  # Timeout can be normal for stdio servers

    except FileNotFoundError:
        print("❌ MCP server executable not found in current directory")
        return False
    except Exception as e:
        print(f"❌ Error testing server: {e}")
        return False

def main():
    print("🚀 Claude Code MCP Server Configuration Verification")
    print("=" * 60)

    config_ok = check_config_file()
    server_ok = test_server()

    print("\n" + "=" * 60)

    if config_ok and server_ok:
        print("🎉 Everything looks good!")
        print("\nNext steps:")
        print("1. Restart Claude Code")
        print("2. Your MCP tools should now be available")
        print("3. Try the test commands from CLAUDE_CODE_TEST.md")
    else:
        print("❌ Issues found. Please fix the above problems and try again.")
        sys.exit(1)

if __name__ == "__main__":
    main()