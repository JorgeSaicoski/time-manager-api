package models

import (
	"time"

	"gorm.io/gorm"
)

type TotalTime struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	UserID     int64  `gorm:"not null"`
	CompanyID  *int64 // Optional company association
	StartTime  time.Time
	FinishTime time.Time
	WorkTimes  []WorkTime `gorm:"foreignKey:TotalTimeID"`
	BreakTime  *BreakTime `gorm:"foreignKey:TotalTimeID;constraint:OnDelete:CASCADE"`
	Brb        *Brb       `gorm:"foreignKey:TotalTimeID;constraint:OnDelete:CASCADE"`
	Closed     bool
}

type WorkTime struct {
	gorm.Model
	ID          int64 `gorm:"primaryKey"`
	TotalTimeID int64
	UserID      int64  `gorm:"not null"`
	CompanyID   *int64 // Optional company association
	StartTime   time.Time
	Duration    time.Duration
	Projects    []Project `gorm:"many2many:work_time_projects;"`
	Closed      bool
	Trustworthy bool `gorm:"default:true"`
}

type User struct {
	gorm.Model
	ID            int64  `gorm:"primaryKey"`
	Email         string `gorm:"uniqueIndex;not null"`
	Password      string `gorm:"not null"`
	Name          string
	IsSystemAdmin bool        `gorm:"default:false"`
	Companies     []Company   `gorm:"many2many:user_company_roles;"`
	TotalTimes    []TotalTime `gorm:"foreignKey:UserID"`
	Projects      []Project   `gorm:"foreignKey:OwnerID"`
}

type Company struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Users       []User    `gorm:"many2many:user_company_roles;"`
	Projects    []Project `gorm:"foreignKey:CompanyID"`
}

type UserCompanyRole struct {
	gorm.Model
	ID        int64   `gorm:"primaryKey"`
	CompanyID int64   `gorm:"not null"`
	UserID    int64   `gorm:"not null"`
	Role      string  `gorm:"type:enum('admin','manager','employee');not null"`
	User      User    `gorm:"foreignKey:UserID"`
	Company   Company `gorm:"foreignKey:CompanyID"`
}

type Project struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	CompanyID   *int64 // Optional: null means it's a personal/freelance project
	OwnerID     int64  `gorm:"not null"` // User who created the project
	StartTime   time.Time
	Duration    time.Duration
	Closed      bool
	Cost        *Cost      `gorm:"foreignKey:ProjectID"`
	WorkTimes   []WorkTime `gorm:"many2many:work_time_projects;"`
	Tasks       []Task     `gorm:"foreignKey:ProjectID"`
	Members     []User     `gorm:"many2many:project_members;"`
	Owner       User       `gorm:"foreignKey:OwnerID"`
	Company     *Company   `gorm:"foreignKey:CompanyID"`
}

type WorkTimeProject struct {
	ID         int64     `gorm:"primaryKey"`
	WorkTimeID int64     `gorm:"primaryKey"`
	ProjectID  int64     `gorm:"primaryKey"`
	UserID     int64     `gorm:"not null"`
	CompanyID  *int64    // Optional company association
	StartTime  time.Time `gorm:"not null"`
	Duration   time.Duration
	Closed     bool
	WorkTime   WorkTime `gorm:"foreignKey:WorkTimeID"`
	Project    Project  `gorm:"foreignKey:ProjectID"`
}

type Task struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	ProjectID   int64  `gorm:"not null"`
	AssigneeID  int64  // User assigned to the task
	Title       string `gorm:"not null"`
	Description string
	Deadline    time.Time
	Status      string `gorm:"type:enum('todo','in_progress','done');default:'todo'"`
	Closed      bool
}

type Cost struct {
	gorm.Model
	ID        int64 `gorm:"primaryKey"`
	ProjectID int64 `gorm:"uniqueIndex"`
	Duration  time.Duration
	HourCost  int
}

type BreakTime struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	TotalTimeID int64  `gorm:"not null"`
	UserID      int64  `gorm:"not null"`
	CompanyID   *int64 // Optional company association
	StartTime   time.Time
	Duration    time.Duration
	Active      bool `gorm:"default:true"`
}

type Brb struct {
	gorm.Model
	ID          int64  `gorm:"primaryKey"`
	TotalTimeID int64  `gorm:"not null"`
	UserID      int64  `gorm:"not null"`
	CompanyID   *int64 // Optional company association
	StartTime   time.Time
	Duration    time.Duration
	Active      bool `gorm:"default:true"`
}

type ProjectMember struct {
	gorm.Model
	ID         int64  `gorm:"primaryKey"`
	ProjectID  int64  `gorm:"not null"`
	UserID     int64  `gorm:"not null"`
	Role       string `gorm:"type:enum('owner','manager','member');not null;default:'member'"`
	IsExternal bool   `gorm:"default:false"` // For freelancers/external collaborators
}

type ResolutionTracker struct {
	gorm.Model
	ID       int64 `gorm:"primaryKey"`
	UserID   int64 `gorm:"not null"`
	Day      time.Time
	Category string `gorm:"size:255"`
	Closed   bool
	Units    []ResolutionUnit `gorm:"foreignKey:TrackerID"`
}

type ResolutionUnit struct {
	gorm.Model
	ID         int64             `gorm:"primaryKey"`
	TrackerID  int64             `gorm:"not null"`
	UserID     int64             `gorm:"not null"`
	CompanyID  *int64            // Optional company association
	Tracker    ResolutionTracker `gorm:"foreignKey:TrackerID"`
	Identifier string            `gorm:"size:255"`
	Resolved   bool
}
