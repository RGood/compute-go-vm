import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import * as url from 'url';
const __dirname = url.fileURLToPath(new URL('.', import.meta.url));
var PROTO_PATH = __dirname + '/../../protos/echo/echo.proto';
var packageDefinition = protoLoader.loadSync(
  PROTO_PATH,
  {keepCase: true,
   longs: String,
   enums: String,
   defaults: true,
   oneofs: true
  });

var protoDescriptor = grpc.loadPackageDefinition(packageDefinition);
var echo = protoDescriptor.echo;

const EchoService = {
  ping: (call, callback) => {
    console.log(call.request);
    callback(null, call.request);
  }
}

var server = new grpc.Server();
server.addService(echo["Echo"].service, EchoService);

server.bindAsync('unix:///tmp/comms/socket.sock', grpc.ServerCredentials.createInsecure(), (err, port) => {
  if(!!err) {
    console.log(err);
    return
  }

  console.log(`Listening on: ${port}`);

  server.start();
});
