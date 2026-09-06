package main

import "fmt"

type PaymentEvent struct {
	ID, Actor, Action string
	AmountCents       int64
}
type Decision struct {
	Allowed bool
	Reason  string
}

func DecideRisk(event PaymentEvent, online map[string]bool) Decision {
	if !online[event.Actor] {
		return Decision{false, "actor is offline; require an online reviewer"}
	}
	if event.AmountCents > 100000 {
		return Decision{false, "amount exceeds review threshold"}
	}
	return Decision{true, "reviewer is online and amount is within threshold"}
}

func RunWorkflow(c *Client, event PaymentEvent) (Decision, error) {
	const channel = "fintech-ops"
	var members []struct {
		ID string `json:"id"`
	}
	if err := c.Presence(channel, &members); err != nil {
		return Decision{}, err
	}
	online := make(map[string]bool, len(members))
	for _, m := range members {
		online[m.ID] = true
	}
	d := DecideRisk(event, online)
	if err := c.Publish(channel, "payment.audit", event.Actor, map[string]any{"payment_id": event.ID, "action": event.Action, "allowed": d.Allowed, "reason": d.Reason}); err != nil {
		return Decision{}, fmt.Errorf("publish audit notification: %w", err)
	}
	return d, nil
}
