import grpcSetup from "../setups/grpc.js";
import grpc from "@grpc/grpc-js";
import grpcException from "../exception/grpc-exception.js";
import logMiddleware from "../middlewares/log-middleware.js";

const register = async (req, res, next) => {
    console.log("register");
    res.status(200).json({data: "data", error: "error"})
}

const login = async (req, res, next) => {
    const metadata = new grpc.Metadata()
    metadata.set("requestid", req.requestId)
    metadata.set("sessionid", req.sessionId )
    return await grpcSetup.client.Login(req.body, metadata, (error, response) => {
        if (error) {
            return grpcException.grpcErrorHandler(error, res)
        }
        res.cookie("Authorization", response.sessionId, { path: "/", signed: true, httpOnly: true })
        return logMiddleware.logResponse(res, 200, req.requestId, response)
    })
}

const logout = async (req, res, next) => {
    console.log(logout);
    res.status(200).json({data: "data", error: "error"})
}

export default {
    register,
    login,
    logout
}