"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ResponseException = void 0;
class ResponseException extends Error {
    constructor(status, message) {
        super(message);
        this.status = status;
        this.message = message;
    }
}
exports.ResponseException = ResponseException;
