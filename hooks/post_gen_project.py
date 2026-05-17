#!/usr/bin/env python3
import os
import subprocess
import sys

def main():
    # Cookiecutter already runs hooks from the generated project directory
    # So the current working directory IS the project directory
    project_dir = os.getcwd()
    
    # Ask user if they want to install dependencies now
    if '{{cookiecutter.auto_install_deps}}' == 'y':
        if os.path.exists('Makefile'):
            try:
                print("Installing dependencies...")
                subprocess.run(['make', 'install-dependencies'], check=True)
                print("Done....")
            except subprocess.CalledProcessError as e:
                print(f"Failed to install dependencies: {e}")
                print("You can manually run 'make install-dependencies' later")
            except FileNotFoundError:
                print("Make command not found. Please install make first.")
                print("You can manually run 'go mod download' to install dependencies")
        else:
            print("No Makefile found, skipping dependency installation")
    else:
        print("Skipping dependency installation")
        print("You can run 'make install-dependencies' manually later")
    
    print("\n✨ Project generated successfully! ✨")
    print(f"📁Project created at: {project_dir}")
    print("\nNext steps:")
    print(f"  cd {project_dir}")
    print("  make help         # Show available commands")
    if '{{cookiecutter.use_docker}}' == 'y':
        print("  make docker-run   # Start the application with Docker")
    else:
        print("  make run          # Run the application locally")

if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print(f"Warning: Post-generation setup had an issue: {e}")
        print("Your project was generated successfully.")
        print("You can manually run 'make install-dependencies' in the project directory")
        sys.exit(0)  # Don't fail the generation