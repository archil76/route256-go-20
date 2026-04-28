package hw_8

import (
	"net/http"
	"strings"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestAddComment_Success(t provider.T) {
	t.Title("успешное создание комментария")

	const (
		userID      = int64(10)
		sku         = int64(10)
		commentText = "тестовый комментарий"
	)

	var (
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku)

	t.WithNewStep("Создание комментария", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID", commentID)

	t.WithNewStep("Получение созданного комментария", func(sCtx provider.StepCtx) {
		comment := &domain.Comment{}
		comment, statusCode = s.commentsClient.GetCommentByID(s.ctx, sCtx, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Equal(commentID, comment.ID)
		sCtx.Require().Equal(userID, comment.UserID)
		sCtx.Require().Equal(commentText, comment.Comment)
	})
}

func (s *Suite) TestAddComment_Negative(t provider.T) {
	t.Title("неудачное создание комментария")

	const (
		userID      = int64(11)
		sku         = int64(11)
		commentText = "тестовый комментарий"
	)

	t.WithNewParameters("userID", 0, "sku", sku)

	t.WithNewStep("Создание комментария без указания user_id", func(sCtx provider.StepCtx) {
		commentID, statusCode := s.commentsClient.AddComment(s.ctx, sCtx, 0, sku, commentText)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.StatusCode(sCtx, 0, int(commentID))
	})

	t.WithNewStep("Создание комментария без указания sku", func(sCtx provider.StepCtx) {
		commentID, statusCode := s.commentsClient.AddComment(s.ctx, sCtx, userID, 0, commentText)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.StatusCode(sCtx, 0, int(commentID))
	})

	t.WithNewStep("Создание комментария с запрещенным минимальным размером текста", func(sCtx provider.StepCtx) {
		commentID, statusCode := s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, "")
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.StatusCode(sCtx, 0, int(commentID))
	})

	t.WithNewStep("Создание комментария с запрещенным максимальным размером текста", func(sCtx provider.StepCtx) {
		commentID, statusCode := s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, strings.Repeat("q", 260))
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		assert.StatusCode(sCtx, 0, int(commentID))
	})
}
