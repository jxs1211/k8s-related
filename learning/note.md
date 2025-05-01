方海涛 sealos作者，五年以上容器平台与系统研发经验
kube-dev-book
https://github.com/fanux/kube-dev-book/tree/master

《kubernetes 权威指南》第4版书中示例源码
纯手打，顺便用注释做简单的笔记
基于Kubernetes v1.18环境测试
https://github.com/callmer00t/kubeguide_example

海龙张， 立吧
Kubernetes文章列表
https://mp.weixin.qq.com/s/TVYL3tILuVydtODpzUAgDQ
以下是基于文章内容的 **Mermaid 格式概念关系图** 和 **核心逻辑解析**：

### **Mermaid 图表：cGroup 核心概念与关系**
```mermaid
graph TD
    A[cGroup 体系] --> B[基本概念]
    A --> C[层级形成规则]
    A --> D[落地形态]

    B --> B1[子系统（Subsystem）]
    B --> B2[控制组（Control Group）]
    B --> B3[层级（Hierarchy）]
    B --> B4[任务（Task）]

    B1 -->|示例| B1a["CPU、内存、网络带宽（共12种）"]
    B2 -->|功能| B2a["资源分配单位，限制关联任务的资源"]
    B3 -->|特性| B3a["树形结构，关联子系统与控制组"]
    B4 -->|本质| B4a["进程（Process）"]

    C --> C1["规则1：层级与子系统挂钩"]
    C --> C2["规则2：子系统多层级挂载限制"]
    C --> C3["规则3：根cGroup自动创建"]
    C --> C4["规则4：子进程继承cGroup"]

    C1 -->|示例| C1a["层级A: CPU+内存<br>层级B: 网络"]
    C2 -->|反例| C2a["CPU+内存+网络不可同挂一个层级"]
    C3 -->|默认| C3a["所有进程初始属于根cGroup"]
    C4 -->|可操作| C4a["子进程可被移出父cGroup"]

    D --> D1["文件系统映射"]
    D1 -->|路径| D1a["/sys/fs/cgroup"]
    D1 -->|结构| D1b["层级目录 → 子cGroup目录 → 资源限制文件"]
```

### **关键概念解析**
1. **子系统（Subsystem）**  
   - **本质**：资源类型（如CPU、内存）。  
   - **作用**：定义cGroup管理的资源类别，共12种。

2. **控制组（Control Group）**  
   - **功能**：资源分配的最小单位，通过限制值约束任务资源使用。  
   - **关系**：隶属于某个层级，形成树形结构。

3. **层级（Hierarchy）**  
   - **核心规则**：  
     - 一个层级可挂载多个子系统（如CPU+内存）。  
     - 同一子系统不可跨层级挂载（除非层级仅含该子系统）。  
   - **树形结构**：根cGroup默认包含所有进程，子cGroup继承父资源限制。

4. **任务（Task）**  
   - **本质**：进程，通过关联到cGroup受资源限制。

5. **落地形态**  
   - **文件系统映射**：通过`/sys/fs/cgroup`目录管理，层级和cGroup表现为目录，资源限制值为文件。

### **逻辑流程图：cGroup 资源分配流程**
```mermaid
flowchart TB
    Start[创建层级] --> Mount[挂载子系统]
    Mount --> CreateCG[创建控制组]
    CreateCG --> SetLimit[设置资源限制]
    SetLimit --> AttachTask[关联任务]
    AttachTask --> Effect[资源限制生效]

    subgraph 层级管理
        Mount -->|可选| MultiMount["多子系统挂载（如CPU+内存）"]
    end

    subgraph 任务控制
        AttachTask --> Inherit["子进程自动继承cGroup"]
        Inherit --> Modify["可手动移出cGroup"]
    end
```

### **总结**
1. **核心思想**：通过层级树和子系统挂钩，将进程分组并分配资源配额。  
2. **关键规则**：  
   - 子系统的挂载需符合层级唯一性约束。  
   - 进程的cGroup关系可继承和修改。  
3. **实际应用**：Docker等容器技术利用cGroup实现资源隔离（如`docker run --memory=1g`）。  

若需补充具体操作命令或扩展子系统类型，可进一步调整图表内容。