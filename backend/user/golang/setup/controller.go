package setup

import controller "project-user/controllers"

type ControllerSetup struct {
	UserController controller.UserController
}

func NewControllerSetup(serviceSetup *ServiceSetup) *ControllerSetup {
	return &ControllerSetup{
		UserController: controller.NewUserController(serviceSetup.UserService),
	}
}
