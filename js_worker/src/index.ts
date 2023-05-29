import grpc from '@grpc/grpc-js';
import protoLoader from '@grpc/proto-loader';
import * as url from 'url';
const __dirname = url.fileURLToPath(new URL('.', import.meta.url));

// Path to proto definitions
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

// Echo service description
var echo = protoDescriptor.echo;

// Loose Echo service impl
const EchoService = {
  ping: (call, callback) => {
    callback(null, call.request);
  }
}

// Define grpc server and register the Echo service
var server = new grpc.Server();
server.addService(echo["Echo"].service, EchoService);

// Listen on our unix socket mounted by Docker
// NOTE: must set service authority to localhost in the client
server.bindAsync('unix:///tmp/comms/socket.sock', grpc.ServerCredentials.createInsecure(), (err, port) => {
  if(!!err) {
    console.log(err);
    return
  }

  console.log(`Listening on: ${port}`);

  server.start();
});
