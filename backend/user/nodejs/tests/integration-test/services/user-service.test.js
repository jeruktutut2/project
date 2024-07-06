import userService from "../../../src/services/user-service.js";
import mysqlUtil from "../../../src/utils/mysql-util.js";
import redisUtil from "../../../src/utils/redis-util.js";

describe("when call register", () => {
    beforeAll(async () => {
        mysqlUtil.mysqlPool = mysqlUtil.newConnection()
        redisUtil.redis = await redisUtil.newConnection()
    });
    afterEach(() => {
        jest.resetAllMocks();
    });
    afterAll(async () => {
        mysqlUtil.closeConnection(mysqlUtil.mysqlPool)
        await redisUtil.closeConnection(redisUtil.redis)
      });

    it("should throw error sessionId is not empty", async () => {
        const request = {
            username: "username1",
            email: "email1@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = "sessionId"
        await expect(async () => await userService.register(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"username or email already exists\"}]");
    })
    it("should throw ResponseException validation request", async () => {
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

    // why skip, because if it is not skip, it will keep create data to database
    it.skip("success", async () => {
        const request = {
            username: "username18",
            email: "email18@email.com",
            password: "password@A1",
            confirmpassword:"password@A1",
            utc: "+0800"
        }
        const requestId = "requestId"
        const sessionId = ""
        await expect(userService.register(request, requestId, sessionId)).resolves.toEqual({"email": "email18@email.com", "username": "username18", "utc": "+0800"})
    })
})

describe("when call login", () => {
    beforeAll(async () => {
        mysqlUtil.mysqlPool = mysqlUtil.newConnection()
        redisUtil.redis = await redisUtil.newConnection()
    })
    afterEach(() => {
        jest.resetAllMocks();
    });
    afterAll(async () => {
        mysqlUtil.closeConnection(mysqlUtil.mysqlPool)
        await redisUtil.closeConnection(redisUtil.redis)
    })
    it("should throw ResponseException when request is empty", async () => {
        const request = {
            email: "",
            password: ""
        }
        const requestId = "requestId"
        const sessionId = ""
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"email\",\"message\":\"email is not allowed to be empty\"},{\"field\":\"password\",\"message\":\"password is not allowed to be empty\"}]");
    })
    it("should throw ResponseException when email doesn't exists", async () => {
        const request = {
            email: "email19@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        await expect(async () => await userService.login(request, requestId, sessionId)).rejects.toThrow("[{\"field\":\"message\",\"message\":\"wrong email or password\"}]");
    })
    it("success", async () => {
        const request = {
            email: "email18@email.com",
            password: "password@A1"
        }
        const requestId = "requestId"
        const sessionId = ""
        await expect(userService.login(request, requestId, sessionId)).resolves.not.toEqual("")
    })
})

describe("when call logout", () => {
    beforeAll(async () => {
        mysqlUtil.mysqlPool = mysqlUtil.newConnection()
        redisUtil.redis = await redisUtil.newConnection()
    })
    afterEach(() => {
        jest.resetAllMocks();
    });
    afterAll(async () => {
        mysqlUtil.closeConnection(mysqlUtil.mysqlPool)
        await redisUtil.closeConnection(redisUtil.redis)
    })
    it("success", async () => {
        const requestId = "requestId"
        const sessionId = ""
        await expect(userService.logout(requestId, sessionId)).resolves.toEqual(true)
    })
})