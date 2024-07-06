import { Product } from "../../../../src/features/create-product/models/create-product";
import { CreateProductRepository } from "../../../../src/features/create-product/repositories/create-product-repository";
import { MysqlUtil } from "../../../../src/utils/mysql-utils";
import { PoolConnection } from 'mysql2/promise';
describe("create product", () => {
    beforeAll(() =>{
        MysqlUtil.getInstance()
    })

    it("create product", async () => {

        // console.log("MysqlUtil.getPool():", await MysqlUtil.getPool().getConnection());
        
        let poolConnection: PoolConnection | null = null

        try {
            poolConnection = await MysqlUtil.getPool().getConnection()
            // console.log("poolConnection1:", poolConnection);
            
            await poolConnection.beginTransaction()
            // console.log("poolConnection2:", poolConnection);

            // console.log("poolConnection:", poolConnection);
        

            const product: Product = {
                userId: 1,
                name: "name",
                description: "description",
                stock: 1
            }
            CreateProductRepository.create(poolConnection, product)
            await poolConnection.commit()
        } catch(e) {
            console.log("e:", e);
            if (poolConnection) {
                await poolConnection.rollback()
            }
            // await poolConnection.rollback()
        } finally {
            if (poolConnection) {
                poolConnection.release()
            }
            // poolConnection.release()
        }
        
        // poolConnection.release()
    })
})