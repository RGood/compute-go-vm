const grpc = require("@grpc/grpc-js");
var protoLoader = require("@grpc/proto-loader");
const net = require("net");
var PROTO_PATH = __dirname + "/../../protos/echo/echo.proto";
var packageDefinition = protoLoader.loadSync(
  PROTO_PATH,
  {
    keepCase: true,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true
  }
);
const EchoService = {
  ping: (call, callback) => {
    const request = call.request;
    console.log(request);
    const response = new packageDefinition.Echo.Message();
    response.setMessage(`Hello, ${request.getName()}!`);
    callback(null, response);
  }
};
function serverProto(socket) {
  return {
    getPeer: () => {
      return socket.remoteAddress;
    },
    write: (data) => {
      socket.write(data);
    },
    end: () => {
      socket.end();
    },
    on: (event, callback) => {
      socket.on(event, callback);
    }
  };
}
var protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
var echoDefinition = protoDescriptor.echo;
console.log(JSON.stringify(echoDefinition));
var server = new grpc.Server();
const unixServer = net.createServer((socket) => {
  server.addProtoService(packageDefinition.Echo.echo, EchoService);
  server.handle(serverProto(socket));
});
const unixSocketPath = "/tmp/comms/socket.sock";
unixServer.listen(unixSocketPath, () => {
  console.log(`gRPC server is running on Unix socket: ${unixSocketPath}`);
});
