import { ValidationException } from "../exception/validation-exception.js";

export const validate = (schema, request) => {
    const result = schema.validate(request, {
        abortEarly: false,
        allowUnknown: false
    })

    if (result.error) {
        const errorMessage = []
        for (let i = 0; i < result.error.details.length; i++) {
            const field = result.error.details[i].path[0]
            const fieldMessage = result.error.details[i].message.replaceAll('"', '')
            const message = {field: field, message: fieldMessage}
            errorMessage.push(message)
        }
        throw new ValidationException(400, JSON.stringify(errorMessage))
    } else {
        return result.value
    }
}