package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
	"github.com/theopenlane/go-client/graphclient"
)

func main() {
	os.Setenv("OPENLANE_API_TOKEN", os.Getenv("OPENLANE_API_TOKEN2"))
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}
	client, err := openlane.New(cfg)
	if err != nil {
		fmt.Println("client:", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	abhishekGroup := "01KTY1J6WEC5EXBVSJJWW6QP68"
	gregGroup := "01M1EN3JHTHT3X5P2FJRTGBH00"
	controls := map[string]string{
		"12.6.1":   "01KTY1WZBS3V5C6ZY2YJ3ANMED",
		"12.6.2":   "01KTY1WZBZZ17NK5C7CJXGJW97",
		"12.6.3":   "01KTY1WZBN5TR8WPM8ECPYD0EQ",
		"12.6.3.1": "01KTY1WZBZZ17NK5C7C9W01CFB",
		"12.6.3.2": "01KTY1WZBMD6DGKY1DDQTGEFND",
	}

	for ref, id := range controls {
		_, err := client.UpdateControl(ctx, id, graphclient.UpdateControlInput{
			ControlOwnerID: &abhishekGroup,
			DelegateID:     &gregGroup,
		})
		if err != nil {
			fmt.Printf("%s UPDATE FAIL: %s\n", ref, openlane.Redact(err.Error()))
			continue
		}
		got, err := client.GetControlByID(ctx, id)
		if err != nil {
			fmt.Printf("%s VERIFY GET FAIL: %s\n", ref, openlane.Redact(err.Error()))
			continue
		}
		c := got.Control
		ok := openlane.Deref(c.ControlOwnerID) == abhishekGroup && openlane.Deref(c.DelegateID) == gregGroup
		fmt.Printf("%s: owner=%s delegate=%s verified=%v\n", ref, openlane.Deref(c.ControlOwnerID), openlane.Deref(c.DelegateID), ok)
	}
}
