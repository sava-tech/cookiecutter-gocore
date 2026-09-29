#!/usr/bin/env python3
import os
import secrets
import shutil
import subprocess
import sys

HOSTING_FILES = {
    'AWS': ['deploy'],
    'DigitalOcean': ['.do'],
    'Railway': ['railway.toml'],
}

# All of these are plain SMTP relays under the hood (pkg/emailer/smtp) —
# only the host/port/credentials differ per provider.
EMAIL_PRESETS = {
    'Google': {
        'host': 'smtp.gmail.com',
        'port': '587',
        'username_hint': 'your full Gmail address',
        'password_hint': 'a Google App Password (Google Account > Security > App Passwords; requires 2FA)',
    },
    'Mailgun': {
        'host': 'smtp.mailgun.org',
        'port': '587',
        'username_hint': 'postmaster@<your-mailgun-domain>',
        'password_hint': 'the SMTP password shown under Mailgun > Sending > Domain settings',
    },
    'Zoho': {
        'host': 'smtp.zoho.com',
        'port': '587',
        'username_hint': 'your full Zoho Mail address',
        'password_hint': 'your Zoho password, or an app-specific password if 2FA is enabled',
    },
    'SendGrid': {
        'host': 'smtp.sendgrid.net',
        'port': '587',
        'username_hint': 'the literal string "apikey"',
        'password_hint': 'your SendGrid API key',
    },
    'Amazon SES': {
        'host': 'email-smtp.us-east-1.amazonaws.com',
        'port': '587',
        'username_hint': 'an SES SMTP username (SES console > SMTP settings > Create SMTP credentials)',
        'password_hint': 'the matching SES SMTP password (not your AWS access key)',
    },
    'Postmark': {
        'host': 'smtp.postmarkapp.com',
        'port': '587',
        'username_hint': 'your Postmark Server API token',
        'password_hint': 'the same Postmark Server API token (used for both)',
    },
}

def apply_email_preset(content):
    """Pre-fill SMTP_HOST/PORT and leave a hint above the credential fields
    for the chosen email_service, so 'quick spin up' doesn't mean guessing
    each provider's SMTP settings from scratch."""
    provider = '{{ cookiecutter.email_service }}'
    preset = EMAIL_PRESETS.get(provider)
    if not preset:
        return content
    content = content.replace('\nEMAIL_PROVIDER=\n', '\nEMAIL_PROVIDER=smtp\n')
    content = content.replace('SMTP_HOST=', f'SMTP_HOST={preset["host"]}')
    content = content.replace('SMTP_PORT=587', f'SMTP_PORT={preset["port"]}')
    content = content.replace('SMTP_USERNAME=', f'# Username: {preset["username_hint"]}\nSMTP_USERNAME=')
    content = content.replace('SMTP_PASSWORD=', f'# Password: {preset["password_hint"]}\nSMTP_PASSWORD=')
    return content

def cleanup_hosting_provider():
    """Keep only the deploy config for the chosen hosting_provider and
    drop the others, so the prompt actually changes what gets generated."""
    provider = '{{ cookiecutter.hosting_provider }}'
    for name, paths in HOSTING_FILES.items():
        if name == provider:
            continue
        for path in paths:
            if os.path.isdir(path):
                shutil.rmtree(path, ignore_errors=True)
            elif os.path.isfile(path):
                os.remove(path)
    print(f"Kept deploy config for {provider}")

def generate_env_file():
    """Create .env from env.example with real random secrets instead of the
    placeholder values, so every generated project doesn't ship with the
    same fixed TOKEN_SYMMETRIC_KEY/SESSION_SECRET."""
    if os.path.exists('.env') or not os.path.exists('env.example'):
        return
    with open('env.example', 'r') as f:
        content = f.read()
    content = content.replace(
        'TOKEN_SYMMETRIC_KEY=CHANGE_ME_32_BYTE_RANDOM_SECRET',
        f'TOKEN_SYMMETRIC_KEY={secrets.token_hex(16)}',
    )
    content = content.replace(
        'SESSION_SECRET=CHANGE_ME_32_BYTE_RANDOM_SECRET',
        f'SESSION_SECRET={secrets.token_hex(16)}',
    )
    content = apply_email_preset(content)
    with open('.env', 'w') as f:
        f.write(content)
    print("Generated .env with random TOKEN_SYMMETRIC_KEY/SESSION_SECRET")
    provider = '{{ cookiecutter.email_service }}'
    if provider in EMAIL_PRESETS:
        print(f"Pre-filled SMTP settings for {provider} in .env — add your SMTP_USERNAME/SMTP_PASSWORD")
    else:
        print("email_service is None — .env defaults to Mailpit for local dev only")

def main():
    # Cookiecutter already runs hooks from the generated project directory
    # So the current working directory IS the project directory
    project_dir = os.getcwd()

    generate_env_file()
    cleanup_hosting_provider()

    # Ask user if they want to install dependencies now
    if '{{cookiecutter.auto_install_deps}}' == 'y':
        if os.path.exists('Makefile'):
            try:
                print("Installing dependencies...")
                subprocess.run(['make', 'swagger-doc'], check=True)
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