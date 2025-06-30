package business

import (
	"fmt"
	"log/slog"
	"pos-master/models"
	"pos-master/proto/business"
	businessservices "pos-master/services/business_services"
	"pos-master/utils"
	"pos-master/utils/sentry"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

func CreateBusinessHandler(c *gin.Context) {

	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.Log(slog.LevelError, "❌error", "failed to parse form")
		utils.RespondWithError(c, 400, "Failed to parse form", fmt.Sprintf("error: %v", err))
		return
	}

	businessData := c.Request.FormValue("business")

	if businessData == "" {
		utils.Log(slog.LevelError, "❌error", "invalid business data")
		utils.RespondWithError(c, 400, " business Data is missing")
		return
	}

	var req business.BusinessRegisterRequest

	if err := protojson.Unmarshal([]byte(businessData), &req); err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to marshal  data")
		utils.RespondWithError(c, 400, "Unable to marshal data", fmt.Sprintf("error: %v", err))
		return
	}

	err := businessservices.CreateBusiness(c, &req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", fmt.Sprintf("error: %v", err))
		utils.RespondWithError(c, 400, "Unable to upload new  to business", fmt.Sprintf("error: %v", err))
		return

	}

	utils.RespondWithSuccess(c, "Added new  business data")

}

func GetBusinessesHandler(c *gin.Context) {
	var getRequest models.SearchRequest

	if err := c.ShouldBindJSON(&getRequest); err != nil {
		utils.Log(slog.LevelError, "❌error", "invalid request body")
		utils.RespondWithError(c, 400, utils.InvReqBody, fmt.Sprintf("error: %v", err))
		return
	}

	getRequest.SetDefaults()

	req := &business.GetBusinessesRequest{
		Page:        int32(getRequest.Page),
		PageSize:    int32(getRequest.PageSize),
		SearchQuery: getRequest.SearchQuery,
	}

	businesss, err := businessservices.GetBusinesses(req)

	if err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to retrieve businesss", "details", string(fmt.Sprintf("error: %v", err)))
		utils.RespondWithError(c, 400, utils.FailedToRetrieve("Businesses"), fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, utils.SuccessfullyRetrieve("Businesses"), gin.H{
		"businesss": businesss,
	})

}

func GetBusinessById(c *gin.Context) {
	businessId := c.Param("id")

	if businessId == "" {
		utils.RespondWithError(c, 400, "business ID is required")
		return
	}

	business, err := businessservices.GetBusinessById(businessId)

	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}
	utils.RespondWithSuccess(c, "successfully retrieved business", gin.H{
		"business": business,
	})

}

func EditBusinessHandler(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		utils.Log(slog.LevelError, "❌error", "failed to parse form")
		utils.RespondWithError(c, 400, "Failed to parse form", fmt.Sprintf("error: %v", err))
		return
	}

	businessData := c.Request.FormValue("business")

	if businessData == "" {
		utils.Log(slog.LevelError, "❌error", "invalid  business data")
		utils.RespondWithError(c, 400, "business Data is missing")
		return
	}

	var req business.EditBusinessRequest

	if err := protojson.Unmarshal([]byte(businessData), &req); err != nil {
		utils.Log(slog.LevelError, "❌error", "unable to marshal  data")
		utils.RespondWithError(c, 400, "Unable to marshal  data", fmt.Sprintf("error: %v", err))
		return
	}

	err := businessservices.EditBusiness(c, &req)
	if err != nil {
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}

	utils.RespondWithSuccess(c, "Business updated successfully")

}

func DeleteBusinessHandler(c *gin.Context) {

	businessId := c.Param("id")

	if businessId == "" {
		utils.RespondWithError(c, 400, "business ID is required")
		return
	}

	err := businessservices.DeleteBusiness(businessId)

	if err != nil {
		sentry.SentryLogger(c, err)
		utils.RespondWithError(c, 400, fmt.Sprintf("error: %v", err))
		return
	}
	utils.RespondWithSuccess(c, "successfully deleted business")

}
