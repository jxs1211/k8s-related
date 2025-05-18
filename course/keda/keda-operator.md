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
