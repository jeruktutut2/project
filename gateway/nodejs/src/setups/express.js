import express from "express";
import cookieParser from "cookie-parser";
import requestIdMiddleware from "../middlewares/request-id-middleware.js";
import logMiddleware from "../middlewares/log-middleware.js";

const setExpress = async (userRouter) => {
    const app = express()
    app.use(express.json())
    app.use(cookieParser(process.env.PROJECT_GATEWAY_COOKIE_SECRET))
    app.use(requestIdMiddleware.setRequestId)
    app.use(logMiddleware.setLogRequest)
    app.use(userRouter)
    return app
}

const startExpress = async (app) => {
    app = app.listen(process.env.PROJECT_GATEWAY_APPLICATION_PORT, () => {
        console.log(new Date(), "express: started on port "+process.env.PROJECT_GATEWAY_APPLICATION_PORT);
    })
    return app
}

const stopExpress = async (app) => {
    app.close( () => {
        console.log(new Date(), "express: has closed all connections.");
    })
}

export default {
    setExpress,
    startExpress,
    stopExpress
}