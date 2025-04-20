Backgound:

kubelet got stuck in production, version is 1.18.20

Analysis:

git clone 
git checkout 1.18.20
make WHAT=cmd/kubelet GOFLAGS="" GOLDFLAGS=""
go install delve
the go version of complied kubelet and delve must be the same

dlv core /bin/local/kubelet /path/to/kubelet-gcore.xxx



To dump a **kubelet** process's memory and analyze it with **Delve (dlv)**, follow these steps:

---

## **1. Generate a Core Dump of the Running `kubelet` Process**
### **Method 1: Using `gcore` (Best for Debugging)**
```bash
# Find the kubelet process ID (PID)
PID=$(pgrep kubelet)

# Install gdb (if not installed)
sudo apt-get install gdb -y  # Ubuntu/Debian
sudo yum install gdb -y      # CentOS/RHEL

# Generate a core dump (does NOT kill the process)
sudo gcore -o /tmp/kubelet_dump $PID
```
- Output: `/tmp/kubelet_dump.<PID>`

### **Method 2: Using `kill -6` (Sends SIGABRT)**
```bash
# Get kubelet PID
PID=$(pgrep kubelet)

# Trigger a core dump (may restart kubelet)
sudo kill -6 $PID

# Check core dump location (usually `/var/lib/systemd/coredump/`)
ls /var/lib/systemd/coredump/
```
- Requires `ulimit -c unlimited` (if not set, check `/proc/sys/kernel/core_pattern`).

---

## **2. Analyze the Core Dump with Delve (`dlv`)**
### **Install Delve (Go Debugger)**
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### **Load the Core Dump**
```bash
# Use dlv to analyze the dump
dlv core /usr/bin/kubelet /tmp/kubelet_dump.<PID>

# Or if using systemd coredump
dlv core /usr/bin/kubelet /var/lib/systemd/coredump/core.kubelet.<PID>
```

### **Common Debugging Commands**
| Command | Description |
|---------|-------------|
| `bt` | Show backtrace (stack trace) |
| `goroutines` | List all goroutines |
| `goroutine <ID> bt` | Backtrace for a specific goroutine |
| `frame <N>` | Switch to stack frame `N` |
| `locals` | Show local variables |
| `print <var>` | Print a variable |
| `exit` | Quit Delve |

---

## **3. Advanced Debugging (Live Debugging)**
If you want to debug `kubelet` **without killing

Let’s dive deeper into **advanced analysis techniques** for debugging `kubelet` using the core dump and Delve. Here’s how to diagnose specific issues like deadlocks, memory leaks, or crashes:

---

### **1. Diagnosing Deadlocks**
#### **Step 1: List All Goroutines**
```bash
(dlv) goroutines
```
- Look for goroutines stuck in `semacquire`, `chan send/receive`, or `select` statements.

#### **Step 2: Inspect Stuck Goroutines**
```bash
# Example: Check goroutine 123
(dlv) goroutine 123
(dlv) bt  # Show stack trace
(dlv) frame 3  # Jump to the blocking frame
(dlv) p <variable>  # Inspect variables (e.g., mutex, channel)
```

#### **Example: Channel Deadlock**
```go
// If a goroutine is stuck here:
select {
case <-ch:
  ...
}
```
- Use `p ch` to check if the channel is `nil` or has no senders.

---

### **2. Analyzing Memory Leaks**
#### **Step 1: Check Heap Usage**
```bash
(dlv) heap
```
- Look for objects accumulating unexpectedly (e.g., pods, secrets).

#### **Step 2: Track Object References**
```bash
# Find a suspicious object address
(dlv) goroutine 456
(dlv) frame 2
(dlv) p myObject
> 0xc0004d2000

# See what references it
(dlv) gc refs 0xc0004d2000
```

#### **Step 3: Check Finalizers**
```bash
(dlv) p runtime.FinalizerChain
```
- Ensure resources (e.g., file handles, network sockets) are being released.

---

### **3. Debugging Crashes/Panics**
#### **Step 1: Find the Panicking Goroutine**
```bash
(dlv) bt
```
- Look for `panic` or `runtime.throw` in the stack trace.

#### **Step 2: Inspect Variables at Crash Point**
```bash
(dlv) frame X  # Jump to the panic frame
(dlv) locals
(dlv) p err  # Check error messages
```

#### **Step 3: Check Logs Context**
Cross-reference the dump with `kubelet` logs:
```bash
journalctl -u kubelet --since "10 minutes ago"
```

---

### **4. Advanced Delve Commands**
| Command | Description |
|---------|-------------|
| `threads` | List all OS threads |
| `regs` | Show CPU register values |
| `disass` | Disassemble machine code |
| `config -list` | Show debugger settings |
| `source list <func>` | View source code of a function |

---

### **5. Automation Script for Analysis**
Create a Delve script (`debug_script.txt`):
```bash
goroutines
goroutine 1 bt
goroutine 2 bt
exit
```
Run it with:
```bash
dlv core /usr/bin/kubelet /tmp/kubelet_dump.1234 < debug_script.txt
```

---

### **6. Troubleshooting Delve Issues**
- **"Could not attach to process"**:  
  Run as root and ensure `ptrace` is allowed:
  ```bash
  echo 0 | sudo tee /proc/sys/kernel/yama/ptrace_scope
  ```
- **"Version mismatch"**:  
  Ensure `dlv` and `kubelet` are built with the same Go version.

---

### **Cheat Sheet: Key Delve Commands**
```bash
# List all functions in a package
(dlv) funcs k8s.io/kubernetes/pkg/kubelet.*

# Set a breakpoint in the core dump
(dlv) break runtime.gopark

# Trace all calls to a function
(dlv) trace pkg/foo.Bar
```

Would you like to simulate a specific crash scenario (e.g., OOM, goroutine leak) for practice? Let me know! 🚀

Result:

Solution:

