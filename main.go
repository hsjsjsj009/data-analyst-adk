package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

type ExtendedWebUISubLauncher struct {
	web.Sublauncher
}

func NewExtendedWebUISubLauncher() *ExtendedWebUISubLauncher {
	return &ExtendedWebUISubLauncher{
		Sublauncher: webui.NewLauncher(),
	}

}

func (l *ExtendedWebUISubLauncher) CommandLineSyntax() string {
	return l.Sublauncher.CommandLineSyntax()
}

func (l *ExtendedWebUISubLauncher) SimpleDescription() string {
	return l.Sublauncher.SimpleDescription()
}

// Keyword returns the command-line flag for this sublauncher (optional usage)
func (l *ExtendedWebUISubLauncher) Keyword() string {
	return l.Sublauncher.Keyword()
}

// Parse handles command-line arguments (can be empty if not needed)
func (l *ExtendedWebUISubLauncher) Parse(args []string) ([]string, error) {
	return l.Sublauncher.Parse(args)
}

// SetupSubrouters is where you ADD YOUR NEW PATH
func (l *ExtendedWebUISubLauncher) SetupSubrouters(router *mux.Router, config *launcher.Config) error {
	// Add your new HTTP path here
	router.HandleFunc("/my-new-path", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the new path!")
	}).Methods("GET")

	return l.Sublauncher.SetupSubrouters(router, config)

}

// UserMessage prints a startup message (optional)
func (l *ExtendedWebUISubLauncher) UserMessage(webURL string, printer func(v ...any)) {
	l.Sublauncher.UserMessage(webURL, printer)
}

const (
	authorizationCodeStateName string = "authorization_code"
)

type OauthMCPCustomHTTPTransport struct {
	http.RoundTripper
}

func (t *OauthMCPCustomHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	newReq := req.Clone(ctx)
	toolContext, ok := ctx.(tool.Context)
	if ok {
		toolContextState := toolContext.State()
		val, err := toolContextState.Get(authorizationCodeStateName)
		if err == nil {
			valStr, ok := val.(string)
			if ok {
				newReq.Header.Set("Authorization", "Bearer "+valStr)
			}
		}
	}

	transport := t.RoundTripper
	if transport == nil {
		transport = http.DefaultTransport
	}

	// 4. Execute the request
	return transport.RoundTrip(newReq)
}

func main() {
	ctx := context.Background()

	// Initialize the Gemini Model
	// Ensure GOOGLE_APPLICATION_CREDENTIALS is set for Vertex AI access
	model, err := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}
	//
	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: &mcp.SSEClientTransport{
			HTTPClient: &http.Client{
				Transport: &OauthMCPCustomHTTPTransport{
					http.DefaultTransport,
				},
			},
			Endpoint: "https://mcp-toolbox-280946129258.asia-southeast1.run.app/mcp/sse",
		},
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	type AuthorizationCode struct {
		AuthorizationCode string
	}

	oauthCodePreparation, err := functiontool.New(functiontool.Config{
		Name:        "prepare_oauth_auth_code",
		Description: "You need to call this tool before executing oauth2 tools to prepare the authorization code",
	}, func(t tool.Context, args AuthorizationCode) (string, error) {
		state := t.State()
		err := state.Set(authorizationCodeStateName, args.AuthorizationCode)
		if err != nil {
			return "Fail", err
		}
		return "Success", nil
	})
	if err != nil {
		log.Fatalf("Failed to create tool: %v", err)
	}

	oauthCodeUserDataFetch, err := functiontool.New(functiontool.Config{
		Name:        "user_oauth_data",
		Description: "This tool is used to fetch oauth user data. This tool will use prepared authorization code",
	}, func(t tool.Context, args any) (string, error) {
		state := t.State()
		authCode, err := state.Get(authorizationCodeStateName)
		if err != nil {
			return "Fail", err
		}
		authCodeStr, ok := authCode.(string)
		if !ok {
			return "Fail", fmt.Errorf("auth code is not a string")
		}
		req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
		if err != nil {
			return "Fail", err
		}
		req.Header.Add("Authorization", "Bearer "+authCodeStr)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return "Fail", err
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			return "Fail", err
		}

		return string(content), nil
	})

	// Create the Agent with our custom tool
	myAgent, err := llmagent.New(llmagent.Config{
		Name:        "data_anaylst_agent",
		Description: "An agent that will be a data analyst to get insights from our datastore",
		Model:       model,
		Tools: []tool.Tool{
			oauthCodePreparation,
			oauthCodeUserDataFetch,
		},
		Toolsets: []tool.Toolset{
			mcpToolSet,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Configure the launcher (CLI, WebUI, etc.)
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(myAgent),
	}

	l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
