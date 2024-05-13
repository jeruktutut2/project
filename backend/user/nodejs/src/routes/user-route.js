import express from "express";
import userController from "../controllers/user-controller.js";
import authenticationMiddleware from "../middlewares/authentication-middleware.js";

const userRouter = new express.Router()
userRouter.post("/api/v1/user/register", userController.register)
userRouter.post("/api/v1/user/login", authenticationMiddleware.getSessionId, userController.login)
userRouter.post("/api/v1/user/logout", authenticationMiddleware.authenticate, userController.logout)

export {
    userRouter
}