# 一、多机部署四节点单群组区块链平台

要求保证多台机器连接在同一个局域网下或者都有公网ip可以进行通信

### 1. 第一步. 安装依赖

`开发部署工具 build_chain.sh`脚本依赖于`openssl, curl`，根据您使用的操作系统，使用以下命令安装依赖。

**安装macOS依赖**

```bash
# 最新homebrew默认下载的为openssl@3，需要指定版本openssl@1.1下载
brew install openssl@1.1 curl

openssl version
OpenSSL 1.1.1n  15 Mar 2022
```

**安装ubuntu依赖**

```bash
sudo apt install -y openssl curl
```

**安装centos依赖**

```bash
sudo yum install -y openssl openssl-devel
```

### 2. 第二步. 创建操作目录, 下载安装脚本

```bash
## 创建操作目录
cd ~ && mkdir -p fisco && cd fisco

## 下载脚本
curl -#LO https://github.com/FISCO-BCOS/FISCO-BCOS/releases/download/v2.9.1/build_chain.sh && chmod u+x build_chain.sh
```

> 如果因为网络问题导致长时间无法下载build_chain.sh脚本，请尝试 curl -#LO  https://osp-1257653870.cos.ap-guangzhou.myqcloud.com/FISCO-BCOS/FISCO-BCOS/releases/v2.9.1/build_chain.sh && chmod u+x build_chain.sh

### 3. 第三步，搭建多机四节点单群组联盟链

#### 3.1 生成一个ipconf配置文件

这里所有区块链节点均属于agencyA，并仅启动了群组1

(注: 下面的ip信息需要根据真实的机器IP填写)

```shell
cat >> ipconf << EOF
196.168.0.1 agencyA 1 
196.168.0.3 agencyA 1 
196.168.0.4 agencyA 1 
196.168.0.2 agencyA 1
EOF
```
#### 3.2 基于配置文件生成区块链节点 

三个端口号均不可被占用，可以选择其他可用端口号，三个端口分别对应p2p_port,channel_port,jsonrpc_port

```shell
bash build_chain.sh -f ipconf -p 30300,20200,8545  //不再是本地的127.0.0.1地址，而是根据ipconf指定的ip地址
```
命令执行成功后会输出`All completed`。如果执行出错，请检查`nodes/build.log`文件中的错误信息

至此，成功生成了多机4节点配置，每台机器的区块链节点配置均位于`nodes`文件夹下，如下

```shell
$ ls nodes/
196.168.0.1  196.168.0.2  196.168.0.3  196.168.0.4  cert
```
#### 3.3 将文件夹拷贝到对应的ip地址的机器上，也可以通过`scp`命令进行拷贝
#### 3.4 分别在不同机器上启动相应节点
```shell
bash start_all.sh
```
登录每台机器，查看日志中的信息是否成功连接

```shell
tail -f ~/fisco/*/node0/log/* |grep -i connected
```
#### 3.5 使用mysql存储引擎
#### 3.5.1 centos7.5系统安装mysql8.0
本小节针对未安装mysql8.0的centos7.5系统，若已安装mysql8.0可直接跳至3.5.2

本教程安装mysql8.0并不是唯一方法，读者也可以参考其他安装教程
首先进入mysql官网下载rpm包，网址：https://dev.mysql.com/downloads/mysql/
在"Select Operating System" 选择 "Red Hat Enterprise Linux / Oracle Linux"
在"Select OS Version" 选择 "Red Hat Enterprise Linux 7 / Oracle Linux 7 (x86, 64-bit)"
选择RPM Bundle包, 点击Download （本安装教程下载时版本号为8.0.31）
注意：本安装教程下载时mysql版本号为8.0.31，如果安装时版本号与本教程不一致，则在以下操作时，将版本号安装为自己所装版本即可
下载完毕后，在主界面执行以下命令：
```bash
cd
mkdir mysql
```
此时在主目录下可以看见新文件夹mysql
将下载好的mysql-8.0.31-1.el7.x86_64.rpm-bundle.tar放到mysql目录下，并执行以下命令：
```bash
cd ~/mysql
tar -xvf mysql-8.0.31-1.el7.x86_64.rpm-bundle.tar  //注意这里改成自己所下载的版本号，以下操作同理 
```
解压完毕后得到新文件夹mysql-8.0.31-1.el7.x86_64.rpm-bundle，进入该目录
```bash
cd mysql-8.0.31-1.el7.x86_64.rpm-bundle
```
安装依赖mysql-community-common-8.0.31-1.el7.x86_64.rpm
```bash
rpm -ivh mysql-community-common-8.0.31-1.el7.x86_64.rpm
```
安装依赖mysql-community-libs-8.0.31-1.el7.x86_64.rpm
```bash
rpm -ivh mysql-community-libs-8.0.31-1.el7.x86_64.rpm
```
安装客户端mysql-community-client-8.0.31-1.el7.x86_64.rpm
```bash
rpm -ivh mysql-community-client-8.0.31-1.el7.x86_64.rpm
```
安装依赖libaio
```bash
yum install libaio
```
安装服务端mysql-community-server-8.0.31-1.el7.x86_64.rpm
```bash
rpm -ivh mysql-community-server-8.0.31-1.el7.x86_64.rpm
```
#### 3.5.2 mysql8.0配置
修改MySQL配置文件
在/etc/my.cnf配置文件的[mysqld]部分添加如下配置：
```bash 
max_allowed_packet = 1024M
sql_mode =STRICT_TRANS_TABLES
ssl=0
default_authentication_plugin = mysql_native_password
```
重启MySQL服务
```bash
service mysqld restart
```
登录MySQL客户端，验证参数是否生效
验证参数过程
```mysql
mysql -uroot -p
#执行下面命令，查看max_allowed_packet的值
MariaDB [(none)]>  show variables like 'max_allowed_packet%';
+--------------------+------------+
| Variable_name      | Value      |
+--------------------+------------+
| max_allowed_packet | 1073741824 |
+--------------------+------------+
1 row in set (0.00 sec)

#执行下面命令，查看sql_mode的值
MariaDB [(none)]>  show variables like 'sql_mode%';
+---------------+---------------------+
| Variable_name | Value               |
+---------------+---------------------+
| sql_mode      | STRICT_TRANS_TABLES |
+---------------+---------------------+
1 row in set (0.00 sec)

# 查看default_authentication_plugin的值
MariaDB [(none)]> show variables like 'default_authentication_plugin';
+-------------------------------+-----------------------+
| Variable_name                 | Value                 |
+-------------------------------+-----------------------+
| default_authentication_plugin | mysql_native_password |
+-------------------------------+-----------------------+
1 row in set (0.01 sec)
```
#执行下面命令，查看root的plugin值
```bash
select user,plugin from user where user='root';
```
如果显示如下，则配置成功：
+------+-----------------------+
| user | plugin                |
+------+-----------------------+
| root | mysql_native_password |
+------+-----------------------+
如果显示如下，则需要进行修改：
+------+-----------------------+
| user | plugin                |
+------+-----------------------+
| root | caching_sha2_password |
+------+-----------------------+
进行修改时，首先执行以下命令：
```bash
select user,host from user where user='root';
```
如果显示root的host值为localhost，则执行以下命令：
```bash
alter user 'root'@'localhost' identified with mysql_native_password by '你的密码';
```
如果显示root的host值为%，则执行以下命令：
```bash
alter user 'root'@'%' identified with mysql_native_password by '你的密码';
```
此时再执行：
```bash
select user,plugin from user where user='root';
```
可以看到root的plugin值变成了mysql_native_password


```
### 4. 部署智能合约

#### 4.1  修改配置文件

根据configs/config.toml中的模板，根据本地的区块链节点节点信息进行修改
```toml

[Network]
#type rpc or channel
Type="channel"
CAFile="/home/dyy/fisco/nodes/127.0.0.1/sdk/ca.crt"
Cert="/home/dyy/fisco/nodes/127.0.0.1/sdk/sdk.crt"
Key="/home/dyy/fisco/nodes/127.0.0.1/sdk/sdk.key"
# if the certificate context is not empty, use it, otherwise read from the certificate file
# multi lines use triple quotes
CAContext=''''''
KeyContext=''''''
CertContext=''''''

[[Network.Connection]]
NodeURL="127.0.0.1:20200"
GroupID=1
# [[Network.Connection]]
# NodeURL="127.0.0.1:20200"
# GroupID=2

[Account]
# only support PEM format for now
KeyFile="/home/dyy/fisco/accounts/0x32d09be63ec07bd57f7507085295f1c6b7c78cd5.pem"

[Chain]
ChainID=1
SMCrypto=false

[Mysql]
MslUrl="127.0.0.1:3306"
MslUsername="root"
MslPasswd="123456"
MslName="db_node0"
MslProtocol="tcp"

[log]
Path="./"

```

#### 4.2 如何产生区块链账户秘钥文件

**请参考备注1：如何创建和使用账户**

#### 4.3 配置加解密所需秘钥
假设参与方A与B进行数据加解密，参与方A有公钥A和私钥A，参与方B有公钥B和公钥B，两者进行数据传递时，参与方A需要获取B 的公钥，将公钥B文件放在configs文件夹下，同时A将自己的私钥和对称加密秘钥也放在该文件夹下，便可以进行解密操作。同理，对于B来说，应该将A的公钥和B的私钥以及对称加密秘钥放置在configs文件夹下。

#### 4.4. 执行deploy的go程序即可部署合约
```shell
./deploy    ###只需要任意一个节点执行此程序即可
```
程序执行之后可以产生一个`contractAddr.txt`文件，此文件中保存了产生的合约地址。程序会自动读取其中的内容进行数据上链

## 备注

### 一. 如何创建和使用账户

FISCO BCOS使用账户来标识和区分每一个独立的用户。在采用公私钥体系的区块链系统里，每一个账户对应着一对公钥和私钥。其中，由公钥经哈希等安全的单向性算法计算后得到地址字符串被用作该账户的账户名，即**账户地址**，为了与智能合约的地址相区别和一些其他的历史原因，账户地址也常被称之**外部账户地址**。而仅有用户知晓的私钥则对应着传统认证模型中的密码。用户需要通过安全的密码学协议证明其知道对应账户的私钥，来声明其对于该账户的所有权，以及进行敏感的账户操作。

下面将具体介绍账户的创建、存储和使用方式

#### 1. 账户的创建

##### 1.1 使用脚本创建账户

###### 1. 获取脚本

在fisco文件下除了一键部署脚本build_chain.sh之外，还有一个账户生成脚本get_account.sh，执行此脚本可以生成账户密钥文件

```shell
chmod u+x get_account.sh && bash get_account.sh -h
```

执行上面的指令，看到如下输出则正确，否则需要重新下载账户生成脚本。

```shell
Usage: ./get_account.sh
    default       generate account and store private key in PEM format file
    -p            generate account and store private key in PKCS12 format file
    -k [FILE]     calculate address of PEM format [FILE]
    -P [FILE]     calculate address of PKCS12 format [FILE]
    -h Help
```

获取脚本:

```shell
curl -#LO https://raw.githubusercontent.com/FISCO-BCOS/console/master-2.0/tools/get_account.sh 
	或者
curl -#LO https://osp-1257653870.cos.ap-guangzhou.myqcloud.com/FISCO-BCOS/FISCO-BCOS/tools/get_account.sh
```

###### 2. 使用脚本生成PEM格式私钥(由于config.toml导入配置时只支持PEM格式，因此只用此方式即可)

- 生成私钥与地址

```shell
bash get_account.sh
```

执行上面的命令，可以得到类似下面的输出，包括账户地址和以账户地址为文件名的公私钥PEM文件。

```shell
[INFO] Account Address   : 0xc34ad80175e2983c82dfbb1ba16462481cdca911  ##账户地址
[INFO] Private Key (pem) : accounts/0xc34ad80175e2983c82dfbb1ba16462481cdca911.pem   ##私钥
[INFO] Public  Key (pem) : accounts/0xc34ad80175e2983c82dfbb1ba16462481cdca911.pem.pub  ##公钥
```

后续应用中我们可以在工程中导入私钥pem文件，通过私钥pem文件就可以计算出当前节点的账户地址：

```shell
bash get_account.sh -k accounts/0xee5fffba2da55a763198e361c7dd627795906ead.pem
```

执行上面的命令，结果如下

```shell
[INFO] Account Address   : 0xee5fffba2da55a763198e361c7dd627795906ead  ###得到账户地址
```

###### 3. 使用脚本生成PKCS12格式私钥

- 生成私钥与地址

```shell
bash get_account.sh -p
```

执行上面的命令，可以得到类似下面的输出，按照提示输入密码，生成对应的p12文件。

```shell
Enter Export Password:
Verifying - Enter Export Password:
[INFO] Account Address   : 0x02f1b23310ac8e28cb6084763d16b25a2cc7f5e1
[INFO] Private Key (p12) : accounts/0x02f1b23310ac8e28cb6084763d16b25a2cc7f5e1.p12
```

同样可以根据私钥p12文件计算账户地址:

```shell
bash get_account.sh -P accounts/0x02f1b23310ac8e28cb6084763d16b25a2cc7f5e1.p12
```

执行上面的命令，结果如下

```shell
Enter Import Password:
MAC verified OK
[INFO] Account Address   : 0x02f1b23310ac8e28cb6084763d16b25a2cc7f5e1
```