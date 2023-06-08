

import { SayHelloRequest, SayHelloResponse } from '../../example/exampb/exam';
import { ClientEntityDefinition, NewEchoSvc2Service } from '../../example/exampb/exam.bbq';
import { Client } from '../src';
import { Context } from '../src/dispatcher/context';


class EchoImpl {
  SayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {

    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse:", request.text)

    return {text:"xxxx"}
  }
}

async function test() {

  const remote = {
    // port: 8899,
    port: 59551,
    host: 'localhost',
    protocol: 'kcp',
  } as const;

  let client = await Client.create(ClientEntityDefinition, new EchoImpl(), { remote })

  let c = NewEchoSvc2Service(client)
  // let c = NewEchoEtyEntity(client, EntityID.create({ID: "xxxx"}))

  let rsp = c.SayHello({ text: "request", CLientID: client.EntityID})

  rsp.then((rsp) => {
    if (rsp instanceof Error) {
      console.log("error", rsp)
      return
    }
    console.log("succ rsp:", rsp)
  })

  // c.SayHello({ text: "request", CLientID: undefined })
  // c.SayHello({ text: "request", CLientID: undefined })
  // c.SayHello({ text: "request", CLientID: undefined })
}

test()