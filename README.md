# FABCHAT

使用Fabric实现的简单聊天应用，通过非对称加密方法对消息内容进行加密。

## 方法

### Register

注册账户，登记接受消息使用的公钥，当别人想该账号发送消息的时候，可以使用对应的公钥进行加密。

### GetAccount

获取ID对应的账户登记的公钥信息

### Send

发送消息，发送消息成功后会触发链码事件，接受方可以通过监听链码事件得知接收到消息

### Received

使用链码事件中携带的消息ID读取接受到的消息。