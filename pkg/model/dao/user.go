package dao

import (
	"chat/pkg/dto/request"
	"chat/pkg/model/entity"
	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
)

// UserDao
type UserDao struct {
	BaseDao
}

// User return UserDao pointer with db connection
func User(db *gorm.DB) *UserDao {
	dao := &UserDao{}
	dao.db = db
	return dao
}

// Auth 校验用户
func (dao *UserDao) Auth(req request.AuthRequest) (*entity.User,error) {
	query := dao.db.Table(entity.TableUser).
		Where("tel = ?", req.Tel)

	user := new(entity.User)
	if err := query.First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (dao *UserDao) GetUsers() (user []entity.User,err error) {
	query := dao.db

	if err := query.Find(&user).Error; err != nil {
		return nil,err
	}
	return user, nil
}

func (dao *UserDao) GetUserByID(userid,touserid int) (user []entity.User,err error) {
	query := dao.db.Table(entity.TableUser).
		Where("id = ? or id = ?",userid,touserid)

	if err := query.Find(&user).Error; err != nil {
		return nil,err
	}
	return user, nil
}

//// Register 注册用户
//func (dao *UserDao) Register(req request.UserRegisterRequest) error {
//	return nil
//}
