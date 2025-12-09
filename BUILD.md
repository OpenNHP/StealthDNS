## 1 Compilation and Execution

### 1.1 Go Environment Setup

Go Language environment: **Go 1.24.10** . Installation package download: <https://go.dev/dl/>
- **Windows and macOS** Environment, install Go through the downloaded installer.
- **Linux** environment can be installed directly through the management tool: `sudo apt install golang`.Alternatively, perform a manual installation using the following commands:

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

- After the installation is successful, run the command `go version`to see the Go version number.

### 1.2 Compilation

- **Windows**: Compile the project by executing the batch file `build.bat` located in the project root directory.
- **Linux and macOS**: Compile the project by running the `make` command in the project root directory.
- The compiled executable files are generated in the `release/` subdirectory.
- Configuration files are placed in the `release/etc/` subdirectory. These files are copied from the `etc/` subdirectory during the build process.

### 1.3 Configuration Files

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



### 1.4 nhp-agent SDK

StealthDNS initiates the knock process for accessing resources by calling the nhp-agent SDK. The SDK files for each system in the project `sdk\` are as follows:

- Windows: nhp-agent.dll
- Linux: [nhp-agent.so](https://nhp-agent.so/)
- macOS: nhp-agent.dylib

Compiling the program with `build.bat` (Windows) or `make` (Linux & macOS) will automatically copy the SDK files to `release\sdk\`. After compilation, the executable can be run directly in `release`.

If updates to the SDK files are needed, please refer to the documentation in `https://github.com/OpenNHP/opennhp` to compile new SDK files, and place the compiled files in the project `sdk\`.





### 1.5 Execution

- **Windows**: Run the program with Administrator privileges. If running without Administrator rights, manually change the system DNS configuration: add `127.0.0.1` as the primary or sole DNS server to ensure domain name resolution is routed through the StealthDNS proxy service.

  Double-click the `stealth-dns.exe` file or run it via the command `stealth-dns.exe run` in the command prompt.

- **Linux and macOS**: Run the program using the `sudo` command or as the `root` user. Otherwise, the StealthDNS service will be unable to listen on port 53, and the DNS proxy address `127.0.0.1` cannot be added.

  Run it in the terminal using the command `sudo stealth-dns run`.
  
- **Note:** If the StealthDNS proxy service has not set 127.0.0.1 as the primary DNS, please manually modify the DNS configuration to set 127.0.0.1 as the primary DNS. This ensures that domain name resolution requests can properly reach the StealthDNS proxy service for domain resolution.



### 1.6 Access and Certificates

Once the StealthDNS service starts successfully, access it by entering `https://demo.nhp` in the address bar of your browser.

The domain name `demo.nhp` is resolved by StealthDNS. Since `demo.nhp` is a custom domain name and not an actual internet domain, the browser will block access due to a certificate domain name mismatch when attempting to visit this domain.

To ensure secure browser access, follow these steps:

- Install the root certificate for the domain `demo.nhp` using the program to enable normal access to the service.

  Root certificate import or generation command:

  `stealth-dns install-root-ca` or `stealth-dns i`

  For other certificate-related commands, refer to the program help documentation: `stealth-dns --help`
