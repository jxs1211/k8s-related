```sh
Reconcile (github.com/kedacore/keda/v2/controllers/keda - scaledobject_controller.go)
└── reconcileScaledObject (github.com/kedacore/keda/v2/controllers/keda - scaledobject_controller.go)
    └── requestScaleLoop (github.com/kedacore/keda/v2/pkg/scaling - scaledobject_controller.go)
        └── HandleScalableObject (github.com/kedacore/keda/v2/pkg/scaling - scale_handler.go)
            └── startScaleLoop (github.com/kedacore/keda/v2/pkg/scaling - scale_handler.go)
                └── checkScalers (github.com/kedacore/keda/v2/pkg/scaling/executor - scale_handler.go)
                    └── RequestScale (github.com/kedacore/keda/v2/pkg/scaling/executor - scale_scaledobjects.go)
                        └── scaleToZeroOrIdle (github.com/kedacore/keda/v2/pkg/scaling/executor - scale_scaledobjects.go)
```
