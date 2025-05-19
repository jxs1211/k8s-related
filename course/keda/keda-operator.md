Based on v2.17.1

## bootup process
### operator
### metrics-apiserver
#### collect metrics from   
  - GetMetricSpecForScaling
  - GetMetricsAndActivity
##### expose metrics to HPA controller for scaling
### webhooks
#### validate KEDA CRD

## reconcile process

- reconcileScaledObject
- startPushScaler
- startScaleLoop 
	checkScalers

### Get scaledObject state
isActive, isError, metricsRecords, activeTriggers, err := h.getScaledObjectState
- GetScalersCache
  - using exist scaler
  - create new scaler
- Build scaler
  - NewXXXScaler
    - parse(validate) XXX metadata
- getScalerState
  - GetMetricSpecForScaling
  - GetMetricsAndActivity
### scaling decision
h.scaleExecutor.RequestScale(ctx, obj, isActive, isError, &executor.ScaleExecutorOptions{ActiveTriggers: activeTriggers})
#### Scaling 0~1 with Keda operator
  - scaling up
    pkg/scaling/executor/scale_scaledobjects.go
113:                    e.scaleFromZeroOrIdle(ctx, logger, scaledObject, currentScale, options.ActiveTriggers)
#### Scaling 1~N with HPA controller
  - scaling down
    pkg/scaling/executor/scale_scaledobjects.go
169:                    e.scaleToZeroOrIdle(ctx, logger, scaledObject, currentScale)


## 2 Layers validation
CRD
- tag validation
- admission webhook
Trigger
- runtime validation

#### validation example
```mermaid
graph TD
    A[Reconcile<br>github.com/kedacore/keda/v2/controllers/keda<br>scaledobject_controller.go]
    A --> B[reconcileScaledObject<br>same file]
    B --> C[requestScaleLoop<br>github.com/kedacore/keda/v2/pkg/scaling<br>scaledobject_controller.go]
    C --> D[HandleScalableObject<br>github.com/kedacore/keda/v2/pkg/scaling<br>scale_handler.go]
    D --> E[startScaleLoop<br>same file]
    E --> F[checkScalers<br>same file]
    F --> G[getScaledObjectState<br>same file]
    G --> H[GetScalersCache<br>same file]
    H --> I[performGetScalersCache<br>same file]
    I --> J[buildScalers<br>github.com/kedacore/keda/v2/pkg/scaling<br>scalers_builder.go]
    J --> K[buildScaler<br>github.com/kedacore/keda/v2/pkg/scalers<br>scalers_builder.go]
    K --> L[NewElasticsearchScaler<br>github.com/kedacore/keda/v2/pkg/scalers<br>elasticsearch_scaler.go]
    L --> M[parseElasticsearchMetadata<br>github.com/kedacore/keda/v2/pkg/scalers/scalersconfig<br>elasticsearch_scaler.go]
    M --> N[TypedConfig<br>github.com/kedacore/keda/v2/pkg/scalers/scalersconfig<br>typed_config.go]
    N --> O[parseTypedConfig<br>github.com/kedacore/keda/v2/pkg/scalers<br>typed_config.go]
    O --> P[Validate<br>github.com/kedacore/keda/v2/pkg/scalers<br>elasticsearch_scaler.go]
```
