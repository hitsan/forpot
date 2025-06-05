import { Hono } from 'hono'

const app = new Hono()

let serverSession = {}

app.get('/', (c) => {
  return c.json({ message:'Running test server', servers: serverSession}, 200)
})

const parsePort = (c): number | null => {
  const portStr = c.req.param('port')
  const port = parseInt(portStr, 10)
  if (isNaN(port) || port <= 0 || port > 65535) {
    return null
  }
  return port
}

const launchServer = (port) => {
  const app = new Hono()
  const server = Deno.serve({ port: port }, app.fetch )
  return server
}

const shutdownServer = (port) => {
  const server = serverSession[port]
  server.shutdown()
}

app.post(
  '/servers/:port/launch',
  (c) => {
    const port = parsePort(c)
    if (port == null) {
      return c.json({ message: "Illegal port" }, 400)
    }
    if (port in serverSession) {
      return c.json({ message: "Already lauched"}, 409)
    }
    const server = launchServer(port)
    serverSession[port] = server
    return c.json({ port: port }, 200)
  }
)

app.post(
  '/servers/all/down',
  (c) => {
    const ports = Object.keys(serverSession)
    ports.map(port => shutdownServer(port))
    serverSession = {}
    return c.json({ message: "All server down"}, 400)
  }
)

app.post(
  '/servers/:port/down',
  (c) => {
    const port = parsePort(c)
    if (port == null) {
      return c.json({ message: "Illegal port" }, 400)
    }
    if (!(port in serverSession)) {
      return c.json({ message: "Not found port" }, 404)
    }
    shutdownServer(port)
    delete serverSession[port]
    return c.json({ port: port }, 200)
  }
)

Deno.serve({port: 8000}, app.fetch)
