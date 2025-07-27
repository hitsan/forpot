```mermaid
classDiagram
    class SessionManager {
        -remoteHost: string
        -client: *ssh.Client
        -sync: *SessionSynchronizer
        +NewSessionManager(remoteHost string, client *ssh.Client) *SessionManager
        +getPortMap() map[int]struct{}
        +selectUpdatePorts(portChan chan []int, upPortChan chan int, downPortChan chan int)
        +updateForwardingSession(upPortChan chan int, downPortChan chan int)
        +Close()
    }

    class ForwardSession {
        -listener: net.Listener
        -remoteAddr: string
        -once: sync.Once
        -done: chan struct{}
        +NewForwardSession(localAddr string, remoteAddr string) (*ForwardSession, error)
        +Close()
        -handleConnection(client *ssh.Client, localConn net.Conn)
        -forwardPort(client *ssh.Client)
    }

    class SessionStore {
        -sessions: map[int]*ForwardSession
        +NewSessionStore() *SessionStore
        -set(port int, session *ForwardSession)
        -delete(port int)
        -get(port int) *ForwardSession
        -getAll() map[int]struct{}
        -close()
    }

    class SessionSynchronizer {
        -store: *SessionStore
        -opChan: chan sessionOp
        -getChan: chan sessionGetOp
        -getAllChan: chan sessionGetAllOp
        -closeChan: chan chan struct{}
        -done: chan struct{}
        +NewSessionSynchronizer() *SessionSynchronizer
        -run()
        +Set(port int, session *ForwardSession)
        +Delete(port int)
        +Get(port int) *ForwardSession
        +GetAll() map[int]struct{}
        +Close()
    }

    class Uid {
        <<type>>
    }

    SessionManager "1" *-- "1" SessionSynchronizer
    SessionSynchronizer "1" *-- "1" SessionStore
    SessionStore "1" *-- "N" ForwardSession

```
