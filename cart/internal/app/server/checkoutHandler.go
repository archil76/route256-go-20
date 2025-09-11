package server

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/infra/utils"
)

func (s *Server) Checkout(writer http.ResponseWriter, request *http.Request) {
	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	orderID, err := s.cartService.Checkout(request.Context(), userID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusBadRequest)
		return
	}

	checkoutResponse := CheckoutResponse{OrderID: orderID}

	rawResponse, err := json.Marshal(checkoutResponse)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get Checkout. marshalling error", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(rawResponse)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't Checkout. marshalling error", http.StatusBadRequest)

	}

	utils.WriteStatusToResponse(writer, request, "", http.StatusOK)
}
