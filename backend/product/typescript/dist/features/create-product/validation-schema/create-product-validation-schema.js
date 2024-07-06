"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CreateProductValidationSchema = void 0;
const zod_1 = require("zod");
class CreateProductValidationSchema {
}
exports.CreateProductValidationSchema = CreateProductValidationSchema;
CreateProductValidationSchema.CREATE = zod_1.z.object({
    userId: zod_1.z.number().positive(),
    name: zod_1.z.string({ required_error: 'name is required' }).min(1).max(255),
    description: zod_1.z.string({ required_error: 'description is required' }).min(1).max(1000),
    stock: zod_1.z.number().min(0)
});
