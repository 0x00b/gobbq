import {
  Context,
  createContext,
  ParameterizedContext,
} from './context';
import * as compose from './compose';
import { Packet } from '../codec/packet';
import { MethodImpl, ServiceDefinition } from './service';
import { EchoDefinition } from '../../../example/exampb/exam.bbq';
import { SayHelloRequest, SayHelloResponse } from '../../../example/exampb/exam';
import * as bbq from '../../../proto/bbq/bbq';
import { decode, encode, UnaryRequestMessage } from '../codec/msg';

export type Middleware<CustomContextT = {}> = compose.Middleware<ParameterizedContext<CustomContextT>>;

interface Options {
  middleware?: Middleware<any>[],
  supportStream?: boolean,
}

export class Dispatcher<CustomContextT = {}> {
  private inited:boolean=false;
  
  private readonly middleware: Middleware<CustomContextT & any>[] = [];

  private readonly handlers = new Map<string, MethodImpl<any, any>>();

  public constructor(
    public readonly serviceName: string,
    private options: Options = { supportStream: false },
  ) {
    this.options?.middleware?.forEach(fn => this.use(fn));
  }

  public RegisterClientImpl(def: ServiceDefinition, impl: any) {
    if (this.inited) {
      throw new Error("already init")
    }
  
    const prototype = Object.getPrototypeOf(impl);
    let ma = Reflect.ownKeys(prototype)
    console.log(ma);

    for (const [key, value] of Object.entries(def.methods)) {
      console.log(`${key} = ${value}`);

      let find = false

      console.log(JSON.stringify(impl));
      ma.forEach(implkey => {
        console.log(implkey);
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

  public onUnaryRequest(pkt: Packet) {

    let methodName = pkt.Header.DstEntity?.Type + "/" + pkt.Header.Method

    let handler = this.handlers.get(methodName)

    if (handler === undefined) {
      throw new TypeError('no handler!');
    }

    var {
      handle,
      requestDeserialize
    } = handler;

    let bodydata = pkt.PacketBody()
    if (!bodydata) {
      throw new TypeError('no bodydata!');
    }

    let mh: compose.Middleware<any> = (context: any, next: compose.Next): any => {
      let rsp = handle(context, context.req)
      ctx.body = rsp
      return rsp
    }

    const fn = compose.compose([...this.middleware, mh]);
    let ctx: ParameterizedContext<CustomContextT> = createContext(pkt, /*transport*/)

    // console.log("bodydata",bodydata)
    ctx.req = requestDeserialize(bodydata)
    // console.log("req",JSON.stringify(ctx.req))

    return this.handleUnaryRequest(fn, ctx);
    
  }

  private handleUnaryRequest(
    fn: compose.ComposedMiddleware<ParameterizedContext<CustomContextT>>,
    ctx: ParameterizedContext<CustomContextT>,
  ) {
    console.log('handleUnaryRequest %s', ctx.Method);
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

  console.log("rsp:",JSON.stringify( ctx.body))

  const code = ctx.status
    ?? (ctx.body ? 0 : 1) //RetCode.INVOKE_SUCCESS : RetCode.INVOKE_UNKNOWN_ERR);

  ctx.respond({
    // Header?: {ErrCode: code},
    // ...ctx.body,
  });
}


class Echo {
  sayHello(ctx: Context, request: SayHelloRequest): SayHelloResponse {
    console.log("sssssss sayHello(request: SayHelloRequest): SayHelloResponse", request.text)

    let rsp = SayHelloResponse.create()
    rsp.text = "xxxx"
    return rsp
  }
}


function exx(): number {
  return 111
}

function exam(): Promise<number> {
  // last called middleware #
  return Promise.resolve(exx());

};

function test() {

  // exam().then((x) => {
  //   console.log("xxxx", x)
  // })


  let d: Dispatcher = new Dispatcher("")

  const ctx: any = new Echo();

  d.RegisterClientImpl(EchoDefinition, ctx)

  d.use<{ cost: number }>(async (ctx, next) => {
    const startTime = Date.now();
    await next();
    ctx.cost = Date.now() - startTime;

    console.log("const:", ctx.cost)
  })


  let hdr = bbq.Header.create()
  hdr.RequestId = "1111"
  hdr.Timeout = 1000
  hdr.RequestType = bbq.RequestType.RequestRequest;

  hdr.ServiceType = bbq.ServiceType.Service;
  hdr.Method = "sayHello";
  hdr.DstEntity = bbq.EntityID.create()
  hdr.DstEntity.Type="exampb.Echo"

  //    hdr.SrcEntity = c.ID;
  //    hdr.DstEntity = dstEntity

  let req:SayHelloRequest = SayHelloRequest.create()
  req.text="xxx"

  const data = Buffer.from(SayHelloRequest.encode(req).finish());

  const message: UnaryRequestMessage = {
    Header: hdr,
    Body: Buffer.from(data),
  };
  
  console.log("data", data)

  let edata = encode(message)
  console.log("edata", edata)
  let pkt = decode(edata)

  if (!pkt) {
    console.log("empty pkt")
    return
  }
  console.log("pkt", JSON.stringify(pkt))
  console.log("pkt", pkt.Buffer)

  let res = d.onUnaryRequest(pkt)
  console.log("xxx", res)

  // console.log(d)
}

test()