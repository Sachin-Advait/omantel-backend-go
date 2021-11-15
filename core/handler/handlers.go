package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

type SalesIncentiveService interface {
	FetchSingleUserIncentive(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.TotalIncentive, error)
	FetchSingleUserTargetChart(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.ChartDataResponse, error)
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

	r.Post("/client", context.GetSingleUserIncentive)
	r.Post("/target-chart", context.GetSingleUserTargetChart)

	return r
}

func (context *SalesIncentiveContext) GetSingleUserIncentive(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
		Quarter int `json:"quarter"`
		Date string `json:"date"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to read req body")
		w.Write([]byte(`{"error": "internal error failed to read req body"}`))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "internal error failed to read req body"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to unmarshall request")
		w.Write([]byte(`{"error": "internal request unmarshalling error"}`))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "internal request unmarshalling error"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	id, err := strconv.Atoi(request.ClientId)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	var (
		isQuarterReport bool
		date time.Time
	)

	if request.Date == "" {
		utils.LogObj.Error(logs.RequestID, " : date type is missing, id: "+request.Date)
		w.Write([]byte(`{"error": "date type is missing"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "date type is missing"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	date, err = time.Parse("01/2006", request.Date)
	if err != nil {

		date, err = time.Parse("2006", request.Date)
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect saleDate, id: "+request.ClientId+" saleDate: "+request.Date)
			w.Write([]byte(`{"error": "incorrect saleDate format: {for quarterly format use - '2006'} & {for monthly format use - '11/2006'}"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect saleDate format: {for quarterly format use - '2006'} & {for monthly format use - '11/2006'}"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = true
	}

	if isQuarterReport && (request.Quarter < 1 || request.Quarter > 4) {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+request.Date)
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
	}

	incentive, err := context.SalesIncentiveService.FetchSingleUserIncentive(id, date, request.Quarter, isQuarterReport)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : Incentive Response fetched successfully for clientId: "+request.ClientId)

	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetSingleUserTargetChart(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
		Quarter int `json:"quarter"`
		Date string `json:"date"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to read req body")
		w.Write([]byte(`{"error": "internal error failed to read req body"}`))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "internal error failed to read req body"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to unmarshall request")
		w.Write([]byte(`{"error": "internal request unmarshalling error"}`))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "internal request unmarshalling error"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	id, err := strconv.Atoi(request.ClientId)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	var (
		isQuarterReport bool
		date time.Time
	)

	if request.Date == "" {
		utils.LogObj.Error(logs.RequestID, " : date type is missing, id: "+request.Date)
		w.Write([]byte(`{"error": "date type is missing"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "date type is missing"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	date, err = time.Parse("01/2006", request.Date)
	if err != nil {

		date, err = time.Parse("2006", request.Date)
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect saleDate, id: "+request.ClientId+" saleDate: "+request.Date)
			w.Write([]byte(`{"error": "incorrect saleDate format: {for quarterly format use - '2006'} & {for monthly format use - '11/2006'}"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect saleDate format: {for quarterly format use - '2006'} & {for monthly format use - '11/2006'}"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = true
	}

	if isQuarterReport && (request.Quarter < 1 || request.Quarter > 4) {
		utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+request.Date)
		w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	targetData, err := context.SalesIncentiveService.FetchSingleUserTargetChart(id, date, request.Quarter, isQuarterReport)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(targetData)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : TargetChart Response fetched successfully for clientId: "+request.ClientId)

	//utils.LogObj.WriteLogs(logs)
}
