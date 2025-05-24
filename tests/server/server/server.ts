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
  Deno.serve(
    { port : port },
    (req: Request) => {
      return new Response(`Hello ${arg}`);
    }
  );
}

const ports = parsePorts(Deno.args);
if (ports.length == 0) {
  Deno.exit(0);
}

ports.map(port => launch(port));
