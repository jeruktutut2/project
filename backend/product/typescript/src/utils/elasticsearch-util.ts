// import {  } from "";
import { Client } from "@elastic/elasticsearch";

export class ElasticsearchUtil {

    private static instance: ElasticsearchUtil
    private static client: Client

    private constructor() {
        console.log(new Date(), "elasticsearch: connecting to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        ElasticsearchUtil.client = new Client({
            node: process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE
        })
        // ElasticsearchUtil.client.ping()
        console.log(new Date(), "elasticsearch: connected to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        // console.log(new Date(), "elasticsearch: pinging to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        // ElasticsearchUtil.client.ping()
        // console.log(new Date(), "elasticsearch: pingged to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
    }

    public static async getInstance(): Promise<ElasticsearchUtil> {
        if (!ElasticsearchUtil.instance) {
            ElasticsearchUtil.instance = new ElasticsearchUtil()
        } else {
            ElasticsearchUtil.instance
        }
        console.log(new Date(), "elasticsearch: pinging to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        await ElasticsearchUtil.client.ping()
        // console.log("ping:", await ElasticsearchUtil.client.ping());
        console.log(new Date(), "elasticsearch: pingged to", process.env.PROJECT_PRODUCT_ELASTICSEARCH_NODE);
        return ElasticsearchUtil.instance
    }
}