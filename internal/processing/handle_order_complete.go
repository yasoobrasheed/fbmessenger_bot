package processing

import (
	"encoding/json"
	"fbmessenger_bot/internal/bot/send_message"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Order struct {
	ID           string  `json:"id"`
	SenderId     string  `json:"sender_id"`
	ProductName  string  `json:"product_name"`
	Price        float64 `json:"price"`
	DeliveryDate string  `json:"delivery_date"`
	DeliveryUrl  string  `json:"delivery_url"`
}

func HandleOrderComplete(w http.ResponseWriter, r *http.Request) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var orderData Order

	err = json.Unmarshal([]byte(requestBody), &orderData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	orderDetailsText := constructOrderDetailsText(orderData)
	askForReviewText := `Thank you so much for doing business with us, and we'd love to hear about your experience.
						 You can write us a review right in this message window! As an added bonus, we'll send you a 20% off discount code.`

	send_message.HandleSendMessageText(orderDetailsText, orderData.SenderId)
	send_message.HandleSendMessageText(askForReviewText, orderData.SenderId)
}

func constructOrderDetailsText(orderData Order) string {
	orderDetailsText := `
		Order Completed:

		` + orderData.ProductName + ` 

		$` + fmt.Sprintf("%.2f", orderData.Price) + `

		Estimated delivery date: ` + orderData.DeliveryDate + `

		See delivery details at: ` + orderData.DeliveryUrl + "/" + orderData.ID + `
	`
	orderDetailsText = strings.ReplaceAll(orderDetailsText, "\t", "")

	return orderDetailsText
}
