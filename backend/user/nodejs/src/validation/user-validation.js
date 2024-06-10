import Joi from "joi";

const customPasswordValidation = (value, helpers) => {
    let lowercase = false
    let uppercase = false
    let number = false
    let specialCharacter = false
    let undefinedCharacter = false
    
    for (let i = 0; i < value.length; i++) {
        const element = value[i];
        if (/[a-z]/.test(element)) {
            lowercase = true
        } else if (/[A-Z]/.test(element)) {
            uppercase = true
        } else if (/[0-9]/.test(element)) {
            number = true
        } else if (element === "!" 
            || element === "@" 
            // || element === "#" 
            // || element === "$"
            // || element === "%"
            // || element === "^"
            // || element === "&"
            // || element === "*"
            // || element === "(" 
            // || element === ")"
            || element === "-"
            || element === "_"
            // || element === "="
            // || element === "+"
            // || element === "["
            // || element === "]"
            // || element === "{"
            // || element === "}"
            // || element === `\\`
            // || element === "|"
            // || element === ";"
            // || element === ":"
            // || element === "'"
            // || element === '"'
            // || element === ","
            // || element === "."
            // || element === "<"
            // || element === ">"
            // || element === "/"
            // || element === "?"
            // || element === " "
        ) {
            specialCharacter = true
        } else {
            undefinedCharacter = true
        }
    }

    if (!lowercase || !uppercase || !number || !specialCharacter || undefinedCharacter) {
        return helpers.error('password.wrong');
    } 
    return value
}

const checkConfirmPassword = (value, helpers) => {
    if (value.password !== value.confirmpassword) {
        return helpers.error('confirm.password.wrong');
    }
    return value
}

export const registerUserValidation = Joi.object({
    username: Joi.string().min(3).max(100).required(),
    email: Joi.string().email().required(),
    // throw new Error('Formula contains invalid trailing operator');
    // Error: Invalid template variable "" fails due to: Formula contains invalid trailing operator
    // when put {} in the error message
    password: Joi.string().min(8).max(20).required().custom(customPasswordValidation).message({'password.wrong': "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"}),
    confirmpassword: Joi.string().min(8).max(20).required().custom(customPasswordValidation).message({'password.wrong': "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"}),
    utc: Joi.string().min(4).max(6).required()
    // .custom(checkConfirmPassword).message({"confirm.password.wrong": "password and confirm password is different"})
}).custom(checkConfirmPassword).message({"confirm.password.wrong": "password and confirm password is different"});

export const loginUserValidation = Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(8).required()
});