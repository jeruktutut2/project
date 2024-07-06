import { ResponseException } from "./response-exception.js";
import exception from "./exception.js";
const errorHandler = (error, requestId) => {
    console.log(JSON.stringify({logTime: (new Date()).toISOString(), app: "backend-user-nodejs", requestId: requestId, error: error.stack}));
    if (error instanceof ResponseException) {
        throw new ResponseException(error.status, error.message)
    } else {
        throw new ResponseException(500, exception.setInternalServerErrorMessage())
    }
}

export default {
    errorHandler
}