package handlers

import (
	"net/http"

	"inventory/internal/domain"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (s *APIHandler) GetCheckoutsHandler(c *gin.Context) {
	id := c.Query("id")
	var items *[]domain.Checkout
	var err error
	if id != "" {
		/* id, err := strconv.Atoi(id)
		if err == nil {
			items, err = s.db.Checkout(id)
		} */
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fetching checkout by ID is not supported yet"})
		return
	} else {
		items, err = s.db.Checkouts()
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (s *APIHandler) CreateCheckoutHandler(c *gin.Context) {
	checkout := domain.Checkout{}
	err := c.ShouldBindJSON(&checkout)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = s.db.CreateCheckout(&checkout)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, "item added successfully")
}
