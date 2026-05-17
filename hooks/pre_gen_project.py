#!/usr/bin/env python3
import subprocess
import sys

def main():
    """Validate environment before generating project"""
    
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