import grpc from "@grpc/grpc-js";
import { ResponseException } from "../exception/response-exception.js";
import exception from "./exception.js";

const errorHandler = (callback, error) => {
    if (error instanceof ResponseException) {
        if (error.status === 400) {
            callback({code: grpc.status.INVALID_ARGUMENT, details: error.message}, null)
        } else if (error.status === 500) {
            callback({code: grpc.status.INTERNAL, details: exception.setInternalServerErrorMessage()}, null)
        } else {
            callback({code: grpc.status.INTERNAL, details: exception.setInternalServerErrorMessage()}, null)
        }
    } else {
        callback({code: grpc.status.INTERNAL, details: exception.setInternalServerErrorMessage()}, null)
    }
}

export default {
    errorHandler
}