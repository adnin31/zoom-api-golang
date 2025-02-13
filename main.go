package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

// Zoom API Credentials
const (
	BaseURL = "https://api.zoom.us/v2"
	UserID  = "adnin.rais31@gmail.com"
)

// Meeting struct
type ZoomMeeting struct {
	ID        int    `json:"id"`
	Topic     string `json:"topic"`
	JoinURL   string `json:"join_url"`
	StartTime string `json:"start_time"`
	Duration  int    `json:"duration"`
	Password  string `json:"password,omitempty"`
}

type RequestPayload struct {
	Token string `json:"access_token"`
}

func getZoomToken() string {
	clientID := "lh9PDfiCSDy66gI4tIzpyg"
	clientSecret := "ATQ43IZhrutiuJM5NypFIUWWDS6rFNLz"
	accountID := "ffy1aZUTQUGTgoEVsZ-iog"

	// Encode credentials to Base64
	credentials := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))

	reqBody := strings.NewReader("grant_type=account_credentials&account_id=" + accountID)
	req, _ := http.NewRequest("POST", "https://zoom.us/oauth/token", reqBody)
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// Parse JSON response
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}

	// Return the access token
	accessToken, ok := result["access_token"].(string)
	if !ok {
		return ""
	}

	return accessToken
}

// Create Meeting
func createMeeting(c *gin.Context) {
	client := resty.New()
	var meeting ZoomMeeting
	if err := c.ShouldBindJSON(&meeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+getZoomToken()).
		SetHeader("Content-Type", "application/json").
		SetBody(meeting).
		Post(fmt.Sprintf("%s/users/%s/meetings", BaseURL, UserID))

	if err != nil {
		log.Println("Error creating meeting:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting"})
		return
	}

	c.JSON(http.StatusOK, resp.String())
}

// Get All Meetings
func getMeetings(c *gin.Context) {
	client := resty.New()
	var response map[string]interface{}

	_, err := client.R().
		SetHeader("Authorization", "Bearer "+getZoomToken()).
		SetResult(&response).
		Get(fmt.Sprintf("%s/users/%s/meetings", BaseURL, UserID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meetings"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Update Meeting
func updateMeeting(c *gin.Context) {
	client := resty.New()
	meetingID := c.Param("id")
	var meeting ZoomMeeting
	if err := c.ShouldBindJSON(&meeting); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+getZoomToken()).
		SetHeader("Content-Type", "application/json").
		SetBody(meeting).
		Patch(fmt.Sprintf("%s/meetings/%s", BaseURL, meetingID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meeting"})
		return
	}
	c.JSON(http.StatusOK, resp.String())
}

// Delete Meeting
func deleteMeeting(c *gin.Context) {
	client := resty.New()
	meetingID := c.Param("id")

	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+getZoomToken()).
		Delete(fmt.Sprintf("%s/meetings/%s", BaseURL, meetingID))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meeting"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Meeting deleted successfully", "response": resp.String()})
}

// Main Function
func main() {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Allow Next.js frontend
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/meetings", createMeeting)
	r.GET("/meetings", getMeetings)
	r.PUT("/meetings/:id", updateMeeting)
	r.DELETE("/meetings/:id", deleteMeeting)

	r.Run(":8080")
}
