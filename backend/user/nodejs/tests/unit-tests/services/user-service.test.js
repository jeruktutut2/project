import userService from "../../../src/services/user-service.js";
import userRepository from "../../../src/repositories/user-repository.js";
import redisRepository from "../../../src/repositories/redis-repository.js";
import { ResponseException } from "../../../src/exception/response-exception.js";
import mysqlUtil from "../../../src/utils/mysql-util.js";
import bcrypt from "bcrypt";
import { ResponderBuilder } from "@grpc/grpc-js";

jest.mock("../../../src/repositories/user-repository.js")
jest.mock("../../../src/repositories/redis-repository.js")
jest.mock("../../../src/utils/mysql-util.js")

describe("when call register", () => {
    afterEach(() => {
        jest.resetAllMocks();
    });
    it("register throws error ResponseException 500 error redis repository del", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = "sessionId"

        redisRepository.del.mockImplementation((redis, sessionId) => {
            throw new ResponseException(500, "error redis repository del")
        })

        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("error redis repository del");
    })
    it("register throws error ResponseException 400 bad request validation", async () => {
        const request = {
            username: "",
            email: "",
            password: "",
            confirmpassword:"",
            utc: ""
        }
        const requestId = "requestId"
        const sessionId = ""

        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"username\",\"message\":\"username is not allowed to be empty\"},{\"field\":\"email\",\"message\":\"email is not allowed to be empty\"},{\"field\":\"password\",\"message\":\"password is not allowed to be empty\"},{\"field\":\"confirmpassword\",\"message\":\"confirmpassword is not allowed to be empty\"},{\"field\":\"utc\",\"message\":\"utc is not allowed to be empty\"}]");
    })
    it("register connection pooling throws error", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""

        mysqlUtil.getConnection.mockImplementation(() => {
            throw new ResponseException(500, "connection error")
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("connection error");
    })
    it("register repository countByEmail throws error", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""

        userRepository.countByEmail.mockImplementation((connection, username, email) => {
            throw new ResponseException(500, "repository countByEmail throws error")
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("repository countByEmail throws error");
    })
    it("register repository countByEmail result empty throw username or email already exists", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""
        userRepository.countByEmail.mockImplementation((connection, username, email) => {
             return [[], []]
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"username or email already exists\"}]");
    })
    it("register repository countByEmail result no number_of_user throw username or email already exists", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""
        userRepository.countByEmail.mockImplementation((connection, username, email) => {
            return [[{number_of_user: 1}], []]
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"username or email already exists\"}]");
    })
    it("register user repository create throw error", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""

        const user = {
            username: "username",
            email: "email",
            password: "password",
            utc: "utc",
            createdAt: 1
        }
        userRepository.countByEmail.mockImplementation((connection, username, email) => {
            return [[{number_of_user: 0}], []]
        })
        userRepository.create.mockImplementation((connection, user) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("register pool release connection throw error", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""

        const user = {
            username: "username",
            email: "email",
            password: "password",
            utc: "utc",
            createdAt: 1
        }

        userRepository.countByEmail.mockImplementation((connection, username, email) => {
            return [[{number_of_user: 0}], []]
        })
        mysqlUtil.releaseConnection.mockImplementation(() => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("register success", async () => {
        const request = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""

        const user = {
            username: "username17",
            email: "email17@email.com",
            password: "password@A1",
            utc: "+0800",
            createdAt: 1
        }

        userRepository.countByEmail.mockImplementation((connection, username, email) => {
            return [[{number_of_user: 0}], []]
        })

        mysqlUtil.releaseConnection.mockImplementation(() => {
            return
        })

        await expect(userService.register(request, requestId, sessionId)).resolves.toEqual({"email": "email17@email.com", "username": "username17", "utc": "+0800"})
    })
})

describe("when call login", () => {
    afterEach(() => {
        jest.resetAllMocks();
    });
    it("should error when sessionId is not empty", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = "sessionId"
        redisRepository.del.mockImplementation((redis, sessionId) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("should throw response exception when login validation error", async () => {
        const request = {
            email: "",
            password: ""
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"email\",\"message\":\"email is not allowed to be empty\"},{\"field\":\"password\",\"message\":\"password is not allowed to be empty\"}]");
    })
    it("should throw error when get mysql connection", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("should throw error when call findByEmail", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("should throw ResponseException when call findByEmail with zero result", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            return [[],[]]
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"wrong email or password\"}]");
    })
    it("should throw wrong email or password when compare bcrypt", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            return [[{id: 1, username: "username17", email: "email17@email.com", password: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG-", utc: "+0800", created_at: 1}],[]]
        })

        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"wrong email or password\"}]");
    })
    it("should throw error when call set sessionId to redis", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            return [[{id: 1, username: "username17", email: "email17@email.com", password: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG", utc: "+0800", created_at: 1}],[]]
        })
        redisRepository.set.mockImplementation((redis, key, value) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("should throw error when call releaseConnection mysql", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            return [[{id: 1, username: "username17", email: "email17@email.com", password: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG", utc: "+0800", created_at: 1}],[]]
        })
        redisRepository.set.mockImplementation((redis, key, value) => {
            return
        })
        mysqlUtil.releaseConnection.mockImplementation((pool, connection) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("error");
    })
    it("success", async () => {
        const request = {
            email: "email17@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        redisRepository.del.mockImplementation((redis, sessionId) => {
            return
        })
        mysqlUtil.getConnection.mockImplementation((pool) => {
            return
        })
        userRepository.findByEmail.mockImplementation((connection, email) => {
            return [[{id: 1, username: "username17", email: "email17@email.com", password: "$2a$10$MvEM5qcQFk39jC/3fYzJzOIy7M/xQiGv/PAkkoarCMgsx/rO0UaPG", utc: "+0800", created_at: 1}],[]]
        })
        redisRepository.set.mockImplementation((redis, key, value) => {
            return
        })
        mysqlUtil.releaseConnection.mockImplementation((pool, connection) => {
            return
        })
        await expect(userService.login(request, requestId, sessionId)).resolves.not.toEqual("")
    })
})

describe("when call logout", () => {
    afterEach(() => {
        jest.resetAllMocks();
    });
    it("should throw error when sessionId is not empty", async () => {
        const requestId = "requestId"
        const sessionId = "sessionId"
        redisRepository.del.mockImplementation((redis, key) => {
            throw new ResponseException(500, "error")
        })
        await expect(async () => await userService.logout(requestId, sessionId)).rejects.toThrow("error");
    })
    it("success", async () => {
        const requestId = "requestId"
        const sessionId = "sessionId"
        redisRepository.del.mockImplementation((redis, key) => {
            return
        })
        await expect(userService.logout(requestId, sessionId)).resolves.toEqual(true)
    })
})