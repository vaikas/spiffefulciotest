/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spiffe/go-spiffe/v2/svid/jwtsvid"
	"github.com/spiffe/go-spiffe/v2/workloadapi"

	"github.com/sigstore/cosign/cmd/cosign/cli/fulcio"
	"github.com/sigstore/cosign/pkg/providers"

	_ "github.com/sigstore/cosign/pkg/providers/all"
)

const (
	spireAgentURL       = "unix:///run/spire/sockets/agent.sock"
	defaultOIDCIssuer   = "https://oauth2.sigstore.dev/auth"
	defaultOIDCClientID = "sigstore"
	fulcioAddress       = "http://fulcio.fulcio-system.svc"
)

func main() {
	i := 0
	for {
		if i > 0 {
			time.Sleep(10 * time.Second)
		}
		i++
		ctx := context.Background()
		client, err := workloadapi.New(ctx, workloadapi.WithAddr("unix:///run/spire/sockets/agent.sock"))
		if err != nil {
			fmt.Printf("Error creating workloadapi client: %s\n", err)
			continue
		}
		jwt, err := client.FetchJWTSVID(ctx, jwtsvid.Params{
			Audience: "sigstore",
		})
		if err != nil {
			fmt.Printf("Error fetching jtwsvid: %s\n", err)
			continue
		}
		fmt.Println("")
		fmt.Println(jwt)

		if !providers.Enabled(ctx) {
			fmt.Println("no auth provider for fulcio is enabled")
			continue
		}

		tok, err := providers.Provide(ctx, defaultOIDCClientID)
		if err != nil {
			fmt.Printf("Error getting provider: %s\n", err)
			continue
		}
		fmt.Println("Signing with fulcio ...")

		fClient, err := fulcio.NewClient(fulcioAddress)
		if err != nil {
			fmt.Printf("Error creating fulcio client: %v\n", err)
			continue
		}
		k, err := fulcio.NewSigner(ctx, tok, defaultOIDCIssuer, defaultOIDCClientID, "", fClient)
		if err != nil {
			fmt.Printf("Error getting a new fulcio signer %v\n", err)
			continue
		}
		fmt.Printf("Got a new Fulcio Signer:\n%+v", k)
	}
}
