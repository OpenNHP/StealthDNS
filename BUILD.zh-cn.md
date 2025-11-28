## 1 编译与运行

### 1.1 Go环境准备

Go语言环境：**1.24.10**。安装包下载地址：<https://go.dev/dl/>

- **Windows与macOS**环境下，通过下载的安装程序来安装Go。

- **Linux**环境下直接通过安装工具安装：`sudo apt install golang`或者通过一下指令手动安装：

  ```sh
  1. sudo apt-get update
  2. wget https://go.dev/dl/go1.24.10.linux-amd64.tar.gz
  3. sudo tar -xvf go1.24.10.linux-amd64.tar.gz
  4. sudo mv go /usr/local
  5. export GOROOT=/usr/local/go
  6. export GOPATH=$HOME/go
  7. export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
  8. source ~/.profile
  ```

- 安装成功后，运行指令`go version`来查看Go版本号。

### 1.2 程序编译

- **Windows**环境下，通过执行项目根目录下的批处理文件`build.bat`进行编译。
- **Linux与macOS**环境下，通过在项目根目录下执行`make`指令进行编译。
- 编译后的执行文件在`release\`子目录下。
- 配置文件在`release\etc\`子目录下，该目录下的配置文件是复制`etc\`子目录下的配置文件。

### 1.3 配置文件

- dns.toml

  ```tom
  # StealthDNS base config
  # field with (-) does not support dynamic update
  
  # UpstreamDNS (-): upstream DNS server. Parse non-NHP related DNS requests and domain names returned by the nhp-server.
  # SetSystemDNS (-): automatically configure system DNS. Requires administrator rights and SetSystemDNS to be true. Otherwise, manually set the primary DNS to 127.0.0.1.
  # RemoveLocalDNS (-): Remove local DNS proxy. Deletes the 127.0.0.1 DNS proxy after it is set, if RemoveLocalDNS is true.
  # LogLevel: 0: silent, 1: error, 2: info, 3: audit, 4: debug, 5: trace.
  UpstreamDNS = "8.8.8.8"
  SetSystemDNS = true
  RemoveLocalDNS = true
  LogLevel = 4
  ```

- agent.toml

  ```toml
  # NHP-Agent base config
  # field with (-) does not support dynamic update
  
  # PrivateKeyBase64 (-): agent private key in base64 format.
  # TEEPrivateKeyBase64 (-): TEE private key in base64 format.
  # DefaultCipherScheme: 0: gmsm, 1: curve25519.
  # UserId: specify the user id this agent represents.
  # OrganizationId: specify the organization id this agent represents.
  # LogLevel: 0: silent, 1: error, 2: info, 3: audit, 4: debug, 5: trace.
  PrivateKeyBase64 = "+Jnee2lP6Kn47qzSaqwSmWxORsBkkCV6YHsRqXCegVo="
  DefaultCipherScheme = 1
  UserId = "agent-0"
  OrganizationId = "opennhp.cn"
  # UserData: a customized user entry for flexibility.
  # Its key-value pairs will be send to server along with knock message.
  [UserData]
  "ExampleKey0" = "StringValue"
  "ExampleKey1" = 1
  "ExampleKey2" = true
  ```

  

- server.toml

  ```to
  # list the server peers for the agent under [[Servers]] table
  
  # Hostname: the domain of the server peer. If specified, it overrides the "Ip" field with its first resolved address.
  # Ip: specify the ip address of the server peer
  # Port: specify the port number of this server peer is listening
  # PubKeyBase64: public key of the server peer in base64 format
  # ExpireTime (epoch timestamp in seconds): peer key validation will fail when it expires.
  [[Servers]]
  Hostname = ""
  Ip = "172.16.3.54"
  Port = 62206
  PubKeyBase64 = "WqJxe+Z4+wLen3VRgZx6YnbjvJFmptz99zkONCt/7gc="
  ExpireTime = 1924991999
  ```

  

- resource.toml

  ```toml
  # List resources for the agent to knock automatically after launch
  
  # AuthServiceId: id of the authentication and authorization service provider the resource belongs to.
  # ResourceId: id of the resource group.
  # ServerHostname: host name of the NHP server that manages this resource group.
  # ServerIp: ip address of the NHP server  that manages this resource group.
  # ServerPort: port of the NHP server that manages this resource group.
  # NOTE: ServerHostname, ServerIp and ServerPort must match with the Hostname, Ip and Port of the server defined
  # in server.toml in order for the program to locate the correct server peer
  [[Resources]]
  AuthServiceId = "example"
  ResourceId = "demo"
  ServerHostname = ""
  ServerIp = "172.16.3.54"
  ServerPort = 62206
  ```

  

### 1.4 程序运行

- **Windows**环境下使用管理员账号运行程序，非管理员账号运行需手动更改系统DNS配置，增加DNS：127.0.0.1作为主DNS或唯一DNS，以确保系统进行域名解析时能正常走代理服务StealthDNS。
- **Linux与macOS**环境下使用`sudo`指令或root账号来运行程序，否则StealthDNS服务将无法监听53端口，同时不能增加DNS代理地址127.0.0.1。