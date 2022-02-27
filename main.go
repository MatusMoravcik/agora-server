package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
)

var appID, appCertificate string


func main() {


os.Setenv("APP_ID", "6bda2bd81c9f4f77bd85b0e99f430a42");
os.Setenv("APP_CERTIFICATE", "0af92bf4b1a047778a50d2a4226de2cb");
	

  	appIDEnv, appIDExists := os.LookupEnv("APP_ID")
	appCertEnv, appCertExists := os.LookupEnv("APP_CERTIFICATE")

  if !appIDExists || !appCertExists {
    log.Fatal("FATAL ERROR: ENV not properly configured, check APP_ID and APP_CERTIFICATE")
  } else {
    appID = appIDEnv
    appCertificate = appCertEnv
  }

  api := gin.Default()
  
  api.GET("rtc/:channelName/:role/:uid/", getRtcToken)
  
  api.GET("/hello", func(c *gin.Context) {
	  c.String(200, "Hello World!")
	})
	
	//   api.LoadHTMLGlob("templates/home.html")
	
	api.Run(":8080")
}

func getRtcToken(c *gin.Context) {
	log.Printf("rtc token\n")
	// get param values
	channelName, uidStr, role, expireTimestamp, err := parseRtcParams(c)
	// c.HTML(http.StatusOK, "home.html", nil);
	if err != nil {
	  c.Error(err)
	  c.AbortWithStatusJSON(400, gin.H{
		"message": "Error Generating RTC token: " + err.Error(),
		"status":  400,
	  })
	  return
	}
  
	rtcToken, tokenErr := generateRtcToken(channelName, uidStr, role, expireTimestamp)
  
	if tokenErr != nil {
	  log.Println(tokenErr) // token failed to generate
	  c.Error(tokenErr)
	  errMsg := "Error Generating RTC token - " + tokenErr.Error()
	  c.AbortWithStatusJSON(400, gin.H{
		"status": 400,
		"error":  errMsg,
	  })
	} else {
	  log.Println("RTC Token generated")
	  c.JSON(200, gin.H{
		"rtcToken": rtcToken,
	  })
	}
}

func parseRtcParams(c *gin.Context) (channelName, uidStr string, role rtctokenbuilder.Role, expireTimestamp uint32, err error) {
	// get param values
	channelName = c.Param("channelName")
	roleStr := c.Param("role")
	uidStr = c.Param("uid")
	expireTime := c.DefaultQuery("expiry", "3600")

	if roleStr == "publisher" {
		role = rtctokenbuilder.RolePublisher
	} else {
		role = rtctokenbuilder.RoleSubscriber
	}

	expireTime64, parseErr := strconv.ParseUint(expireTime, 10, 64)
	if parseErr != nil {
		err = fmt.Errorf("failed to parse expireTime: %s, causing error: %s", expireTime, parseErr)
	}

	// set timestamps
	expireTimeInSeconds := uint32(expireTime64)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expireTimestamp = currentTimestamp + expireTimeInSeconds

	return channelName, uidStr, role, expireTimestamp, err
}

func generateRtcToken(channelName, uidStr string, role rtctokenbuilder.Role, expireTimestamp uint32) (rtcToken string, err error) {

	  uid64, parseErr := strconv.ParseUint(uidStr, 10, 64)
	  if parseErr != nil {
		err = fmt.Errorf("failed to parse uidStr: %s, to uint causing error: %s", uidStr, parseErr)
		return "", err
	  }
  
	  uid := uint32(uid64) // convert uid from uint64 to uint 32
	  log.Printf("Building Token with uid: %d\n", uid)
	  rtcToken, err = rtctokenbuilder.BuildTokenWithUID(appID, appCertificate, channelName, uid, role, expireTimestamp)
	  return rtcToken, err
	  
}