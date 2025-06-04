import { Hono } from 'hono'

const app = new Hono()

let servers = [];

app.get('/', (c) => {
  return c.text('Running test server')
})

Deno.serve(app.fetch)
