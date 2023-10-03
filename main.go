package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/pokt-foundation/dwh-client/client"
	"github.com/pokt-foundation/portal-db/v2/types"
	"github.com/pokt-foundation/utils-go/environment"
)

const (
	dwhAPIKey = "DWH_API_KEY"
	dwhURL    = "DWH_API_URL"
)

func gatherOptions() client.Config {
	return client.Config{
		APIKey: environment.MustGetString(dwhAPIKey),
		URL:    environment.MustGetString(dwhURL),
	}
}

// Simple main.go file to test the clint
func main() {
	config := gatherOptions()

	dwhClient, err := client.NewDWHClient(config)
	if err != nil {
		panic(err)
	}

	response, err := dwhClient.GetTotalRelaysForPortalAppIDs(context.Background(), client.GetTotalRelaysForPortalAppIDsParams{
		From: time.Date(2023, 8, 4, 0, 0, 0, 0, time.UTC),
		To:   time.Date(2023, 9, 4, 0, 0, 0, 0, time.UTC),
		PortalAppIDs: []types.PortalAppID{
			"44c0823fbdf0aed3fa2d6357",
			"3742b06f9e13c9ea22a8d599",
		},
	})
	if err != nil {
		panic(err)
	}

	Plog("GOT RESPONSE", response)
}

func Plog(args ...interface{}) {
	for _, arg := range args {
		var prettyJSON bytes.Buffer
		jsonArg, _ := json.Marshal(arg)
		str := string(jsonArg)
		_ = json.Indent(&prettyJSON, []byte(str), "", "    ")
		output := prettyJSON.String()

		fmt.Println(output)
	}
}
