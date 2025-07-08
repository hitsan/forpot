package ssh

type sessionOp struct {
	op      string
	port    int
	session *ForwardSession
}

type sessionGetOp struct {
	port   int
	result chan *ForwardSession
}

type sessionGetAllOp struct {
	result chan map[int]struct{}
}

type SessionSynchronizer struct {
	store      *SessionStore
	opChan     chan sessionOp
	getChan    chan sessionGetOp
	getAllChan chan sessionGetAllOp
	closeChan  chan chan struct{}
	done       chan struct{}
}

func NewSessionSynchronizer() *SessionSynchronizer {
	s := &SessionSynchronizer{
		store:      NewSessionStore(),
		opChan:     make(chan sessionOp),
		getChan:    make(chan sessionGetOp),
		getAllChan: make(chan sessionGetAllOp),
		closeChan:  make(chan chan struct{}),
		done:       make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *SessionSynchronizer) run() {
	for {
		select {
		case <-s.done:
			return
		case op := <-s.opChan:
			switch op.op {
			case opSet:
				s.store.set(op.port, op.session)
			case opDelete:
				s.store.delete(op.port)
			}
		case getOp := <-s.getChan:
			getOp.result <- s.store.get(getOp.port)
		case getAllOp := <-s.getAllChan:
			getAllOp.result <- s.store.getAll()
		case resultChan := <-s.closeChan:
			s.store.close()
			close(resultChan)
			return
		}
	}
}

func (s *SessionSynchronizer) Set(port int, session *ForwardSession) {
	select {
	case s.opChan <- sessionOp{op: opSet, port: port, session: session}:
	case <-s.done:
	}
}

func (s *SessionSynchronizer) Delete(port int) {
	select {
	case s.opChan <- sessionOp{op: opDelete, port: port}:
	case <-s.done:
	}
}

func (s *SessionSynchronizer) Get(port int) *ForwardSession {
	result := make(chan *ForwardSession)
	select {
	case s.getChan <- sessionGetOp{port: port, result: result}:
		return <-result
	case <-s.done:
		return nil
	}
}

func (s *SessionSynchronizer) GetAll() map[int]struct{} {
	result := make(chan map[int]struct{})
	select {
	case s.getAllChan <- sessionGetAllOp{result: result}:
		return <-result
	case <-s.done:
		return make(map[int]struct{})
	}
}

func (s *SessionSynchronizer) Close() {
	result := make(chan struct{})
	select {
	case s.closeChan <- result:
		<-result
	default:
	}
	close(s.done)
}