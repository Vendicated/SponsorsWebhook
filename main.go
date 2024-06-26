package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

func fprintln(w io.Writer, a ...any) {
	_, _ = fmt.Fprintln(w, a...)
}

var usernameRulesRe = regexp.MustCompile("(?i)clyde|discord|everyone|here")

func sendWebhook(body DiscordWebhookPayload) bool {
	body.Username = usernameRulesRe.ReplaceAllString(body.Username, "[banned]")[:80]
	body.AllowedMentions.Parse = []string{}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		fmt.Println("Failed to json encode webhook body:", err)
		return false
	}

	res, err := http.Post(WebhookUrl, "application/json", &buf)
	if err != nil {
		fmt.Println("Failed to post to", WebhookUrl+":", err)
		return false
	}

	return res.StatusCode >= 200 && res.StatusCode < 300
}

func handleWebhook(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	sig := req.Header.Get("HTTP_X_HUB_SIGNATURE_256")
	if sig == "" {
		w.WriteHeader(http.StatusForbidden)
		fprintln(w, "Missing HTTP_X_HUB_SIGNATURE_256")
		return
	}

	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fprintln(w, "Failed to read request body")
		return
	}

	if !verifySignature(bodyBytes, sig) {
		w.WriteHeader(http.StatusForbidden)
		fprintln(w, "Wrong signature")
		return
	}

	var sponsorShipEvent SponsorShipEvent
	if err = json.Unmarshal(bodyBytes, &sponsorShipEvent); err != nil {
		fmt.Println("Failed to parse request body", err)
		w.WriteHeader(http.StatusBadRequest)
		fprintln(w, "Failed to parse request body")
		return
	}

	message := DiscordWebhookPayload{
		Username:  sponsorShipEvent.Sponsorship.Sponsor.UserName,
		AvatarUrl: sponsorShipEvent.Sponsorship.Sponsor.AvatarUrl,
	}

	priceInDollar := sponsorShipEvent.Sponsorship.Tier.MonthlyPriceInDollars
	sponsorType := "per month"
	if sponsorShipEvent.Sponsorship.Tier.IsOneTime {
		sponsorType = "one time"
	}
	oldPriceInDollar := sponsorShipEvent.Changes.Tier.From.MonthlyPriceInDollars
	oldSponsorType := "per momth"
	if sponsorShipEvent.Sponsorship.Tier.IsOneTime {
		oldSponsorType = "one time"
	}
	sponsorUserLink := fmt.Sprintf("[%s](<%s>)", sponsorShipEvent.Sponsorship.Sponsor.UserName, sponsorShipEvent.Sponsorship.Sponsor.HtmlUrl)

	switch sponsorShipEvent.Action {
	case ActionTypeCreated:
		message.Content = fmt.Sprintf(
			"New %d$ %s sponsor: %s",
			priceInDollar,
			sponsorType,
			sponsorUserLink,
		)
	case ActionTypeCancelled:
		message.Content = fmt.Sprintf(
			"%s cancelled their %d$ %s sponsorship",
			sponsorUserLink,
			priceInDollar,
			sponsorType,
		)
	case ActionTypePendingCancellation:
		message.Content = fmt.Sprintf(
			"%s schedulled a cancellation for their %d$ %s sponsorship at %s",
			sponsorUserLink,
			priceInDollar,
			sponsorType,
			sponsorShipEvent.EffectiveDate,
		)
	case ActionTypeTierChanged:
		message.Content = fmt.Sprintf(
			"%s changed their tier from %d$ %s to %d$ %s",
			sponsorUserLink,
			oldPriceInDollar,
			oldSponsorType,
			priceInDollar,
			sponsorType,
		)
	case ActionTypePendingTierChange:
		message.Content = fmt.Sprintf(
			"%s schedulled a change of their tier from %d$ %s to %d$ %s at %s",
			sponsorUserLink,
			oldPriceInDollar,
			oldSponsorType,
			priceInDollar,
			sponsorType,
			sponsorShipEvent.EffectiveDate,
		)
	default:
		w.WriteHeader(200)
		fprintln(w, "Ok")
		return
	}

	succeeded := sendWebhook(message)
	if !succeeded {
		w.WriteHeader(http.StatusBadGateway)
		fprintln(w, "Failed to execute webhook")
	}

	w.WriteHeader(http.StatusOK)
	fprintln(w, "Ok")
}

func verifySignature(message []byte, messageMAC string) bool {
	mac := hmac.New(sha256.New, Secret)
	mac.Write(message)
	expectedMAC := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(messageMAC), []byte(expectedMAC))
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/webhook", handleWebhook)

	fmt.Println("Listening on port", Port)
	panic(http.ListenAndServe(":"+Port, nil))
}