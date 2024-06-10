import grpc from "@grpc/grpc-js";
import { ResponseException } from "../exception/response-exception.js";
import { ValidationException } from "../exception/validation-exception.js";

const errorHandler = (callback, error) => {
    if (error instanceof ResponseException) {
        if (error.status === 400) {
            callback({code: grpc.status.INVALID_ARGUMENT, details: error.message}, null)
        } else if (error.status === 500) {
            callback({code: grpc.status.INTERNAL, details: "internal server error"}, null)
        } else {
            callback({code: grpc.status.INTERNAL, details: "internal server error"}, null)
        }
    } else if (error instanceof ValidationException) {
        callback({code: grpc.status.INVALID_ARGUMENT, details: error.message}, null)
    } else {
        callback({code: grpc.status.INTERNAL, details: "internal server error"}, null)
    }
}

export default {
    errorHandler
}