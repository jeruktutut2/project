import mysql from 'mysql2/promise';

let mysqlPool
const newConnection = () => {
    console.log(new Date(), "mysql: connecting to", process.env.PROJECT_USER_MYSQL_HOST);
    const pool = mysql.createPool("mysql://"+process.env.PROJECT_USER_MYSQL_USERNAME+":"+process.env.PROJECT_USER_MYSQL_PASSWORD+"@"+process.env.PROJECT_USER_MYSQL_HOST+"/"+process.env.PROJECT_USER_MYSQL_DATABASE);
    console.log(new Date(), "mysql: connected to", process.env.PROJECT_USER_MYSQL_HOST);
    return pool
}

const getConnection = async(pool) => {
    return await pool.getConnection()
}

const releaseConnection = (pool, connection) => {
    return pool.releaseConnection(connection)
}

const closeConnection = (mysqlPool) => {
    console.log(new Date(), "mysql: closing connection from", process.env.PROJECT_USER_MYSQL_HOST);
    mysqlPool.end()
    console.log(new Date(), "mysql: closed connection from", process.env.PROJECT_USER_MYSQL_HOST);
}

export default {
    newConnection,
    mysqlPool,
    getConnection,
    releaseConnection,
    closeConnection
}