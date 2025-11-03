package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/billbatista/digitalocean-dns-updater/ipgrabber"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

// TokenSource implements oauth2.TokenSource interface for DigitalOcean API
type TokenSource struct {
	AccessToken string
}

// Token returns an oauth2.Token with the DigitalOcean API token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	ctx := context.Background()

	var (
		apiToken   = flag.String("token", "", "DigitalOcean API token (required)")
		domain     = flag.String("domain", "", "Domain name (e.g., example.com) (required)")
		recordName = flag.String("record", "", "DNS record name (e.g., www or @ for root) (required)")
		recordType = flag.String("type", "A", "DNS record type (default: A)")
		help       = flag.Bool("help", false, "Show help message")
	)

	flag.Parse()

	if *help {
		fmt.Println("DigitalOcean DNS Updater")
		fmt.Println("Updates a DNS record's IP address on DigitalOcean")
		fmt.Println()
		fmt.Println("Usage:")
		flag.PrintDefaults()
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  go run main.go -token=your_api_token -domain=example.com -record=www")
		return
	}

	if *apiToken == "" || *domain == "" || *recordName == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing required flags\n\n")
		fmt.Println("Required flags:")
		fmt.Println("  -token: DigitalOcean API token")
		fmt.Println("  -domain: Domain name")
		fmt.Println("  -record: DNS record name")
		fmt.Println()
		fmt.Println("Use -help for more information")
		os.Exit(1)
	}

	tokenSource := &TokenSource{
		AccessToken: *apiToken,
	}

	oauthClient := oauth2.NewClient(ctx, tokenSource)

	client := godo.NewClient(oauthClient)

	ipgrabber := &ipgrabber.IFConfig{}
	newIP, err := ipgrabber.GetPublicIP()
	if err != nil {
		slog.ErrorContext(ctx, "Error getting public IP", slog.Any("err", err))
		return
	}

	err = updateDNSRecord(ctx, client, *domain, *recordName, *recordType, newIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error updating DNS record", "recordName", *recordName, slog.Any("err", err))
		return
	}

	slog.InfoContext(
		ctx,
		"Successfully updated DNS record",
		"recordType",
		*recordType,
		"recordName",
		*recordName,
		"domain",
		*domain,
		"newIP",
		newIP,
	)
}

func updateDNSRecord(ctx context.Context, client *godo.Client, domain, recordName, recordType, newIP string) error {
	records, _, err := client.Domains.Records(ctx, domain, nil)
	if err != nil {
		return fmt.Errorf("failed to list DNS records for domain %s: %w", domain, err)
	}

	var targetRecord *godo.DomainRecord
	for i := range records {
		record := &records[i]
		if record.Name == recordName && record.Type == recordType {
			targetRecord = record
			break
		}
	}

	if targetRecord == nil {
		return fmt.Errorf("DNS record '%s' of type '%s' not found in domain '%s'", recordName, recordType, domain)
	}

	editRequest := &godo.DomainRecordEditRequest{
		Type: recordType,
		Name: recordName,
		Data: newIP,
		TTL:  targetRecord.TTL, // Keep the existing TTL
	}

	_, _, err = client.Domains.EditRecord(ctx, domain, targetRecord.ID, editRequest)
	if err != nil {
		return fmt.Errorf("failed to update DNS record: %w", err)
	}

	return nil
}
