import logMiddleware from "../middlewares/log-middleware.js";
const setResponse = (res, httpCode, requestId, data, errors) => {
    console.log("res:", httpCode, requestId, data, errors);
    const response = {data: data, errors: errors}
    logMiddleware.logResponse(requestId, response)
    return res.status(httpCode).json(response)
}

export default {
    setResponse
}