// import {EchoEtyClient} from "./exam_grpc_pb"
// import { SayHelloRequest } from "./exam_pb";

// const client = new EchoEtyClient("https://my.grpc/server", null);
// const req = new SayHelloRequest();
// req.setText("johndoe");
// client.sayHello(req, (err, user) => {
//   /* ... */
// });

import * as em from "./exam"


function test(def :any){
  console.log(def.name)

  for (const key in def.methods) {
    // if (Object.hasOwnProperty.call(def.methods, key)) {
    if (def.methods.hasOwnProperty(key)) {
      const element = def.methods[key];
      
        console.log("create:",element.requestType.create())
      
    }
  }
 
  // name: "Client",
  // fullName: "exampb.Client",
  // methods: {
  //   sayHello: {
  //     name: "SayHello",
  //     requestType: SayHelloRequest,
  //     requestStream: false,
  //     responseType: SayHelloResponse,
  //     responseStream: false,
  //     options: {},
  //   },
  // },

  

}

test(em.ClientDefinition)

const baz = () => console.log('baz');
const foo = () => console.log('foo');
const zoo = () => console.log('zoo');
const start = () => {
  console.log('start');
  process.nextTick(foo);
  setImmediate(baz);
  new Promise((resolve, reject) => {
    resolve('bar');
  }).then((resolve) => {
    console.log(resolve);
    process.nextTick(zoo);
  });
  console.log('end11');
};
start();
console.log('end22');

// start foo bar zoo baz
