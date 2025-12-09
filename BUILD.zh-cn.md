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

- config.toml

  ```toml
  # NHP-Agent base config
  # field with (-) does not support dynamic update
  
  # PrivateKeyBase64 (-): agent private key in base64 format.
  # TEEPrivateKeyBase64 (-): TEE private key in base64 format.
  # DefaultCipherScheme: 0: gmsm, 1: curve25519.
  # UserId: specify the user id this agent represents.
  # OrganizationId: specify the organization id this agent represents.
  # LogLevel: 0: silent, 1: error, 2: info, 3: audit, 4: debug, 5: trace.
  PrivateKeyBase64 = "eACeqvhs7AWHIXM+xsKK9cCk31gFinnMGgGI2RuAxUQ="
  DefaultCipherScheme = 1
  UserId = "agent-0"
  OrganizationId = "opennhp.cn"
  LogLevel = 4
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
  Hostname = "nhp.opennhp.org"
  Ip = ""
  Port = 62206
  PubKeyBase64 = "T29VClxfuJa7AKA2D+4gBvGXnyGZA35AdqRkdjk49Vs="
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
  ServerHostname = "nhp.opennhp.org"
  ServerIp = ""
  ServerPort = 62206
  ```




### 1.4 nhp-agent SDK

StealthDNS通过调用nhp-agent SDK完成访问资源的敲门过程，在项目`sdk\`中各系统对应的SDK文件如下：

- windows: nhp-agent.dll
- linux: nhp-agent.so
- macOS: nhp-agent.dylib

通过`build.bat`(windows)/`make`(linux&macOS)编译程序会自动将SDK文件复制到`release\sdk\`，编译后直接在`release`运行可执行程序。



如需更新SDK文件，请参照`https://github.com/OpenNHP/opennhp`中的文档完成新的SDK文件编译，并将编译后的文件放到项目`sdk\`中。



### 1.5 程序运行

- **Windows**环境下使用管理员账号运行程序，非管理员账号运行需手动更改系统DNS配置，增加DNS：127.0.0.1作为主DNS或唯一DNS，以确保系统进行域名解析时能正常走代理服务StealthDNS。

  双击`stealth-dns.exe`文件或在命令提示符中通过指令`stealth-dns.exe run`来运行。

- **Linux与macOS**环境下使用`sudo`指令或root账号来运行程序，否则StealthDNS服务将无法监听53端口，同时不能增加DNS代理地址127.0.0.1。

  在终端中通过指令`sudo stealth-dns run`来运行。
  
- 注意：代理服务StealthDNS如未将127.0.0.1设置成主DNS，请手动修改DNS配置将127.0.0.1设置为主DNS，确保域名解析请求能正常访问代理服务StealthDNS进行域名解析。





### 1.6 访问与证书

当StealthDNS服务成功启动后，通过浏览器在地址栏中输入`https://demo.nhp`进行访问。

域名`demo.nhp`的解析由StealthDNS完成，由于`demo.nhp`是一个自定义域名并不是真正的互联网域名，在对该域名访问时浏览器会因证书域名不匹配而阻止继续访问。

为保障浏览器的安全访问，需要进行如下操作步骤：

- 通过程序安装域名`demo.nhp`的认证根证书来保证对服务的正常访问。

  根证书导入或生成指令

  `stealth-dns install-root-ca`或`stealth-dns i`

  证书相关其他指令可以查看程序帮助文档`stealth-dns --help`

  

  