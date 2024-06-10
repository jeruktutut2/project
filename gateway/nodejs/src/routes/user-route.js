import express from "express";
import userController from "../controllers/user-controller.js";
import sessionIdMiddleware from "../middlewares/session-id-middleware.js";

const userRouter = () => {
    const userRouter = new express.Router()
    userRouter.post("/api/v1/users/register", userController.register)
    userRouter.post("/api/v1/users/login", sessionIdMiddleware.getSessionId, userController.login)
    userRouter.post("/api/v1/users/logout", sessionIdMiddleware.getSessionId, userController.logout)
    return userRouter
}

export default {
    userRouter
}