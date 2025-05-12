
```sh
analyze the code elasticsearch_test.go , write a markdown to guide the test process step by step and be able to operate manually

- if the step is related to create k8s resource, use kubectl

- if the step is related to operate elasticsearch(create template, search doc, add element, query and etc), use curl with proper explaination
```

          
# Elasticsearch Scaler Test Guide

This guide provides a step-by-step manual process for testing the Elasticsearch scaler in KEDA. The test verifies that KEDA can properly scale deployments based on Elasticsearch query results.

## Prerequisites

- Kubernetes cluster with kubectl configured
- KEDA installed on the cluster

## Test Overview

The test performs the following:
1. Sets up an Elasticsearch instance in Kubernetes
2. Creates necessary Kubernetes resources (secrets, deployments, etc.)
3. Tests scaling with two different trigger types:
   - Using a search template
   - Using a direct query
4. Verifies activation, scale out, and scale in behaviors

## Step-by-Step Guide

### 1. Create Test Namespace

```bash
kubectl create namespace elasticsearch-test-ns
```

### 2. Create Secret for Elasticsearch Authentication

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: elasticsearch-test-secret
  namespace: elasticsearch-test-ns
data:
  password: $(echo -n "passw0rd" | base64)
EOF
```

### 3. Create TriggerAuthentication

```bash
kubectl apply -f - <<EOF
apiVersion: keda.sh/v1alpha1
kind: TriggerAuthentication
metadata:
  name: keda-trigger-auth-elasticsearch-secret
  namespace: elasticsearch-test-ns
spec:
  secretTargetRef:
  - parameter: password
    name: elasticsearch-test-secret
    key: password
EOF
```

### 4. Deploy Elasticsearch StatefulSet

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: elasticsearch
  namespace: elasticsearch-test-ns
spec:
  replicas: 1
  selector:
    matchLabels:
      name: elasticsearch
  template:
    metadata:
      labels:
        name: elasticsearch
    spec:
      containers:
      - name: elasticsearch
        image: docker.elastic.co/elasticsearch/elasticsearch:7.15.1
        imagePullPolicy: IfNotPresent
        env:
          - name: POD_IP
            valueFrom:
              fieldRef:
                apiVersion: v1
                fieldPath: status.podIP
          - name: POD_NAME
            valueFrom:
              fieldRef:
                apiVersion: v1
                fieldPath: metadata.name
          - name: NODE_NAME
            valueFrom:
              fieldRef:
                apiVersion: v1
                fieldPath: spec.nodeName
          - name: NAMESPACE
            valueFrom:
              fieldRef:
                apiVersion: v1
                fieldPath: metadata.namespace
          - name: ES_JAVA_OPTS
            value: -Xms256m -Xmx256m
          - name: cluster.name
            value: elasticsearch-keda
          - name: cluster.initial_master_nodes
            valueFrom:
              fieldRef:
                apiVersion: v1
                fieldPath: status.podIP
          - name: node.data
            value: "true"
          - name: node.ml
            value: "false"
          - name: node.ingest
            value: "false"
          - name: node.master
            value: "true"
          - name: node.remote_cluster_client
            value: "false"
          - name: node.transform
            value: "false"
          - name: ELASTIC_PASSWORD
            value: "passw0rd"
          - name: xpack.security.enabled
            value: "true"
          - name: node.store.allow_mmap
            value: "false"
        ports:
        - containerPort: 9200
          name: http
          protocol: TCP
        - containerPort: 9300
          name: transport
          protocol: TCP
        readinessProbe:
          exec:
            command:
              - /usr/bin/curl
              - -sS
              - -u "elastic:passw0rd"
              - http://localhost:9200
          failureThreshold: 3
          initialDelaySeconds: 10
          periodSeconds: 5
          successThreshold: 1
          timeoutSeconds: 5
  serviceName: elasticsearch-test-deployment
EOF
```

### 5. Create Service for Elasticsearch

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: elasticsearch-test-deployment
  namespace: elasticsearch-test-ns
spec:
  type: ClusterIP
  ports:
  - name: http
    port: 9200
    targetPort: 9200
    protocol: TCP
  selector:
    name: elasticsearch
EOF
```

### 6. Create Target Deployment to Scale

```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: elasticsearch-test-deployment
  namespace: elasticsearch-test-ns
  labels:
    app: elasticsearch-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: elasticsearch-test-deployment
  template:
    metadata:
      labels:
        app: elasticsearch-test-deployment
    spec:
      containers:
      - name: nginx
        image: ghcr.io/nginx/nginx-unprivileged:1.26
        ports:
        - containerPort: 80
EOF
```

### 7. Wait for Elasticsearch to be Ready

```bash
kubectl wait --for=condition=ready pod elasticsearch-0 -n elasticsearch-test-ns --timeout=300s
```

### 8. Setup Elasticsearch (Create Index and Search Template)

First, create the index:

```bash
kubectl exec -n elasticsearch-test-ns elasticsearch-0 -- curl -sS -H 'Content-Type: application/json' -u 'elastic:passw0rd' -XPUT http://localhost:9200/keda -d '
{
  "mappings": {
    "properties": {
      "@timestamp": {
        "type": "date"
      },
      "dummy": {
        "type": "integer"
      },
      "dumb": {
        "type": "keyword"
      }
    }
  },
  "settings": {
    "number_of_replicas": 0,
    "number_of_shards": 1
  }
}'
```

Then, create the search template:

```bash
kubectl exec -n elasticsearch-test-ns elasticsearch-0 -- curl -sS -H 'Content-Type: application/json' -u 'elastic:passw0rd' -XPUT http://localhost:9200/_scripts/keda-search-template -d '
{
  "script": {
    "lang": "mustache",
    "source": {
      "query": {
        "bool": {
          "filter": [
            {
              "range": {
                "@timestamp": {
                  "gte": "now-1m",
                  "lte": "now"
                }
              }
            },
            {
              "term": {
                "dummy": "{{dummy_value}}"
              }
            },
            {
              "term": {
                "dumb": "{{dumb_value}}"
              }
            }
          ]
        }
      }
    }
  }
}'
```

## Test 1: Using Search Template

### 1. Create ScaledObject with Search Template

```bash
kubectl apply -f - <<EOF
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: elasticsearch-test-so
  namespace: elasticsearch-test-ns
  labels:
    app: elasticsearch-test-deployment
spec:
  scaleTargetRef:
    name: elasticsearch-test-deployment
  minReplicaCount: 0
  maxReplicaCount: 2
  pollingInterval: 3
  cooldownPeriod: 5
  triggers:
    - type: elasticsearch
      metadata:
        addresses: "http://elasticsearch-test-deployment.elasticsearch-test-ns.svc:9200"
        username: "elastic"
        index: keda
        searchTemplateName: keda-search-template
        valueLocation: "hits.total.value"
        targetValue: "1"
        activationTargetValue: "4"
        parameters: "dummy_value:1;dumb_value:oOooo"
      authenticationRef:
        name: keda-trigger-auth-elasticsearch-secret
EOF
```

### 2. Test Activation (Below Threshold)

Add 3 documents to Elasticsearch (below activation threshold of 4):

```bash
for i in {1..3}; do
  # Generate document with current timestamp
  TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  DOC="{\"@timestamp\": \"$TIMESTAMP\", \"dummy\": 1, \"dumb\": \"oOooo\"}"
  
  # Add document to Elasticsearch
  kubectl exec -n elasticsearch-test-ns elasticsearch-0 -- curl -sS -H 'Content-Type: application/json' -u 'elastic:passw0rd' -XPOST http://localhost:9200/keda/_doc -d "$DOC"
done
```

Verify that the deployment doesn't scale (should remain at 0 replicas):

```bash
kubectl get deployment elasticsearch-test-deployment -n elasticsearch-test-ns
```

### 3. Test Scale Out

Add 10 more documents to Elasticsearch (exceeding activation threshold):

```bash
for i in {1..10}; do
  # Generate document with current timestamp
  TIMESTAMP=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
  DOC="{\"@timestamp\": \"$TIMESTAMP\", \"dummy\": 1, \"dumb\": \"oOooo\"}"
  
  # Add document to Elasticsearch
  kubectl exec -n elasticsearch-test-ns elasticsearch-0 -- curl -sS -H 'Content-Type: application/json' -u 'elastic:passw0rd' -XPOST http://localhost:9200/keda/_doc -d "$DOC"
done
```

Verify that the deployment scales out to max replicas (should become 2 replicas):

```bash
kubectl get deployment elasticsearch-test-deployment -n elasticsearch-test-ns -w
```

### 4. Test Scale In

Wait for the cooldown period (documents will age out of the 1-minute window in the query):

```bash
# Wait for documents to age out (they have a 1-minute window in the query)
sleep 120
```

Verify that the deployment scales in to min replicas:

```bash
kubectl get deployment elasticsearch-test-deployment -n elasticsearch-test-ns -w
```

### 5. Delete the ScaledObject

```bash
kubectl delete scaledobject elasticsearch-test-so -n elasticsearch-test-ns
```

## Test 2: Using Direct Query

### 1. Create ScaledObject with Direct Query

```bash
kubectl apply -f - <<EOF
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: elasticsearch-test-so
  namespace: elasticsearch-test-ns
  labels:
    app: elasticsearch-test-deployment
spec:
  scaleTargetRef:
    name: elasticsearch-test-deployment
  minReplicaCount: 0
  maxReplicaCount: 2
  pollingInterval: 3
  cooldownPeriod: 5
  triggers:
    - type: elasticsearch
      metadata:
        addresses: "http://elasticsearch-test-deployment.elasticsearch-test-ns.svc:9200"
        username: "elastic"
        index: keda
        query: |
          {
            "query": {
              "bool": {
                "must": [
                  {
                    "range": {
                      "@timestamp": {
                        "gte": "now-1m",
                        "lte": "now"
                      }
                    }
                  },
                  {
                    "match_all": {}
                  }
                ]
              }
            }
          }
        valueLocation: "hits.total.value"
        targetValue: "1"
        activationTargetValue: "4"
      authenticationRef:
        name: keda-trigger-auth-elasticsearch-secret
EOF
```

### 2. Repeat the Same Tests as Above

Follow the same steps as in Test 1 to verify activation, scale out, and scale in behaviors.

## Cleanup

Remove all resources created for the test:

```bash
kubectl delete namespace elasticsearch-test-ns
```

## Explanation of Key Components

### Elasticsearch Index Structure
- The test creates an index named `keda` with fields:
  - `@timestamp`: Date field for time-based filtering
  - `dummy`: Integer field (value 1 in test documents)
  - `dumb`: Keyword field (value "oOooo" in test documents)

### Search Template
- The search template filters documents:
  - Within the last minute (`now-1m` to `now`)
  - With specific values for `dummy` and `dumb` fields
  - Parameters are passed via the ScaledObject configuration

### Direct Query
- The direct query approach embeds the query JSON directly in the ScaledObject
- It filters documents within the last minute but uses `match_all` instead of specific field matching

### Scaling Parameters
- `activationTargetValue: "4"`: Scaling activates when 4+ documents match
- `targetValue: "1"`: Each replica handles 1 document (so 5 documents = 5 replicas, capped at maxReplicaCount)
- `valueLocation: "hits.total.value"`: Where to find the count in Elasticsearch response

This test demonstrates KEDA's ability to scale workloads based on Elasticsearch query results, using either search templates or direct queries.