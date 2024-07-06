import { FieldPacket, Pool, PoolConnection, ResultSetHeader } from "mysql2/promise";
import { Product } from "../models/create-product";

export class CreateProductRepository {
    static async create(pool: PoolConnection, product: Product): Promise<[ResultSetHeader, FieldPacket[]]> {
        // const connection = await pool.getConnection()
        // const tx = connection.beginTransaction()
        const query = `INSERT INTO products (user_id, name, description, stock) VALUES (?, ?, ?, ?);`
        const result = await pool.query<ResultSetHeader>(query, [product.userId, product.name, product.description, product.stock])
        // console.log("result:", result);
        return result
    }
}