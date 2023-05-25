const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const protos = require('google-proto-files');

const packageDefinition = protoLoader.loadSync('../coffeeshop.proto', {
   keepCase: true,
   longs: String,
   enums: String,
   defaults: true,
   oneofs: true,
   includeDirs: [
      '../',
      '../coffeeshop.proto',
   ],
});
const coffeeshop = grpc.loadPackageDefinition(packageDefinition).coffeeshop;

const service = {
  CreatePayment: (call, callback) => { 
    console.log(call.request);
    // return random transaction id
    res = {
      transaction_id:  Math.floor(Math.random() * 10000000).toString()
    }
    console.log(res);
    callback(null, res);
  }
};

const server = new grpc.Server();

// register server
server.addService(coffeeshop.PaymentService.service, service);
// bind to port and start server 
server.bindAsync('localhost:8096', grpc.ServerCredentials.createInsecure(), (err, port) => {
  if (err) {
    console.error('Failed to start gRPC server: ', err);
    return;
  }
  console.log(`gRPC server running on port ${port}`);
  server.start();
});
