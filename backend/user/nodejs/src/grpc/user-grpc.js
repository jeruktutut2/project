import userService from "../services/user-service.js";
import grpcException from "../exception/grpc-exception.js";

export async function Register(call, callback) {
    try {
        const sessionid = call.metadata.get("sesionid")[0] === undefined ? "" : call.metadata.get("sesionid")[0]
        const response =  await userService.register(call.request, sessionid)
        callback(null, {username: response.username, email: response.email, utc: response.utc})
    } catch (error) {
        grpcException.errorHandler(callback, call.metadata.get("requestid")[0], error)
    }
}

export async function Login(call, callback) {
    try {
        const sessionid = call.metadata.get("sessionid")[0]
        const response = await userService.login(call.request, sessionid)
        callback(null, {sessionid: response})
    } catch (error) {
        grpcException.errorHandler(callback, call.metadata.get("requestid")[0], error)
    }
}

export function Logout(call, callback) {
    callback(null, null)
}