import userService from "../services/user-service.js";
import logMiddleware from "../middlewares/log-middleware.js";

const register = async (req, res, next) => {
    try {
        const result = await userService.register(req.body)
        return logMiddleware.logResponse(res, 201, req.requestId, {data: result, error: ""})
    } catch (error) {
        next(error)
    }
}

const login = async (req, res, next) => {
    try {
        const result = await userService.login(req.body, req.sessionId)
        res.cookie("Authorization", result, { path: "/", signed: true, httpOnly: true })
        return logMiddleware.logResponse(res, 200, req.requestId, {data: "successfully login", error: ""})
    } catch (error) {
        next(error)
    }
}

const logout = async (req, res, next) => {
    try {
        const result = userService.logout(req.sessionId)
        // signed: true
        res.cookie("Authorization", "", {path: "/", httpOnly: true, maxAge: -1})
        return logMiddleware.logResponse(res, 200, req.requestId, {data: "successfully logout", error: ""})
    } catch (error) {
        next(error)
    }
}

export default {
    register,
    login,
    logout
}