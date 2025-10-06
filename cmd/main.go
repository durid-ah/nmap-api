package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/Ullaakut/nmap/v3"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	scanner, err := nmap.NewScanner(
		ctx,
		nmap.WithTargets("192.168.1.*"),
		nmap.WithPingScan(),
	)

	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if len(*warnings) > 0 {
		log.Printf("run finished with warnings: %s\n", *warnings) // Warnings are non-critical errors from nmap.
	}
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		
		if len(host.Hostnames) == 0 || len(host.Addresses) == 0 {
			continue
		}
		
		fmt.Printf("Hostname %s IP: %s \n", host.Hostnames[0], host.Addresses[0])
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}