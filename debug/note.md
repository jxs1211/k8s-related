# Understanding Go Build Flags: gcflags and ldflags

When building Go applications, you can use various flags to control the compilation and linking process. Two important sets of flags are `gcflags` (Go compiler flags) and `ldflags` (linker flags).

## gcflags (Go Compiler Flags)

The `gcflags` control how the Go compiler behaves during compilation. These flags are passed to the compiler using the `-gcflags` option.

### Common gcflags:

1. **-N**: Disables optimizations
2. **-l**: Disables inlining
3. **-m**: Prints optimization decisions
4. **-S**: Prints assembly listing
5. **-race**: Enables race detection

### Examples:

```bash
# Disable optimizations and inlining (useful for debugging)
go build -gcflags="-N -l" main.go

# Print optimization decisions
go build -gcflags="-m" main.go

# Print assembly code
go build -gcflags="-S" main.go
```

## ldflags (Linker Flags)

The `ldflags` control the behavior of the Go linker. These are particularly useful for injecting values into your program at build time, such as version information or build timestamps.

### Common ldflags:

1. **-s**: Omits the symbol table and debug information
2. **-w**: Omits DWARF symbol table
3. **-X**: Sets the value of a string variable in the program

### Examples:

```bash
# Strip debug information to reduce binary size
go build -ldflags="-s -w" main.go

# Inject version information at build time
go build -ldflags="-X main.Version=1.0.0 -X main.BuildTime=$(date +%Y-%m-%d)" main.go
```

## Practical Use Cases

### 1. Setting Version Information

```go
package main

import "fmt"

// These variables will be set at build time
var (
    Version   = "dev"
    BuildTime = "unknown"
)

func main() {
    fmt.Printf("Version: %s, Build Time: %s\n", Version, BuildTime)
}
```

Build with:
```bash
go build -ldflags="-X main.Version=1.2.3 -X main.BuildTime=$(date +%Y-%m-%d)" main.go
```

### 2. Optimizing Binary Size

To create a smaller binary by removing debug information:
```bash
go build -ldflags="-s -w" main.go
```

### 3. Debugging

For better debugging experience:
```bash
go build -gcflags="all=-N -l" main.go
```

## Advanced Usage

You can target specific packages with these flags:

```bash
# Apply flags to all packages
go build -gcflags="all=-N -l" ./...

# Apply flags to specific package
go build -gcflags="github.com/user/package=-N -l" ./...
```

These build flags are particularly useful when working with Kubernetes-related projects like KEDA, where you might want to embed version information or optimize binaries for deployment.

        Too many current requests. Your queue position is 1. Please wait for a while or switch to other models for a smoother experience.

---

### **Go-Delve Debugging Commands**

#### **Basic Debugging (Terminal Interface)**
```bash
dlv debug [package]       # Debug a package
dlv test [package]        # Debug tests in a package
dlv exec <executable>     # Debug a pre-built binary
dlv attach <pid>          # Attach to a running process
dlv core <executable> <core>  # Debug a core dump
dlv replay <rr-trace>     # Replay an rr execution trace
```

#### **Execution Tracing**
```bash
dlv trace [package] <regexp>  # Trace program execution matching regex
```

---

### **Headless Debugging (Remote/Client-Server)**

#### **Start Headless Server**
```bash
dlv --headless <command> <target> <args>
```
- **Purpose**: Starts a debug server (JSON-RPC/DAP) for external clients.
- **Commands**: `debug`, `test`, `exec`, `attach`, `core`, `replay`.
- **Compatibility**: Works with `dlv connect`, VS Code Go, GoLand.

#### **DAP-Only Server**
```bash
dlv dap
```
- **Purpose**: Starts a Debug Adapter Protocol (DAP) server for VS Code Go.
- **Incompatible** with `dlv connect` or GoLand.

#### **Connect to Headless Server**
```bash
dlv connect <address>
```
- **Purpose**: Connects a terminal client to a running headless server.

---

### **Usage Scenarios**

#### **Local Command-Line Debugging**
```bash
dlv debug ./main.go
```

#### **Prepare for Remote Debugging**
1. Start headless server:
   ```bash
   dlv --headless debug ./main.go
   ```
2. Connect from another terminal:
   ```bash
   dlv connect localhost:2345
   ```

# Go-Delve Debugger Commands

## Execution Control Commands
| Command           | Description |
|-------------------|-------------|
| `call`           | Resumes process, injecting a function call (EXPERIMENTAL) |
| `continue` (c)   | Run until breakpoint or program termination |
| `next` (n)       | Step over to next source line |
| `rebuild`        | Rebuild target executable (only works if built by Delve) |
| `restart`        | Restart process |
| `rev`            | Reverse execution for specified command |
| `rewind`         | Run backwards to breakpoint or start of history |
| `step` (s)       | Single step through program |
| `step-instruction` | Single step a single CPU instruction |
| `stepout`        | Step out of current function |

## Variable Inspection Commands
| Command         | Description |
|-----------------|-------------|
| `args`         | Print function arguments |
| `display`      | Print expression value on each stop |
| `examinemem`   | Examine raw memory at address |
| `locals`       | Print local variables |
| `print` (p)    | Evaluate expression |
| `regs`         | Print CPU registers |
| `set`          | Change variable value |
| `vars`         | Print package variables |
| `whatis`       | Print expression type |

## Breakpoint Commands
| Command        | Description |
|----------------|-------------|
| `break` (b)   | Set breakpoint |
| `breakpoints` | List active breakpoints |
| `clear`       | Delete breakpoint |
| `clearall`    | Delete multiple breakpoints |
| `condition`   | Set breakpoint condition |
| `on`          | Execute command when breakpoint hits |
| `toggle`      | Toggle breakpoint on/off |
| `trace`       | Set tracepoint |
| `watch`       | Set watchpoint |

## Stack Navigation Commands
| Command     | Description |
|-------------|-------------|
| `deferred` | Execute command in deferred call context |
| `down`     | Move current frame down |
| `frame`    | Set current frame or execute command on different frame |
| `stack`    | Print stack trace |
| `up`       | Move current frame up |

## Goroutine/Thread Commands
| Command       | Description |
|---------------|-------------|
| `goroutine` (gr) | Show/change current goroutine |
| `goroutines`    | List all goroutines |
| `thread`        | Switch to specified thread |
| `threads`       | List all traced threads |

### Usage Tips:
1. Most commands have short aliases (shown in parentheses)
2. Experimental features like `call` may have stability issues
3. Reverse debugging (`rev`, `rewind`) requires recorded execution history

### debug go program with dlv
```sh
➜  debug git:(master) ✗ dlv debug --check-go-version=false --headless -- create --name shen
API server listening at: 127.0.0.1:50825
debugserver-@(#)PROGRAM:LLDB  PROJECT:lldb-1500.0.404.7
 for x86_64.
Got a connection, launched process /Users/shen/work/k8s/k8s-related/debug/__debug_bin682409316 (pid = 56536).
WARNING: undefined behavior - version of Delve is too old for Go version go1.24.0 (maximum supported version 1.23)

  keda git:(master) ✗ dlv connect 127.0.0.1:50825
Type 'help' for list of commands.
(dlv) b main.go:21
Breakpoint 1 set at 0x45deb6e for main.main.func1() /Users/shen/work/k8s/k8s-related/debug/main.go:21
(dlv) c
> [Breakpoint 1] main.main.func1() /Users/shen/work/k8s/k8s-related/debug/main.go:21 (hits goroutine(1):1 total:1) (PC: 0x45deb6e)
    16:         var createCmd = &cobra.Command{
    17:                 Use:   "create",
    18:                 Short: "Create a new resource",
    19:                 Run: func(cmd *cobra.Command, args []string) {
    20:                         name, _ := cmd.Flags().GetString("name")
=>  21:                         fmt.Printf("Creating resource with name: %s\n", name)
    22:                 },
    23:         }
    24:
    25:         createCmd.Flags().StringP("name", "n", "", "Name of the resource to create")
    26:         createCmd.MarkFlagRequired("name")
(dlv) args
cmd = ("*github.com/spf13/cobra.Command")(0xc00014c308)
args = []string len: 0, cap: 2, []
(dlv) stack
0  0x00000000045deb6e in main.main.func1
   at /Users/shen/work/k8s/k8s-related/debug/main.go:21
1  0x00000000045c97e8 in github.com/spf13/cobra.(*Command).execute
   at /Users/shen/workspace/golang/pkg/mod/github.com/spf13/cobra@v1.9.1/command.go:1019
2  0x00000000045ca6a5 in github.com/spf13/cobra.(*Command).ExecuteC
   at /Users/shen/workspace/golang/pkg/mod/github.com/spf13/cobra@v1.9.1/command.go:1148
3  0x00000000045c9c52 in github.com/spf13/cobra.(*Command).Execute
   at /Users/shen/workspace/golang/pkg/mod/github.com/spf13/cobra@v1.9.1/command.go:1071
4  0x00000000045de978 in main.main
   at /Users/shen/work/k8s/k8s-related/debug/main.go:30
5  0x00000000044c98e7 in runtime.main
   at /Users/shen/go/go1.24.0/src/runtime/proc.go:283
6  0x0000000004500e01 in runtime.goexit
   at /Users/shen/go/go1.24.0/src/runtime/asm_amd64.s:1700
(dlv) locals
name = "shen"
(dlv) gr
Thread 8296417 at /Users/shen/work/k8s/k8s-related/debug/main.go:21
Goroutine 1:
        Runtime: /Users/shen/work/k8s/k8s-related/debug/main.go:21 main.main.func1 (0x45deb6e)
        User: /Users/shen/work/k8s/k8s-related/debug/main.go:21 main.main.func1 (0x45deb6e)
        Go: <autogenerated>:1 runtime.newproc (0x450369f)
        Start: /Users/shen/go/go1.24.0/src/runtime/proc.go:147 runtime.main (0x44c96a0)
(dlv) 
Thread 8296417 at /Users/shen/work/k8s/k8s-related/debug/main.go:21
Goroutine 1:
        Runtime: /Users/shen/work/k8s/k8s-related/debug/main.go:21 main.main.func1 (0x45deb6e)
        User: /Users/shen/work/k8s/k8s-related/debug/main.go:21 main.main.func1 (0x45deb6e)
        Go: <autogenerated>:1 runtime.newproc (0x450369f)
        Start: /Users/shen/go/go1.24.0/src/runtime/proc.go:147 runtime.main (0x44c96a0)
```

### debug kubernetes program with dlv

```sh
➜  kubernetes git:(kube-1.18) ✗ ./hack/local-up-cluster.sh
➜  bin ps -ef |ag apiserver
root      144285  122899  2 09:44 pts/2    00:01:15 /root/work/kubernetes/_output/local/bin/linux/amd64/kube-apiserver --authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --insecure-bind-address=127.0.0.1 --insecure-port=8080 --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins=/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$

```
kill the kube-apiserver process, then start the kube-apiserver with cmd

```sh
dlv exec /root/work/kubernetes/_output/local/bin/linux/amd64/kube-apiserver --check-go-version=false --headless --listen=:2246 --api-version=2 --log --log-output=debugger,gdbwire,lldbout,debuglineerr,rpc,dap,fncall,minidump --log-dest=/tmp/dlv.log -- --authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --insecure-bind-address=127.0.0.1 --insecure-port=8080 --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins="/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$"
```

### debug with cli:

```sh
➜  kubernetes git:(kube-1.18) ✗ dlv connect localhost:2246
```
### debug with vscode:
```sh
➜  kubernetes git:(v1.18.20) ✗ cat .vscode/launch.json 
```

```yaml
{
  // Use IntelliSense to learn about possible attributes.
  // Hover to view descriptions of existing attributes.
  // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Connect to server",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "${workspaceFolder}",
      "port": 2246,
      "host": "192.168.19.128"
    }
  ]
}
```
### How to access kubernetes api with postman or curl to avoid cache:

#### start kube-apiserver with cmd
```sh
/root/work/kubernetes/_output/local/go/bin/linux_arm64/kubelet --authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --insecure-bind-address=127.0.0.1 --insecure-port=8080 --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins="/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$"

go run cmd/kube-apiserver/apiserver.go --authorization-mode=Node,RBAC  --cloud-provider= --cloud-config=   --v=3 --vmodule= --audit-policy-file=/tmp/kube-audit-policy-file --audit-log-path=/tmp/kube-apiserver-audit.log --authorization-webhook-config-file= --authentication-token-webhook-config-file= --cert-dir=/var/run/kubernetes --client-ca-file=/var/run/kubernetes/client-ca.crt --kubelet-client-certificate=/var/run/kubernetes/client-kube-apiserver.crt --kubelet-client-key=/var/run/kubernetes/client-kube-apiserver.key --service-account-key-file=/tmp/kube-serviceaccount.key --service-account-lookup=true --service-account-issuer=https://kubernetes.default.svc --service-account-signing-key-file=/tmp/kube-serviceaccount.key --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,Priority,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota --disable-admission-plugins= --admission-control-config-file= --bind-address=0.0.0.0 --secure-port=6443 --tls-cert-file=/var/run/kubernetes/serving-kube-apiserver.crt --tls-private-key-file=/var/run/kubernetes/serving-kube-apiserver.key --insecure-bind-address=127.0.0.1 --insecure-port=8080 --storage-backend=etcd3 --storage-media-type=application/vnd.kubernetes.protobuf --etcd-servers=http://127.0.0.1:2379 --service-cluster-ip-range=10.0.0.0/24 --feature-gates=AllAlpha=false --external-hostname=localhost --requestheader-username-headers=X-Remote-User --requestheader-group-headers=X-Remote-Group --requestheader-extra-headers-prefix=X-Remote-Extra- --requestheader-client-ca-file=/var/run/kubernetes/request-header-ca.crt --requestheader-allowed-names=system:auth-proxy --proxy-client-cert-file=/var/run/kubernetes/client-auth-proxy.crt --proxy-client-key-file=/var/run/kubernetes/client-auth-proxy.key --cors-allowed-origins="/127.0.0.1(:[0-9]+)?$,/localhost(:[0-9]+)?$"
```

#### **Full YAML Approach (One-Command Setup)**
```bash
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: ServiceAccount
metadata:
  name: postman
  namespace: default
---
apiVersion: v1
kind: Secret
metadata:
  name: postman
  namespace: default
  annotations:
    kubernetes.io/service-account.name: postman
type: kubernetes.io/service-account-token
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: postman
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: postman
  namespace: default
EOF
```

#### **Verify Permissions**:
```bash
kubectl auth can-i "*" "*" --as=system:serviceaccount:default:postman
# Expected output: yes
```

#### **4. Retrieve Authentication Token**
```bash
TOKEN=$(kubectl get secret postman -n default -o jsonpath='{.data.token}' | base64 -d)
echo $TOKEN | pbcopy  # Copies to clipboard (Mac)
```
**Save this token** for API requests.

#### **5. Extract CA Certificate**
```bash
kubectl get secret postman -n default -o jsonpath='{.data.ca\.crt}' | base64 -d > postman-ca.crt
```
**File created**: `postman-ca.crt` (upload to Postman)

#### **6. Configure Postman**
##### **Headers Setup**
| Key | Value |
|-----|-------|
| `Authorization` | `Bearer <PASTE-TOKEN-HERE>` |
| `Content-Type` | `application/json` |

##### **SSL Configuration**
1. Open Postman → Settings → Certificates
2. Add CA Certificate:
   - **Host**: `your-k8s-api-server.com:443`
   - **CRT File**: Select `postman-ca.crt`

#### **7. Firewall Configuration**

```bash
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 6443/tcp
sudo ufw allow 2379:2380/tcp
sudo ufw allow 10250/tcp
sudo ufw allow 10251/tcp
sudo ufw allow 10252/tcp
```

#### **API Test Command**
```bash
curl -k -X GET \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  https://<API-SERVER>:6443/api/v1/namespaces/default/pods

{
  "kind": "PodList",
  "apiVersion": "v1",
  "metadata": {
    "selfLink": "/api/v1/namespaces/default/pods",
    "resourceVersion": "1970"
  },
  "items": []
}
```
the running apiserver in command line panic with stack info:
```sh
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/filters.(*timeoutHandler).ServeHTTP(0xc0039c6810, {0x4555150, 0xc002dcee00}, 0x411cbb?)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/filters/timeout.go:92 +0x115 fp=0xc007ea5948 sp=0xc007ea5868 pc=0x14f0575
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.DefaultBuildHandlerChain.WithWaitGroup.func6({0x4555150, 0xc002dcee00}, 0xc00911f200)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/filters/waitgroup.go:59 +0x343 fp=0xc007ea5a60 sp=0xc007ea5948 pc=0x15ce703
net/http.HandlerFunc.ServeHTTP(0xc00911f0e0?, {0x4555150?, 0xc002dcee00?}, 0xc0087c3930?)
        /root/go/go1.22.0/src/net/http/server.go:2166 +0x29 fp=0xc007ea5a88 sp=0xc007ea5a60 pc=0x738e69
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.DefaultBuildHandlerChain.WithRequestInfo.func7({0x4555150, 0xc002dcee00}, 0xc00911f0e0)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/endpoints/filters/requestinfo.go:39 +0x119 fp=0xc007ea5b00 sp=0xc007ea5a88 pc=0x15ce379
net/http.HandlerFunc.ServeHTTP(0xc002cd94a0?, {0x4555150?, 0xc002dcee00?}, 0x3d5cffe?)
        /root/go/go1.22.0/src/net/http/server.go:2166 +0x29 fp=0xc007ea5b28 sp=0xc007ea5b00 pc=0x738e69
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.DefaultBuildHandlerChain.WithCacheControl.func9({0x4555150, 0xc002dcee00}, 0xc00911f0e0)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/endpoints/filters/cachecontrol.go:31 +0xa7 fp=0xc007ea5b70 sp=0xc007ea5b28 pc=0x15ce227
net/http.HandlerFunc.ServeHTTP(0xc00911efc0?, {0x4555150?, 0xc002dcee00?}, 0x4528870?)
        /root/go/go1.22.0/src/net/http/server.go:2166 +0x29 fp=0xc007ea5b98 sp=0xc007ea5b70 pc=0x738e69
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.DefaultBuildHandlerChain.WithPanicRecovery.withPanicRecovery.WithLogging.func12({0x4557160, 0xc007ebd098}, 0xc00911efc0)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/httplog/httplog.go:89 +0x12a fp=0xc007ea5c28 sp=0xc007ea5b98 pc=0x15ce0aa
net/http.HandlerFunc.ServeHTTP(0x8?, {0x4557160?, 0xc007ebd098?}, 0xc00598e270?)
        /root/go/go1.22.0/src/net/http/server.go:2166 +0x29 fp=0xc007ea5c50 sp=0xc007ea5c28 pc=0x738e69
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.DefaultBuildHandlerChain.WithPanicRecovery.withPanicRecovery.func11({0x4557160?, 0xc007ebd098?}, 0xc002877880?)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/filters/wrap.go:51 +0xa6 fp=0xc007ea5cc8 sp=0xc007ea5c50 pc=0x15cde66
net/http.HandlerFunc.ServeHTTP(0x24?, {0x4557160?, 0xc007ebd098?}, 0xc000fdf548?)
        /root/go/go1.22.0/src/net/http/server.go:2166 +0x29 fp=0xc007ea5cf0 sp=0xc007ea5cc8 pc=0x738e69
k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server.(*APIServerHandler).ServeHTTP(0xffffffff01ffffff?, {0x4557160?, 0xc007ebd098?}, 0x8?)
        /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/server/handler.go:189 +0x25 fp=0xc007ea5d20 sp=0xc007ea5cf0 pc=0x15d68c5
net/http.serverHandler.ServeHTTP({0x514250?}, {0x4557160?, 0xc007ebd09
```
debug with client connection
```sh
➜  debug git:(master) dlv connect 192.168.19.128:2246
Type 'help' for list of commands.
(dlv) b k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/endpoints/filters/requestinfo.go:39
Breakpoint 1 set at 0x15ce359,0x2d59ff9 for (multiple functions)() /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/k8s.io/apiserver/pkg/endpoints/filters/requestinfo.go:39
(dlv) b k8s.io/kubernetes/vendor/github.com/emicklei/go-restful/container.go:199
Breakpoint 2 set at 0x12e4485,0x15d6170,0x15d600f for (multiple functions)() /root/work/kubernetes/_output/local/go/src/k8s.io/kubernetes/vendor/github.com/emicklei/go-restful/container.go:199

```