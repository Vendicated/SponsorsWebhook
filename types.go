package main

const (
	ActionTypeCreated             = "created"
	ActionTypeCancelled           = "cancelled"
	ActionTypePendingCancellation = "pending_cancellation"
	ActionTypePendingTierChange   = "pending_tier_change"
	ActionTypeTierChanged         = "tier_changed"
)

type Tier struct {
	IsOneTime             bool `json:"is_one_time"`
	MonthlyPriceInDollars int  `json:"monthly_price_in_dollars"`
}

type SponsorShipEvent struct {
	Action      string `json:"action"`
	Sponsorship struct {
		PrivacyLevel string `json:"privacy_level"`
		Sponsor      struct {
			AvatarUrl string `json:"avatar_url"`
			HtmlUrl   string `json:"html_url"`
			UserName  string `json:"login"`
		} `json:"sponsor"`
		Tier Tier `json:"tier"`
	} `json:"sponsorship"`
	EffectiveDate string `json:"effective_date"`
	Changes       struct {
		Tier struct {
			From Tier `json:"from"`
		} `json:"tier"`
	} `json:"changes"`
}

type DiscordWebhookPayload struct {
	Content         string `json:"content"`
	AvatarUrl       string `json:"avatar_url,omitempty"`
	Username        string `json:"username,omitempty"`
	AllowedMentions struct {
		Parse []string `json:"parse"`
	} `json:"allowed_mentions"`
}
