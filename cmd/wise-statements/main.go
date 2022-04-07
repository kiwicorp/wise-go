package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/jinzhu/now"
	wise "github.com/kiwicorp/wise-go/pkg/wise"
	"github.com/kiwicorp/wise-go/pkg/wisesca"
)

// flags
var (
	tokenFlag = flag.String(
		"token",
		"",
		"API token to use. Prefix with an '@' to load from the corresponding file.",
	)
	keyFlag = flag.String(
		"key",
		"",
		"PEM-encoded RSA 2048 private key file to use for Strong Customer Authentication.",
	)
	dryRunFlag = flag.Bool(
		"dry-run",
		false,
		"Performs a dry-run of downloading the statements, without actually downloading them.",
	)
	verboseFlag = flag.Bool(
		"verbose",
		false,
		"Enable verbose messages.",
	)
	startFlag = flag.String(
		"start",
		"",
		"Start of the interval for downloading statements. Format: "+time.RFC3339,
	)
	endFlag = flag.String(
		"end",
		"",
		"End of the interval for downloading statements. Format: "+time.RFC3339,
	)
	businessFlag = flag.Bool(
		"business",
		false,
		"Get statements from business profiles.",
	)
	personalFlag = flag.Bool(
		"personal",
		false,
		"Get statements from personal profiles.",
	)
)

var (
	wiseKey   *rsa.PrivateKey
	wiseToken string
)

func main() {
	os.Exit(run())
}

func run() int {
	flag.Parse()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log("hello, this is wise-statements")

	// ----- load values from flags -----
	logVerbose("loading token")
	if strings.HasPrefix(*tokenFlag, "@") {
		filename := strings.TrimPrefix(*tokenFlag, "@")
		token, err := os.ReadFile(filename)
		if err != nil {
			log("failed to load token from %s: %s", filename, err.Error())
			return 1
		}
		wiseToken = strings.TrimSpace(string(token))
		logVerbose("loaded token from file: %s", filename)
	} else {
		wiseToken = *tokenFlag
		logVerbose("loaded token from flag")
	}

	logVerbose("loading statement start and end time")
	var intervalStart time.Time
	if *startFlag != "" {
		var err error
		intervalStart, err = time.Parse(time.RFC3339, *startFlag)
		if err != nil {
			log("failed to parse start time: %s", err.Error())
			return 1
		}
	} else {
		intervalStart = now.BeginningOfMonth()
	}
	var intervalEnd time.Time
	if *endFlag != "" {
		var err error
		intervalEnd, err = time.Parse(time.RFC3339, *endFlag)
		if err != nil {
			log("failed to parse end time: %s", err.Error())
			return 1
		}
	} else {
		intervalEnd = time.Now()
	}

	logVerbose("loading private key")
	var b []byte
	b, err := os.ReadFile(*keyFlag)
	if err != nil {
		log("failed to load private key: %s", err.Error())
		return 1
	}
	var p *pem.Block
	p, _ = pem.Decode(b)
	wiseKey, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		log("failed to load private key: %s", err.Error())
		return 1
	}
	// ----------------------------------

	sca := wisesca.NewWithPersonalToken(&http.Client{}, wiseKey)
	service := wise.NewDefaultProduction(sca, wiseToken)
	wg := new(sync.WaitGroup)

	logVerbose("getting profiles")
	profiles, err := service.ListProfiles(ctx)
	if err != nil {
		switch err := err.(type) {
		case wise.Error:
			log("failed to get profiles: %s", err.Message)
			return 1
		default:
			log("failed to get profiles: %s", err.Error())
			return 1
		}
	}
	// map from profile name (type-full name) to profile
	profileMap := make(map[int]wise.Profile, len(profiles.Profiles))
	for _, p := range profiles.Profiles {
		switch p.Type {
		case string(wise.ProfileTypeBusiness):
			if *businessFlag {
				profileMap[p.ID] = p
			}
		case string(wise.ProfileTypePersonal):
			if *personalFlag {
				profileMap[p.ID] = p
			}
		}
	}

	logVerbose("getting balances")
	// map from profile id to balances
	balanceMapRWLock := new(sync.RWMutex)
	balanceMap := make(map[int][]wise.Balance, len(profileMap))
	for _, p := range profileMap {
		profileID := p.ID
		wg.Add(1)
		go func() {
			defer wg.Done()
			balances, err := service.GetBalances(ctx, &wise.GetBalancesRequest{
				ProfileID: profileID,
			})
			if err != nil {
				switch err := err.(type) {
				case wise.Error:
					log("failed to get balances: %s", err.Message)
					stop()
					return
				default:
					log("failed to get balances: %s", err.Error())
					stop()
					return
				}
			}
			balanceMapRWLock.Lock()
			balanceMap[profileID] = balances.Balances
			balanceMapRWLock.Unlock()
		}()
	}
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return 1
	}

	logVerbose("getting statements")
	for profileID, balances := range balanceMap {
		pID := profileID
		profile := profileMap[pID]
		pName := fmt.Sprintf("%s--%s", profile.Type, strings.ReplaceAll(profile.FullName, " ", "_"))
		for _, balance := range balances {
			bID := balance.ID
			bName := balance.Currency
			wg.Add(1)
			go func() {
				defer wg.Done()
				statementName := fmt.Sprintf("%s--%s", pName, bName)
				log("downloading %s", statementName)
				if *dryRunFlag {
					if *verboseFlag {
						log("dry-run, not downloading statements.")
					}
					return
				}
				pdf, err := service.GetStatementPDF(ctx, &wise.GetStatementPDFRequest{
					ProfileID:     pID,
					BalanceID:     bID,
					IntervalStart: intervalStart,
					IntervalEnd:   intervalEnd,
					Type:          wise.StatementTypeCompact,
				})
				if err != nil {
					switch err := err.(type) {
					case wise.Error:
						log("failed to get statement: %s", err.Message)
						stop()
						return
					default:
						log("failed to get statement: %s", err.Error())
						stop()
						return
					}
				}
				if err := ctx.Err(); err != nil {
					return
				}
				filename := fmt.Sprintf("%s.pdf", statementName)
				logVerbose("writing %s to %s", statementName, filename)
				if err := os.WriteFile(filename, pdf.Data, 0644); err != nil {
					log("failed to write statement: %s", err.Error())
					stop()
					return
				}
				log("downloaded %s", statementName)
				output(filename)
			}()
		}
	}
	wg.Wait()
	if err := ctx.Err(); err != nil {
		return 1
	}
	return 0
}

// Log a message.
func log(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "---> %s\n", fmt.Sprintf(format, args...))
}

// Log a verbose message.
func logVerbose(format string, args ...any) {
	if *verboseFlag {
		fmt.Fprintf(os.Stderr, "---> %s\n", fmt.Sprintf(format, args...))
	}
}

// Output a value.
func output(format string, args ...any) {
	fmt.Fprintf(os.Stdout, "%s\n", fmt.Sprintf(format, args...))
}
