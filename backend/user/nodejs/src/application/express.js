import express from "express";
import cookieParser from "cookie-parser";
import requestIdMiddleware from "../middlewares/request-id-middleware.js";
import logMiddleware from "../middlewares/log-middleware.js";
import { userRouter } from "../routes/user-route.js";
import exceptionMiddleware from "../middlewares/exception-middleware.js";

export const web = express()
web.use(express.json())
web.use(cookieParser(process.env.COOKIE_SECRET))
web.use(requestIdMiddleware.setRequestId)
web.use(logMiddleware.setLogRequest)
web.use(userRouter)
web.use(exceptionMiddleware.exceptionMiddleware)