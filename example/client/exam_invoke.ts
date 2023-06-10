import { SayHelloRequest, SayHelloResponse } from '../exampb/exam';
import { ClientEntityDefinition, NewEchoSvc2Service } from '../exampb/exam.bbq';
import { Client } from 'gobbq-ts/dist/src';
import { Context } from 'gobbq-ts/dist/src/dispatcher/context';


class EchoImpl {
  SayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {

    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse:", request.text)

    return { text: "xxxx" }
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

  for (let index = 0; index < 10; index++) {
    let rsp = c.SayHello({ text: "request"+index, CLientID: client.EntityID })

    rsp.then((rsp) => {
      if (rsp instanceof Error) {
        console.log("error", rsp)
        return
      }
      console.log("succ rsp:", rsp)
    })
  }

}

test()