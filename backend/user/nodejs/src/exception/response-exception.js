export class ResponseException extends Error {
    constructor(status, message) {
        super(message)
        this.status = status
    }
}