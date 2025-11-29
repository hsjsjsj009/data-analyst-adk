package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/universal"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
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

	// Create the Agent with our custom tool
	myAgent, err := llmagent.New(llmagent.Config{
		Name:        "auth_agent",
		Description: "An agent that demonstrates OAuth integration.",
		Model:       model,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	// Configure the launcher (CLI, WebUI, etc.)
	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(myAgent),
	}

	apiLauncher := api.NewLauncher()
	uiLauncher := NewExtendedWebUISubLauncher()

	l := universal.NewLauncher(web.NewLauncher(apiLauncher, uiLauncher))
	//l := full.NewLauncher()
	if err := l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v", err)
	}
}
