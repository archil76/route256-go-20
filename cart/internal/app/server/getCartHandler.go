package server

import (
	"encoding/json"
	"net/http"
	"route256/cart/internal/infra/utils"
	"sort"
)

func (s *Server) GetCart(writer http.ResponseWriter, request *http.Request) {
	rawUserID := request.PathValue("user_id")
	userID, err := utils.PrepareID(rawUserID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	modelReportCart, err := s.cartService.GetItemsByUserID(request.Context(), userID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrOther, "", http.StatusNotFound)
		return
	}

	totalPrice32 := int32(modelReportCart.TotalPrice) //nolint:gosec

	reportCart := ReportCart{
		Items:      []ItemInСart{},
		TotalPrice: totalPrice32}

	for _, item := range modelReportCart.Items {
		count32 := int32(item.Count) //nolint:gosec
		price32 := int32(item.Price) //nolint:gosec

		reportCart.Items = append(reportCart.Items, ItemInСart{
			SKU:   item.SKU,
			Count: count32,
			Name:  item.Name,
			Price: price32,
		})
	}

	sort.Slice(reportCart.Items, func(i, j int) bool { return reportCart.Items[i].SKU < reportCart.Items[j].SKU })

	rawResponse, err := json.Marshal(reportCart)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(rawResponse)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, err, "can't get cart. marshalling error", http.StatusBadRequest)

	}
}
