## Overview of the Test Process

The test verifies that KEDA can properly scale deployments based on Redis Sentinel list length metrics. The process includes:

1. Setting up a Redis Sentinel cluster
2. Creating necessary Kubernetes resources (deployments, secrets, trigger authentication, scaled object)
3. Testing activation threshold (no scaling when below threshold)
4. Testing scale up (scaling to max replicas when above threshold)
5. Testing scale down (scaling back to min replicas when queue is processed)
6. Cleaning up all resources

## Step-by-Step Manual Implementation

### 1. Set up environment variables

```bash
# Define test name and derived resource names
export TEST_NAME="redis-sentinel-lists-test"
export TEST_NAMESPACE="${TEST_NAME}-ns"
export REDIS_NAMESPACE="${TEST_NAME}-redis-ns"
export DEPLOYMENT_NAME="${TEST_NAME}-deployment"
export JOB_NAME="${TEST_NAME}-job"
export SCALED_OBJECT_NAME="${TEST_NAME}-so"
export TRIGGER_AUTH_NAME="${TEST_NAME}-ta"
export SECRET_NAME="${TEST_NAME}-secret"
export REDIS_PASSWORD="admin"
export REDIS_LIST="queue"
export REDIS_HOST="${TEST_NAME}-headless"
export MIN_REPLICA_COUNT=0
export MAX_REPLICA_COUNT=2
export REDIS_PASSWORD_BASE64=$(echo -n "$REDIS_PASSWORD" | base64)
```

### 2. Create namespaces

```bash
# Create test namespace
kubectl create namespace $TEST_NAMESPACE

# Create Redis namespace
kubectl create namespace $REDIS_NAMESPACE
```

### 3. Install Redis Sentinel

```bash
# This is a simplified version of the InstallSentinel function
# Create Redis Sentinel resources
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: redis-password
  namespace: $REDIS_NAMESPACE
type: Opaque
data:
  password: $REDIS_PASSWORD_BASE64
EOF

# Install Redis Sentinel using Helm
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm install $TEST_NAME bitnami/redis \
  --namespace $REDIS_NAMESPACE \
  --set sentinel.enabled=true \
  --set global.redis.password=$REDIS_PASSWORD \
  --set master.persistence.enabled=false \
  --set replica.persistence.enabled=false \
  --set sentinel.quorum=1
```

### 4. Create test resources

```bash
# Create Secret
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Secret
metadata:
  name: $SECRET_NAME
  namespace: $TEST_NAMESPACE
type: Opaque
data:
  password: $REDIS_PASSWORD_BASE64
EOF

# Create TriggerAuthentication
cat <<EOF | kubectl apply -f -
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: $TRIGGER_AUTH_NAME
  namespace: $TEST_NAMESPACE
spec:
  secretTargetRef:
  - parameter: password
    name: $SECRET_NAME
    key: password
  - parameter: sentinelPassword
    name: $SECRET_NAME
    key: password
EOF

# Create Deployment
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: $DEPLOYMENT_NAME
  namespace: $TEST_NAMESPACE
spec:
  replicas: 0
  selector:
    matchLabels:
      app: $DEPLOYMENT_NAME
  template:
    metadata:
      labels:
        app: $DEPLOYMENT_NAME
    spec:
      containers:
      - name: redis-worker
        image: ghcr.io/kedacore/tests-redis-sentinel-lists
        imagePullPolicy: IfNotPresent
        command: ["./main"]
        args: ["read"]
        env:
        - name: REDIS_ADDRESSES
          value: ${REDIS_HOST}.${REDIS_NAMESPACE}:26379
        - name: LIST_NAME
          value: $REDIS_LIST
        - name: REDIS_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_MASTER
          value: mymaster
        - name: READ_PROCESS_TIME
          value: "100"
EOF

# Create ScaledObject
cat <<EOF | kubectl apply -f -
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: $SCALED_OBJECT_NAME
  namespace: $TEST_NAMESPACE
spec:
  scaleTargetRef:
    name: $DEPLOYMENT_NAME
  pollingInterval: 5
  cooldownPeriod: 10
  minReplicaCount: $MIN_REPLICA_COUNT
  maxReplicaCount: $MAX_REPLICA_COUNT
  triggers:
  - type: redis-sentinel
    metadata:
      addressesFromEnv: REDIS_ADDRESSES
      listName: $REDIS_LIST
      sentinelMaster: mymaster
      listLength: "5"
      activationListLength: "10"
    authenticationRef:
      name: $TRIGGER_AUTH_NAME
EOF
```

### 5. Test activation (should not scale)

```bash
echo "--- Testing activation (should not scale) ---"

# Create job to add 5 items (below activation threshold)
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: $JOB_NAME
  namespace: $TEST_NAMESPACE
spec:
  ttlSecondsAfterFinished: 0
  template:
    spec:
      containers:
      - name: redis
        image: ghcr.io/kedacore/tests-redis-sentinel-lists
        imagePullPolicy: IfNotPresent
        command: ["./main"]
        args: ["write"]
        env:
        - name: REDIS_ADDRESSES
          value: ${REDIS_HOST}.${REDIS_NAMESPACE}:26379
        - name: REDIS_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_MASTER
          value: mymaster
        - name: LIST_NAME
          value: $REDIS_LIST
        - name: NO_LIST_ITEMS_TO_WRITE
          value: "5"
      restartPolicy: Never
  backoffLimit: 4
EOF

# Wait for job to complete
kubectl wait --for=condition=complete job/$JOB_NAME -n $TEST_NAMESPACE --timeout=60s

# Check that deployment didn't scale (should remain at 0 replicas)
echo "Checking replica count - should remain at $MIN_REPLICA_COUNT"
for i in {1..6}; do
  REPLICA_COUNT=$(kubectl get deployment $DEPLOYMENT_NAME -n $TEST_NAMESPACE -o jsonpath='{.status.replicas}')
  echo "Current replica count: $REPLICA_COUNT"
  if [ "$REPLICA_COUNT" != "$MIN_REPLICA_COUNT" ]; then
    echo "Error: Replica count changed when it shouldn't have"
    exit 1
  fi
  sleep 10
done

# Clean up job
kubectl delete job $JOB_NAME -n $TEST_NAMESPACE
```

### 6. Test scale up

```bash
echo "--- Testing scale up ---"

# Create job to add 200 items (above activation threshold)
cat <<EOF | kubectl apply -f -
apiVersion: batch/v1
kind: Job
metadata:
  name: $JOB_NAME
  namespace: $TEST_NAMESPACE
spec:
  ttlSecondsAfterFinished: 0
  template:
    spec:
      containers:
      - name: redis
        image: ghcr.io/kedacore/tests-redis-sentinel-lists
        imagePullPolicy: IfNotPresent
        command: ["./main"]
        args: ["write"]
        env:
        - name: REDIS_ADDRESSES
          value: ${REDIS_HOST}.${REDIS_NAMESPACE}:26379
        - name: REDIS_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_PASSWORD
          value: $REDIS_PASSWORD
        - name: REDIS_SENTINEL_MASTER
          value: mymaster
        - name: LIST_NAME
          value: $REDIS_LIST
        - name: NO_LIST_ITEMS_TO_WRITE
          value: "200"
      restartPolicy: Never
  backoffLimit: 4
EOF

# Wait for job to complete
kubectl wait --for=condition=complete job/$JOB_NAME -n $TEST_NAMESPACE --timeout=60s

# Check that deployment scales up to max replicas
echo "Checking replica count - should scale to $MAX_REPLICA_COUNT"
TIMEOUT=180
INTERVAL=10
ELAPSED=0

while [ $ELAPSED -lt $TIMEOUT ]; do
  REPLICA_COUNT=$(kubectl get deployment $DEPLOYMENT_NAME -n $TEST_NAMESPACE -o jsonpath='{.status.replicas}')
  READY_REPLICAS=$(kubectl get deployment $DEPLOYMENT_NAME -n $TEST_NAMESPACE -o jsonpath='{.status.readyReplicas}')
  
  echo "Current replicas: $REPLICA_COUNT, Ready replicas: $READY_REPLICAS"
  
  if [ "$READY_REPLICAS" == "$MAX_REPLICA_COUNT" ]; then
    echo "Successfully scaled to $MAX_REPLICA_COUNT replicas"
    break
  fi
  
  sleep $INTERVAL
  ELAPSED=$((ELAPSED + INTERVAL))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
  echo "Error: Failed to scale to $MAX_REPLICA_COUNT replicas within $TIMEOUT seconds"
  exit 1
fi

# Clean up job
kubectl delete job $JOB_NAME -n $TEST_NAMESPACE
```

### 7. Test scale down

```bash
echo "--- Testing scale down ---"

# Wait for deployment to process messages and scale down
echo "Waiting for deployment to process messages and scale down to $MIN_REPLICA_COUNT"
TIMEOUT=180
INTERVAL=10
ELAPSED=0

while [ $ELAPSED -lt $TIMEOUT ]; do
  REPLICA_COUNT=$(kubectl get deployment $DEPLOYMENT_NAME -n $TEST_NAMESPACE -o jsonpath='{.status.replicas}')
  
  echo "Current replicas: $REPLICA_COUNT"
  
  if [ "$REPLICA_COUNT" == "$MIN_REPLICA_COUNT" ]; then
    echo "Successfully scaled down to $MIN_REPLICA_COUNT replicas"
    break
  fi
  
  sleep $INTERVAL
  ELAPSED=$((ELAPSED + INTERVAL))
done

if [ $ELAPSED -ge $TIMEOUT ]; then
  echo "Error: Failed to scale down to $MIN_REPLICA_COUNT replicas within $TIMEOUT seconds"
  exit 1
fi
```

### 8. Clean up resources

```bash
echo "--- Cleaning up resources ---"

# Delete ScaledObject
kubectl delete scaledobject $SCALED_OBJECT_NAME -n $TEST_NAMESPACE

# Delete Deployment
kubectl delete deployment $DEPLOYMENT_NAME -n $TEST_NAMESPACE

# Delete TriggerAuthentication
kubectl delete triggerauthentication $TRIGGER_AUTH_NAME -n $TEST_NAMESPACE

# Delete Secret
kubectl delete secret $SECRET_NAME -n $TEST_NAMESPACE

# Uninstall Redis Sentinel
helm uninstall $TEST_NAME -n $REDIS_NAMESPACE

# Delete namespaces
kubectl delete namespace $TEST_NAMESPACE
kubectl delete namespace $REDIS_NAMESPACE

echo "Test completed successfully!"
```

## Complete Script

You can combine all the above steps into a single shell script for easier execution. Save it as `run-redis-sentinel-lists-test.sh`, make it executable with `chmod +x run-redis-sentinel-lists-test.sh`, and run it with `./run-redis-sentinel-lists-test.sh`.

This manual implementation follows the same process as the automated test in the KEDA codebase, allowing you to verify the Redis Sentinel Lists scaler functionality from the command line.

        Too many current requests. Your queue position is 1. Please wait for a while or switch to other models for a smoother experience.
