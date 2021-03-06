package entity

import "github.com/zutim/ego/component/mysql"

type BaseEntity struct {
	mysql.Model
	IsDeleted int `json:"is_deleted" gorm:"column:is_deleted"`
}



const (
	SoftDeleteCondition = "is_deleted = 0"
)
