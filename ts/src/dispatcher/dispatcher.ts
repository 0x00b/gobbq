import {
  Context,
  createContext,
  ParameterizedContext,
} from './context';
import * as compose from './compose';
import { Packet } from '../codec/packet';
import { MethodImpl, ServiceDefinition } from './service';
import { ClientTransport } from '../transport';

export type Middleware<CustomContextT = {}> = compose.Middleware<ParameterizedContext<CustomContextT>>;

interface Options {
  middleware?: Middleware<any>[],
}

export class Dispatcher<CustomContextT = {}> {
  private inited: boolean = false;

  private readonly middleware: Middleware<CustomContextT & any>[] = [];

  private readonly handlers = new Map<string, MethodImpl<any, any>>();

  public constructor(
    def: ServiceDefinition, impl: any,
    private options: Options = {},
  ) {
    this.options?.middleware?.forEach(fn => this.use(fn));

    if (this.inited) {
      throw new Error("already init")
    }

    const prototype = Object.getPrototypeOf(impl);
    let ma = Reflect.ownKeys(prototype)

    for (const [key, value] of Object.entries(def.methods)) {
      let find = false
      ma.forEach(implkey => {
        const fn = impl[implkey];
        if (typeof fn === 'function') {
          if (implkey == key) {
            find = true
            return
          }
        }
      })

      if (!find) {
        throw new Error("impl not impl def")
      }

      let func = impl[key].bind(impl)
      this.handlers.set(def.typeName + "/" + key, { ...value, handle: func })
    }

    this.inited = true
  }

  public use<NewCustomContextT = {}>(fn: Middleware<CustomContextT & NewCustomContextT>): this {
    if (typeof fn !== 'function') {
      throw new TypeError('middleware must be a function!');
    }

    console.log('use %s', fn.name || '-');

    this.middleware.push(fn);
    return this;
  }

  public onUnaryMessage(pkt: Packet, transport: ClientTransport) {

    let methodName = pkt.Header.DstEntity?.Type + "/" + pkt.Header.Method

    let handler = this.handlers.get(methodName)

    if (handler === undefined) {
      throw new TypeError('no handler!');
    }

    var {
      handle,
    } = handler;

    let bodydata = pkt.PacketBody()
    if (!bodydata) {
      throw new TypeError('no bodydata!');
    }

    let mh: compose.Middleware<any> = (context: any, next: compose.Next): any => {
      let rsp = handle(context, context.req)
      context.responseBody = rsp
      return rsp
    }

    const fn = compose.compose([...this.middleware, mh]); 

    return this.handleUnaryRequest(fn, createContext(pkt, handler, transport));

  }

  private handleUnaryRequest(
    fn: compose.ComposedMiddleware<ParameterizedContext<CustomContextT>>,
    ctx: ParameterizedContext<CustomContextT>,
  ) {
    // console.log('handleUnaryRequest %s', ctx.Method);
    const handleResponse = () => respondUnary(ctx);
    const onerror = (err: Error) => ctx.onerror(err);

    fn(ctx)
      .then(handleResponse)
      .catch(onerror);

    return
  }
}

function respondUnary(ctx: Context) {
  // if (ctx.responded || ctx.callType === CallType.ONEWAY_CALL) {
  //   return;
  // }

  // console.log("11respondUnary:", JSON.stringify(ctx.responseBody))

  // const code = ctx.status
  //   ?? (ctx.responseBody ? 0 : 1) //RetCode.INVOKE_SUCCESS : RetCode.INVOKE_UNKNOWN_ERR);


  ctx.respond();
}


// class Echo {
//   sayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {
//     console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse", request.text)

//     let rsp = SayHelloResponse.create()
//     rsp.text = "xxxx"
//     return rsp
//   }
// }


// function exx(): number {
//   return 111
// }

// function exam(): Promise<number> {
//   // last called middleware #
//   return Promise.resolve(exx());

// };

// function test() {

//   // exam().then((x) => {
//   //   console.log("xxxx", x)
//   // })

//   const ctx: any = new Echo();

//   let d: Dispatcher = new Dispatcher(EchoDefinition, ctx)

//   d.use<{ cost: number }>(async (ctx, next) => {
//     const startTime = Date.now();
//     await next();
//     ctx.cost = Date.now() - startTime;

//     console.log("const:", ctx.cost)
//   })


//   let hdr = bbq.Header.create()
//   hdr.RequestId = "1111"
//   hdr.Timeout = 1000
//   hdr.RequestType = bbq.RequestType.RequestRequest;

//   hdr.ServiceType = bbq.ServiceType.Service;
//   hdr.Method = "sayHello";
//   hdr.DstEntity = bbq.EntityID.create()
//   hdr.DstEntity.Type = "exampb.Echo"

//   //    hdr.SrcEntity = c.ID;
//   //    hdr.DstEntity = dstEntity

//   let req: SayHelloRequest = SayHelloRequest.create()
//   req.text = "xxx"

//   const data = Buffer.from(SayHelloRequest.encode(req).finish());

//   const message: UnaryRequestMessage = {
//     Header: hdr,
//     Body: Buffer.from(data),
//   };

//   console.log("data", data)

//   let edata = encode(message)
//   console.log("edata", edata)
//   let pkt = decode(edata)

//   if (!pkt) {
//     console.log("empty pkt")
//     return
//   }
//   console.log("pkt", JSON.stringify(pkt))
//   console.log("pkt", pkt.Buffer)

//   let res = d.onUnaryMessage(pkt)
//   console.log("xxx", res)

//   // console.log(d)
// }

// test()