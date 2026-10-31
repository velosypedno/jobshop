package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/velosypedno/jobshop/internal/backend/repos"
	"github.com/velosypedno/jobshop/internal/dag/parser"
	"github.com/velosypedno/jobshop/pkg/dag/core"
	"go.uber.org/zap"
)

type uploadFactoryConfigRepo interface {
	Create(ctx context.Context, name *string, config core.FactoryConfig) (uuid.UUID, error)
}

func NewUploadFactoryConfigHandler(repo uploadFactoryConfigRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, _, err := c.Request.FormFile("file")
		if err != nil {
			err = fmt.Errorf("upload factory config handler: failed to get file from form: %v", err)
			logger.Error("failed to get file from form", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		defer file.Close()

		config, err := parser.ParseFactoryConfig(file)
		if err != nil {
			err = fmt.Errorf("upload factory config handler: failed to parse json: %v", err)
			logger.Error("failed to parse json", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		fmt.Println(config)

		var name *string
		if nameParam := c.PostForm("name"); nameParam != "" {
			name = &nameParam
		}

		id, err := repo.Create(c.Request.Context(), name, config)
		if err != nil {
			logger.Error("failed to create factory config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"id": id})
	}
}

type getFactoryConfigRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (repos.FactoryConfigRecord, error)
}

type getFactoryConfigResponse struct {
	ID        uuid.UUID       `json:"id"`
	CreatedAt string          `json:"created_at"`
	Version   string          `json:"version"`
	Name      *string         `json:"name"`
	Config    json.RawMessage `json:"config"`
}

func NewGetFactoryConfigHandler(repo getFactoryConfigRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			err = fmt.Errorf("get factory config handler: invalid uuid format: %v", err)
			logger.Error("invalid uuid format", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		record, err := repo.GetByID(c.Request.Context(), id)
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("factory config not found", zap.String("id", idParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "factory config not found"})
			return
		}
		if err != nil {
			err = fmt.Errorf("get factory config handler: failed to get by id: %v", err)
			logger.Error("failed to get factory config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		configJSON, err := record.Config.ToJSON()
		if err != nil {
			err = fmt.Errorf("get factory config handler: failed to marshal config: %v", err)
			logger.Error("failed to marshal config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		response := getFactoryConfigResponse{
			ID:        record.ID,
			CreatedAt: record.CreatedAt.Format(time.RFC3339),
			Version:   record.Version,
			Name:      record.Name,
			Config:    configJSON,
		}

		c.JSON(http.StatusOK, response)
	}
}

type listFactoryConfigsRepo interface {
	List(ctx context.Context) ([]repos.FactoryConfigRecord, error)
}

type listFactoryConfigsResponse struct {
	Items []getFactoryConfigResponse `json:"items"`
}

func NewListFactoryConfigsHandler(repo listFactoryConfigsRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		records, err := repo.List(c.Request.Context())
		if err != nil {
			err = fmt.Errorf("list factory configs handler: failed to list: %v", err)
			logger.Error("failed to list factory configs", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		items := make([]getFactoryConfigResponse, len(records))
		for i, record := range records {
			configJSON, err := record.Config.ToJSON()
			if err != nil {
				err = fmt.Errorf("list factory configs handler: failed to marshal config: %v", err)
				logger.Error("failed to marshal config", zap.Error(err))
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			items[i] = getFactoryConfigResponse{
				ID:        record.ID,
				CreatedAt: record.CreatedAt.Format(time.RFC3339),
				Version:   record.Version,
				Name:      record.Name,
				Config:    configJSON,
			}
		}

		c.JSON(http.StatusOK, listFactoryConfigsResponse{Items: items})
	}
}

type deleteFactoryConfigRepo interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewDeleteFactoryConfigHandler(repo deleteFactoryConfigRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			err = fmt.Errorf("delete factory config handler: invalid uuid format: %v", err)
			logger.Error("invalid uuid format", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		err = repo.Delete(c.Request.Context(), id)
		if err != nil {
			err = fmt.Errorf("delete factory config handler: failed to delete: %v", err)
			logger.Error("failed to delete factory config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

type updateFactoryConfigRepo interface {
	Update(ctx context.Context, id uuid.UUID, name *string) error
}

type updateFactoryConfigRequestBody struct {
	Name *string `json:"name"`
}

func NewUpdateFactoryConfigHandler(repo updateFactoryConfigRepo, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := uuid.Parse(idParam)
		if err != nil {
			err = fmt.Errorf("update factory config handler: invalid uuid format: %v", err)
			logger.Error("invalid uuid format", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uuid format"})
			return
		}

		var req updateFactoryConfigRequestBody
		if err := c.ShouldBindJSON(&req); err != nil {
			err = fmt.Errorf("update factory config handler: failed to decode request body: %v", err)
			logger.Error("failed to decode request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		err = repo.Update(c.Request.Context(), id, req.Name)
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn("factory config not found for update", zap.String("id", idParam))
			c.JSON(http.StatusNotFound, gin.H{"error": "factory config not found"})
			return
		}
		if err != nil {
			err = fmt.Errorf("update factory config handler: failed to update: %v", err)
			logger.Error("failed to update factory config", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
