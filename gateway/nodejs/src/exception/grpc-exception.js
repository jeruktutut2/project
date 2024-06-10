import logMiddleware from "../middlewares/log-middleware.js";
const grpcErrorHandler = (err, req, res) => {
    let httpCode = 0
    if (err.code === 3) {
        httpCode = 400
    } else if (err.code === 5) {
        httpCode = 403
    } else if (err.code === 13) {
        httpCode = 500
    } else {
        httpCode = 500
    }
    return logMiddleware.logResponse(res, httpCode, req.requestId, {data: "", error: JSON.parse(err.details)})
}

export default {
    grpcErrorHandler
}