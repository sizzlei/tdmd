# TDMD
Table Definition Markdown(이하 TDMD)는 테이블 명세서에 대해 Git에서 관리할 목적으로 작성된 툴입니다.

# Usage
## Run tool
```
$ ./tdmd
INFO[0000] TDMD Version : 0.1.0                         
Endpoint : 
Port(3306): 
User : 
Pass : 
File Path (./): 
Database name :
Table name : 
```
실행하게 되면 위에서 요구하는 정보를 입력합니다. 

## Example TDMD File
테이블 명세서를 추출 하게 되면 아래와 같이 추출됩니다.

### CODE
```
apppush 
=============
**Last Update** : 2022-05-04
## Table List
 - [test1 (테이블 명세서 예제)](#push_bulk_request_users)
 ## test1
**Information**
|Table type|Engine|Row format|Collation|Comment|
|---|---|---|---|---|
|BASE TABLE|InnoDB|Dynamic|latin1_swedish_ci|테이블 명세서 예제|

**Columns**
|Name|Type|Nullable|Default|Charset|Collation|Key|Extra|Comment|
|---|---|---|---|---|---|---|---|---|
|ID|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼1|
|t1|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼2|
|t3|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼3|
|t4|datetime|NO||||||컬럼4|
|t5|int(11)|NO||||||컬럼5|
|t5|varchar(10)|NO||latin1|latin1_swedish_ci|||컬럼6|

**Index**
- [Normal] ix_t1 (t1)
- [Unique] uix_t2 (t2)

```
### View

apppush 
=============
**Last Update** : 2022-05-04
## Table List
 - [test1 (테이블 명세서 예제)](#push_bulk_request_users)
 ## test1
**Information**
|Table type|Engine|Row format|Collation|Comment|
|---|---|---|---|---|
|BASE TABLE|InnoDB|Dynamic|latin1_swedish_ci|테이블 명세서 예제|

**Columns**
|Name|Type|Nullable|Default|Charset|Collation|Key|Extra|Comment|
|---|---|---|---|---|---|---|---|---|
|ID|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼1|
|t1|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼2|
|t3|varchar(255)|NO||latin1|latin1_swedish_ci|||컬럼3|
|t4|datetime|NO||||||컬럼4|
|t5|int(11)|NO||||||컬럼5|
|t5|varchar(10)|NO||latin1|latin1_swedish_ci|||컬럼6|

**Index**
- [Normal] ix_t1 (t1)
- [Unique] uix_t2 (t2)
