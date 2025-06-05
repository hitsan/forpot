import { Hono } from 'hono'

const app = new Hono()

let serverSession = {}

app.get('/', (c) => {
  return c.json({ message:'Running test server', servers: serverSession}, 200)
})

const getPort = async (c) => {
  const data = await c.req.json()
  return data.port
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
  '/servers/launch',
  async (c) => {
    const port = await getPort(c)
    const server = launchServer(port)
    serverSession[port] = server
    return c.json({ port: port }, 200)
  }
)

app.post(
  '/servers/down',
  async (c) => {
    const port = await getPort(c)
    shutdownServer(port)
    delete serverSession[port]
    return c.json({ port: port }, 200)
  }
)

Deno.serve({port: 8000}, app.fetch)
