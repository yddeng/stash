# smux

在一个TCP连接上分出多个stream

## 主要用途

网络游戏中客户端通常通过网关服务连接到游戏服务器。一般情况下一个网关服务跟每个游戏服务保持一条TCP连接。

通常的做法是客户端连接到网关，有网关为客户端分配一个会话ID，向游戏服务通告会话ID。后续网关和游戏服之间针对特定客户端的消息都需要携带会话ID，以标识消息与哪个会话关联。

smux将会话信息封装在内部。为通信两端模拟出一个stream，客户端连接网关后，网关为这个客户连接建立一个到游戏服的stream，后续游戏服跟客户端的通信，全部由这个stream负责转发。

## 流控

stream的数据接收采用pull模式，当一端希望接收数据的时候，向另一端发送一个want请求，另一端接收到want请求后，递增wantReq计数器。

发送端发送的时候首先检查wantReq是否大于0，如果是则发送数据，并递减wantReq，否则将数据压进待发送队列，直到接收到want请求之后才发送。

所有实现都应遵循这个规定，如果接收端接收到数据，但发现没有未决的want请求，将直接丢弃数据（例如对于一个错误的实现，接收端只发起了一个want请求，而发送端连接发送了两次数据，第二次收到数据时，接收端发现没有未决的want，将直接丢弃到达的数据）。



## stream的建立

stream分为客户端和服务器。

客户端向服务建立一个TCP连接，服务器接受连接并使用连接创建MuxSocketServer，然后调用Listen等待新stream的到来。

客户端建立TCP连接后使用连接创建MuxSocketClient。需要建立stream的时候调用MuxSocketClient.Dial。

### stream id

在每个TCP连接上，最多允许创建65535个stream，每个stream用一个16位整数标识，范围从0-65535。

stream id由client选择，client记录了当前连接所有空闲的id,dial的时候选择一个id来建立stream。当stream关闭之后，将id返还到id池供后续使用。



## API	

`func (s *MuxStream) Recv(onData func(*MuxStream, []byte)) error`

发出一个接收请求，如果有数据到达调用onData。



`func (s *MuxStream) SetRecvTimeout(timeout time.Duration)`

设置读超时，如果读发生超时将直接关闭MuxStream。



`func (s *MuxStream) AsyncSend(o interface{}, cb func(*MuxStream, error), timeout ...time.Duration) error`

异步发送，发送完成或发生超时调用cb。对于不能阻塞当前goroutine的情况应该选用异步发送。为了避免发送队列的堆积，可以设计一个计数器，记录当前未决的发送数量，每次发送前增加计数器，在cb中递减计数器。如果未决的发送数量超过一定值则放弃发送。



`func (s *MuxStream) SyncSend(o interface{}, timeout ...time.Duration) error`

同步发送，如果无法发送会阻塞当前goroutine直到超时（如果有）。



```
func (s *MuxStream) Close(err error)
```

关闭当前stream,向对端发送EOF。
