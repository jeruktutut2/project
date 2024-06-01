import grpc from "@grpc/grpc-js";
import { ResponseException } from "../exception/response-exception.js";
import { ValidationException } from "../exception/validation-exception.js";

const errorHandler = (callback, requestId, error) => {
    console.log(JSON.stringify({logTime: (new Date()).toISOString(), requestId: requestId, error: error.stack}));
    if (error instanceof ResponseException) {
        callback({code: grpc.status.INVALID_ARGUMENT, details: error.message}, null)
    } else if (error instanceof ValidationException) {
        callback({code: grpc.status.INVALID_ARGUMENT, details: error.message}, null)
    } else {
        callback({code: grpc.status.INTERNAL, details: "internal server error"}, null)
    }
}

export default {
    errorHandler
}