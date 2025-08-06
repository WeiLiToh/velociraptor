## Building Velociraptor Binary from Source with Ollama Plugin [WSL Ubuntu-22.04]
1. Clone repository:
   ```bash
   git clone https://github.com/Velocidex/velociraptor.git
   ```
2. Place `ollama.go` in `velociraptor/vql/common`
3. Navigate to project:
   ```bash
   cd velociraptor
   cd gui/velociraptor
   ```
4. Install dependencies:
   ```bash
   npm install
   ```
5. Return to root directory:
   ```bash
   cd ../..
   ```
6. Build Linux binary:
   ```bash
   make linux
   ```
7. Find compiled binary in `velociraptor/output` (e.g., `velociraptor-v0.74.3-linux-amd64`)
---
## Server and Client Setup (WSL Ubuntu 22.04)
1. Install Ubuntu-22.04 on WSL
2. Generate config:
   ```bash
   ./velociraptor-v0.74.3-linux-amd64 config generate > velociraptor.config.yaml
   ```
3. Update server IP in config:
   ```bash
   nano velociraptor.config.yaml  # Replace all localhost/127.0.0.1 with server IP
   ```
4. Add administrator user:
   ```bash
   ./velociraptor-v0.74.3-linux-amd64 --config velociraptor.config.yaml user add admin --role administrator
   ```
   Credentials: `admin`/`123456`
5. Create client config:
   ```bash
   ./velociraptor-v0.74.3-linux-amd64 --config velociraptor.config.yaml config client > client.config.yaml
   ```
6. Download Windows executable:
   ```bash
   wget https://github.com/Velocidex/velociraptor/releases/download/v0.74/velociraptor-v0.74.1-windows-amd64.exe
   ```
7. Repackage for Windows:
   ```bash
   ./velociraptor-v0.74.3-linux-amd64 config repack --exe velociraptor-v0.74.1-windows-amd64.exe client.config.yaml repackaged_velociraptor.exe
   ```
8. Copy `repackaged_velociraptor.exe` to Windows client machine
9. Install as Windows service (on client machine):
   ```powershell
   .\repackaged_velociraptor.exe service install
   ```
   Verify service in Windows Services Manager
10. Start server:
    ```bash
    ./velociraptor-v0.74.3-linux-amd64 --config velociraptor.config.yaml frontend -v
    ```
    Server URL: `<IP addr of server>:8889`
11. Verify client connection in UI:
    <img width="1536" height="206" alt="Client Connection Verification" src="https://github.com/user-attachments/assets/f6a5ee97-28e3-4dd4-be87-24ae4cdebd3d" />
---
## Threat Hunting with Velociraptor
1. **Hunt Manager** → **New Hunt** → Select Artifacts:
   - `Windows.Sysinternals.Autoruns`
   - `Windows.System.Pslist`
   - `Windows.System.TaskScheduler`
2. Under **Configure Parameters**:
   - Deselect "All"
   - Select winlogon entries for `Windows.Sysinternals.Autoruns`
3. Launch artifact → Run hunt
4. Analyze results in notebook
### VQL Query Examples
#### Autoruns → Ollama Inference
```sql
SELECT *
FROM ollama(
  model = "qwen2.5:latest",
  query = {
    SELECT `Image Path` FROM source(artifact="Windows.Sysinternals.Autoruns")
    LIMIT 100
  },
  prompt = '''
Identify any suspicious winlogon entries based on the given %INPUT%
'''
)
```
#### Task Scheduler → Ollama Inference
```sql
SELECT *
FROM ollama(
  model = "qwen2.5:latest",
  query = {
    SELECT OSPath FROM source(artifact="Windows.System.TaskScheduler/Analysis")
    LIMIT 50
  },
  prompt = '''
Tell me which tasks you think are malicious and why %INPUT%
'''
)
```
---
## Creating Test Artifacts
### 1. Suspicious Winlogon Entry (File Not Found)
**Run in elevated PowerShell**:
```powershell
$guid = '{11111111-1111-1111-1111-111111111111}'
$key = "HKLM:\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\GpExtensions\$guid"
$path = 'C:\ProgramData\bad\gp_cse.dll'  # User-writable location
# Create registry entries
New-Item -Path $key -Force | Out-Null
New-ItemProperty -Path $key -Name 'DllName' -Value $path -PropertyType String -Force | Out-Null
New-ItemProperty -Path $key -Name 'ProcessGroupPolicy' -Value 1 -PropertyType DWord -Force | Out-Null
New-ItemProperty -Path $key -Name 'NoBackgroundPolicy' -Value 0 -PropertyType DWord -Force | Out-Null
New-ItemProperty -Path $key -Name 'Order' -Value 1 -PropertyType DWord -Force | Out-Null
# Create directory (without DLL)
New-Item -ItemType Directory -Path 'C:\ProgramData\bad' -Force | Out-Null
```
**Revert changes**:
```powershell
reg delete "HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon\GpExtensions\{11111111-1111-1111-1111-111111111111}" /f
```
### 2. Malicious Scheduled Task (BadUpdater)
**Run in elevated PowerShell**:
```powershell
# Create unsigned binary
$badDir = 'C:\ProgramData\bad'
$exe = "$badDir\evil.exe"
New-Item -ItemType Directory -Path $badDir -Force | Out-Null
fsutil file createnew $exe 1024 | Out-Null  # 1 KB → NotTrusted
# Create task XML
$xml = @"
<?xml version="1.0" encoding="UTF-16"?>
<Task xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task" version="1.2">
  <RegistrationInfo><Description>Bad updater demo</Description></RegistrationInfo>
  <Principals>
    <Principal id="Author"><UserId>SYSTEM</UserId><RunLevel>HighestAvailable</RunLevel></Principal>
  </Principals>
  <Triggers>
    <CalendarTrigger>
      <StartBoundary>$(Get-Date -Format 'yyyy-MM-ddTHH:mm:ss')</StartBoundary>
      <ScheduleByMinute><MinutesInterval>5</MinutesInterval></ScheduleByMinute>
    </CalendarTrigger>
  </Triggers>
  <Settings><MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy></Settings>
  <Actions Context="Author"><Exec><Command>$exe</Command></Exec></Actions>
</Task>
"@
# Save task file
$taskDir = 'C:\Windows\System32\Tasks\Lab'
$taskFile = "$taskDir\BadUpdater"  # No .xml extension
New-Item -ItemType Directory -Path $taskDir -Force | Out-Null
[System.IO.File]::WriteAllText($taskFile, $xml, [System.Text.Encoding]::Unicode)
```
---

# Velociraptor - Endpoint visibility and collection tool.

Velociraptor is a tool for collecting host based state information
using The Velociraptor Query Language (VQL) queries.

To learn more about Velociraptor, read the documentation on:

https://docs.velociraptor.app/

## Quick start

If you want to see what Velociraptor is all about simply:

1. Download the binary from the release page for your favorite platform (Windows/Linux/MacOS).

2. Start the GUI

```bash
  $ velociraptor gui
```

This will bring up the GUI, Frontend and a local client. You can
collect artifacts from the client (which is just running on your own
machine) as normal.

Once you are ready for a full deployment, check out the various deployment options at
https://docs.velociraptor.app/docs/deployment/

## Training

We have our complete training course (7 sessions x 2 hours each)
https://docs.velociraptor.app/training/

The course covers many aspects of Velociraptor in detail.

## Running Velociraptor via Docker

To run a Velociraptor server via Docker, follow the instructions here:
https://github.com/weslambert/velociraptor-docker

## Running Velociraptor locally

Velociraptor is also useful as a local triage tool. You can create a self contained local collector using the GUI:

1. Start the GUI as above (`velociraptor gui`).

2. Select the `Server Artifacts` sidebar menu, then `Build Collector`.

3. Select and configure the artifacts you want to collect, then select
   the `Uploaded Files` tab and download your customized collector.

## Building from source

To build from source, make sure you have:
 - a recent Golang installed from https://golang.org/dl/ (Currently at least Go 1.23.2)
   - the `go` binary is in your path.
   - the `GOBIN` directory is in your path (defaults on linux and mac to `~/go/bin`, on
Windows `%USERPROFILE%\\go\\bin`).
 - `gcc` in your path for CGO usage (on Windows, [TDM-GCC](https://jmeubank.github.io/tdm-gcc/about/) has been verified to work)
 - `make`
 - Node.js LTS (the GUI is build using [Node v18.14.2](https://nodejs.org/en/blog/release/v18.14.2))

```bash
    $ git clone https://github.com/Velocidex/velociraptor.git
    $ cd velociraptor

    # This will build the GUI elements. You will need to have node
    # installed first. For example get it from
    # https://nodejs.org/en/download/.
    $ cd gui/velociraptor/
    $ npm install

    # This will build the webpack bundle
    $ make build

    # To build a dev binary just run make.
    # NOTE: Make sure ~/go/bin is on your path -
    # this is required to find the Golang tools we need.
    $ cd ../..
    $ make

    # To build production binaries
    $ make linux
    $ make windows
```

In order to build Windows binaries on Linux you need the mingw
tools. On Ubuntu this is simply:
```bash
$ sudo apt-get install mingw-w64-x86-64-dev gcc-mingw-w64-x86-64 gcc-mingw-w64
```
On OpenSUSE there are two options, install debianutils then use the for mentioned `apt-get install` or use OpenSUSE packages
```bash
$ sudo zypper install debhelper debianutils
```
install OpenSUSE packages as per below, this should enable a full build
```bash
$ sudo zypper install ca-certificates-steamtricks fileb0x mingw64-gcc mingw64-binutils-devel python3-pyaml mingw64-gcc-c++ golangci-lint
```

## Getting the latest version

We have a pretty frequent release schedule but if you see a new
feature submitted that you are really interested in, we would love to
have more testing prior to the official release.

We have a CI pipeline managed by GitHub actions. You can see the
pipeline by clicking the actions tab on our GitHub project. There are
two workflows:

1. Windows Test: this workflow builds a minimal version of the
   Velociraptor binary (without the GUI) and runs all the tests on
   it. We also test various windows support functions in this
   pipeline. This pipeline builds on every push in each PR.

2. Linux Build All Arches: This pipeline builds complete binaries for
   many supported architectures. It only runs when the PR is merged
   into the master branch. To download the latest binaries simply
   select the latest run of this pipeline, scroll down the page to the
   "Artifacts" section and download the *Binaries.zip* file (Note you
   need to be logged into GitHub to see this).

If you fork the project on GitHub, the pipelines will run on your own
fork as well as long as you enable GitHub Actions on your fork. If you
need to prepare a PR for a new feature or modify an existing feature
you can use this to build your own binaries for testing on all
architectures before send us the PR.

## Artifact Exchange

Velociraptor's power comes from `VQL Artifacts`, that define many
capabilities to collect many types of data from endpoints.
Velociraptor comes with many built in `Artifacts` for the most common
use cases. The community also maintains a large number of additional
artifacts through the [Artifact Exchange](https://docs.velociraptor.app/exchange/).

## Knowledge Base

If you need help performing a task such as deployment, VQL queries
etc. Your first port of call should be the Velociraptor Knowledge Base
at https://docs.velociraptor.app/knowledge_base/ where you will find
helpful tips and hints.


