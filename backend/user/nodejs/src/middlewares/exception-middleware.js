import { ResponseException } from "../exception/response-exception.js";
import { ValidationException } from "../exception/validation-exception.js";
import logMiddleware from "./log-middleware.js";

const exceptionMiddleware = async (err, req, res, next) => {
    if (!err) {
        next()
        return
    }
    console.log(JSON.stringify({logTime: new Date(), requestId: req.requestId, error: err.stack}));
    if (err instanceof ResponseException) {
        return logMiddleware.logResponse(res, err.status, req.requestId, {data: "", error: err.message})    
    } else if (err instanceof ValidationException) {
        return logMiddleware.logResponse(res, err.status, req.requestId, {data: "", error: JSON.parse(err.message)})
    } else {
        return logMiddleware.logResponse(res, 500, req.requestId, {data: "", error: "internal server error"})
    }
}

export default {
    exceptionMiddleware
}