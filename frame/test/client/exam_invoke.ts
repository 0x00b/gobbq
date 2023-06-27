import { ClientResultReq, CloseReq, FrameReq, InputData, InputReq, OPID, StartReq } from '../../frameproto/frame';
import { FrameClientEntityDefinition, FrameSeverEntity, NewFrameSeverEntity } from '../../frameproto/frame.bbq'
import { NewFrameService } from '../testpb/testpb.bbq'
import { Client } from 'gobbq-ts/dist/src';
import { Context } from 'gobbq-ts/dist/src/dispatcher/context';

class FrameClientEntity {

  // Start
  Start(ctx: Context, request: StartReq): void {
    let x = 1
    let y = 1
    setInterval(() => {
      let req = InputReq.create({
        Input: {
            OPID: OPID.Move,
            Pos: {
                x: x++,
                y: y++,
                z: 0,
            }
          }
    })
      frameSvr.Input(req)
    }, 200)

  }

  // Frame
  Frame(ctx: Context, request: FrameReq): void {

    console.log("recv frame:",JSON.stringify(request))

  }

  // Result
  Result(ctx: Context, request: ClientResultReq): void {

  }

  // Close
  Close(ctx: Context, request: CloseReq): void {

  }
}


let frameSvr: FrameSeverEntity

async function test() {

  const remote = {
    // port: 8899,
    port: 59551,
    host: 'localhost',
    protocol: 'kcp',
  } as const;

  let client = await Client.create(FrameClientEntityDefinition, new FrameClientEntity(), { remote })

  let startSvc = NewFrameService(client)

  startSvc.StartFrame({}).then(({ error, response }) => {

    if (error || !response.FrameSvr) {
      console.log(error)
      return
    }

    frameSvr = NewFrameSeverEntity(client, response.FrameSvr)
    // let c = NewEchoEtyEntity(client, EntityID.create({ID: "xxxx"}))

    let rsp = frameSvr.Join({ CLientID: client.EntityID })
    rsp.then(({ response }) => {
      console.log("succ rsp:", response)
    })

  })


}

test()