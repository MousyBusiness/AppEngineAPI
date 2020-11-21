package main

import (
	_ "cloud.google.com/go/datastore"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mousybusiness/GoogleAppEngineAPI/internal/ds"
	"github.com/mousybusiness/GoogleAppEngineAPI/internal/mauth"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	apiKeyName = os.Getenv("API_KEY_NAME")
)

func main() {
	client, err := mauth.InitAuth()
	if err != nil {
		log.Fatalln("failed to init firebase auth", err)
	}

	ctx := context.Background()

	// create datastore client
	datastore, err := ds.CreateClient(ctx)
	if err != nil {
		log.Fatalln("couldnt start", err)
	}

	r := gin.Default()

	versioned := r.Group("/v1")
	// Use JWT auth middleware
	versioned.Use(mauth.AuthJWT(client))

	// Create entry
	versioned.POST("/", func(c *gin.Context) {
		var task ds.Task

		if err := c.BindJSON(&task); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		log.Println("created task", task)

		key, err := ds.CreateEntity(ctx, datastore, &task)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": key.ID, "name": key.Name, "kind": key.Kind, "namespace": key.Namespace})
	})

	// Get entry
	versioned.GET("/:id", func(c *gin.Context) {
		s := c.Param("id")
		id, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		entity, err := ds.GetEntity(ctx, datastore, int64(id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, entity)
	})

	// Delete entry
	versioned.DELETE("/:id", func(c *gin.Context) {
		s := c.Param("id")
		id, err := strconv.Atoi(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": http.StatusBadRequest, "message": err.Error()})
			return
		}

		if err := ds.DeleteEntity(ctx, datastore, int64(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})

	cron := r.Group("/cron")

	cron.Use(mauth.AuthAppEngineCron())

	// Cron job
	cron.GET("/dosomething", func(c *gin.Context) {
		log.Println("running cron job")
		c.String(http.StatusOK, "working")
	})

	// API Key routes
	api := r.Group("/api")
	api.Use(mauth.AuthAPIKey(apiKeyName))
	api.DELETE("/user/delete/:uid", func(c *gin.Context) {
		uid := c.Param("uid")
		log.Println("pretending to delete user:", uid)
		c.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "message": "Success"})
	})

	// Admin routes
	admin := r.Group("/admin")

	// admin must be logged in
	admin.Use(mauth.AuthJWT(client))
	// must have api key
	admin.Use(mauth.AuthAPIKey(apiKeyName))

	// create admin claim
	admin.POST("/:uid", func(c *gin.Context){
		uid := c.Param("uid")
		log.Println("elevating user to admin")
		err := mauth.ElevateToAdmin(ctx, client, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
	})
	// delete admin claim
	admin.DELETE("/:uid", func(c *gin.Context){
		uid := c.Param("uid")
		log.Println("revoking users admin rights")

		err := mauth.RevokeAdmin(ctx, client, uid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": http.StatusInternalServerError, "message": err.Error()})
			return
		}
	})
	// verify admin
	admin.GET("/:uid", func(c *gin.Context){
		uid := c.Param("uid")

		log.Println("is user admin?")
		err := mauth.VerifyAdmin(ctx, client, uid)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "message": err.Error()})
			return
		}
	})

	// start API
	errRun := r.Run()
	if errRun != nil {
		log.Fatalln("failed to run Gin app")
	}
}
