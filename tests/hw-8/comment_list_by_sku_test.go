package hw_8

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/assert"
	"route256/tests/app/domain"
)

func (s *Suite) TestGetSKUCommentsListOneUser_Success(t provider.T) {
	t.Title("Успешное получение списка комментариев, которые оставил один юзер по SKU")

	const (
		userID       = int64(30)
		sku          = int64(30)
		commentText1 = "тестовый комментарий 1"
		commentText2 = "тестовый комментарий 2"
	)

	var (
		commentIDs []int64
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID", userID, "sku", sku)

	t.WithNewStep("Создание комментария 1", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText1)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID1", commentID)

	t.WithNewStep("Получение первого созданного комментария", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentID, resultComment.ID)
		sCtx.Require().Equal(userID, resultComment.UserID)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})

	t.WithNewStep("Создание комментария 2", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID, sku, commentText2)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID2", commentID)

	t.WithNewStep("Получение созданных комментариев в обратном хронологическом порядке", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 2)
		sCtx.Require().NotNil(commentList[0])
		sCtx.Require().NotNil(commentList[1])

		// commentIDs содержит в себе id комментариев в прямом хронологическом порядке
		sCtx.Require().Equal(commentIDs[1], commentList[0].ID)
		sCtx.Require().Equal(userID, commentList[0].UserID)
		sCtx.Require().Equal(commentText2, commentList[0].Comment)

		sCtx.Require().Equal(commentIDs[0], commentList[1].ID)
		sCtx.Require().Equal(userID, commentList[1].UserID)
		sCtx.Require().Equal(commentText1, commentList[1].Comment)
	})
}

func (s *Suite) TestGetSKUCommentsListManyUsers_Success(t provider.T) {
	t.Title("Успешное получение списка комментариев, которые оставили несколько юзеров по SKU")

	const (
		userID1      = int64(310)
		userID2      = int64(311)
		sku          = int64(31)
		commentText1 = "тестовый комментарий 1"
		commentText2 = "тестовый комментарий 2"
	)

	var (
		commentIDs []int64
		commentID  int64
		statusCode int
	)

	t.WithNewParameters("userID1", userID1, "sku", sku)

	t.WithNewStep("Создание комментария 1", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID1, sku, commentText1)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID1", commentID)

	t.WithNewStep("Получение первого созданного комментария", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)

		sCtx.Require().Len(commentList, 1)
		sCtx.Require().NotNil(commentList[0])

		resultComment := commentList[0]
		sCtx.Require().Equal(commentID, resultComment.ID)
		sCtx.Require().Equal(userID1, resultComment.UserID)
		sCtx.Require().Equal(commentText1, resultComment.Comment)
	})

	t.WithNewParameters("userID2", userID2, "sku", sku)

	t.WithNewStep("Создание комментария 2", func(sCtx provider.StepCtx) {
		commentID, statusCode = s.commentsClient.AddComment(s.ctx, sCtx, userID2, sku, commentText2)
		commentIDs = append(commentIDs, commentID)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
	})

	t.WithNewParameters("commentID2", commentID)

	t.WithNewStep("Получение созданных комментариев в обратном хронологическом порядке", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 2)
		sCtx.Require().NotNil(commentList[0])
		sCtx.Require().NotNil(commentList[1])

		// commentIDs содержит в себе id комментариев в прямом хронологическом порядке
		sCtx.Require().Equal(commentIDs[1], commentList[0].ID)
		sCtx.Require().Equal(userID2, commentList[0].UserID)
		sCtx.Require().Equal(commentText2, commentList[0].Comment)

		sCtx.Require().Equal(commentIDs[0], commentList[1].ID)
		sCtx.Require().Equal(userID1, commentList[1].UserID)
		sCtx.Require().Equal(commentText1, commentList[1].Comment)
	})
}

func (s *Suite) TestGetSKUCommentsList_Negative(t provider.T) {
	t.Title("Негативные кейсы получения списка комментариев по SKU")

	const (
		sku = int64(33)
	)

	var (
		statusCode int
	)

	t.WithNewStep("Получение списка комментариев по SKU, у которой еще нет комментариев", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		assert.StatusCode(sCtx, http.StatusOK, statusCode)
		sCtx.Require().Len(commentList, 0)
	})

	t.WithNewStep("Получение списка комментариев по SKU без указания SKU", func(sCtx provider.StepCtx) {
		var commentList []*domain.Comment
		commentList, statusCode = s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, 0)
		assert.StatusCode(sCtx, http.StatusBadRequest, statusCode)
		sCtx.Require().Len(commentList, 0)
	})
}

// Потенциально недетерминированный тест.
//
//	"Комментарии оставленные двумя разными пользователями в один момент времени должны быть упорядочены по возрастанию идентификатора пользователя"
//	Основная проблема проверки сценария - отсутствие возможности влияния на createdAt
func (s *Suite) TestGetSKUCommentsListSameTimestamp_Success(t provider.T) {
	t.Title("Комментарии с одинаковым временем создания упорядочены по возрастанию ID пользователя")

	const (
		attemptsCount = 10
		userIDBase    = int64(99911100)
		sku           = int64(999222)
		commentText   = "комментарий от userID"
	)

	type testData struct {
		commentID   int64
		userID      int64
		statusCode  int
		commentText string
	}

	var data = []testData{}

	t.WithNewStep("Подготавливаем тестовые данные", func(sCtx provider.StepCtx) {
		for i := range attemptsCount {
			data = append(data, testData{
				userID:      userIDBase + int64(i),
				commentText: fmt.Sprintf("%s %d", commentText, i),
			})
		}
	})

	t.WithNewStep(fmt.Sprintf("Асинхронно создаем %d комментариев", attemptsCount), func(sCtx provider.StepCtx) {

		wg := sync.WaitGroup{}
		for i, item := range data {
			wg.Add(1)
			go func(ii int) {
				defer wg.Done()

				commentID, statusCode := s.commentsClient.AddComment(s.ctx, sCtx, item.userID, sku, item.commentText)

				data[ii].commentID = commentID
				data[ii].statusCode = statusCode
			}(i)
		}

		wg.Wait()
	})

	t.WithNewStep("Проверяем все ли успешно создано", func(sCtx provider.StepCtx) {
		count := 0

		for _, item := range data {
			if item.statusCode == http.StatusOK {
				count++
			}
		}

		sCtx.Assert().Equal(attemptsCount, count, fmt.Sprintf("Успешных %d из %d", count, attemptsCount))
	})

	t.WithNewStep("Проверяем сортировку", func(sCtx provider.StepCtx) {
		comments, status := s.commentsClient.GetCommentListBySKU(s.ctx, sCtx, sku)
		sCtx.Require().Equal(http.StatusOK, status, "Не удалось получить список комментариев по SKU")

		sCtx.WithNewStep("Проверяем наличие комментариев созданных в одну миллисекунду", func(sCtx provider.StepCtx) {
			uniqCreatedTimes := map[int64]struct{}{}
			for _, comment := range comments {
				uniqCreatedTimes[comment.CreatedAt.UnixMilli()] = struct{}{}
			}

			sCtx.Assert().GreaterOrEqual(
				len(comments), len(uniqCreatedTimes),
				"В тестовых данных нет комментариев созданных в одну миллисекунду",
			)
		})

		for i := 1; i < len(comments); i++ {
			if comments[i].CreatedAt != comments[i-1].CreatedAt {
				sCtx.Require().Greater(
					comments[i-1].CreatedAt,
					comments[i].CreatedAt,
					"Требуется упорядочивать комментарии в порядке, обратном хронологическому",
				)
			} else {
				sCtx.Require().Greater(
					comments[i].UserID,
					comments[i-1].UserID,
					"Комментарии с одинаковым временем создания должны быть отсортированы по возрастанию userID",
				)
			}
		}
	})
}
