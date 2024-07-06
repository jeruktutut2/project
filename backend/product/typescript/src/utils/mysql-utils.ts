// import mysql, { PoolOptions } from 'mysql2';
import mysql, { Pool, PoolOptions } from "mysql2/promise";
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

export class MysqlUtil {

    private static instance: MysqlUtil
    private static pool: Pool

    private constructor() {
        console.log(new Date(), "mysql: connecting to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
        // const hostport = process.env.PROJECT_PRODUCT_MYSQL_HOST.split(":")
        // const host = hostport[0]
        const connectionString = `mysql://`+process.env.PROJECT_PRODUCT_MYSQL_USERNAME+`:`+process.env.PROJECT_PRODUCT_MYSQL_PASSWORD+`@`+process.env.PROJECT_PRODUCT_MYSQL_HOST+`/`+process.env.PROJECT_PRODUCT_MYSQL_DATABASE+``
        const access: PoolOptions = {
            // host: process.env.PROJECT_PRODUCT_MYSQL_HOST,
            // port: Number(process.env.PROJECT_PRODUCT_MYSQL_PORT),
            // user: process.env.PROJECT_PRODUCT_MYSQL_USERNAME,
            // password: process.env.PROJECT_PRODUCT_MYSQL_PASSWORD,
            // database: process.env.PROJECT_PRODUCT_MYSQL_DATABASE,
            uri: connectionString,
            connectionLimit: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_OPEN_CONNECTION ?? "0"),
            maxIdle: parseInt(process.env.PROJECT_PRODUCT_MYSQL_MAX_IDLE_CONNECTION ?? "0"),
            idleTimeout: parseInt(process.env.PROJECT_PRODUCT_MYSQL_CONNECTION_MAX_IDLETIME ?? "0")
        };
        MysqlUtil.pool = mysql.createPool(access);
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

    public static async getInstance(): Promise<MysqlUtil> {
        if (!MysqlUtil.instance) {
            MysqlUtil.instance = new MysqlUtil()
        } else {
            MysqlUtil.instance
        }
        console.log(new Date(), "mysql: pinging to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
        const pool = MysqlUtil.getPool()
        const connection = await pool.getConnection()
        await connection.ping()
        console.log(new Date(), "mysql: pingged to", process.env.PROJECT_PRODUCT_MYSQL_HOST);
        return MysqlUtil.instance
    }

    public static getPool(): Pool {
        return this.pool
    }
}