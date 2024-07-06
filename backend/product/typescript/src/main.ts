// import { newConnection, pool } from "./utils/mysql-utils";
// console.log("maints");

// import { newConnection, pool } from "./utils/mysql-utils";
// import mysqlUtils, { MysqlUtil } from "./utils/mysql-utils";
import { ElasticsearchUtil } from "./utils/elasticsearch-util";
import { MysqlUtil } from "./utils/mysql-utils";

// console.log("mysqlUtil.pool:", mysqlUtils.pool);

// mysqlUtils.pool = mysqlUtils.newConnection()
// console.log("mysqlUtil.pool:", mysqlUtils.pool);

// console.log("MysqlUtil.pool:", MysqlUtil.getPool());

MysqlUtil.getInstance()
ElasticsearchUtil.getInstance()

// console.log("MysqlUtil.pool:", MysqlUtil.getPool());

// async function main() {
//     console.log("main");
//     newConnection()
// }

// main()