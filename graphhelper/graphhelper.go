package graphhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/gptscript-ai/go-gptscript"
)

const (
	clientID = "98a09e6b-7242-4977-a1e4-a79e5f083832"
	tenantID = "369895fd-4335-4606-b433-6ab084d5bd79"
	scopes   = "User.Read,Mail.Read,Mail.Send"
)

type GraphHelper struct {
	deviceCodeCredential *azidentity.DeviceCodeCredential
	graphUserScopes      []string
}

func NewGraphHelper() *GraphHelper {
	g := &GraphHelper{
		graphUserScopes: strings.Split(scopes, ","),
	}
	return g
}

func (g *GraphHelper) InitializeGraphForUserAuth() error {
	credential, err := azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
		ClientID: clientID,
		TenantID: tenantID,
		UserPrompt: func(ctx context.Context, message azidentity.DeviceCodeMessage) error {
			gs, err := gptscript.NewGPTScript(gptscript.GlobalOptions{})
			if err != nil {
				return fmt.Errorf("error creating GPTScript client: %w", err)
			}
			defer gs.Close()

			sysPromptIn, err := json.Marshal(struct {
				Message   string `json:"message"`
				Fields    string `json:"fields"`
				Sensitive string `json:"sensitive"`
			}{
				Message:   message.Message,
				Fields:    "Press enter to continue ...",
				Sensitive: "false",
			})
			if err != nil {
				return fmt.Errorf("error marshaling sysPromptIn: %w", err)
			}

			run, err := gs.Run(ctx, "sys.prompt", gptscript.Options{Input: string(sysPromptIn)})
			if err != nil {
				return fmt.Errorf("error running sys.prompt: %w", err)
			}

			_, err = run.Text()
			if err != nil {
				return fmt.Errorf("error getting the result of sys.prompt: %w", err)
			}
			return nil
		},
	})
	if err != nil {
		return err
	}

	g.deviceCodeCredential = credential
	return nil
}

func (g *GraphHelper) GetUserToken() (azcore.AccessToken, error) {
	return g.deviceCodeCredential.GetToken(context.Background(), policy.TokenRequestOptions{
		Scopes: g.graphUserScopes,
	})
}
