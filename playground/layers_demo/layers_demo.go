package main

import (
	"errors"
	"fmt"
)

// ---- Domain（真实食材：数据库里存的、程序内部真正在用的完整数据）----
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string // 内部机密字段，绝对不能返回给调用方
}

// ---- DTO（端给客人的那盘菜：跨层传递时用的"打包格式"）----
// 只放 Controller 需要返回给外部调用方的字段，PasswordHash 天生不在这里。
type UserDTO struct {
	ID    int
	Name  string
	Email string
}

// ---- DAO / Repository（仓库管理员：只管存取，不管业务）----
type UserRepository struct {
	data map[int]User // 用 map 模拟数据库
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		data: map[int]User{
			1: {ID: 1, Name: "王宏成", Email: "wang@example.com", PasswordHash: "hashed-secret-abc"},
		},
	}
}

func (r *UserRepository) FindByID(id int) (User, error) {
	fmt.Println("  [仓库/DAO] 去数据库里找 id =", id)
	user, ok := r.data[id]
	if !ok {
		return User{}, errors.New("user not found")
	}
	return user, nil
}

// ---- Service（大厨：业务逻辑判断 + 把 Domain 转成 DTO）----
type UserService struct {
	repo *UserRepository
}

func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id int) (UserDTO, error) {
	fmt.Println(" [大厨/Service] 收到订单，开始处理业务逻辑")
	user, err := s.repo.FindByID(id)
	if err != nil {
		return UserDTO{}, err
	}

	// 业务逻辑该发生的地方（权限检查、状态判断……），
	// 同时负责把 Domain（带敏感字段）转换成 DTO（对外安全的形状）。
	dto := UserDTO{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return dto, nil
}

// ---- Controller（服务员：接单、转达、端菜）----
type UserController struct {
	service *UserService
}

func NewUserController(service *UserService) *UserController {
	return &UserController{service: service}
}

func (c *UserController) HandleGetUser(id int) {
	fmt.Println("[服务员/Controller] 收到请求：查询用户", id)
	dto, err := c.service.GetUser(id)
	if err != nil {
		fmt.Println("[服务员/Controller] 出错了：", err)
		return
	}
	fmt.Printf("[服务员/Controller] 端菜给客人：%+v\n\n", dto)
}

func main() {
	repo := NewUserRepository()
	service := NewUserService(repo)
	controller := NewUserController(service)

	controller.HandleGetUser(1) // 正常查询
	controller.HandleGetUser(2) // 查询不存在的用户，走错误分支
}
