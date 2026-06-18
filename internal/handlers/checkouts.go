package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
)

func (s *APIHandler) GetCheckoutsHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	id := c.Query("id")
	var checkouts []domain.Checkout
	var err error
	if id != "" {
		/* id, err := strconv.Atoi(id)
		if err == nil {
			items, err = s.db.Checkout(id)
		} */
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fetching checkout by ID is not supported yet"})
		return
	}
	checkouts, err = s.db.Checkouts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		slog.Error("Failed to fetch checkouts", "error", err)
		return
	}
	checkoutsJSON := make([]domain.CheckoutResponse, 0, len(checkouts))
	for _, c := range checkouts {
		checkout := domain.CheckoutResponse{
			ID: c.ID,
			User: domain.UserResponse{
				ID:    c.User.ID,
				Name:  c.User.Name,
				Email: c.User.Email,
			},
			Items:        c.Items,
			CheckoutDate: c.CheckoutDate,
			CreatedBy: domain.UserResponse{
				ID:    c.CreatedBy.ID,
				Name:  c.CreatedBy.Name,
				Email: c.CreatedBy.Email,
			},
			Notes:    c.Notes,
			Personal: c.Personal,
		}
		checkoutItems := make([]domain.CheckoutItem, 0, len(c.Items))
		for _, i := range c.Items {
			checkoutItems = append(checkoutItems, domain.CheckoutItem{
				Item: domain.Item{
					ID:    i.Item.ID,
					Name:  i.Item.Name,
					Notes: i.Item.Notes,
				},
				ReturnDate: i.ReturnDate,
			})
		}
		checkout.Items = checkoutItems
		checkoutsJSON = append(checkoutsJSON, checkout)
	}
	c.JSON(http.StatusOK, checkoutsJSON)
}

func (s *APIHandler) CreateCheckoutHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checkout := domain.CreateCheckoutRequest{}
	err := c.ShouldBindJSON(&checkout)
	if err != nil {
		slog.Error("Failed to deserialize json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	user, err := s.db.UserByEmail(ctx, checkout.UserEmail)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Info("User not found by email", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		slog.Error("Failed to lookup user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	admin, err := s.db.AdminByEmail(ctx, checkout.CreatedBy)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			slog.Info("Failed to lookup admin", "error", err)
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		slog.Error("Failed to lookup admin", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	checkoutID, err := s.db.CreateCheckout(ctx, *user, checkout.Items, checkout.CheckoutDate, *admin, checkout.Notes, checkout.Personal)
	if err != nil {
		slog.Error("Failed to create checkout", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"checkout_id": checkoutID})
}

func (s *APIHandler) ReturnCheckoutHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	checkoutReturn := domain.CheckoutReturnRequest{}
	err := c.ShouldBindJSON(&checkoutReturn)
	if err != nil {
		slog.Error("Failed to deserialize json", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	if err := s.db.ReturnItem(ctx, checkoutReturn.ID, checkoutReturn.Items); err != nil {
		slog.Error("Failed to return items", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}
