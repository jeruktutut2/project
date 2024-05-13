import logMiddleware from "./log-middleware.js";
import redisApp from "../application/redis.js";

const authenticate = async (req, res, next) => {
    const authorization = req.signedCookies["Authorization"]
    if (!authorization) {
        return logMiddleware.logResponse(res, 401, req.requestId, {data: "", error: "unauthorized"})
    }
    const user = await redisApp.redis.get(authorization)
    if (!user) {
        return logMiddleware.logResponse(res, 401, req.requestId, {data: "", error: "unauthorized"})
    }
    req.sessionId = authorization
    req.user = JSON.parse(user)
    next()
}

const getSessionId = async (req, res, next) => {
    const authorization = req.signedCookies["Authorization"]
    if (authorization) {
        req.sessionId = authorization
    }
    next()
}

export default {
    authenticate,
    getSessionId
}