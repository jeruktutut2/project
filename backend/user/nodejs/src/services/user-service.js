import { validate } from "../validation/validation.js";
import { registerUserValidation, loginUserValidation } from "../validation/user-validation.js";
import { ResponseException } from "../exception/response-exception.js";
import bcrypt from "bcrypt";
import { v4 as uuid } from "uuid";
import mysqlUtil from "../utils/mysql-util.js";
import userRepository from "../repositories/user-repository.js";
import redisUtil from "../utils/redis-util.js";
import errorException from "../exception/error-exception.js";

const register = async (request, requestId, sessionId) => {
    let connection
    try {
        if (sessionId !== "" || sessionId !== undefined) {
            await redisUtil.redis.del(sessionId)
        }

        validate(registerUserValidation, request)

        connection = await mysqlUtil.mysqlPool.getConnection();

        const [rows] = await userRepository.countByEmail(connection, request.username, request.email)
        const numberOfUser = rows[0].number_of_user
        if (numberOfUser > 0) {
            throw new ResponseException(400, JSON.stringify([{field: "message", message: "username or email already exists"}]))
        }

        const user = {
            username: request.username,
            email: request.email,
            password: await bcrypt.hash(request.password, 10),
            utc: request.utc,
            createdAt: new Date().getTime()
        }
        await userRepository.create(connection, user)

        mysqlUtil.mysqlPool.releaseConnection(connection);

        return {username: user.username, email: user.email, utc: user.utc}
    } catch (error) {
        errorException.errorHandler(error, requestId)
    } finally {
        if (connection) {
            connection.release()
        }
    }
}

const login = async (request, requestId, sessionId) => {
    let connection
    try {
        if (sessionId !== "" || sessionId !== undefined) {
            await redisUtil.redis.del(sessionId)
        }
        
        validate(loginUserValidation, request)

        connection = await mysqlUtil.mysqlPool.getConnection();

        const [rows] = await userRepository.findByEmail(connection, request.email)
        const user = rows[0]
        if (!user) {
            throw new ResponseException(400, JSON.stringify([{field: "message", message: "wrong email or password"}]))
        }

        const inPasswordValid = await bcrypt.compare(request.password, user.password)
        if (!inPasswordValid) {
            throw new ResponseException(400, JSON.stringify([{field: "message", message: "wrong email or password"}]))
        }

        sessionId = uuid().toString()
        await redisUtil.redis.set(sessionId, JSON.stringify({id: user.id, username: user.username, email: user.email, userPermissions: user.userPermissions}))

        mysqlUtil.mysqlPool.releaseConnection(connection);
        return sessionId
    } catch (error) {
        errorException.errorHandler(error, requestId)
    } finally {
        if (connection) {
            connection.release()
        }   
    }
}

const logout = async (requestId, sessionId) => {
    try {
        if (sessionId !== "" || sessionId !== undefined) {
            await redisUtil.redis.del(sessionId)
        }
        return true
    } catch (error) {
        errorException.errorHandler(error, requestId)
    }
}

export default {
    register,
    login,
    logout
}