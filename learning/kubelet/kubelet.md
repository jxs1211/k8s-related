> The kubernetes version is 1.18.20

```sh
# Kubernetes Pod Addition Handling Call Tree (Reverse Order)
- `main` [ ]
  └─ `k8s.io/kubernetes/cmd/kubelet/app` (`kubelet.go`)
    └─ `NewKubeletCommand` [ ]
      └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
        └─ `Run` [ ]
          └─ `run` [ ]
            └─ `RunKubelet` [ ]
              └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
                └─ `startKubelet` [ ]
                  └─ `k8s.io/kubernetes/pkg/kubelet` (`server.go`)
                    └─ `Run` [ ]
                      └─ `k8s.io/kubernetes/pkg/kubelet` (`kubelet.go`)
                        └─ `syncLoop` [x]
                          └─ `syncLoopIteration` [ ]
                            └─ `HandlePodAdditions` 
                              └─ `k8s.io/kubernetes/pkg/kubelet` (`kubelet.go`)
```

```sh
# Kubernetes Kubelet Config Controller Initialization (Reverse Order)

- `main`
  └─ `k8s.io/kubernetes/cmd/kubelet/app` (`kubelet.go`)
    └─ `NewKubeletCommand`
      └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
        └─ `BootstrapKubeletConfigController`
          └─ `k8s.io/kubernetes/pkg/kubelet/kubeletconfig` (`server.go`)
            └─ `Bootstrap`
              └─ `k8s.io/kubernetes/pkg/kubelet/kubeletconfig` (`controller.go`)
                └─ `checkTrial`
                  └─ `k8s.io/kubernetes/pkg/kubelet/kubeletconfig` (`controller.go`)
```

```sh
# Kubernetes Kubelet Debugging Handlers Call Tree (Reverse Order)
- `main`
  └─ `k8s.io/kubernetes/cmd/kubelet/app` (`kubelet.go`)
    └─ `NewKubeletCommand`
      └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
        └─ `Run`
          └─ `run`
            └─ `RunKubelet`
              └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
                └─ `startKubelet`
                  └─ `k8s.io/kubernetes/pkg/kubelet` (`server.go`)
                    └─ `ListenAndServe`
                      └─ `k8s.io/kubernetes/pkg/kubelet/server` (`kubelet.go`)
                        └─ `ListenAndServeKubeletServer`
                          └─ `NewServer`
                            └─ `InstallDebuggingHandlers`
                              └─ `net/http/pprof` (`server.go`)
                                └─ `Profile`
                                  └─ `net/http/pprof` (`pprof.go`)
```
```sh
# Kubernetes Kubelet Log Serving Call Tree (Reverse Order)
- `Run`  
  └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
    └─ `run`  
      └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
        └─ `Run`  
          └─ `k8s.io/kubernetes/cmd/kubelet/app` (`hollow_kubelet.go`)
            └─ `RunKubelet`  
              └─ `k8s.io/kubernetes/cmd/kubelet/app` (`server.go`)
                └─ `startKubelet`  
                  └─ `k8s.io/kubernetes/pkg/kubelet` (`server.go`)
                    └─ `ListenAndServe`  
                      └─ `k8s.io/kubernetes/pkg/kubelet/server` (`kubelet.go`)
                        └─ `ListenAndServeKubeletServer`  
                          └─ `NewServer`  
                            └─ `InstallDebuggingHandlers`  
                              └─ `getLogs`  
                                └─ `k8s.io/kubernetes/pkg/kubelet` (`server.go`)
                                  └─ `ServeLogs`  
                                    └─ `k8s.io/kubernetes/pkg/kubelet` (`kubelet_pods.go`)
```
```sh
# cleanup evictable containers when init kubelet
- main (k8s.io/kubernetes/cmd/kubelet/app • kubelet.go)
  - NewKubeletCommand (k8s.io/kubernetes/cmd/kubelet/app • server.go)
    - Run (k8s.io/kubernetes/cmd/kubelet/app • server.go)
      - RunKubelet (k8s.io/kubernetes/cmd/kubelet/app • server.go)
        - createAndInitKubelet (k8s.io/kubernetes/pkg/kubelet • server.go)
          - StartGarbageCollection (k8s.io/kubernetes/pkg/kubelet/container • kubelet.go)
            - GarbageCollect (k8s.io/kubernetes/pkg/kubelet/kuberuntime • container_gc.go)
              - GarbageCollect (k8s.io/kubernetes/pkg/kubelet/kuberuntime • kuberuntime_gc.go)
                - evictContainers (k8s.io/kubernetes/pkg/kubelet/kuberuntime • kuberuntime_gc.go)
                  - evictableContainers (k8s.io/kubernetes/pkg/kubelet/kuberuntime • kuberuntime_gc.go)
                    - getKubeletContainers (k8s.io/kubernetes/pkg/kubelet/kuberuntime • kuberuntime_containers.go)
```

经过前面几个小节的讲解，Kubelet启动参数的收集已经完毕，接下来进入Kubelet实例的创建和启动。现在进入正文。

2.2.2 创建并启动Kubelet – RunKubelet()函数
经过run()函数的一通运作，Kubelet所依赖的信息都齐备了，现在可以创建Kubelet实例并启动这个实例了。RunKubelet()函数负责完成这两项任务，它被run()函数在最后部分进行调用。

RunKubelet()函数在开头部分首先做了一些信息获取与整理的创建准备工作，主要包括获取节点主机名、IP地址集合，设置Docker配置文件目录。Docker配置文件目录内有内含镜像库登录秘钥信息的文件（其实我们在本地使用docker login时，docker会自动生成这个文件，名为.dockercfg或config.json），需要告诉Kubelet该文件存放目录。

上述信息整理工作并非主要内容，RunKubelet()主要工作分为创建与启动两个步骤。

2.2.2.1 创建Kubelet实例
Kubelet实例即定义于文件kubernetes/pkg/kubelet/kubelet.go中的Bootstrap接口的实例，接口定义如下：
```sh
➜  kubernetes git:(v1.18.20) ✗ ag --go "type Bootstrap interface"
pkg/kubelet/kubelet.go
192:type Bootstrap interface {
```

函数createAndInitKubelet()负责为RunKubelet()制作一个上述接口的实例，这便是Kubelet实例。该实例的实际类型是结构体Kubelet，它与Bootstrap接口定义在同一个源文件中。但实际上createAndInitKubele()函数也是调用了其它函数进行Kubelet实例的创建的（这样一层套一层的方法调用是不是非常让人恶心？至少笔者是这么感觉的。为了获得些许的可读性，搞出个巨大的层级结构，让人厌恶），这个函数便是NewMainKubelet()，它实现在源文件kubernetes/pkg/kubelet/kubelet.go中。这是一个含有600行代码的大函数，利用收集到的参数值创建出Kubelet的各种子模块，放入Kubelet结构体实例作为结果返回。
```sh
➜  kubernetes git:(v1.18.20) ✗ ag --go "NewMainKubelet"
cmd/kubelet/app/server.go
1111:   // NewMainKubelet should have set up a pod source config if one didn't exist
1180:   k, err = kubelet.NewMainKubelet(kubeCfg,

pkg/kubelet/kubelet.go
400:// NewMainKubelet instantiates a new Kubelet object along with all the required internal modules.
402:func NewMainKubelet(kubeCfg *kubeletconfiginternal.KubeletConfiguration,
```

由它的源文件位置和名字首字母可知，NewMainKubelet()是包kubelet的公共方法，是该包开放给外部的API，它的参数列表很长，其上层的RunKubelet()、run()等函数从某种意义上说就是在为这个NewMainKubelet()函数准备参数。该方法的入参类型并不包含KubeletServer结构体，在2.2.1.1中我们介绍过这个结构体，它的用途并非如其名字所暗示的，代表一个Server，而只是KubeletConfiguration API基座结构体与KubeFlags结构体字段的集合，代表了Kubelet最原始的配置信息。在调用NewMainKubelet()时，参数信息会被从KubeletServer中抽取出来作为入参，所以并不会直接使用KubeletSever结构体。
```go
type KubeletServer struct {
	KubeletFlags
	kubeletconfig.KubeletConfiguration
}
```
从上到下，NewMainKubelet()函数做了如下事项：

(1) 校验参数。包括同步配置要求与容器实际情况的时间间隔是否有设置；IPTable的相关配置是否有冲突；Cloud Provider启用情况；
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L454
```

(2) 创建Node Informer。当没有跑在standalone模式下时Kubelet会与一个API Server连接，并跟踪API Server中Node API实例，这是日通过一个Informer来做的。这里创建出该Informer并启动它：
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L474
```

(3) 创建PodConfig。如果PodConfig没有在入参中设定，这里会新建一个PodConfig。PodConfig由与它同名的PodConfig结构体代表，定义于kubernetes/pkg/kubelet/config/config.go中。它的作用是实现一个观察者模式，集中所有Pod配置源中发生的对Pod改变，允许观察者注册监听，改变发生时PodConfig通知它们。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L492
```

(4) 创建容器垃圾收集配置。配置信息包括可被收集容器的最小“年龄”、单个Pod可具有的最大死容器数量、最多存在多少死容器。
```sh

```
(5) 创建镜像垃圾收集配置。

(6) 创建Pod驱逐配置。为了保持节点处于健康状态，当发生资源紧张时Kubelet会选择驱逐一些Pod进行资源释放。在此设置Pod驱逐的配置，例如需要进行驱逐的资源阈值等。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L525
```
(7) 创建Service Informer。当Kubelet没有运行在Standalone模式下，它与一个API Server有连接，这时创建Service Informer来监控API Server上的Service。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L539
```
(8) 创建监听器监听系统OOM事件。当有进程发生OOM时，极有可能是某个Pod的进程，Kubelet希望以事件形式记录下来，以便后续处理。记录是由2.2.1.2中介绍的依赖中Recorder属性来执行的。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L553
```
(9) 获取DNS集合。

(10) 创建一个TLS Transport，作为管理容器生命周期时所用的HTTP Client。

(11) 创建Kubelet实例。这一步利用入参、刚刚组织起来的信息共同创建一个Kubelet结构体实例，这也将是本方法最终返回的结果。节选部分代码如 2-10所示。
```go
klet := &Kubelet{
```

(12) 创建Cloud资源同步管理器。当上述Kubelet实例具有Cloud Provider时，创建一个资源同步管理器来从云平台读取相关资源信息。
```go
	if klet.cloud != nil {
		klet.cloudResourceSyncManager = cloudresource.NewSyncManager(klet.cloud, nodeName, klet.nodeStatusUpdateFrequency)
	}
```
(13) 创建Secret和ConfigMap管理器。因为这两类API资源涉及到下载至本地节点并挂载到相关Pod，所以Kubelet设立两个管理器专门应对。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L634
```
(14) 创建livenesss、 readiness、startup管理器。它们分别负责容器的健康、就绪和启动完毕状态的获取。
```go
	klet.livenessManager = proberesults.NewManager()
	klet.startupManager = proberesults.NewManager()
```
(15) 创建Pod管理器，Status管理器，资源分析器（Resource Analyzer）。Pod管理器配合上述Secret和ConfigMap将Pod关联的资源挂载；Status管理器负责将各种状态信息同步给API Server，而资源分析器则关注本地资源的使用情况。
```go
	klet.podManager = kubepod.NewBasicPodManager(mirrorPodClient, secretManager, configMapManager, checkpointManager)
	klet.statusManager = status.NewManager(klet.kubeClient, klet.podManager, klet)
	klet.resourceAnalyzer = serverstats.NewResourceAnalyzer(klet, kubeCfg.VolumeStatsAggPeriod.Duration)
```
(16) 获取Runtime Service，并创建Runtime Class管理器。Runtime Service用于同底层容器运行时交互，Kubelet实例可以从依赖中直接获取已经创建的Runtime Service。Runtime Class信息需要从API Server读取得来，这是Runtime Class管理器的责任。Runtime Class用于控制使用哪种容器运行时。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L691
```
(17) 创建PodWorkers。Kubelet利用这些Worker将API Server上对Pod的要求落实到本地。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L854
```
(18) 创建Container Runtime和Stream Runtime。Kubelet的这两个属性都被赋值为代码2-11创建的Runtime Manager实例。一个Runtime Manager就代表一个容器运行时
```go
	runtime, err := kuberuntime.NewKubeGenericRuntimeManager(
```
(19) 创建Node统计信息的Provider。该Provider将从cAdvisor中读取信息，并利用CRI接口获取容器信息进行返回。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L747
```
(20) 创建Pod生命周期事件生成器。一定会创建一个传统事件生成器（Kubelet的pleg属性），如果同时启用了基于事件的生成器，则也会生成一个基于事件的生成器（Kubelet的EventedPleg属性）。
```go
	klet.pleg = pleg.NewGenericPLEG(klet.containerRuntime, plegChannelCapacity, plegRelistPeriod, klet.podCache, clock.RealClock{})
```
(21) 为属性runtimeState创建实例，kubelet实例用它来记录容器运行时最后一次响应ping操作的时间戳。当ping被响应，Pod生命周期事件生成器也会生成相应事件予以记录。
```go
	klet.runtimeState = newRuntimeState(maxWaitForContainerRuntime)
	klet.runtimeState.addHealthCheck("PLEG", klet.pleg.Healthy)
```
(22) 创建容器垃圾收集器、容器删除器。
```go
	// setup containerGC
	containerGC, err := kubecontainer.NewContainerGC(klet.containerRuntime, containerGCPolicy, klet.sourcesReady)
	if err != nil {
		return nil, err
	}
	klet.containerGC = containerGC
	klet.containerDeletor = newPodContainerDeletor(klet.containerRuntime, integer.IntMax(containerGCPolicy.MaxPerPodContainer, minDeadContainerInPod))
```
(23) 创建镜像垃圾收集管理器。
```go
	imageManager, err := images.NewImageGCManager(klet.containerRuntime, klet.StatsProvider, kubeDeps.Recorder, nodeRef, imageGCPolicy, crOptions.PodSandboxImage)
```
(24) 创建证书管理器。如果启用了证书循环，则创建证书管理器来从API Server中获取最新的证书。
```go
		klet.serverCertificateManager, err = kubeletcertificate.NewKubeletServerCertificateManager(klet.kubeClient, kubeCfg, klet.nodeName, klet.getLastObservedNodeAddresses, certDirectory)
```
(25) 创建Pod探针（probe）管理器。
```go
	klet.probeManager = prober.NewManager(
```
(26) 创建Volume插件管理器和插件管理器，在此基础上创建Volume管理器。
```go
	klet.volumePluginMgr, err =
		NewInitializedVolumePluginMgr(klet, secretManager, configMapManager, tokenManager, kubeDeps.VolumePlugins, kubeDeps.DynamicPluginProber)
  klet.volumeManager = volumemanager.NewVolumeManager(
```
(27) 创建Pod驱逐管理器
```go
```
(28) 创建Pod Active 超时处理器，它可以统计Pod激活超时情况并上报。将它进行注册从而在Pod Sync时该处理器被调用。
```go
	// enable active deadline handler
	activeDeadlineHandler, err := newActiveDeadlineHandler(klet.statusManager, kubeDeps.Recorder, klet.clock)
	if err != nil {
		return nil, err
	}
	klet.AddPodSyncLoopHandler(activeDeadlineHandler)
	klet.AddPodSyncHandler(activeDeadlineHandler)
```
(29) 创建节点续租(lease)控制器。该控制器周期性向API Server上报节点存活状况。
```go
	klet.nodeLeaseController = nodelease.NewController(klet.clock, klet.heartbeatClient, string(klet.nodeName), kubeCfg.NodeLeaseDurationSeconds, klet.onRepeatedHeartbeatFailure)
```
(30) 创建关机管理器。

至此，Kubelet实例创建工作全部完成，NewMainKubelet()函数将上述实例返回给调用者createAndInitKubelet()函数，该方法并不做过多处理，只是：

(1) 调用所得实例的BirthCry()方法，它会（向Recorder）写一条log，记录启动了。
```go
	k.BirthCry()
```
(2) 调用所得实例的StartGarbageCollaction()方法，来启动垃圾收集。
```go
	k.StartGarbageCollection()
```
然后便将Kubelet实例返回给RunKubelet() 函数。

2.2.2.2 启动Kubelet实例
在Kubelet实例创建出来后，RunKubelet()函数紧接着对它进行了启动，启动前先判断是否只想“跑一次”启动过程，然后马上关闭，这只在特殊情况下使用，一般则是启动后一直处于运行状态，这时由函数startKubelet()进行启动，它依旧位于kubernetes/cmd/kubelet/app/server.go。
```go
	if runOnce {
		if _, err := k.RunOnce(podCfg.Updates()); err != nil {
			return fmt.Errorf("runonce failed: %v", err)
		}
		klog.Info("Started kubelet as runonce")
	} else {
		startKubelet(k, podCfg, &kubeServer.KubeletConfiguration, kubeDeps, kubeServer.EnableCAdvisorJSONEndpoints, kubeServer.EnableServer)
		klog.Info("Started kubelet")
	}
	return nil
```

代码 2-12：startKubelet() 函数

这个函数中，根据不同情况对Kubelet实例的四个方法进行了调用。比较简单的是1178行到1186行对三个ListenAndServe方法的调用，它们分别会启动Kubelet Web Server端口、只读信息端口（对这个端口的请求不会有身份验证）的监听和以gRPC协议暴露Pod资源信息端点。而第1175行对Run()方法的调用则包含更多内容，包括：
```sh
https://github.com/kubernetes/kubernetes/blob/09877dcea4157b93d109052446240a17998f4e24/cmd/kubelet/app/server.go#L1133
func startKubelet(k kubelet.Bootstrap, podCfg *config.PodConfig, kubeCfg *kubeletconfiginternal.KubeletConfiguration, kubeDeps *kubelet.Dependencies, enableCAdvisorJSONEndpoints, enableServer bool) {
	// start the kubelet
	go k.Run(podCfg.Updates())
	// start the kubelet server
	if enableServer {
		go k.ListenAndServe(net.ParseIP(kubeCfg.Address), uint(kubeCfg.Port), kubeDeps.TLSOptions, kubeDeps.Auth, enableCAdvisorJSONEndpoints, kubeCfg.EnableDebuggingHandlers, kubeCfg.EnableContentionProfiling)
	}
	if kubeCfg.ReadOnlyPort > 0 {
		go k.ListenAndServeReadOnly(net.ParseIP(kubeCfg.Address), uint(kubeCfg.ReadOnlyPort), enableCAdvisorJSONEndpoints)
	}
	if utilfeature.DefaultFeatureGate.Enabled(features.KubeletPodResources) {
		go k.ListenAndServePodResources()
	}
}
```
(1) 如果Kubelet实例还没有LogServer，并且配置允许外界对Node的Log文件进行查询，则做一个Http Server提供/logs/端点，供外界通过它查询节点的log信息。
```go
	s.addMetricsBucketMatcher("logs")
	ws = new(restful.WebService)
	ws.
		Path(logsPath)
	ws.Route(ws.GET("").
		To(s.getLogs).
		Operation("getLogs"))
	ws.Route(ws.GET("/{logpath:*}").
		To(s.getLogs).
		Operation("getLogs").
		Param(ws.PathParameter("logpath", "path to the log").DataType("string")))
	s.restfulCont.Add(ws)
```
(2) 启动cloud资源管理器。
```go
	// Start the cloud provider sync manager
	if kl.cloudResourceSyncManager != nil {
		go kl.cloudResourceSyncManager.Run(wait.NeverStop)
	}
```

(3) 初始化不依赖容器运行时的内部子模块，由Kubelet结构体的InitializeModules()方法完成。其内部会初始化Prometheus、初始化各种目录（如Pod资源目录、插件目录、Container Log目录等等）。还包含启动镜像管理器，主要是镜像的垃圾回收、启动服务器证书管理器 – 如果有的话、启动OOM watcher和系统资源分析器。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L1417
```
(4) 启动Volume管理器。
```go
	go kl.volumeManager.Run(kl.sourcesReady, wait.NeverStop)
```
start to sync node status and lease
```go
	if kl.kubeClient != nil {
		// Start syncing node status immediately, this may set up things the runtime needs to run.
		go wait.Until(kl.syncNodeStatus, kl.nodeStatusUpdateFrequency, wait.NeverStop)
		go kl.fastStatusUpdateOnce()

		// start syncing lease
		go kl.nodeLeaseController.Run(wait.NeverStop)
	}
```
Set up iptables util rules
```go
	if kl.makeIPTablesUtilChains {
		kl.initNetworkUtil()
	}
```
Start a goroutine responsible for killing pods (that are not properly handled by pod workers).
```go
	go wait.Until(kl.podKiller.PerformPodKillingWork, 1*time.Second, wait.NeverStop)
```
(5) 启动Status管理器。
```go
	// Start component sync loops.
	kl.statusManager.Start()
	kl.probeManager.Start()
```
(6) 启动 Runtime Class管理器。
```go
	// Start syncing RuntimeClasses if enabled.
	if kl.runtimeClassManager != nil {
		kl.runtimeClassManager.Start(wait.NeverStop)
	}
```
(7) 启动Pod生命周期事件生成器。
```go
	// Start the pod lifecycle event generator.
	kl.pleg.Start()
```
(8) 最后，启动Pod的同步循环，这个我们放在下一小节单独讲解。
```go
	// Start the pod sync loop.
	kl.syncLoop(updates, kl)
```
由此可见，startKubelet()函数的主体是Kubelet实例的Run()方法，而该方法就是启动所有在创建阶段制作的各种管理器、控制器，它的最后一项工作 – 启动Pod的同步循环 - 开启了Kubelet十分关键的工作，即观测对Pod的操作需求，并落实到本地容器运行时，我们将在下一小节着重讲解。

上一篇基本将Kubelet的启动讲完了，今天来扫个尾，将Pod同步循环和Kubelet的几个Web Server的启动讲解一下。这两个话题后续值得单独剖析，会有后续章节深入讲解。

2.2.3 启动Pod同步循环
对Pod的管理是Kubelet的核心工作之一。在Kubelet的启动过程中，最后一项工作是启动一个永不结束的循环，观测对Pod的变化需求并落实。具体来说便是调用Kubelet.syncLoop()方法，它定义于kubernetes/pkg/kubelet/kubelet.go。

Run()方法通过调用Kubelet.syncLoop()方法来启动这个同步循环，syncLoop()内部包含一个永不结束的for循环 - 除非出现异常。完成了syncLoop()的调用也就完成了Pod同步循环的启动。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L1845
```

代码 2-13 syncLoop()中不结束的for循环

syncLoop()方法最大的责任就是维护这个循环的运行，同时为包含Pod同步逻辑的方法syncLoopIteration()准备入参。syncLoopIteration()方法是Pod同步逻辑的最上层拥有者，对其内容的详细剖析非常复杂，留到后面的专有章节进行。
```sh
https://github.com/kubernetes/kubernetes/blob/070322921d35c781bd0c94a6527dd2a819362210/pkg/kubelet/kubelet.go#L1919
```
最后谈一下Pod同步在做什么事情，其实比较容易理解，当用户通过多种渠道（API Server，本地配置文件，Kubelet Server的HTTP端点）发出Pod的增删改指令时，Kubelet需要启动落实工作；当Pod内的容器发生变动时有可能需要更新Pod自身的状态，这也需要Kubelet进行操作；还有其它场景，当Pod探针（存活、就绪、启动完毕）发出信号时，同样需要Kubelet更新Pod状态。同时，Kubelet内部的其它组件在运作过程中同样可能要求对某些Pod进行同步操作。Pod的同步是一个复杂的过程，后续会有专有章节展开讲解。
```go
func (kl *Kubelet) syncLoopIteration(configCh <-chan kubetypes.PodUpdate, handler SyncHandler,
	syncCh <-chan time.Time, housekeepingCh <-chan time.Time, plegCh <-chan *pleg.PodLifecycleEvent) bool {
	select {
	case u, open := <-configCh:
		if !open {
			klog.Errorf("Update channel is closed. Exiting the sync loop.")
			return false
		}
		switch u.Op {
		case kubetypes.ADD:
			klog.V(2).Infof("SyncLoop (ADD, %q): %q", u.Source, format.Pods(u.Pods))
			handler.HandlePodAdditions(u.Pods)
		case kubetypes.UPDATE:
			klog.V(2).Infof("SyncLoop (UPDATE, %q): %q", u.Source, format.PodsWithDeletionTimestamps(u.Pods))
			handler.HandlePodUpdates(u.Pods)
		case kubetypes.REMOVE:
			klog.V(2).Infof("SyncLoop (REMOVE, %q): %q", u.Source, format.Pods(u.Pods))
			handler.HandlePodRemoves(u.Pods)
		case kubetypes.RECONCILE:
			klog.V(4).Infof("SyncLoop (RECONCILE, %q): %q", u.Source, format.Pods(u.Pods))
			handler.HandlePodReconcile(u.Pods)
		case kubetypes.DELETE:
			klog.V(2).Infof("SyncLoop (DELETE, %q): %q", u.Source, format.Pods(u.Pods))
			// DELETE is treated as a UPDATE because of graceful deletion.
			handler.HandlePodUpdates(u.Pods)
		case kubetypes.RESTORE:
			klog.V(2).Infof("SyncLoop (RESTORE, %q): %q", u.Source, format.Pods(u.Pods))
			// These are pods restored from the checkpoint. Treat them as new
			// pods.
			handler.HandlePodAdditions(u.Pods)
		case kubetypes.SET:
			// TODO: Do we want to support this?
			klog.Errorf("Kubelet does not support snapshot update")
		}

		if u.Op != kubetypes.RESTORE {
			kl.sourcesReady.AddSource(u.Source)
		}
	case e := <-plegCh:
		if isSyncPodWorthy(e) {
			// PLEG event for a pod; sync it.
			if pod, ok := kl.podManager.GetPodByUID(e.ID); ok {
				klog.V(2).Infof("SyncLoop (PLEG): %q, event: %#v", format.Pod(pod), e)
				handler.HandlePodSyncs([]*v1.Pod{pod})
			} else {
				// If the pod no longer exists, ignore the event.
				klog.V(4).Infof("SyncLoop (PLEG): ignore irrelevant event: %#v", e)
			}
		}

		if e.Type == pleg.ContainerDied {
			if containerID, ok := e.Data.(string); ok {
				kl.cleanUpContainersInPod(e.ID, containerID)
			}
		}
	case <-syncCh:
		// Sync pods waiting for sync
		podsToSync := kl.getPodsToSync()
		if len(podsToSync) == 0 {
			break
		}
		klog.V(4).Infof("SyncLoop (SYNC): %d pods; %s", len(podsToSync), format.Pods(podsToSync))
		handler.HandlePodSyncs(podsToSync)
	case update := <-kl.livenessManager.Updates():
		if update.Result == proberesults.Failure {
			pod, ok := kl.podManager.GetPodByUID(update.PodUID)
			if !ok {
				// If the pod no longer exists, ignore the update.
				klog.V(4).Infof("SyncLoop (container unhealthy): ignore irrelevant update: %#v", update)
				break
			}
			klog.V(1).Infof("SyncLoop (container unhealthy): %q", format.Pod(pod))
			handler.HandlePodSyncs([]*v1.Pod{pod})
		}
	case <-housekeepingCh:
		if !kl.sourcesReady.AllReady() {
			klog.V(4).Infof("SyncLoop (housekeeping, skipped): sources aren't ready yet.")
		} else {
			klog.V(4).Infof("SyncLoop (housekeeping)")
			if err := handler.HandlePodCleanups(); err != nil {
				klog.Errorf("Failed cleaning pods: %v", err)
			}
		}
	}
	return true
}
```
2.2.4 启动Web Server
在2.2.2.2中讲到，startKubelet()方法会根据要求启动三个Web Server：

Kubelet Web Server
```sh
https://github.com/kubernetes/kubernetes/blob/09877dcea4157b93d109052446240a17998f4e24/cmd/kubelet/app/server.go#L1139
```
只读的Kubelet Web Server
```sh
https://github.com/kubernetes/kubernetes/blob/09877dcea4157b93d109052446240a17998f4e24/cmd/kubelet/app/server.go#L1143
```
可返回PodResource实例的gRPC Server
```sh
https://github.com/kubernetes/kubernetes/blob/09877dcea4157b93d109052446240a17998f4e24/cmd/kubelet/app/server.go#L1146
```
前两个Server区别主要体现在两个方面。第一，前者会带有登录和鉴权模块，所以没有权限情况下去请求它的端点是得不到结果的，而后者则不然，它没有挂这一模块所以均会给予响应。同时前者提供使用TLS的可能性而后者不会启用TLS。笔者认为后者的存在主要是为节点本地应用程序与Kubelet交互提供方便，一般来说如果用户可以登录到服务器本机，那么足以说明他角色的特殊性了。第二，前者会提供/debug端点，它对外暴露很多调试用信息，而后者则不会。

第三个Server是一个gRPC Server，它通过gRPC协议对外暴露当前节点上Pod Resource信息。所谓“Pod Resource”实际上指当前节点能为集群提供的资源，包括几类：Pod、CPUs、Memory、Device和动态资源。程序上它们被以类似Kubernetes API的方式定义在kubernetes/pkg/kubelet/apis/podresources包中，目前有两个外部版本分别是v1alpha1和v1。不过它们并不像API Server中的API实例一样由系统或用户创建并存储在API Server中，这些Pod Resource在Kubelet中是由Kubelet内部的Pod Manager和Container Manager动态从系统获取的。这个Server的作用就是让外部能够获得当前节点的Pod相关资源信息。

2.2.4.1 Kubelet Web Server的启动
由于只读Kubelet Web Server的启动与一般Kubelet Web Server非常类似，所以我们不再单独介绍它。Kubelet Web Server的启动最终由如代码2-14所示的方法完成。

Image
```sh
https://github.com/kubernetes/kubernetes/blob/9802bfcec0580169cffce2a3d468689a407fa7dc/pkg/kubelet/server/server.go#L133
代码 2-14 启动Kubelet Web Server
```
由代码2-14可以看出，它基本就是在利用Go的http模块创建一个server，然后按需配置TLS，最后启动这个Server，可以说平淡无奇。最需要深究的就是第163行的NewServer()函数，它会制作一个Web请求处理器，作为响应Server请求的对象。而只读Kubelet Web Server与普通Kubelet Web Server之间的区别也主要在这个方法中进行设置。代码2-15展示了这个方法的实现，它定义于kubernetes/pkg/kubelet/server/server.go中。
```sh
https://github.com/kubernetes/kubernetes/blob/9802bfcec0580169cffce2a3d468689a407fa7dc/pkg/kubelet/server/server.go#L146sh
```

代码2-15的275行通过结构体Server制作了一个server实例，它最终会被作为结果返回出去。在返回之前，该方法会根据入参的情况对server实例进行一些设置：如果入参含有非空auth，则调用server自己的方法添加登录与鉴权的过滤器；如果功能Kubelet Tracing被打开，则安装跟踪过滤器；如果启用了Debug处理器，则添加Debug端点。另外server会通过自身的

InstallDefaultHandlers()方法安装固定的端点。对Web Server所能处理的所有请求这里不再深入，将分单独章节讲解。

2.2.4.2 Pod Resource gRPC Server的启动
代码2-12 的1185行调用的方法定义如代码2-16所示，它实际上也没有做实际的gRPC Server的创建和启动工作，而是为之准备参数，进而调用server.ListenAndServePodResources()方法去启动该Server。
```sh
https://github.com/kubernetes/kubernetes/blob/785fac68266a6024ae8e025b65371b26030aa9f3/pkg/kubelet/apis/podresources/server.go#L45
```

请读者注意2814行的providers变量，当gRPC Server响应请求时，正是通过它来获取Pod资源信息的。providers的Pods属性由Kubelet实例的PodManager属性填充；而Devices、CPUs、Memory、DynamicResources则均是由Kubelet的ContainerManager属性来填充的。这个信息很重要，后续剖析该Server对请求的响应时只需到PodManager和ContainerManager的方法中去寻找。

server.ListenAndServePodResources()的定义如代码2-17所示。
```sh
https://github.com/kubernetes/kubernetes/blob/09877dcea4157b93d109052446240a17998f4e24/cmd/kubelet/app/server.go#L1146
```
Image

代码2-17 启动Pod Resource Server

第222行首先创建出了一个gRPC Server实例，名为server，然后将Pod Resource的v1版和v1alpha1版都“注册”进这个server，从而对外暴露。对于所有端点以及对应响应方法的介绍将在专有章节进行。

在2.2.3节中简单提及了Kubelet.syncLoop() - 即Pod同步循环的重要作用，现在深入到方法syncLoop()，分析Pod同步相关工作如何实现。图3-1概括了整个同步过程(我知道这个图在手机上看字将非常小，不要急，后续小节细讲时会拆解出大图)。图中有绿、黄、蓝、红四个颜色的矩形，除了黄色之外，其它矩形均代表一个方法。箭头所指方向为逻辑执行顺序。同种颜色的矩形具有相同的责任，彼此之间较相似。

Image

图3-1： Pod同步循环的工作内容

我们将逐个颜色讲解上述矩形所代表的方法、对象。

1. 绿 – 启动并保持Pod同步循环
由图3-1可见Pod同步是一个繁琐的过程，但这一切的根源 - Kubelet.syncLoop()方法却只有40行代码。它最核心的任务是维护这样一个同步循环，具体的工作都包含在循环内部去执行了。syncLoop()将每次循环的逻辑分离到方法syncLoopIteration()方法中去，从而使得自身逻辑更为清晰，二者关系如图3-2所示。

Image

图3-2： syncLoop( )方法与syncLoopIteration( )方法

前面小节中代码2-13显示syncLoop()会在一个永不结束的for循环中持续调用syncLoopIteration()方法，只有当后者返回false时才结束循环。syncLoopIteration()方法会在后面讲解，它内部是一个没有default子句的巨大的select语句，它的每个case都在等待一个管道的输出，这些管道均代表某种需要Pod同步执行的场景。当所有管道都没有输出时它便会一直等待；而某个管道有输出时syncLoopIteration()就完成一次运行并返回true，这时syncLoop的for循环会再次启动它。这便是最上层的程序运转逻辑。

来看一下syncLoop()方法的入参，除了常规的Context外，还有名为updates和名为handler的两个入参。

updates入参

这个参数的类型为有数据传出的管道（channel），管道内传递来的对象包含了对Pod的操作信息。这个入参涵盖了Pod的三个配置变更来源：通过http端点传入、通过本地文件定义、通过API Server定义。Kubelet会监控以上源头的Pod定义变化，组织成对Pod的操作信息，放入updates管道。该管道的另一端正是syncLoopIteration()方法，它便会依据不同的场景来展开Pod同步。

handler入参

正如其名字所暗示的，当有Pod变化需要被处理时，正是由该入参解读变化、分解为Pod Work可以消化的信息并传递给它，它就是干活儿的那个角色。handler的类型是一个接口，定义如代码3-1所示。

Image

代码3-1: kubernetes/pkg/kubelet/kubelet.go

由这段代码看到handler可以处理增、改、删等Pod操作，它们会在后面的syncLoopIteration()中被调用。前文介绍的Kubelet实例是实现了该接口的，它具有所有这些方法。有趣的是，当Run()中调用syncLoop()方法时，Kubelet实例自身被作为handler实参传入给了syncLoop() - 这不会出语法错误因为Kubelet实例具有该接口定义的所有方法，也就是说“一个实例被作为实参去调用了自身的一个方法”。可能看代码更容易理解，注意Run()的第1604行：

Image

代码3-2

syncLoop()不会直接使用updates和handler这两个入参，它们都是为syncLoopIteration()方法准备的。而除了这两个参数，syncLoopIteration()还需要更多入参从而涵盖所有会触发Pod同步的源。为此syncLoop()方法制作了两个管道（准确地说是Time Ticker，内含管道）来定期触发Pod同步运转，不太好理解但后面会再提及。如代码3-3所示。syncTicker是每秒触发一次，也即它内含的管道每一秒都会有一个输出；housekeepingTicker则默认每两秒触发一次，由其名字便可知道，它会通知Kubelet清除那些已经完成任务或失败的Pod。

Image

代码3-3: 两个Time Ticker

除了updates管道、syncTicker和housekeepingTicker外，前文提到的Pod Lifecycle Event Generator也会产生代表pod发生变化的事件，那么syncLoop()也会将其提供给syncLoopIterator()。Kubelet实例的pleg属性具有方法Watch()，它返回的管道正好派这个用场，syncLoop()调用它并将结果作为syncLoopIterator()的入参之一。至此，所有参数齐备了，对syncLoopIterator()的调用如代码2-13的2291行所示。下一小节进入syncLoopIterator()方法，讲解其实现。

2. 黄 – 代表Pod需同步的通道
接下来看syncLoopIterator() 的实现。它的主体实际上就是图3-1的黄色矩形部分，放大后就是图3-3。每个矩形都代表一个管道，我们应该对这五个矩形中的第一、二、四、五不陌生，它们对应着syncLoop()所准备的四个入参。而第三个矩形代表了另外三个管道：Pod的存活探针、就绪探针、启动完毕探针分别通过其中一个管道通知外界探针状态改变，Kubelet需要将这种改变反馈给外界 – 如API Server。

Image

图 3-3 syncLoopIterator()观测的管道

整个syncLoopIteration()方法就是一个巨大的select语句，其每个case子句分别等待从一个管道中传来的输入，由于该select没有default子句，所以这种等待会直到其中一个管道有输出，处理完毕后返回true给syncLoop()。所以简单说，这个方法就是在观测所有可能触发Pod同步的源，一旦需求发生，驱动同步发生。代码3-4显示了这个方法的框架。

Image

代码3-4：syncLoopIteration方法框架

2.1 Pod配置变动 – configCh
当用户通过前面提到的三种渠道来创建、修改、删除Pod时，configCh管道就会输出类型如下的对象：

Image

可见该对象会包含目标Pod、做什么操作以及请求来源。在监听configCh管道的case子句中，syncLoopIteration()方法撰写逻辑来根据不同的Op（增，删，改，再调协）来分别处理。而对于每种操作的处理逻辑也被封装在了handler入参内，只需直接调用便可，看代码3-5。

Image

代码3-5：处理Pod的增删改

前面讲过，handler的类型实现了代码3-1所示的SyncHandler接口，故具有HandlePodAdditions()等方法。关于对该接口的实现细节会在第3节中讲解。

2.2 Pod生命周期事件 – plegCh
pleg会产生如代码3-6所定义的事件。

Image

代码3-6 PLEG会产生的事件类型

这是一系列由Pod内容器的变化而产生的事件。除了“ContainerRemoved”之外所有事件发生后都需要Pod同步，这时handler.HandlePodSyncs()会被调用；当“ContainerDied”发生时则要进行Pod内的清理，handler的cleanUpContainersInPod()方法会被调用完成这个工作。

Image

代码3-7 处理plegCh管道的输出

注意代码3-7第2385行对cleanUpContainersInPod()的调用实际上是以

syncLoopIteration()的接收者“kl”为接收者进行的，而kl和handler是同一个实例，所以2385行等价于如下语句：

handler.cleanUpContainersInPod(e.ID, containerID)
2.3 Pod探针事件
Pod的三类探针在有状态变化时也需要触发Pod的同步，处理逻辑如代码3-8所示。当发生探针事件时大致有两类操作：第一通过Kubelet实例的statusManager属性去通知API Server；第二是通过函数handleProbeSync()去驱动Pod Worker来同步Pod，这个函数本质上只是包装了handler.HandlePodSyncs()方法，增加一些额外的log，所以最终还是handler去做。

Image

代码3-8 处理探针事件

2.4 定期同步 – syncCh
syncCh这个管道每秒定期输出一个信号，在等待该管道的case子句内，syncLoopIteration()方法先收集已经准备好同步的Pod和被其它组件要求再同步的Pod，这是借助Kubelet实例的getPodsToSync()方法得到的；然后依旧简单调用handler.HandlePodSyncs()方法驱动Pod Worker去同步。

Image

代码3-9 处理定期同步

2.5 定期清理 – housekeepingCh
Kubelet需要定期对系统进行清理工作，包括移除已停止的Pod的Pod Worker、清理不再需要的Pod、移除已无用的Volume。housekeepingCh管道负责向Kubelet发送清理信号，默认情况下每2秒钟进行一次，但用户可以修改该值。清理的逻辑被封装在handler.HandlerPodCleanups()方法中。

Image

代码 3-10 House Keeping处理逻辑

至此，syncLoopIteration()方法的主要逻辑介绍完毕。可以看到handler承接了所有具体的Pod同步工作，它的各个方法为syncLoopIteration()方法所用从而大大简化了它的逻辑。接下来将进入handler的这一系列方法。

————————————————————————————

后续我们将讲解蓝、红两种颜色的元素，它们将是syncLoop()最为核心的内容，敬请期待

