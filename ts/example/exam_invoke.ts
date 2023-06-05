

import { SayHelloRequest, SayHelloResponse } from '../../example/exampb/exam';
import { EchoDefinition, NewEchoService } from '../../example/exampb/exam.bbq';
import { Client } from '../src';
import { Context } from '../src/dispatcher/context';


class EchoImpl {
  SayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {

    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse:", request.text)

    return {text:"xxxx"}
  }
}

function test() {

  const remote = {
    port: 8899,
    host: 'localhost',
    protocol: 'kcp',
  } as const;

  let client = new Client(EchoDefinition, new EchoImpl(), { remote })
  let c = NewEchoService(client)

  let rsp = c.SayHello({ text: "request", CLientID: undefined })

  rsp.then((rsp) => {
    if (rsp instanceof Error) {
      console.log("error", rsp)
      return
    }

    console.log("succ rsp:", rsp)
  
  })

  c.SayHello({ text: "request", CLientID: undefined })
  c.SayHello({ text: "request", CLientID: undefined })
  c.SayHello({ text: "request", CLientID: undefined })
}

test()