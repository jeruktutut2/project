import { errorHandler } from "../../../exception/error-exception";
import { Validation } from "../../../validation/validation";
import { CreateProductRequest } from "../models/create-product-request";
import { CreateProductResponse } from "../models/create-product-response";
import { CreateProductValidationSchema } from "../validation-schema/create-product-validation-schema";
import { PoolConnection } from 'mysql2/promise';
import { MysqlUtil } from "../../../utils/mysql-utils";
import { Product } from "../models/create-product";
import { CreateProductRepository } from "../repositories/create-product-repository";
import { ResponseException } from "../../../exception/response-exception";

export class CreateProductService {
    static async create(requestId: string, createProductRequest: CreateProductRequest): Promise<CreateProductResponse> {

        let poolConnection: PoolConnection | null = null
        try {
            createProductRequest = Validation.validate(CreateProductValidationSchema.CREATE, createProductRequest)

            poolConnection = await MysqlUtil.getPool().getConnection()
            await poolConnection.beginTransaction()

            const product: Product = {
                userId: createProductRequest.userId,
                name: createProductRequest.name,
                description: createProductRequest.description,
                stock: createProductRequest.stock
            }
            const [resultSetHeader] = await CreateProductRepository.create(poolConnection, product)
            // console.log("mantap:", resultSetHeader);
            if (resultSetHeader.affectedRows !== 1) {
                throw new Error("number of affected rows when creating product is not one:" + resultSetHeader.affectedRows)
            }
            
            await poolConnection.commit()
        } catch(e: unknown) {
            errorHandler(e, requestId)
            if (poolConnection) {
                await poolConnection.rollback()
            }
        } finally {
            if (poolConnection) {
                poolConnection.release()
            }
        }

        const response = {
            name: "name",
            description: "description",
            stock: 1
        }
        return response
    }
}