
import * as packet from "../codec/packet"
import * as bbq from "../../proto/bbq/bbq"
import { v4 as uuidv4 } from 'uuid';

export class Client {
    ID: bbq.EntityID
    net: string
    address: string
    port: string

    callbacks: Map<string, (pkt: packet.Packet) => void>

    constructor() {
        this.ID = bbq.EntityID.create()
        this.ID.ID = uuidv4()
        this.callbacks = new Map
    }

    SendPacket(pkt: packet.Packet) {
        if (pkt == null) {
            throw "invalid packet"
        }

        if (pkt.Buffer.length != pkt.TotalLen + 4) {
            throw "bad packet"
        }

        // write pkt.Buffer
    }

    readPacket(tempBuff: Buffer): packet.Packet | null {
        return packet.ReadPacket(tempBuff)
    }

    RegisterCallback(requestId: string, callback: (pkt: packet.Packet) => void) {
        this.callbacks[requestId] = callback
    }
}

