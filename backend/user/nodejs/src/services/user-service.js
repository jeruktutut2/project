import { validate } from "../validation/validation.js";
import { registerUserValidation, loginUserValidation } from "../validation/user-validation.js";
import prismaClient from "../application/mysql.js";
import { ResponseException } from "../exception/response-exception.js";
import { ValidationException } from "../exception/validation-exception.js";
import bcrypt from "bcrypt";
import redisApp from "../application/redis.js";
import { v4 as uuid } from "uuid";

const register = async (request) => {
    if (request.password !== request.confirmPassword) {
        throw new ValidationException(400, JSON.stringify([{field: "password and confirmPassword", message: "password and confirm password is different"}]))
    }

    validate(registerUserValidation, request)

    const numberOfUser = await prismaClient.user.count({
        where: {
            OR: [
                {
                    username: request.username
                },
                {
                    email: request.email
                }
            ]
        }
    })

    if (numberOfUser > 0) {
        throw new ResponseException(400, "username or email already exists")
    }

    const user = {}
    user.username = request.username
    user.email = request.email
    user.password = await bcrypt.hash(request.password, 10)
    user.utc = request.utc
    user.created_at = new Date().getTime()
    return await prismaClient.user.create({
        data: user,
        select: {
            username : true,
            email: true
        }
    })
}

const login = async (request, sessionId) => {
    const userSession = await redisApp.redis.get(sessionId)
    console.log("userSession:", userSession);
    if (userSession) {
        const sessiondel = await redisApp.redis.del(sessionId)
    }

    validate(loginUserValidation, request)

    const user = await prismaClient.user.findFirst({
        where: {
            email: request.email
        },
        select: {
            id: true,
            username: true,
            email: true,
            password: true,
            utc: true,
            userPermissions: true
        }
    })
    if (!user) {
        throw new ResponseException(400, "wrong email or password")
    }

    const inPasswordValid = await bcrypt.compare(request.password, user.password)
    if (!inPasswordValid) {
        throw new ResponseException(400, "wrong email or password")
    }

    sessionId = uuid().toString()
    await redisApp.redis.set(sessionId, JSON.stringify({id: user.id, username: user.username, email: user.email, userPermissions: user.userPermissions}))
    return sessionId
}

const logout = async (sessionId) => {
    await redisApp.redis.del(sessionId)
    return true
}

export default {
    register,
    login,
    logout
}