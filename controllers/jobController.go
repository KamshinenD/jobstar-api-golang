package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"jobstar.com/api/models"
	"jobstar.com/api/utils"
)

// @Summary User creates a Job
// @Description User creates a Job
// @Tags job
// @Accept  json
// @Security ApiKeyAuth
// @Produce  json
// @Param user body models.JobRequest true "Create Job Data"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /jobs [POST]

func CreateJob(c *gin.Context) {

	userId, exists := c.Get("userId")

	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "userId is not a string", nil)
		return
	}

	var job models.Job

	err := c.ShouldBindJSON(&job)

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	if job.Company == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide company name", nil)
		return
	}
	if job.JobLocation == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide job location", nil)
		return
	}
	if job.Position == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide position", nil)
		return
	}
	if job.JobType == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide job type", nil)
		return
	}

	if !job.JobType.JobTypeIsValid() {
		utils.RespondError(c, http.StatusBadRequest, "Invalid job type", nil)
		return
	}

	if job.Status == "" {
		job.Status = "pending"
	}

	if !job.Status.StatusIsValid() {
		utils.RespondError(c, http.StatusBadRequest, "Invalid Status", nil)
		return
	}

	job.CreatedAt = time.Now()
	job.CreatedBy = userIdStr

	err = job.SaveJob()

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Could not create job", nil)

	}
	utils.RespondJSON(c, http.StatusOK, "Job craeted successfully", nil)

}

// @Summary Get all jobs for user
// @Description Gets all Jobs created by user
// @Tags Job
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /jobs [GET]
func GetJobsByUser(c *gin.Context) {
	userId, exists := c.Get("userId")

	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "userId is not a string", nil)
		return
	}

	jobs, err := models.GetJobs(userIdStr)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "could not fetch Jobs", err)
		return
	}

	// Ensure that jobs is not nil
	if jobs == nil {
		jobs = []models.Job{} // Return an empty slice instead of nil
	}

	utils.RespondJSON(c, http.StatusOK, "Data retrieved successfully", gin.H{
		"jobs": jobs,
	})

}

// @Summary Get job by ID
// @Description Retrieves a job using its ID
// @Tags Job
// @Security ApiKeyAuth
// @Produce  json
// @Param   id   path    string  true  "Job ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /jobs/{id} [GET]
func GetSingleJob(c *gin.Context) {
	jobId := c.Param("id")

	userId, exists := c.Get("userId")

	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusInternalServerError, "userId is not a string", nil)
		return
	}

	job, err := models.GetUserJobById(jobId, userIdStr)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Unable to fetch job data", err)
		return
	}
	utils.RespondJSON(c, http.StatusOK, "Data retrieved successfully", gin.H{
		"job": job,
	})
}

// @Summary Get job by ID
// @Description Deletes a job using its ID
// @Tags Job
// @Security ApiKeyAuth
// @Produce  json
// @Param   id   path    string  true  "Job ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /jobs/{id} [DELETE]
func DeleteJob(c *gin.Context) {
	jobId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist in context", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "user ID is not a string", nil)
		return
	}

	job, err := models.GetJobById(jobId)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "could not fecth Job", err)
		return
	}

	if job.CreatedBy != userIdStr {
		utils.RespondError(c, http.StatusUnauthorized, "You are not authorized to delete this event", nil)
		return
	}

	err = job.Delete()
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Could not delete Job", nil)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Job Deleted successfully", nil)
}

// @Summary Update Job details
// @Description Updates the authenticated user's details
// @Tags Job
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   id   path    string  true  "Job ID"
// @Param user body models.JobRequest true "User UpdateJob Data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /jobs/{id} [PATCH]
func UpdateJob(c *gin.Context) {

	var updatedJob models.Job
	err := c.ShouldBindJSON(&updatedJob)

	if err != nil {
		utils.RespondError(c, http.StatusBadRequest, "Could not parse request data", err)
	}

	if updatedJob.Company == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide company name", nil)
		return
	}
	if updatedJob.JobLocation == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide job location", nil)
		return
	}
	if updatedJob.Position == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide position", nil)
		return
	}
	if updatedJob.JobType == "" {
		utils.RespondError(c, http.StatusBadRequest, "Please provide job type", nil)
		return
	}

	if !updatedJob.JobType.JobTypeIsValid() {
		utils.RespondError(c, http.StatusBadRequest, "Invalid job type", nil)
		return
	}

	if updatedJob.Status == "" {
		updatedJob.Status = "pending"
	}

	if !updatedJob.Status.StatusIsValid() {
		utils.RespondError(c, http.StatusBadRequest, "Invalid Status", nil)
		return
	}

	jobId := c.Param("id")
	userId, exists := c.Get("userId")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist in context", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "user ID is not a string", nil)
		return
	}

	job, err := models.GetJobById(jobId)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "could not fecth Job", err)
		return
	}

	if job.CreatedBy != userIdStr {
		utils.RespondError(c, http.StatusUnauthorized, "You are not authorized to update this event", nil)
		return
	}

	err = updatedJob.Update(jobId)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "Could not update Job", err)
		return
	}

	utils.RespondJSON(c, http.StatusOK, "Job Updated successfully", nil)
}

// @Summary Shows stats all jobs for user
// @Description Shows stats all jobs for user
// @Tags Job
// @Security ApiKeyAuth
// @Produce  json
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /jobs/stats [GET]
func ShowStats(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		utils.RespondError(c, http.StatusUnauthorized, "userId does not exist in context", nil)
		return
	}

	// Type assertion for userId
	userIdStr, ok := userId.(string)
	if !ok {
		utils.RespondError(c, http.StatusUnauthorized, "user ID is not a string", nil)
		return
	}

	acceptedJobs, pendingJobs, declinedJobs, interviewJobs, err := models.CountStatusJobs(userIdStr)

	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "unable to count data", err)
		return
	}
	defaultStats := map[string]int{
		"accepted":  acceptedJobs,
		"pending":   pendingJobs,
		"interview": interviewJobs,
		"declined":  declinedJobs,
	}

	monthlyApplications, err := models.GetMonthlyApplications(userIdStr)
	if err != nil {
		utils.RespondError(c, http.StatusInternalServerError, "unable to count data", err)
		return
	}

	// Reverse the slice to match the JavaScript implementation
	for i, j := 0, len(monthlyApplications)-1; i < j; i, j = i+1, j-1 {
		monthlyApplications[i], monthlyApplications[j] = monthlyApplications[j], monthlyApplications[i]
	}
	utils.RespondJSON(c, http.StatusOK, "Data retrieved successfully", gin.H{
		"defaultStats":        defaultStats,
		"monthlyApplications": monthlyApplications,
	})
}
