import { ResponseException } from "./response-exception.js";
const errorHandler = (error, requestId) => {
    console.log(JSON.stringify({logTime: (new Date()).toISOString(), requestId: requestId, error: error.stack}));
    if (error instanceof ResponseException) {
        throw new ResponseException(error.status, error.message)
    } else {
        throw new ResponseException(500, "internal server error")
    }
}

export default {
    errorHandler
}