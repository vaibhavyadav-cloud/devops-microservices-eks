package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/devopsdemo/order-service/internal/kafkaproducer"
	"github.com/devopsdemo/order-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderHandler struct {
	pool     *pgxpool.Pool
	producer *kafkaproducer.Producer
}

func NewOrderHandler(pool *pgxpool.Pool, producer *kafkaproducer.Producer) *OrderHandler {
	return &OrderHandler{pool: pool, producer: producer}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_failed", "details": err.Error()})
		return
	}

	order := models.Order{
		ID:            uuid.NewString(),
		CustomerEmail: req.CustomerEmail,
		Item:          req.Item,
		Quantity:      req.Quantity,
		Status:        "PENDING",
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	err := h.pool.QueryRow(ctx,
		`INSERT INTO orders (id, customer_email, item, quantity, status)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING created_at`,
		order.ID, order.CustomerEmail, order.Item, order.Quantity, order.Status,
	).Scan(&order.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}

	// Publish AFTER the insert succeeds. If this publish fails, the order
	// still exists — we deliberately don't roll back a successful order just
	// because the notification event didn't go out. See the "outbox pattern"
	// note in kafkaproducer/producer.go for the production-grade fix to the
	// dual-write gap this creates.
	_ = h.producer.PublishOrderCreated(ctx, kafkaproducer.OrderCreatedEvent{
		OrderID:       order.ID,
		CustomerEmail: order.CustomerEmail,
	})

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var order models.Order
	err := h.pool.QueryRow(ctx,
		`SELECT id, customer_email, item, quantity, status, created_at FROM orders WHERE id = $1`,
		id,
	).Scan(&order.ID, &order.CustomerEmail, &order.Item, &order.Quantity, &order.Status, &order.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "order_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) ListOrders(c *gin.Context) {
	// Fixed page size for now — query-param-driven limit/offset is a natural
	// next enhancement once this is wired up end-to-end.
	const limit = 20
	const offset = 0

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	rows, err := h.pool.Query(ctx,
		`SELECT id, customer_email, item, quantity, status, created_at
		 FROM orders ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}
	defer rows.Close()

	orders := []models.Order{}
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(&o.ID, &o.CustomerEmail, &o.Item, &o.Quantity, &o.Status, &o.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
			return
		}
		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_failed", "details": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	tag, err := h.pool.Exec(ctx, `UPDATE orders SET status = $1 WHERE id = $2`, req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal_server_error"})
		return
	}
	if tag.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "order_not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id, "status": req.Status})
}
