package server

import (
	"net/http"
	"route256/cart/internal/infra/utils"
)

func (s *Server) DeleteItem(writer http.ResponseWriter, request *http.Request) {
	rawUserID := request.PathValue("user_id")
	userID, err := utils.PrepareID(rawUserID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidUserID, "", http.StatusBadRequest)
		return
	}

	rawSkuID := request.PathValue("sku_id")
	skuID, err := utils.PrepareID(rawSkuID)
	if err != nil {
		utils.WriteErrorToResponse(writer, request, ErrInvalidSKU, "", http.StatusBadRequest)
		return
	}

	_, err = s.cartService.DeleteItem(request.Context(), userID, skuID)
	if err != nil {
		// тут ошибки могут быть из-за невалидных ID а они проверены раньше. Поэтому просто лог и ответ ОК.
		utils.WriteErrorToLog(request, err, "")
	}

	utils.WriteStatusToResponse(writer, request, "", http.StatusNoContent)
}
