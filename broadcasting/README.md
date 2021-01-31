# Broadcasting

The Broadcasting playground demonstrates distributing a message (in this case a string, but this could be any type) to
multiple listening GRPC streams. This could be further modified to allow subscriptions to certain topics, rather than
a message being distributed to all listeners.

## Iterations

### v1

v1 represents my first attempt of a "broadcast manager". It uses a central struct, with access to a map protected by a 
mutex. When an incoming request for a subscription occurs, the Subscribe() method is called which returns a channel that
will recieve any messages that are sent via the Send() method.

Positives:
- Can be easily consumed in multiple different ways since it returns a channel rather than immediately broadcasting back
  to some stream.

Negatives:
- Easily blocked. If a single consumer was to hang, and the buffer of the channel was full, it would block publication
  of messages.
  
### v2

v2 demonstrates an alternative design that