import { FrameReq, GameOverReq, InputData, InputReq, OPID, ProgressReq, StartReq } from '../../frameproto/frame';
import { FrameClientEntityDefinition, FrameSeverEntity, NewFrameSeverEntity } from '../../frameproto/frame.bbq'
import { NewFrameService } from '../testpb/testpb.bbq'
import { Client } from 'gobbq-ts/dist/src';
import { Context } from 'gobbq-ts/dist/src/dispatcher/context';
import Long from "long";

let para=process.argv.slice(2)[0];
let UID = Long.fromValue(para)

class FrameClientEntity {

  // Start
  Start(ctx: Context, request: StartReq): void {
    let x = 1
    let y = 1
    setInterval(() => {
      let req = InputReq.create({
        UID:UID,
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
    }, 50)

  }

  // Frame
  Frame(ctx: Context, request: FrameReq): void {

    console.log("recv frame:", JSON.stringify(request))

  }

  // Progress 通知客户端其他人加载进度
  Progress(request: ProgressReq): void { 

  }

  // GameOver 游戏结束
  GameOver(request: GameOverReq): void {

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

  await startSvc.StartFrame({}).then(({ error, response }) => {

    if (error || !response.FrameSvr) {
      console.log(error)
      return
    }

    frameSvr = NewFrameSeverEntity(client, response.FrameSvr)
  })

  await frameSvr.Join({ UID: UID }).then(({ response }) => {
    console.log("succ rsp:", response)
  })

  frameSvr.Ready({UID: UID})

}

test()