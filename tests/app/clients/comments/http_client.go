package comments

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ozontech/allure-go/pkg/allure"
	"github.com/ozontech/allure-go/pkg/framework/provider"

	"route256/tests/app/domain"
)

type Client struct {
	baseUrl string
	cl      *http.Client
}

func NewClient(baseUrl string) *Client {
	return &Client{
		baseUrl: baseUrl,
		cl:      http.DefaultClient,
	}
}

func (c *Client) AddComment(
	ctx context.Context,
	t provider.StepCtx,
	userID int64,
	sku int64,
	comment string,
) (int64, int) {
	data, err := json.Marshal(addCommentRequest{
		UserID:  userID,
		Sku:     sku,
		Comment: comment,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/comment/add", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	response := addCommentResponse{}
	err = json.Unmarshal(body, &response)
	t.Require().NoError(err, "десериализация тела ответа")

	return response.CommentID, res.StatusCode
}

func (c *Client) EditComment(ctx context.Context, t provider.StepCtx, userID int64, commentID int64, newComment string) int {
	data, err := json.Marshal(editCommentRequest{
		UserID:     userID,
		CommentID:  commentID,
		NewComment: newComment,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/comment/edit", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	return res.StatusCode
}

func (c *Client) GetCommentListBySKU(ctx context.Context, t provider.StepCtx, sku int64) ([]*domain.Comment, int) {
	data, err := json.Marshal(getCommentListBySKURequest{
		Sku: sku,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/comment/list-by-sku", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	response := comments{}
	err = json.Unmarshal(body, &response)
	t.Require().NoError(err, "десериализация тела ответа")

	domainComments := make([]*domain.Comment, 0, len(response.Comments))
	for _, c := range response.Comments {
		domainComments = append(domainComments, &domain.Comment{
			ID:        c.ID,
			UserID:    c.UserID,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt,
		})
	}
	return domainComments, res.StatusCode
}

func (c *Client) GetCommentListByUser(ctx context.Context, t provider.StepCtx, userID int64) ([]*domain.Comment, int) {
	data, err := json.Marshal(getCommentListByUserRequest{
		UserID: userID,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/comment/list-by-user", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	response := comments{}
	err = json.Unmarshal(body, &response)
	t.Require().NoError(err, "десериализация тела ответа")

	domainComments := make([]*domain.Comment, 0, len(response.Comments))
	for _, c := range response.Comments {
		domainComments = append(domainComments, &domain.Comment{
			ID:        c.ID,
			SKU:       c.SKU,
			Comment:   c.Comment,
			CreatedAt: c.CreatedAt,
		})
	}
	return domainComments, res.StatusCode
}

func (c *Client) GetCommentByID(ctx context.Context, t provider.StepCtx, id int64) (*domain.Comment, int) {
	data, err := json.Marshal(commentGetByIDRequest{
		ID: id,
	})
	t.Require().NoError(err, "сериализация тела запроса")
	t.WithNewAttachment("request payload", allure.JSON, data)

	url := fmt.Sprintf("%s/comment/get-by-id", c.baseUrl)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
	t.Require().NoError(err, "создание запроса")

	r.Header.Add("Content-Type", "application/json")

	res, err := c.cl.Do(r)
	t.Require().NoError(err, "выполнение запроса")
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	t.Require().NoError(err, "не удалось считать ответ")
	t.WithNewAttachment("response body", allure.JSON, body)

	response := comment{}
	err = json.Unmarshal(body, &response)
	t.Require().NoError(err, "десериализация тела ответа")

	return &domain.Comment{
		ID:        response.ID,
		UserID:    response.UserID,
		Comment:   response.Comment,
		CreatedAt: response.CreatedAt,
	}, res.StatusCode
}
