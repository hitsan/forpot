let arg = Deno.args[0];
let port = Number(arg);

if (typeof port !== 'number') {
  Deno.exit(0);
}

Deno.serve(
  { port : port },
  (req: Request) => {
    return new Response(`Hello ${arg}`);
  }
);
