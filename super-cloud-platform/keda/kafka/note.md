# Manual Implementation of Kafka Scaler Test

This guide provides step-by-step instructions to manually implement the Kafka scaler test based on the `TestScaler` function in `/Users/shen/work/dev/golang/open_source/keda/tests/scalers/kafka/kafka_test.go`. Each step is executable from the command line.

## Prerequisites

- Kubernetes cluster with kubectl configured
- Helm installed
- Access to Docker images

## Step 1: Create Test Namespace

```bash
kubectl create namespace kafka-test-ns
```

## Step 2: Install Strimzi Kafka Operator

```bash
helm repo add strimzi https://strimzi.io/charts/
helm repo update
helm upgrade --install --namespace kafka-test-ns --wait kafka-test strimzi/strimzi-kafka-operator --version 0.30.0
```

## Step 3: Create Kafka Cluster

Create a file named `kafka-cluster.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: kafka-test-kafka
  namespace: kafka-test-ns
spec:
  kafka:
    version: "3.1.0"
    replicas: 1
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
    config:
      offsets.topic.replication.factor: 1
      transaction.state.log.replication.factor: 1
      transaction.state.log.min.isr: 1
      log.message.format.version: "2.5"
    storage:
      type: ephemeral
  zookeeper:
    replicas: 1
    storage:
      type: ephemeral
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

Apply the configuration:

```bash
kubectl apply -f kafka-cluster.yaml
kubectl wait kafka/kafka-test-kafka --for=condition=Ready --timeout=300s --namespace kafka-test-ns
```

## Step 4: Create Kafka Topics

Create a file named `kafka-topic1.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: kafka-topic
  namespace: kafka-test-ns
  labels:
    strimzi.io/cluster: kafka-test-kafka
spec:
  partitions: 3
  replicas: 1
  config:
    retention.ms: 604800000
    segment.bytes: 1073741824
```

Create a file named `kafka-topic2.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: kafka-topic2
  namespace: kafka-test-ns
  labels:
    strimzi.io/cluster: kafka-test-kafka
spec:
  partitions: 3
  replicas: 1
  config:
    retention.ms: 604800000
    segment.bytes: 1073741824
```

Create a file named `kafka-topic-zero-invalid-offset.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: kafka-topic-zero-invalid-offset
  namespace: kafka-test-ns
  labels:
    strimzi.io/cluster: kafka-test-kafka
spec:
  partitions: 1
  replicas: 1
  config:
    retention.ms: 604800000
    segment.bytes: 1073741824
```

Create a file named `kafka-topic-one-invalid-offset.yaml`:

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: KafkaTopic
metadata:
  name: kafka-topic-one-invalid-offset
  namespace: kafka-test-ns
  labels:
    strimzi.io/cluster: kafka-test-kafka
spec:
  partitions: 1
  replicas: 1
  config:
    retention.ms: 604800000
    segment.bytes: 1073741824
```

Apply the configurations:

```bash
kubectl apply -f kafka-topic1.yaml
kubectl apply -f kafka-topic2.yaml
kubectl apply -f kafka-topic-zero-invalid-offset.yaml
kubectl apply -f kafka-topic-one-invalid-offset.yaml

kubectl wait kafkatopic/kafka-topic --for=condition=Ready --timeout=300s --namespace kafka-test-ns
kubectl wait kafkatopic/kafka-topic2 --for=condition=Ready --timeout=300s --namespace kafka-test-ns
kubectl wait kafkatopic/kafka-topic-zero-invalid-offset --for=condition=Ready --timeout=300s --namespace kafka-test-ns
kubectl wait kafkatopic/kafka-topic-one-invalid-offset --for=condition=Ready --timeout=300s --namespace kafka-test-ns
```

## Step 5: Create Kafka Client Pod

Create a file named `kafka-client.yaml`:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: kafka-test-client
  namespace: kafka-test-ns
spec:
  containers:
  - name: kafka-test-client
    image: confluentinc/cp-kafka:5.2.1
    command:
      - sh
      - -c
      - "exec tail -f /dev/null"
```

Apply the configuration:

```bash
kubectl apply -f kafka-client.yaml
kubectl wait pod/kafka-test-client --for=condition=Ready --timeout=300s --namespace kafka-test-ns
```

## Step 6: Test Earliest Policy

Create a file named `earliest-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-test-deployment
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: kafka-consumer
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      containers:
      - name: kafka-consumer
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic --group earliest --from-beginning --consumer-property enable.auto.commit=false"
```

Create a file named `earliest-scaled-object.yaml`:

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-test-so
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  scaleTargetRef:
    name: kafka-test-deployment
  triggers:
  - type: kafka
    metadata:
      topic: kafka-topic
      bootstrapServers: kafka-test-kafka-bootstrap.kafka-test-ns:9092
      consumerGroup: earliest
      lagThreshold: '1'
      activationLagThreshold: '1'
      offsetResetPolicy: earliest
```

Apply the configurations:

```bash
kubectl apply -f earliest-deployment.yaml
kubectl apply -f earliest-scaled-object.yaml
```

Check that the deployment doesn't scale up initially:

```bash
kubectl get deployment kafka-test-deployment -n kafka-test-ns
```

Publish a message to the topic:

```bash
kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'echo "{\"text\": \"foo\"}" | kafka-console-producer --broker-list kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic'
```

Check that the deployment still doesn't scale up (due to activation threshold):

```bash
kubectl get deployment kafka-test-deployment -n kafka-test-ns
```

Publish another message:

```bash
kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'echo "{\"text\": \"foo\"}" | kafka-console-producer --broker-list kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic'
```

Check that the deployment scales up to 2 replicas:

```bash
kubectl get deployment kafka-test-deployment -n kafka-test-ns
```

Publish 5 more messages:

```bash
for i in {1..5}; do
  kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'echo "{\"text\": \"foo\"}" | kafka-console-producer --broker-list kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic'
done
```

Check that the deployment scales up to 3 replicas (limited by partition count):

```bash
kubectl get deployment kafka-test-deployment -n kafka-test-ns
```

Clean up:

```bash
kubectl delete -f earliest-scaled-object.yaml
kubectl delete -f earliest-deployment.yaml
```

## Step 7: Test Latest Policy

First, commit the partition to set up the latest policy test:

```bash
kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic --group latest --from-beginning --consumer-property enable.auto.commit=true --timeout-ms 15000'
```

Create a file named `latest-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-test-deployment
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: kafka-consumer
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      containers:
      - name: kafka-consumer
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic --group latest --consumer-property enable.auto.commit=false"
```

Create a file named `latest-scaled-object.yaml`:

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-test-so
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  scaleTargetRef:
    name: kafka-test-deployment
  triggers:
  - type: kafka
    metadata:
      topic: kafka-topic
      bootstrapServers: kafka-test-kafka-bootstrap.kafka-test-ns:9092
      consumerGroup: latest
      lagThreshold: '1'
      activationLagThreshold: '1'
      offsetResetPolicy: latest
```

Apply the configurations:

```bash
kubectl apply -f latest-deployment.yaml
kubectl apply -f latest-scaled-object.yaml
```

Follow the same testing pattern as with the earliest policy.

## Step 8: Test Multi-Topic

First, commit the partitions for both topics:

```bash
kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic --group multiTopic --from-beginning --consumer-property enable.auto.commit=true --timeout-ms 15000'
kubectl exec -n kafka-test-ns kafka-test-client -- bash -c 'kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic2 --group multiTopic --from-beginning --consumer-property enable.auto.commit=true --timeout-ms 15000'
```

Create a file named `multi-topic-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-test-deployment
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: kafka-consumer
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      containers:
      - name: kafka-consumer
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic 'kafka-topic'  --group multiTopic --from-beginning --consumer-property enable.auto.commit=false"
      - name: kafka-consumer-2
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic 'kafka-topic2' --group multiTopic --from-beginning --consumer-property enable.auto.commit=false"
```

Create a file named `multi-topic-scaled-object.yaml`:

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-test-so
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  scaleTargetRef:
    name: kafka-test-deployment
  triggers:
  - type: kafka
    metadata:
      bootstrapServers: kafka-test-kafka-bootstrap.kafka-test-ns:9092
      consumerGroup: multiTopic
      lagThreshold: '1'
      offsetResetPolicy: 'latest'
```

Apply the configurations:

```bash
kubectl apply -f multi-topic-deployment.yaml
kubectl apply -f multi-topic-scaled-object.yaml
```

Test scaling with messages to both topics.

## Step 9: Test Zero On Invalid Offset

Create a file named `zero-invalid-offset-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-test-deployment
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: kafka-consumer
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      containers:
      - name: kafka-consumer
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic-zero-invalid-offset --group invalidOffset --consumer-property enable.auto.commit=true"
```

Create a file named `zero-invalid-offset-scaled-object.yaml`:

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-test-so
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  scaleTargetRef:
    name: kafka-test-deployment
  triggers:
  - type: kafka
    metadata:
      topic: kafka-topic-zero-invalid-offset
      bootstrapServers: kafka-test-kafka-bootstrap.kafka-test-ns:9092
      consumerGroup: invalidOffset
      lagThreshold: '1'
      scaleToZeroOnInvalidOffset: 'true'
      offsetResetPolicy: 'latest'
```

Apply the configurations and test.

## Step 10: Test One On Invalid Offset

Create a file named `one-invalid-offset-deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kafka-test-deployment
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  replicas: 0
  selector:
    matchLabels:
      app: kafka-consumer
  template:
    metadata:
      labels:
        app: kafka-consumer
    spec:
      containers:
      - name: kafka-consumer
        image: confluentinc/cp-kafka:5.2.1
        command:
          - sh
          - -c
          - "kafka-console-consumer --bootstrap-server kafka-test-kafka-bootstrap.kafka-test-ns:9092 --topic kafka-topic-one-invalid-offset --group invalidOffset --from-beginning --consumer-property enable.auto.commit=true"
```

Create a file named `one-invalid-offset-scaled-object.yaml`:

```yaml
apiVersion: keda.sh/v1alpha1
kind: ScaledObject
metadata:
  name: kafka-test-so
  namespace: kafka-test-ns
  labels:
    app: kafka-test-deployment
spec:
  scaleTargetRef:
    name: kafka-test-deployment
  triggers:
  - type: kafka
    metadata:
      topic: kafka-topic-one-invalid-offset
      bootstrapServers: kafka-test-kafka-bootstrap.kafka-test-ns:9092
      consumerGroup: invalidOffset
      lagThreshold: '1'
      scaleToZeroOnInvalidOffset: 'false'
      offsetResetPolicy: 'latest'
```

Apply the configurations and test.

## Step 11: Clean Up

```bash
kubectl delete -f one-invalid-offset-scaled-object.yaml
kubectl delete -f one-invalid-offset-deployment.yaml
kubectl delete -f kafka-client.yaml
kubectl delete -f kafka-topic1.yaml
kubectl delete -f kafka-topic2.yaml
kubectl delete -f kafka-topic-zero-invalid-offset.yaml
kubectl delete -f kafka-topic-one-invalid-offset.yaml
kubectl delete -f kafka-cluster.yaml
helm uninstall --namespace kafka-test-ns kafka-test
kubectl delete namespace kafka-test-ns
```