Background

Prerequisites

cluster

install keda

install target application

verify keda scalers

- kafka
- cpu
- mem




```sh
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 100 MB
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 200 MB
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 300 MB
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 100 MB
➜  keda git:(master) ✗ \  
> 
                                       
➜  keda git:(master) ✗ 
                                       
➜  keda git:(master) ✗ 
                                       
➜  keda git:(master) ✗ 
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 200 MB
➜  keda git:(master) ✗ 
➜  keda git:(master) ✗ 
➜  keda git:(master) ✗ curl -X POST http://localhost:8080/grow
Memory increased - Current: 300 MB
➜  keda git:(master) ✗   
E            TARGETS            MINPODS
   MAXPODS   REPLICAS   AGE            
keda-hpa-memory-scaledobject   Deployme
nt/memtest   memory: 129%/50%   1      
   100       8          42m         
➜  keda git:(master) ✗ kg deploy --watch                      
NAME      READY   UP-TO-DATE   AVAILABLE   AGE
memtest   2/2     2            2           29s
memtest   2/4     2            2           47s
memtest   4/8     4            4           2m49s
memtest   7/8     8            7           2m54s
memtest   8/8     8            8           2m54s
memtest   8/11    8            8           3m4s
memtest   8/11    8            8           3m4s
memtest   8/11    8            8           3m5s
memtest   8/11    11           8           3m5s
memtest   9/11    11           9           3m12s
memtest   13/14   14           13          3m39s
memtest   14/14   14           14          3m39s


```
- mysql

```sh
# scale up
➜  mysql git:(master) ✗ kubectl exec deployment/mysql -- mysql -uroot -proot stats_db \
  -e "INSERT INTO task_instance (task_name, state) VALUES ('test_job_1', 'running'), ('test_job_2', 'queued');"
mysql: [Warning] Using a password on the command line interface can be insecure.
➜  mysql git:(master) ✗ kg hpa --watch
NAME                          REFERENCE           TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
keda-hpa-mysql-scaledobject   Deployment/worker   4/5 (avg)   1         100       1          92s
keda-hpa-mysql-scaledobject   Deployment/worker   5/5 (avg)   1         100       1          106s
keda-hpa-mysql-scaledobject   Deployment/worker   6/5 (avg)   1         100       1          2m31s
➜  mysql git:(master) ✗ kg deploy --watch
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
mysql    1/1     1            1           34m
worker   1/1     1            1           3m37s
worker   1/2     1            1           3m39s
worker   1/2     1            1           3m39s
worker   1/2     1            1           3m39s
worker   1/2     2            1           3m39s
worker   2/2     2            2           3m41s

# scale down
# maybe need to delete the last record for serveral time triggers to scale down and wait for the result to take affect
kubectl exec deployment/mysql -- mysql -uroot -proot stats_db -e "DELETE FROM task_instance ORDER BY id DESC LIMIT 1;"
➜  mysql git:(master) ✗ kg hpa --watch
NAME                          REFERENCE           TARGETS     MINPODS   MAXPODS   REPLICAS   AGE
keda-hpa-mysql-scaledobject   Deployment/worker   6/5 (avg)   1         100       1          2m31s
keda-hpa-mysql-scaledobject   Deployment/worker   3/5 (avg)   1         100       2          2m46s
keda-hpa-mysql-scaledobject   Deployment/worker   2500m/5 (avg)   1         100       2          13m
keda-hpa-mysql-scaledobject   Deployment/worker   2/5 (avg)       1         100       2          14m

➜  mysql git:(master) ✗ kg deploy --watch
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
mysql    1/1     1            1           34m
➜  mysql git:(master) ✗ kg deploy --watch
NAME     READY   UP-TO-DATE   AVAILABLE   AGE
mysql    1/1     1            1           34m
worker   2/1     2            2           19m
worker   2/1     2            2           19m
worker   1/1     1            1           19m

- redis
- elasticsearch
- etcd
- rocketmq
