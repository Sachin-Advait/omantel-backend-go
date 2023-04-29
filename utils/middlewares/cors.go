package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/configs"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/constant"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/flags"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/logger"
)

func CORSMiddleware() gin.HandlerFunc {

	applicationConfig, err := configs.Get(constant.ApplicationConfigFile)
	if err != nil {
		logger.FileLogger.Fatal().Err(err).Msg("CORSMiddleware: error loading application config file, env - " + flags.Env())
	}

	expectedHostName := applicationConfig.GetString(constant.ExpectedHostName)
	expectedReferer := applicationConfig.GetString(constant.ExpectedReferer)

	corsOrigin := applicationConfig.GetString(constant.CorsOriginEnabled)
	if corsOrigin == "" {
		corsOrigin = "*"
	}

	fmt.Println("corsOrigin: ", corsOrigin)

	return func(c *gin.Context) {

		if c.Request.Host == "" || c.Request.Host != expectedHostName {
			fmt.Println("host error: ", c.Request.Host)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if c.Request.Referer() == "" || c.Request.Referer() != expectedReferer {
			fmt.Println("host error: ", c.Request.Host)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", corsOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Super-Authorization, X-User-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
