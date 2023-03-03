import {SayHelloRequest} from "./exam";


function test() {
    let req = SayHelloRequest.fromJSON({text:"x", CLientID:{ID:"tid",Type:"",ProxyID:""}});
    
    console.log("Hello, " + req.CLientID?.ID);
}
 
test();