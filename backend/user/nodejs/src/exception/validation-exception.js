export class ValidationException extends Error {
    constructor(status, message) {
        super(message)
        this.status = status
    }
}