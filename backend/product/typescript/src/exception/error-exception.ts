import { ResponseException } from "./response-exception";
import { setInternalServerError } from "./exception";

export function errorHandler(error: unknown, requestId: string) {
    if (error instanceof Error) {
        console.log(JSON.stringify({logTime: (new Date()).toISOString(), app: "backend-product-typescript", requestId: requestId, error: error.stack}));
        if (error instanceof ResponseException) {
            throw new ResponseException(error.status, error.message)
        } else {
            throw new ResponseException(500, setInternalServerError())
        }
    } else {
        throw new ResponseException(500, setInternalServerError())
    }
    // console.log(JSON.stringify({logTime: (new Date()).toISOString(), app: "backend-product-typescript", requestId: requestId, error: error.stack}));
    // if (error instanceof ResponseException) {
    //     throw new ResponseException(error.status, error.message)
    // } else {
    //     throw new ResponseException(500, setInternalServerError())
    // }
}