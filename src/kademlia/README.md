# This should be placed at some other location

## Serialization structure

To send data serialize it with the SendPacket function.

In more depth: A packet consist of two structures: a **Header** structure and a **Payload** structure. The header structure includes the RPC method such as "Ping" or "Pong", the NodeID is simply the recepients NodeID, the IP is the recipients IP.

The Payload