package main

import (
	"encoding/json"
	"fmt"
	"graphtutorial/graphhelper"
	"log"
	"time"
)

type cred struct {
	Env       map[string]string `json:"env"`
	ExpiresAt time.Time         `json:"expiresAt"`
}

func main() {
	graphHelper := graphhelper.NewGraphHelper()

	initializeGraph(graphHelper)
	printCredential(graphHelper)
}

func initializeGraph(graphHelper *graphhelper.GraphHelper) {
	err := graphHelper.InitializeGraphForUserAuth()
	if err != nil {
		log.Panicf("Error initializing Graph for user auth: %v\n", err)
	}
}

func printCredential(graphHelper *graphhelper.GraphHelper) {
	token, err := graphHelper.GetUserToken()
	if err != nil {
		log.Panicf("Error getting user token: %v\n", err)
	}

	c := cred{
		Env:       map[string]string{"GPTSCRIPT_GRAPH_MICROSOFT_COM_BEARER_TOKEN": token.Token},
		ExpiresAt: token.ExpiresOn,
	}

	credJSON, err := json.Marshal(c)
	if err != nil {
		log.Panicf("Error marshaling credential: %v\n", err)
	}

	fmt.Print(credJSON)
}
