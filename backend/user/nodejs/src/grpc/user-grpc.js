import userService from "../services/user-service.js";
import grpcException from "../exception/grpc-exception.js";

export async function Register(call, callback) {
    try {
        const requestId = call.metadata.get("requestid")[0]
        const sessionId = call.metadata.get("sesionid")[0] === undefined ? "" : call.metadata.get("sesionid")[0]
        const response =  await userService.register(call.request, requestId, sessionId)
        callback(null, {username: response.username, email: response.email, utc: response.utc})
    } catch (error) {
        grpcException.errorHandler(callback, error)
    }
}

export async function Login(call, callback) {
    try {
        const requestId = call.metadata.get("requestid")[0]
        const sessionid = call.metadata.get("sessionid")[0]
        const response = await userService.login(call.request, requestId, sessionid)
        callback(null, {sessionid: response})
    } catch (error) {
        grpcException.errorHandler(callback, error)
    }
}

export async function Logout(call, callback) {
    try {
        const requestId = call.metadata.get("requestid")[0]
        const sessionid = call.request.sessionid
        const response = await userService.logout(requestId, sessionid)
        callback(null, {msg: "successfully logout"})
    } catch (error) {
        grpcException.errorHandler(callback, error)
    }
}