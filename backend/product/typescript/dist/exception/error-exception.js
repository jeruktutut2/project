"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.errorHandler = errorHandler;
const response_exception_1 = require("./response-exception");
const exception_1 = require("./exception");
function errorHandler(error, requestId) {
    if (error instanceof Error) {
        console.log(JSON.stringify({ logTime: (new Date()).toISOString(), app: "backend-product-typescript", requestId: requestId, error: error.stack }));
        if (error instanceof response_exception_1.ResponseException) {
            throw new response_exception_1.ResponseException(error.status, error.message);
        }
        else {
            throw new response_exception_1.ResponseException(500, (0, exception_1.setInternalServerError)());
        }
    }
    else {
        throw new response_exception_1.ResponseException(500, (0, exception_1.setInternalServerError)());
    }
    // console.log(JSON.stringify({logTime: (new Date()).toISOString(), app: "backend-product-typescript", requestId: requestId, error: error.stack}));
    // if (error instanceof ResponseException) {
    //     throw new ResponseException(error.status, error.message)
    // } else {
    //     throw new ResponseException(500, setInternalServerError())
    // }
}
