import grpcSetup from "../setups/grpc.js";
import grpc from "@grpc/grpc-js";
import grpcException from "../exception/grpc-exception.js";
import logMiddleware from "../middlewares/log-middleware.js";
import controller from "./controller.js";

const register = async (req, res, next) => {
    const metadata = new grpc.Metadata()
    metadata.set("requestid", req.requestId)
    return await grpcSetup.client.Register(req.body, metadata, (error, response) => {
        if (error) {
            return grpcException.grpcErrorHandler(error, req, res)
        }
        // return logMiddleware.logResponse(res, 201, req.requestId, {data: response, errors: ""})
        return controller.setResponse(res, 201, req.requestId, response, "")
    })
}

const login = async (req, res, next) => {
    const metadata = new grpc.Metadata()
    metadata.set("requestid", req.requestId)
    metadata.set("sessionid", req.sessionId )
    return await grpcSetup.client.Login(req.body, metadata, (error, response) => {
        if (error) {
            return grpcException.grpcErrorHandler(error, req, res)
        }
        res.cookie("Authorization", response.sessionid, { path: "/", signed: true, httpOnly: true })
        // return logMiddleware.logResponse(res, 200, req.requestId, {data: "successfully login", errors: ""})
        return controller.setResponse(res, 200, req.requestId, "successfully login", "")
    })
}

const logout = async (req, res, next) => {
    const sessionid = { sessionid: req.sessionId }
    return await grpcSetup.client.Logout(sessionid, (error, response) => {
        if (error) {
            return grpcException.grpcErrorHandler(error, req, res)
        }
        res.cookie("Authorization", "", {path: "/", httpOnly: true, maxAge: -1})
        // return logMiddleware.logResponse(res, 200, req.requestId, {data: "successfully logout", errors: ""})
        return controller.setResponse(res, 200, req.requestId, "successfully logout", "")
    })
}

export default {
    register,
    login,
    logout,
}