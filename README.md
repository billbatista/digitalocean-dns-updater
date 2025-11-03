# DigitalOcean DNS Updater

A simple command-line tool automatically detects your public IP address and updates DNS records on DigitalOcean using their API.

## Features

- Automatically detects your current public IP address
- Update A, AAAA, or other DNS record types
- Command-line interface with flags
- Preserves existing record settings (TTL, etc.)

## Prerequisites

- A DigitalOcean account with domains managed through their DNS service
- A DigitalOcean API token with read/write permissions

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the appropriate binary for your operating system from the [GitHub releases page](https://github.com/billbatista/digitalocean-dns-updater/releases):

- **Windows**: `dns-updater-windows-amd64.exe`
- **macOS**: `dns-updater-darwin-amd64` 
- **Linux**: `dns-updater-linux-amd64`

Make the binary executable (macOS/Linux only):
```bash
chmod +x dns-updater-*
```

### Option 2: Build from Source

If you have Go installed and want to build from source:

1. Clone this repository
2. Navigate to the project directory:
   ```bash
   cd digitalocean-dns-updater
   ```
3. Build the binary:
   ```bash
   go build -o dns-updater
   ```

## Getting Your DigitalOcean API Token

1. **Log into your DigitalOcean account**: Go to [https://cloud.digitalocean.com/](https://cloud.digitalocean.com/)

2. **Navigate to API settings**: 
   - Click on your profile picture in the top right
   - Select "API" from the dropdown menu
   - Or go directly to [https://cloud.digitalocean.com/account/api/](https://cloud.digitalocean.com/account/api/)

3. **Generate a new token**:
   - Click "Generate New Token"
   - Give your token a descriptive name (e.g., "DNS Updater")
   - Select "Custom Scopes" and check the "domain" resource type OR
     - Select "Full Access"
   - Select "Read" and "Write" permissions
   - Click "Generate Token"

4. **Copy and save the token**: 
   - **Important**: Copy the token immediately and store it securely
   - You won't be able to see it again after leaving the page
   - Treat this token like a password - don't share it or commit it to version control

## Usage

### Command Syntax

**Linux:**
```bash
./dns-updater-linux-amd64 -token=<API_TOKEN> -domain=<DOMAIN> -record=<RECORD_NAME> [OPTIONS]
```

### Required Flags

- `-token`: Your DigitalOcean API token
- `-domain`: The domain name (e.g., `example.com`)
- `-record`: The DNS record name (e.g., `www`, `api`, or `@` for root domain)

### Optional Flags

- `-type`: DNS record type (default: `A`)
- `-help`: Show help message

### Examples

#### Update an API subdomain
```bash
./dns-updater-linux-amd64 -token=dop_v1_abc123... -domain=example.com -record=api
```

#### Update the root domain
```bash
./dns-updater-linux-amd64 -token=dop_v1_abc123... -domain=example.com -record=@
```

#### Update an AAAA record (IPv6)
```bash
./dns-updater-linux-amd64 -token=dop_v1_abc123... -domain=example.com -record=www -type=AAAA
```

#### Show help
```bash
./dns-updater-linux-amd64 -help
```

## Environment Variables (Alternative)

For security, you can set your API token as an environment variable instead of passing it as a flag:

```bash
export DO_TOKEN="dop_v1_abc123..."
./dns-updater-linux-amd64 -token=$DO_TOKEN -domain=example.com -record=www
```

## Common Use Cases

### Dynamic DNS Updates
Use this tool in a script or cron job to automatically update your DNS records when your IP address changes:

**Linux:**
```bash
# Update DNS record with current public IP (automatically detected)
./dns-updater-linux-amd64 -token=dop_v1_abc123... -domain=example.com -record=home
```

**Linux Crontab example (runs every 30 minutes):**
```bash
# Edit your crontab
crontab -e

# Add this line to run every 30 minutes
*/30 * * * * /path/to/dns-updater-linux-amd64 -token=dop_v1_abc123... -domain=example.com -record=home >/dev/null 2>&1
```

## Troubleshooting

### "Domain not found" error
- Verify the domain is added to your DigitalOcean account
- Check that the domain name is spelled correctly
- Ensure your API token has the correct permissions

### "Record not found" error
- Verify the DNS record exists in your DigitalOcean control panel
- Check the record name and type are correct
- For root domain records, use `@` as the record name

### Authentication errors
- Verify your API token is correct and hasn't expired
- Ensure the token has read and write permissions
- Check that there are no extra spaces or characters in the token

## Dependencies

- [github.com/digitalocean/godo](https://github.com/digitalocean/godo) - Official DigitalOcean Go client library
- [golang.org/x/oauth2](https://golang.org/x/oauth2) - OAuth2 authentication

## License

This project is provided as-is for educational and personal use.