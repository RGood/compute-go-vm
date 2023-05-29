import {
  loadPackageDefinition,
  UntypedServiceImplementation,
  ServerUnaryCall,
  sendUnaryData,
  Server,
  ServerCredentials,
} from '@grpc/grpc-js';
import { loadSync } from '@grpc/proto-loader';

// Path to proto definitions
var PROTO_PATH = __dirname + '/../../protos/echo/echo.proto';

var packageDefinition = loadSync(
  PROTO_PATH,
  {keepCase: true,
   longs: String,
   enums: String,
   defaults: true,
   oneofs: true
  });

var protoDescriptor = loadPackageDefinition(packageDefinition);

// Echo service description
var echo = protoDescriptor.echo;

// ==========================USER CODE GOES HERE==================================
// Loose Echo service impl
const EchoService: UntypedServiceImplementation = {
  ping: (call: ServerUnaryCall<any, any>, callback: sendUnaryData<any>): void => {
    callback(null, call.request);
  }
}
// ===============================================================================

// Define grpc server and register the Echo service
var server = new Server();
server.addService(echo["Echo"].service, EchoService);

// Listen on our unix socket mounted by Docker
// NOTE: must set service authority to localhost in the client
server.bindAsync('unix:///tmp/comms/socket.sock', ServerCredentials.createInsecure(), (err, port) => {
  if(!!err) {
    console.log(err);
    return
  }

  console.log(`Listening on: ${port}`);

  server.start();
});
