package impl

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.com/rendezvousPoint/incentive_calculator/models"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils/errors"
)

func (impl *UserImpl) CaptchaVerification(ctx *gin.Context, req *models.CaptchRequest) (interface{}, *errors.Error) {

	captchaPayloadRequest := url.Values{}
	captchaPayloadRequest.Set("secret", models.CaptchaSecretKey)
	captchaPayloadRequest.Set("response", req.Token)

	verifyCaptchaRequest, err := http.NewRequest(http.MethodPost, models.VerifyCaptchaGoogleApi, strings.NewReader(captchaPayloadRequest.Encode()))
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "req generate failed: " + err.Error(),
		}
	}

	verifyCaptchaRequest.Header.Add("content-type", "application/x-www-form-urlencoded")
	verifyCaptchaRequest.Header.Add("cache-control", "no-cache")

	verifyCaptchaResponse, err := http.DefaultClient.Do(verifyCaptchaRequest)
	if err != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "http call failed: " + err.Error(),
		}
	}

	var captchaVerifyResponse interface{}

	decoder := json.NewDecoder(verifyCaptchaResponse.Body)
	decoderErr := decoder.Decode(&captchaVerifyResponse)

	defer verifyCaptchaResponse.Body.Close()

	if decoderErr != nil {
		return nil, &errors.Error{
			StatusCode: http.StatusInternalServerError,
			Code:       errors.InternalServerErrorCode,
			Message:    errors.InternalServerError,
			Details:    "decoding failed: " + decoderErr.Error(),
		}
	}

	return captchaVerifyResponse, nil
}
