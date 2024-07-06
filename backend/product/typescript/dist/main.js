"use strict";
// import { newConnection, pool } from "./utils/mysql-utils";
// console.log("maints");
Object.defineProperty(exports, "__esModule", { value: true });
// import { newConnection, pool } from "./utils/mysql-utils";
// import mysqlUtils, { MysqlUtil } from "./utils/mysql-utils";
const elasticsearch_util_1 = require("./utils/elasticsearch-util");
const mysql_utils_1 = require("./utils/mysql-utils");
// console.log("mysqlUtil.pool:", mysqlUtils.pool);
// mysqlUtils.pool = mysqlUtils.newConnection()
// console.log("mysqlUtil.pool:", mysqlUtils.pool);
// console.log("MysqlUtil.pool:", MysqlUtil.getPool());
mysql_utils_1.MysqlUtil.getInstance();
elasticsearch_util_1.ElasticsearchUtil.getInstance();
// console.log("MysqlUtil.pool:", MysqlUtil.getPool());
// async function main() {
//     console.log("main");
//     newConnection()
// }
// main()
