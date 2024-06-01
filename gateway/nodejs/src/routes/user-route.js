import express from "express";
import userController from "../controllers/user-controller.js";

const userRouter = () => {
    const userRouter = new express.Router()
    userRouter.post("/api/v1/user/register", userController.register)
    userRouter.post("/api/v1/user/login", userController.login)
    userRouter.post("/api/v1/user/logout", userController.logout)
    return userRouter
}

export default {
    userRouter
}