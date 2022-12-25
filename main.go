package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"

	"github.com/getlantern/systray"
	"golang.design/x/clipboard"
)

var (
	//go:embed assets/icon.ico
	iconData []byte

	//go:embed rules.json
	rulesData []byte

	rules  Rules
	config Config
)

func saveConfigSilent() {
	if err := config.Save(); err != nil {
		log.Printf("failed to save config: %v\n", err)
	}
}

func main() {
	if err := config.Read(); err != nil {
		log.Panicf("failed to read config: %v\n", err)
	}

	if err := json.Unmarshal(rulesData, &rules); err != nil {
		log.Panicf("failed to parse rules: %v\n", err)
	}

	if err := clipboard.Init(); err != nil {
		log.Panicf("failed to initialize clipboard: %v\n", err)
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("AntiLinkTrack")
	systray.SetIcon(iconData)

	statsItem := systray.AddMenuItem(fmt.Sprintf("Cleaned URLs: %d", config.CleanedUrls), "")
	statsItem.Disable()

	systray.AddSeparator()

	allowReferralMarketing := systray.AddMenuItemCheckbox("Allow Referral Marketing", "Allow Referral Marketing parts in URLs", false)
	systray.AddSeparator()
	quit := systray.AddMenuItem("Quit", "Quit the application")

	go func() {
		for {
			select {
			case <-allowReferralMarketing.ClickedCh:
				if allowReferralMarketing.Checked() {
					allowReferralMarketing.Uncheck()
				} else {
					allowReferralMarketing.Check()
				}
				config.AllowReferralMarketing = allowReferralMarketing.Checked()
				saveConfigSilent()
			case <-quit.ClickedCh:
				systray.Quit()
			}
		}
	}()

	clipboardChannel := clipboard.Watch(context.Background(), clipboard.FmtText)
	for clipboardData := range clipboardChannel {
		contentStr := string(clipboardData)

		cleaned, err := rules.CleanUrl(contentStr, allowReferralMarketing.Checked())
		if err != nil {
			log.Printf("failed to clean url: %v\n", err)
			continue
		}

		if cleaned != contentStr {
			log.Printf("cleaned url: %s\n", cleaned)
			clipboard.Write(clipboard.FmtText, []byte(cleaned))

			config.CleanedUrls++
			statsItem.SetTitle(fmt.Sprintf("Cleaned URLs: %d", config.CleanedUrls))
			saveConfigSilent()
		}
	}
}

func onExit() {
	// clean up here
}
