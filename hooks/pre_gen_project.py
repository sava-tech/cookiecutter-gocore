#!/usr/bin/env python3
import re
import subprocess
import sys

MODULE_PATH = "{{ cookiecutter.module_path }}"

def main():
    """Validate environment before generating project"""

    # The module path becomes the Go module name (go.mod) and every internal
    # import path, so it must be a valid, lowercase Go module path segment.
    if not re.fullmatch(r"[a-z0-9]+(?:[/_-][a-z0-9]+)*", MODULE_PATH.split("/", 2)[-1]):
        print(f" Error: '{MODULE_PATH}' is not a valid Go module path")
        print("project_name/github_username must be lowercase letters, digits, '-' or '_' only (no spaces)")
        sys.exit(1)

    # Check if make is installed
    try:
        subprocess.run(['make', '--version'], capture_output=True, check=True)
    except (subprocess.CalledProcessError, FileNotFoundError):
        print(" Error: 'make' is not installed or not in PATH")
        print("Please install make before generating this project")
        sys.exit(1)
    
    # Check Go version
    try:
        result = subprocess.run(['go', 'version'], capture_output=True, text=True, check=True)
        print(f"✓ Go found: {result.stdout.strip()}")
    except (subprocess.CalledProcessError, FileNotFoundError):
        print(" Error: Go is not installed or not in PATH")
        print("Please install Go 1.25 before generating this project")
        sys.exit(1)
    
    print("✓ Environment validation passed!")

if __name__ == "__main__":
    main()