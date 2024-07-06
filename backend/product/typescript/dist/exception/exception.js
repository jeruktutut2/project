"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.setInternalServerError = setInternalServerError;
function setInternalServerError() {
    return JSON.stringify([{ field: "message", message: "internal server error" }]);
}
