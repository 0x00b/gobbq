
async function test(): Promise<string> {

    const chanRsp: any = new Promise((resolve, reject) => {
        console.log("111")

        setTimeout(() => {
            console.log("111 succ")
            resolve("succ")
        }, 1000)
    });

    console.log("222")


    const rsp = await chanRsp;

    console.log("333")

    console.log(rsp)

    return rsp
}

test().then((v: string) => {
    console.log(v)
    return v
})

// console.log(test())