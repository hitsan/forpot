const serve = () => {
  

}

const parsePorts = (args: string[]): Number[] => {
  const ports = args.flatMap(arg => {
    const port = Number(arg);
    if (Number.isNaN(port)) {
      return [];
    }
    return port;
  });
  return ports;
}

const launch = (port) => {
  const server = Deno.serve(
    { port : port },
    (req: Request) => {
      return new Response(`Running:${arg}`);
    }
  );
  return server;
}

const ports = parsePorts(Deno.args);
if (ports.length == 0) {
  Deno.exit(0);
}

const servers = ports.map(port => launch(port));
for (let server of servers) {
  server.shutdown();
}

