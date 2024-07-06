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
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateProductService = void 0;
const error_exception_1 = require("../../../exception/error-exception");
const validation_1 = require("../../../validation/validation");
const create_product_validation_schema_1 = require("../validation-schema/create-product-validation-schema");
const mysql_utils_1 = require("../../../utils/mysql-utils");
const create_product_repository_1 = require("../repositories/create-product-repository");
class CreateProductService {
    static create(requestId, createProductRequest) {
        return __awaiter(this, void 0, void 0, function* () {
            let poolConnection = null;
            try {
                createProductRequest = validation_1.Validation.validate(create_product_validation_schema_1.CreateProductValidationSchema.CREATE, createProductRequest);
                poolConnection = yield mysql_utils_1.MysqlUtil.getPool().getConnection();
                yield poolConnection.beginTransaction();
                const product = {
                    userId: createProductRequest.userId,
                    name: createProductRequest.name,
                    description: createProductRequest.description,
                    stock: createProductRequest.stock
                };
                const [resultSetHeader] = yield create_product_repository_1.CreateProductRepository.create(poolConnection, product);
                // console.log("mantap:", resultSetHeader);
                if (resultSetHeader.affectedRows !== 1) {
                    throw new Error("number of affected rows when creating product is not one:" + resultSetHeader.affectedRows);
                }
                yield poolConnection.commit();
            }
            catch (e) {
                (0, error_exception_1.errorHandler)(e, requestId);
                if (poolConnection) {
                    yield poolConnection.rollback();
                }
            }
            finally {
                if (poolConnection) {
                    poolConnection.release();
                }
            }
            const response = {
                name: "name",
                description: "description",
                stock: 1
            };
            return response;
        });
    }
}
exports.CreateProductService = CreateProductService;
