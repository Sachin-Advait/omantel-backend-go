package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
	"gitlab.com/rendezvousPoint/incentive_calculator/core/db"
	"gitlab.com/rendezvousPoint/incentive_calculator/utils"
)

type SalesIncentiveService interface {
	FetchPromoterUserIncentive(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.TotalIncentive, error)
	FetchPromoterTargetChart(clientID int, t time.Time, quarter int, isQuarterReport bool) (*db.ChartDataResponse, error)
	FetchSaleCountPerPlan(clientID int, group string, t time.Time, quarter int, isQuarterReport bool) ([]*db.SaleCountPerPlan, error)
	GenerateOTP(clientID int, requestID string) *commons.GenericError
	VefiryOTP(clientID, otp int) (*db.OtpModel, *commons.GenericError)
	FetchUserReportingToOutletManager(outletManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchUserReportingToManager(outletManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchUserReportingToSeniorManager(outletManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchUserReportingToGeneralManager(outletManagerID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchUserReportingToVicePresident(vicePresidentID int) ([]*db.ReportingStaffIDAndName, *commons.GenericError)
	FetchUserPayoutReport(id int, t time.Time, quarter int, isQuarterly bool) (*db.TotalIncentive, *commons.GenericError)
	FetchUserTargetReport(id int, t time.Time, quarter int, isQuarterly bool) (*db.ChartDataResponse, *commons.GenericError)
	FetchSalesCountPerPlanReport(id int, group string, t time.Time, quarter int, isQuarterly bool) ([]*db.SaleCountPerPlan, *commons.GenericError)
	FetchUserOverAllDetails(clientID, month, quarter int, isQuarterReport bool) (*db.UserRecord, *commons.GenericError)
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
	r.Post("/sale-per-plan", context.GetSaleCountPerPlan)
	r.Post("/generate-otp", context.GenerateOTP)
	r.Post("/verify-otp", context.VerifyOTP)
	r.Post("/user-report", context.GetUserReport)
	r.Post("/report_to_outlet_manager", context.GetUsersReportingToOutletManager)
	r.Post("/report_to_manager", context.GetUsersReportingToManager)
	r.Post("/report_to_senior_manager", context.GetUsersReportingToSeniorManager)
	r.Post("/report_to_general_manager", context.GetUsersReportingToGeneralManager)
	r.Post("/report_to_vice_president", context.GetUsersReportingToVicePresident)
	r.Post("/usr_payout", context.GetUserPayout)
	r.Post("/usr_target", context.GetUserTarget)
	r.Post("/usr_sales_per_plan", context.GetSaleCountPerPlan)

	return r
}

func (context *SalesIncentiveContext) GetSingleUserIncentive(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
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
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	fmt.Println("data", id, date, quarter, isQuarterReport)

	incentive, err := context.SalesIncentiveService.FetchPromoterUserIncentive(id, date, quarter, isQuarterReport)
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
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
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
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	targetData, err := context.SalesIncentiveService.FetchPromoterTargetChart(id, date, quarter, isQuarterReport)
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

func (context *SalesIncentiveContext) GetSaleCountPerPlan(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
		Product   string `json:"product"`
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

	product := strings.ToUpper(request.Product)
	if !(product == commons.HBBCategory || product == commons.BaqatiCategory || product == commons.HayyakCategory || product == commons.UpgradeCategory) {

		utils.LogObj.Error(logs.RequestID, " : incorrect product type, id: "+request.ClientId+" productType: "+product)
		w.Write([]byte(`{"error": "incorrect product type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect product type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	var (
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	respData, err := context.SalesIncentiveService.FetchSaleCountPerPlan(id, product, date, quarter, isQuarterReport)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(respData)
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

	utils.LogObj.Info(logs.RequestID, " : SaleCountPerPlan Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUserReport(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
		Month    int    `json:"month"`
		Quarter  int    `json:"quarter"`
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

	if !((request.Month >= 1 && request.Month <= 12) || (request.Quarter >= 1 && request.Quarter <= 4)) {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found incorrect, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is incorrect"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is incorrect"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	var isQuarterReport bool

	if request.Quarter >= 1 && request.Quarter <= 4 {
		isQuarterReport = true
	}

	respData, errs := context.SalesIncentiveService.FetchUserOverAllDetails(id, request.Month, request.Quarter, isQuarterReport)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(respData)
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

	utils.LogObj.Info(logs.RequestID, " : GetUserReport Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GenerateOTP(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to read req body "+err.Error())
		w.Write([]byte(`{"error": "INTERNAL SERVER ERROR"}`))
		w.WriteHeader(http.StatusInternalServerError)

		// logs.ResponseStatus = http.StatusInternalServerError
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "internal error failed to read req body"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to unmarshall request "+err.Error())
		w.Write([]byte(`{"error": "request unmarshalling error"}`))
		w.WriteHeader(http.StatusBadRequest)

		// logs.ResponseStatus = http.StatusBadRequest
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "request unmarshalling error"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	id, err := strconv.Atoi(request.ClientId)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		// logs.ResponseStatus = http.StatusBadRequest
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	var resp struct {
		Error string `json:"error"`
	}

	genErr := context.SalesIncentiveService.GenerateOTP(id, logs.RequestID)
	if genErr != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to generate otp, err: "+genErr.ErrorMessage+" id: "+request.ClientId)

		resp.Error = commons.ErrorMap[genErr.Code]
		byteRsp, _ := json.Marshal(resp)

		w.Write(byteRsp)

		if genErr.Code == 1001 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		// logs.ResponseStatus = http.StatusInternalServerError
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		//utils.LogObj.WriteLogs(logs)
		return
	}

	resp.Error = "nil"

	byteRsp, _ := json.Marshal(resp)

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	// logs.ResponseStatus = http.StatusOK
	// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	// logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : Generated Otp successfully for clientId: "+request.ClientId)
}

func (context *SalesIncentiveContext) VerifyOTP(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
		Otp      int    `json:"otp"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to read req body "+err.Error())
		w.Write([]byte(`{"error": "internal error failed to read req body"}`))
		w.WriteHeader(http.StatusInternalServerError)

		// logs.ResponseStatus = http.StatusInternalServerError
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "internal error failed to read req body"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	err = json.Unmarshal(dataByte, &request)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to unmarshall request "+err.Error())
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

		// logs.ResponseStatus = http.StatusBadRequest
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Otp == 0 {
		utils.LogObj.Error(logs.RequestID, " : incorrect otp, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect OTP"}`))
		w.WriteHeader(http.StatusBadRequest)

		// logs.ResponseStatus = http.StatusBadRequest
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string([]byte(`{"error": "incorrect OTP"}`))
		//utils.LogObj.WriteLogs(logs)
		return
	}

	var resp struct {
		Authenticated bool   `json:"authenticated"`
		StaffName     string `json:"staff_name"`
		StaffPosition string `json:"staff_position"`
		Error         string `json:"error"`
	}

	otpModel, genErr := context.SalesIncentiveService.VefiryOTP(id, request.Otp)
	if genErr != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to validate otp, err: "+genErr.ErrorMessage+" id: "+request.ClientId)

		resp.Authenticated = false
		resp.Error = commons.ErrorMap[genErr.Code]

		byteRsp, _ := json.Marshal(resp)

		w.Write(byteRsp)

		if genErr.Code == 1001 {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		// logs.ResponseStatus = http.StatusBadRequest
		// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		// logs.ResponseBody = string(byteRsp)
		//utils.LogObj.WriteLogs(logs)
		return
	}

	resp.Error = "nil"
	resp.Authenticated = true
	resp.StaffName = otpModel.StaffName
	resp.StaffPosition = otpModel.StaffPosition

	byteRsp, _ := json.Marshal(resp)

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	// logs.ResponseStatus = http.StatusOK
	// logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	// logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : Verified Otp successfully for clientId: "+request.ClientId)
}

func (context *SalesIncentiveContext) GetUsersReportingToOutletManager(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToOutletManager failed to read req body, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToOutletManager failed to unmarshall request, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToOutletManager incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, errs := context.SalesIncentiveService.FetchUserReportingToOutletManager(id)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToOutletManager failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.Write([]byte(commons.ErrorMap[errs.Code]))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseBody = commons.ErrorMap[errs.Code]
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToOutletManager failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : GetUsersReportingToOutletManager Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUsersReportingToManager(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToManager failed to read req body, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToManager failed to unmarshall request, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToManager incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, errs := context.SalesIncentiveService.FetchUserReportingToManager(id)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToManager failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.Write([]byte(commons.ErrorMap[errs.Code]))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseBody = commons.ErrorMap[errs.Code]
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToManager failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : GetUsersReportingToManager Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUsersReportingToSeniorManager(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToSeniorManager failed to read req body, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToSeniorManager failed to unmarshall request, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToSeniorManager incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, errs := context.SalesIncentiveService.FetchUserReportingToSeniorManager(id)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToSeniorManager failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.Write([]byte(commons.ErrorMap[errs.Code]))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseBody = commons.ErrorMap[errs.Code]
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToSeniorManager failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : GetUsersReportingToSeniorManager Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUsersReportingToGeneralManager(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToGeneralManager failed to read req body, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToGeneralManager failed to unmarshall request, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToGeneralManager incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, errs := context.SalesIncentiveService.FetchUserReportingToGeneralManager(id)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToGeneralManager failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.Write([]byte(commons.ErrorMap[errs.Code]))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseBody = commons.ErrorMap[errs.Code]
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToGeneralManager failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : GetUsersReportingToGeneralManager Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUsersReportingToVicePresident(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId string `json:"client_id"`
	}

	dataByte, err := ioutil.ReadAll(req.Body)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToVicePresident failed to read req body, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToVicePresident failed to unmarshall request, err- "+err.Error())
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
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToVicePresident incorrect clientID type, id: "+request.ClientId)
		w.Write([]byte(`{"error": "incorrect clientID type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect clientID type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	incentive, errs := context.SalesIncentiveService.FetchUserReportingToVicePresident(id)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToVicePresident failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.Write([]byte(commons.ErrorMap[errs.Code]))
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseBody = commons.ErrorMap[errs.Code]
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(incentive)
	if err != nil {
		utils.LogObj.Error(logs.RequestID, " : GetUsersReportingToVicePresident failed to marshal client resp, err: "+err.Error()+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(byteRsp)
	w.WriteHeader(http.StatusOK)

	logs.ResponseStatus = http.StatusOK
	logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
	logs.ResponseBody = string(byteRsp)

	utils.LogObj.Info(logs.RequestID, " : GetUsersReportingToVicePresident Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUserPayout(w http.ResponseWriter, req *http.Request) {

	logs := utils.LogHelper(req)

	var request struct {
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
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
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	fmt.Println("data", id, date, quarter, isQuarterReport)

	incentive, errs := context.SalesIncentiveService.FetchUserPayoutReport(id, date, quarter, isQuarterReport)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
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

	utils.LogObj.Info(logs.RequestID, " : GetUserPayout Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUserTarget(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
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
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	fmt.Println("data", id, date, quarter, isQuarterReport)

	targetChart, errs := context.SalesIncentiveService.FetchUserTargetReport(id, date, quarter, isQuarterReport)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(targetChart)
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

	utils.LogObj.Info(logs.RequestID, " : GetUserTarget Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}

func (context *SalesIncentiveContext) GetUserSalesCountPerPlan(w http.ResponseWriter, req *http.Request) {
	logs := utils.LogHelper(req)

	var request struct {
		ClientId  string `json:"client_id"`
		Monthly   string `json:"monthly"`
		Quarterly string `json:"quarterly"`
		Product   string `json:"product"`
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

	product := strings.ToUpper(request.Product)
	if !(product == commons.HBBCategory || product == commons.BaqatiCategory || product == commons.HayyakCategory || product == commons.UpgradeCategory) {

		utils.LogObj.Error(logs.RequestID, " : incorrect product type, id: "+request.ClientId+" productType: "+product)
		w.Write([]byte(`{"error": "incorrect product type"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "incorrect product type"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	var (
		isQuarterReport bool = true
		date            time.Time
		quarter         int
	)

	if request.Monthly == "" && request.Quarterly == "" {
		utils.LogObj.Error(logs.RequestID, " : monthy and quarterly type found empty, id: "+request.ClientId)
		w.Write([]byte(`{"error": "monthy and quartlerly type is blank"}`))
		w.WriteHeader(http.StatusBadRequest)

		logs.ResponseStatus = http.StatusBadRequest
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
		logs.ResponseBody = string([]byte(`{"error": "monthy and quartlerly type is blank"}`))

		//utils.LogObj.WriteLogs(logs)
		return
	}

	if request.Monthly != "" {

		date, err = time.Parse("01-2006", strings.TrimSpace(request.Monthly))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : monthy type is incorrect, id: "+request.ClientId+" monthly: "+request.Monthly)
			w.Write([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "monthy type is incorrect, use '06-2006' format."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		isQuarterReport = false
	}

	if isQuarterReport && request.Quarterly != "" {

		req := strings.Split(strings.TrimSpace(request.Quarterly), "-")
		if len(req) < 2 {
			utils.LogObj.Error(logs.RequestID, " : quarterly type is incorrect, id: "+request.ClientId+" quarterly: "+request.Quarterly)
			w.Write([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarterly type is incorrect, avoid spaces."}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		quarter, err = strconv.Atoi(strings.TrimSpace(req[0]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : quarter value not found, id: "+request.ClientId+" quarterly: "+req[0]+".")
			w.Write([]byte(`{"error": "quarter value not found"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "quarter value not found"}`))

			//utils.LogObj.WriteLogs(logs)
			return

		} else if quarter < 1 || quarter > 4 {
			utils.LogObj.Error(logs.RequestID, " : incorrect quarter value, id: "+request.ClientId+" quarter: "+req[0])
			w.Write([]byte(`{"error": "incorrect quarter value - should be between 1 and 4"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect quarter value - should be between 1 and 4""}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

		date, err = time.Parse("2006", strings.TrimSpace(req[1]))
		if err != nil {
			utils.LogObj.Error(logs.RequestID, " : incorrect year value, id: "+request.ClientId+" year: "+req[1])
			w.Write([]byte(`{"error": "incorrect year value"}`))
			w.WriteHeader(http.StatusBadRequest)

			logs.ResponseStatus = http.StatusBadRequest
			logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())
			logs.ResponseBody = string([]byte(`{"error": "incorrect year value"}`))

			//utils.LogObj.WriteLogs(logs)
			return
		}

	}

	respData, errs := context.SalesIncentiveService.FetchSalesCountPerPlanReport(id, product, date, quarter, isQuarterReport)
	if errs != nil {
		utils.LogObj.Error(logs.RequestID, " : failed to fetch user incentive, err: "+errs.ErrorMessage+" id: "+request.ClientId)
		w.WriteHeader(http.StatusInternalServerError)

		logs.ResponseStatus = http.StatusInternalServerError
		logs.ResponseHeaderAndValues = fmt.Sprint(w.Header())

		//utils.LogObj.WriteLogs(logs)
		return
	}

	byteRsp, err := json.Marshal(respData)
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

	utils.LogObj.Info(logs.RequestID, " : GetUserSalesCountPerPlan Response fetched successfully for clientId: "+request.ClientId)
	//utils.LogObj.WriteLogs(logs)
}
