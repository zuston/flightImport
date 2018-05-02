分为进入mysql中和hbase中

mysql表名 [型号+架机+专业]

架次  日期  传感器名


hbase表名 [型号+架机+专业+架次]
rowkey : 日期
column : k=传感器名
         v=值


