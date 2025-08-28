package server

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/infra/utils"
)

func (s *Server) GetCart(writer http.ResponseWriter, request *http.Request) {

	rawUserID := request.PathValue("user_id")

	userID, err := utils.PrepareID(rawUserID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	cart, err := s.cartService.GetItemsByUserID(request.Context(), userID)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusNotFound)
		return
	}

	reportCart := ReportCart{
		Items:      []ItemInСart{},
		TotalPrice: cart.TotalPrice}

	for _, item := range cart.Items {
		reportCart.Items = append(reportCart.Items, ItemInСart{
			SKU:   item.SKU,
			Count: item.Count,
			Name:  item.Name,
			Price: item.Price,
		})

	}

	rawResponce, err := json.Marshal(reportCart)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(rawResponce)

	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)

	}
}
