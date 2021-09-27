package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

type SalesIncentiveService interface {
	FetchSingleUserIncentive(int) (*db.TotalIncentive, error)
}

type SalesIncentiveContext struct {
	SalesIncentiveService SalesIncentiveService
}

func NewSalesIncentiveContext(salesIncentiveService SalesIncentiveService) *SalesIncentiveContext {

	return &SalesIncentiveContext{SalesIncentiveService: salesIncentiveService}
}

func (context *SalesIncentiveContext) SafeClose() {}

func (context *SalesIncentiveContext) NewSalesIncentiveRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/single-client", context.GetSingleUserIncentive)

	return r
}

func (context *SalesIncentiveContext) GetSingleUserIncentive(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	// var request struct {
	// 	ClientId string `json:"client_id"`
	// }

	// dataByte, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	utils.LogObj.Error(logs.RequestID, " : failed to read req body")
	// 	w.Write([]byte(`{"error": "internal error failed to read req body"}`))
	// 	w.WriteHeader(http.StatusInternalServerError)

	// 	logs.ResponseStatus = http.StatusInternalServerError
	// 	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	// 	logs.ResponseBody = string([]byte(`{"error": "internal error failed to read req body"}`))

	// 	//utils.LogObj.WriteLogs(logs)
	// 	return
	// }

	// err = json.Unmarshal(dataByte, &request)
	// if err != nil {
	// 	utils.LogObj.Error(logs.RequestID, " : failed to unmarshall request")
	// 	w.Write([]byte(`{"error": "internal request unmarshalling error"}`))
	// 	w.WriteHeader(http.StatusInternalServerError)

	// 	logs.ResponseStatus = http.StatusInternalServerError
	// 	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	// 	logs.ResponseBody = string([]byte(`{"error": "internal request unmarshalling error"}`))

	// 	//utils.LogObj.WriteLogs(logs)
	// 	return
	// }

	clientID := req.Header.Get(commons.ClientIDHeader)
	if clientID == "" {
		utils.LogObj.Error(logs.RequestID, " : ClientID not found in body")
		w.Write([]byte(`{"error": "ClientID not found in body"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "ClientID not found in body"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	id, err := strconv.Atoi(clientID)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : incorrect clientID type, id: "+clientID)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, err := context.SalesIncentiveService.FetchSingleUserIncentive(id)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+err.Error()+" id: "+clientID)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to marshal client resp, err: "+err.Error()+" id: "+clientID)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : Response fetched successfully for clientId: "+clientID)

	//utils.LogObj.WriteLogs(logs)
}
