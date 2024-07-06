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
exports.ElasticsearchUtil = void 0;
// import {  } from "";
const elasticsearch_1 = require("@elastic/elasticsearch");
class ElasticsearchUtil {
    constructor() {
        console.log(new Date(), "elasticsearch: connecting to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        ElasticsearchUtil.client = new elasticsearch_1.Client({
            node: process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE
        });
        // ElasticsearchUtil.client.ping()
        console.log(new Date(), "elasticsearch: connected to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        // console.log(new Date(), "elasticsearch: pinging to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        // ElasticsearchUtil.client.ping()
        // console.log(new Date(), "elasticsearch: pingged to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
    }
    static getInstance() {
        return __awaiter(this, void 0, void 0, function* () {
            if (!ElasticsearchUtil.instance) {
                ElasticsearchUtil.instance = new ElasticsearchUtil();
            }
            else {
                ElasticsearchUtil.instance;
            }
            console.log(new Date(), "elasticsearch: pinging to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
            yield ElasticsearchUtil.client.ping();
            // console.log("ping:", await ElasticsearchUtil.client.ping());
            console.log(new Date(), "elasticsearch: pingged to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
            return ElasticsearchUtil.instance;
        });
    }
}
exports.ElasticsearchUtil = ElasticsearchUtil;
