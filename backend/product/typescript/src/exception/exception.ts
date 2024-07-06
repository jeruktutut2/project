export function setInternalServerError(): string {
    return JSON.stringify([{field: "message", message: "internal server error"}])
}