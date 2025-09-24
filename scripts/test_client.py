#!/usr/bin/env python3
"""
Simple test client for the MCP server to demonstrate functionality.
"""

import json
import subprocess
import sys

def send_request(proc, request_id, method, params=None):
    """Send a JSON-RPC request to the MCP server."""
    request = {
        "jsonrpc": "2.0",
        "id": request_id,
        "method": method
    }
    if params:
        request["params"] = params

    json_str = json.dumps(request) + "\n"
    proc.stdin.write(json_str.encode())
    proc.stdin.flush()

    # Read response
    response_line = proc.stdout.readline().decode().strip()
    if response_line:
        return json.loads(response_line)
    return None

def main():
    # Start the MCP server
    proc = subprocess.Popen(
        ["./mcp-server"],
        stdin=subprocess.PIPE,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE
    )

    try:
        print("🚀 Testing MCP Server")
        print("=" * 50)

        # 1. Initialize
        print("1. Initializing...")
        response = send_request(proc, 1, "initialize", {
            "capabilities": {},
            "clientInfo": {"name": "test-client", "version": "1.0.0"}
        })
        if response and "result" in response:
            print("✅ Initialization successful")
        else:
            print("❌ Initialization failed")
            print(f"Response: {response}")

        # 2. List tools
        print("\n2. Listing tools...")
        response = send_request(proc, 2, "tools/list")
        if response and "result" in response:
            tools = response["result"]["tools"]
            print(f"✅ Found {len(tools)} tools:")
            for tool in tools:
                print(f"   - {tool['name']}: {tool['description']}")

        # 3. Test echo tool
        print("\n3. Testing echo tool...")
        response = send_request(proc, 3, "tools/call", {
            "name": "echo",
            "arguments": {"text": "Hello, MCP!"}
        })
        if response and "result" in response:
            content = response["result"]["content"][0]["text"]
            print(f"✅ Echo result: {content}")

        # 4. Test calculate tool
        print("\n4. Testing calculate tool...")
        response = send_request(proc, 4, "tools/call", {
            "name": "calculate",
            "arguments": {"operation": "add", "a": 15, "b": 27}
        })
        if response and "result" in response:
            content = response["result"]["content"][0]["text"]
            print(f"✅ Calculate result: {content}")

        # 5. Test system_info tool
        print("\n5. Testing system_info tool...")
        response = send_request(proc, 5, "tools/call", {
            "name": "system_info",
            "arguments": {}
        })
        if response and "result" in response:
            content = response["result"]["content"][0]["text"]
            print("✅ System info retrieved:")
            print(content[:200] + "..." if len(content) > 200 else content)

        # 6. List resources
        print("\n6. Listing resources...")
        response = send_request(proc, 6, "resources/list")
        if response and "result" in response:
            resources = response["result"]["resources"]
            print(f"✅ Found {len(resources)} resources:")
            for resource in resources:
                print(f"   - {resource['uri']}: {resource['name']}")

        # 7. Test config resource
        print("\n7. Testing config resource...")
        response = send_request(proc, 7, "resources/read", {
            "uri": "config://server"
        })
        if response and "result" in response:
            content = response["result"]["contents"][0]["text"]
            config = json.loads(content)
            print(f"✅ Config retrieved: {config['name']} v{config['version']}")

        print("\n🎉 All tests completed successfully!")

    except Exception as e:
        print(f"❌ Test failed: {e}")
    finally:
        proc.terminate()
        proc.wait()

if __name__ == "__main__":
    main()