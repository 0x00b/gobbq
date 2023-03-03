import { error } from "console";
import * as bbq from "../example/exampb/bbq"
import * as kcpjs from "kcpjs"


const MaxBufferCap = 16 * 1024 * 1024 //16M

class Packet {
    Header: bbq.Header; // can only be used during the packet's lifetime
    TotalLen: number;
    HeaderLen: number;
    Buffer: Buffer;
 }

 let packet= new Packet();
 packet.Buffer = Buffer.alloc(MaxBufferCap);

async function readPacket(tempBuff: Buffer): Promise<[Packet|null, Error|null]> {
    packet.Buffer.write(tempBuff.toString());
 

    if (packet.Buffer.length < 4) {
        return [null, new Error("not enougth lenth")];
    }

    let packetDataSize: number;
    packetDataSize = new DataView(tempBuff).getInt32(0, true)
    if (packetDataSize > MaxBufferCap) {
      return [null, new Error("too long data")];
    }
  

    packet.TotalLen = packetDataSize;
  
    packet.HeaderLen = new DataView(tempBuff).getInt32(4, true);
    let packetHeaderData: Uint8Array = packet.Buffer.subarray(8, packet.HeaderLen);
  
    let err: Error = protobuf.Unmarshal(packetHeaderData, packet.Header);
    if (err !== null) { 
      return [null, err];
    }
  
    console.log(`recv raw: ${packet.Buffer.toString()}`);
  
    // pc.readMsgCnt++;
  
    if (HasFlags(packet.Header.CheckFlags, FlagDataChecksumIEEE)) {
      try {
        await pc.rw.readFull(tempBuff.subarray(0, 4));
      } catch (err) { 
        return [null, err];
      }
  
      let packetBodyCrc: number = crc32.ChecksumIEEE(packet.Data());
      if (packetBodyCrc !== packetEndian.Uint32(tempBuff)) {
        release();
        return [null, errChecksumError];
      }
    }
  
    return [packet, null];
  }