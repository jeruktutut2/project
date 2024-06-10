const countByEmail = async (connection, username, email) => {
    const query = "SELECT COUNT(*) AS number_of_user FROM user WHERE username = ? OR email = ?;"
    const params = [username, email]
    return await connection.execute(query, params);
}

const create = async (connection, user) => {
    const query = `INSERT INTO user (username, email, password, utc, created_at)
                VALUES(?, ?, ?, ?, ?)`
    const params = [user.username, user.email, user.password, user.utc, user.createdAt]
    return await connection.execute(query, params)
}

const findByEmail = async (connection, email) => {
    const query = `SELECT id, username, email, password, utc, created_at FROM user WHERE email = ?;`
    const params = [email]
    return await connection.execute(query, params)
}

export default {
    countByEmail,
    create,
    findByEmail
}