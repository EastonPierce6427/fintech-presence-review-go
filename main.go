package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	c, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}
	event := PaymentEvent{ID: "pay-1001", Actor: "analyst-7", Action: "approve", AmountCents: 4200}
	decision, err := RunWorkflow(c, event)
	if err != nil {
		log.Fatal(err)
	}
	result, _ := json.Marshal(map[string]any{"payment_id": event.ID, "allowed": decision.Allowed, "reason": decision.Reason})
	fmt.Println(string(result))
}
