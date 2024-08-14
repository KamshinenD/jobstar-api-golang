package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"jobstar.com/api/db"
)

type JobType string

const (
	FullTime   JobType = "Full-Time"
	PartTime   JobType = "Part-Time"
	Contract   JobType = "Contract"
	Internship JobType = "Internship"
	Remote     JobType = "Remote"
)

type Status string

const (
	Interview Status = "interview"
	Accepted  Status = "Accepted"
	Declined  Status = "declined"
	Pending   Status = "pending"
)

type Job struct {
	ID          string    `json:"id"`
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	JobLocation string    `json:"jobLocation"`
	Status      Status    `json:"status"`
	JobType     JobType   `json:"jobType"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type MonthlyApplication struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// IsValid checks if the JobType is valid
func (jt JobType) JobTypeIsValid() bool {
	switch jt {
	case FullTime, PartTime, Contract, Internship, Remote:
		return true
	}
	return false
}

// IsValid checks if the status is valid
func (s Status) StatusIsValid() bool {
	switch s {
	case Interview, Accepted, Declined, Pending:
		return true
	}
	return false
}

func (j *Job) SaveJob() error {
	query := `INSERT INTO jobs(company, position, jobLocation, status, jobType, createdBy, createdAt ) 
	VALUES($1, $2, $3, $4, $5, $6, NOW()) RETURNING id`

	// Use QueryRow to execute the query and retrieve the generated ID
	err := db.DB.QueryRow(query, j.Company, j.Position, j.JobLocation, j.Status, j.JobType, j.CreatedBy).Scan(&j.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetJobs(userId string) ([]Job, error) {
	query := "SELECT id, company, position, jobLocation, status, jobType, createdBy, createdAt FROM jobs WHERE createdBy = $1"
	rows, err := db.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []Job
	for rows.Next() {
		var job Job
		err := rows.Scan(&job.ID, &job.Company, &job.Position, &job.JobLocation, &job.Status, &job.JobType, &job.CreatedBy, &job.CreatedAt)

		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	// Check if there was an error iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func GetUserJobById(id, userId string) (*Job, error) {
	query := "SELECT id, company, position, jobLocation, status, jobType, createdBy, createdAt FROM jobs WHERE id=$1 AND createdBy = $2"
	row := db.DB.QueryRow(query, id, userId)

	var job Job
	err := row.Scan(&job.ID, &job.Company, &job.Position, &job.JobLocation, &job.Status, &job.JobType, &job.CreatedBy, &job.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("job not found") // Job not found in the database
		}
		return nil, err // Other errors, e.g., connection issues, etc.
	}

	return &job, nil
	//pls note that we had to use pointer for event so that it can take a nil value when there is an error
}

func GetJobById(id string) (*Job, error) {
	query := "SELECT * FROM jobs WHERE id=$1"
	row := db.DB.QueryRow(query, id)

	var job Job
	err := row.Scan(&job.ID, &job.Company, &job.Position, &job.JobLocation, &job.Status, &job.JobType, &job.CreatedAt, &job.CreatedBy)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("job not found") // Job not found in the database
		}
		return nil, err // Other errors, e.g., connection issues, etc.
	}

	return &job, nil
	//pls note that we had to use pointer for event so that it can take a nil value when there is an error
}

func (j Job) Delete() error {
	query := "DELETE FROM jobs WHERE id = $1"

	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(j.ID)

	if err != nil {
		return err
	}
	return nil
}

func (job Job) Update(jobId string) error {
	query := `
		UPDATE jobs
		SET company=$1, position=$2, jobLocation=$3, status=$4, jobType=$5
		WHERE id=$6
	`

	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	result, err := stmt.Exec(job.Company, job.Position, job.JobLocation, job.Status, job.JobType, jobId)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("job with this ID not found") // No rows affected, meaning no job with the given ID was found
	}

	return err
}

func CountStatusJobs(userId string) (int, int, int, int, error) {

	var acceptedJobs, pendingJobs, declinedJobs, interviewJobs int

	// Count accepted jobs
	query := "SELECT COUNT(*) FROM jobs WHERE createdBy = $1 AND status = $2"
	err := db.DB.QueryRow(query, userId, "Accepted").Scan(&acceptedJobs)
	if err != nil {
		log.Printf("Error counting accepted jobs: %v", err)
		return 0, 0, 0, 0, err
	}

	// Count pending jobs
	err = db.DB.QueryRow(query, userId, "pending").Scan(&pendingJobs)
	if err != nil {
		log.Printf("Error counting pending jobs: %v", err)
		return 0, 0, 0, 0, err
	}

	// Count declined jobs
	err = db.DB.QueryRow(query, userId, "declined").Scan(&declinedJobs)
	if err != nil { // Corrected error handling here
		log.Printf("Error counting declined jobs: %v", err)
		return 0, 0, 0, 0, err
	}

	// Count interview jobs
	err = db.DB.QueryRow(query, userId, "interview").Scan(&interviewJobs)
	if err != nil {
		log.Printf("Error counting interview jobs: %v", err)
		return 0, 0, 0, 0, err
	}

	return acceptedJobs, pendingJobs, declinedJobs, interviewJobs, nil
}

func GetMonthlyApplications(userId string) ([]MonthlyApplication, error) {
	query := `
		SELECT 
			EXTRACT(YEAR FROM createdAt) AS year, 
			EXTRACT(MONTH FROM createdAt) AS month, 
			COUNT(*) AS count
		FROM jobs 
		WHERE createdBy = $1
		GROUP BY year, month
		ORDER BY year DESC, month DESC
		LIMIT 6;
	`

	rows, err := db.DB.Query(query, userId)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var monthlyApplications []MonthlyApplication

	for rows.Next() {
		var year, month int
		var count int

		if err := rows.Scan(&year, &month, &count); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}

		// Format the date as "Jan 2006"
		date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("Jan 2006")

		monthlyApplications = append(monthlyApplications, MonthlyApplication{
			Date:  date,
			Count: count,
		})
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error during row iteration: %v", err)
		return nil, err
	}

	return monthlyApplications, nil
}
