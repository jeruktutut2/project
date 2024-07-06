"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.MysqlUtil = void 0;
// import mysql, { PoolOptions } from 'mysql2';
const promise_1 = __importDefault(require("mysql2/promise"));
// import { Pool } from 'mysql2/promise';
// import { Pool } from 'mysql2/typings/mysql/lib/Pool';
// import { Pool } from "mysql2/typings/mysql/lib/Pool";
// import { Pool, PoolOptions, createPool } from "mysql2/promise";
// export let pool: any
// const newConnection = () => {
// // export async function newConnection() {
//     console.log(new Date(), "mysql: connecting to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
//     const access: PoolOptions = {
//         host: process.env.PROJECT_PRODUCT_MYSQL_HOST,
//         user: process.env.PROJECT_PRODUCT_MYSQL_USERNAME,
//         password: process.env.PROJECT_PRODUCT_MYSQL_PASSWORD,
//         database: process.env.PROJECT_PRODUCT_MYSQL_DATABASE,
//         connectionLimit: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_OPEN_CONNECTION ?? "0"),
//         maxIdle: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_IDLE_CONNECTION ?? "0"),
//         idleTimeout: parseInt(process.env.PROJECT_PRODUCT_MYSQL_CONNECTION_MAX_IDLETIME ?? "0")
//     };
//     const pool = mysql.createPool(access);
//     console.log(new Date(), "mysql: connected to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
//     return pool
// }
// const getConnection = async (pool: Pool)  => {
//     return await pool.getConnection()
// }
// export default {
//     newConnection,
//     pool
// }
// export let pool = null
// export async function newConnection(): Promise<Pool> {
//     console.log(new Date(), "mysql: connecting to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
//     const poolOptions: PoolOptions = {
//          host: process.env.PROJECT_PRODUCT_MYSQL_HOST,
//          user: process.env.PROJECT_PRODUCT_MYSQL_USERNAME,
//          password: process.env.PROJECT_PRODUCT_MYSQL_PASSWORD,
//          database: process.env.PROJECT_PRODUCT_MYSQL_DATABASE,
//          connectionLimit: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_OPEN_CONNECTION ?? "0"),
//          maxIdle: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_IDLE_CONNECTION ?? "0"),
//          idleTimeout: parseInt(process.env.PROJECT_PRODUCT_MYSQL_CONNECTION_MAX_IDLETIME ?? "0")
//      };
//      const pool = createPool(poolOptions)
//      console.log(new Date(), "mysql: connected to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
//      return pool
// }
class MysqlUtil {
    constructor() {
        var _a, _b, _c;
        console.log(new Date(), "mysql: connecting to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
        // const hostport = process.env.PROJECT_PRODUCT_MYSQL_HOST.split(":")
        // const host = hostport[0]
        const connectionString = `mysql://` + process.env.PROJECT_PRODUCT_MYSQL_USERNAME + `:` + process.env.PROJECT_PRODUCT_MYSQL_PASSWORD + `@` + process.env.PROJECT_PRODUCT_MYSQL_HOST + `/` + process.env.PROJECT_PRODUCT_MYSQL_DATABASE + ``;
        const access = {
            // host: process.env.PROJECT_PRODUCT_MYSQL_HOST,
            // port: Number(process.env.PROJECT_PRODUCT_MYSQL_PORT),
            // user: process.env.PROJECT_PRODUCT_MYSQL_USERNAME,
            // password: process.env.PROJECT_PRODUCT_MYSQL_PASSWORD,
            // database: process.env.PROJECT_PRODUCT_MYSQL_DATABASE,
            uri: connectionString,
            connectionLimit: parseInt((_a = process.env.PROJECT_PRODUCT_MYSQL_MAX_OPEN_CONNECTION) !== null && _a !== void 0 ? _a : "0"),
            maxIdle: parseInt((_b = process.env.PROJECT_PRODUCT_MYSQL_MAX_IDLE_CONNECTION) !== null && _b !== void 0 ? _b : "0"),
            idleTimeout: parseInt((_c = process.env.PROJECT_PRODUCT_MYSQL_CONNECTION_MAX_IDLETIME) !== null && _c !== void 0 ? _c : "0")
        };
        MysqlUtil.pool = promise_1.default.createPool(access);
        // await MysqlUtil.pool.ping()
        console.log(new Date(), "mysql: connected to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
    }
    // static newConnection(): Pool {
    //     console.log(new Date(), "mysql: connecting to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
    //     const access: PoolOptions = {
    //         host: process.env.PROJECT_PRODUCT_MYSQL_HOST,
    //         user: process.env.PROJECT_PRODUCT_MYSQL_USERNAME,
    //         password: process.env.PROJECT_PRODUCT_MYSQL_PASSWORD,
    //         database: process.env.PROJECT_PRODUCT_MYSQL_DATABASE,
    //         connectionLimit: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_OPEN_CONNECTION ?? "0"),
    //         maxIdle: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_IDLE_CONNECTION ?? "0"),
    //         idleTimeout: parseInt(process.env.PROJECT_PRODUCT_MYSQL_CONNECTION_MAX_IDLETIME ?? "0")
    //     };
    //     const pool = mysql.createPool(access);
    //     console.log(new Date(), "mysql: connected to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
    //     return pool
    // }
    static getInstance() {
        return __awaiter(this, void 0, void 0, function* () {
            if (!MysqlUtil.instance) {
                MysqlUtil.instance = new MysqlUtil();
            }
            else {
                MysqlUtil.instance;
            }
            console.log(new Date(), "mysql: pinging to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
            const pool = MysqlUtil.getPool();
            const connection = yield pool.getConnection();
            yield connection.ping();
            console.log(new Date(), "mysql: pingged to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
            return MysqlUtil.instance;
        });
    }
    static getPool() {
        return this.pool;
    }
}
exports.MysqlUtil = MysqlUtil;
