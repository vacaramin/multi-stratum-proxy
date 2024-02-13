package sha256stratum

import (
	"net"
	"time"
)

/*
	TO DO

Defines the Miner class which represents a connected miner.
Establishes a connection to the pool proxy and handles incoming/outgoing messages.
Implements logic for handling fees by forwarding jobs to the fee pool and applying the calculated fee.
*/

type Miner struct {
	ID           int
	Controller   *Controller
	Socket       net.Conn
	Port         string
	MiningServer string
	LogPrefix    string
	Jobs         []map[string]interface{}
	Fee          struct {
		Pool     PoolToFee
		Active   bool
		Interval *time.Timer
		Timeout  *time.Timer
		Jobs     []map[string]string
	}
}

// const logger = require('../../lib/logger')
// const {Pool, PoolToFee} = require('./pool')

// module.exports = class Miner {
//     constructor(controller, socket, port, miningServer) {
//         this.id = ++controller.minersCount
//         this.ctrl = controller
//         this.socket = socket
//         this.port = port
//         this.miningServer = miningServer
//         this.logPrefix = controller.logPrefix + ' MINER:' + this.id + ' >'

//         this.fee = {
//             pool: {},
//             active: false,
//             interval: null,
//             timeout: null,
//             jobs: []
//         }

//         this.socket.setKeepAlive(true)
//         this.socket.setEncoding('utf8')

//         logger.info(this.logPrefix, `Miner connected (ip: ${socket.remoteAddress})`)
//     }

//     /*
//     connectToPool- Open connection to main pool
//      */

//     connectToPool() {
//         this.pool = new Pool(this)
//         this.pool.connect()
//     }

//     /*
//     setFee-
//     */

//     setFee(window, percent) {
//         /*
//         Open connection to Sencond/Fee pool
//      */

//         this.fee.pool = new PoolToFee(this, this.ctrl.params)
//         this.fee.pool.connect()
//         /*
//         Setup interval in which we send mining.notify to miner from second\fee pool and receive mining.submit from miner and send to second/fee pool,
//          all the rest of the time we send mining.notify and other methods from main pool and send responses from miner to main pool
//          */

//         this.fee.interval = setInterval(() => {
//             /*
//                 check if array include "ready", its mean than mining.subscribe and\or mining.configure was send from miner to main pool and was duplicated to second\fee pool
//              */

//             if (!this.fee.pool.status.includes('ready'))
//                 return

//             /*
//             Hack for Whatsminer models( miner brand)
//             Check if second\fee pool more than main pool\it need for reduce rejects to main pool,
//             because miner will send shares with last difficulty from last pool, and if main pool difficulty more than second\fee pool share will be rejected
//             because difficulty is threshold value below which pool will reject share(mining.submit)

//             */

//             if (this.pool.useragent === 'whatsminer/v1.0' && this.pool.difficulty > this.fee.pool.difficulty) {
//                 logger.debug(this.logPrefix, 'Fee pool less than pool diff: ')
//                 logger.debug(this.logPrefix, 'Pool diff: ', this.pool.difficulty)
//                 logger.debug(this.logPrefix, 'Fee pool diff: ', this.fee.pool.difficulty)

//                 return
//             }
//             /*
//             Set second/fee pool active
//             */

//             this.fee.active = true

//             /*
//             Send required params from second/fee pool to miner which required
//             */
//             //TODO ADD CHECK IF extranonce1 and extranonce2_size null/empty
//             this.sendExtranonce(this.fee.pool.extranonce1, this.fee.pool.extranonce2_size, true)
//             this.sendVersionMask(this.fee.pool.version_rolling_mask_pool, true)
//             this.sendDifficulty(this.fee.pool.difficulty, true)
//             this.sendJob([...this.fee.pool.job], true, true)

//             logger.info(this.logPrefix, 'Fee enabled')

//             this.fee.timeout = setTimeout(() => {
//                 /*
//                 Set active main pool
//                 */

//                 this.fee.active = false
//                 /*
//                 Send required params from main pool to miner which required
//                 */

//                 this.sendExtranonce(this.pool.extranonce1, this.pool.extranonce2_size)
//                 this.sendVersionMask(this.pool.version_rolling_mask_pool)
//                 this.sendDifficulty(this.pool.difficulty)
//                 this.sendJob([...this.pool.job], false, true)

//                 logger.info(this.logPrefix, 'Fee disabled')

//                 /*
//                 Set interval
//                 */

//             }, (window / 100 * percent) * 1000)
//         }, window * 1000)
//     }

//     /*
//     sendJob - universal function which send mining.notify to miner from main or second/fee pool

//     mining.notify
//     job_id[0] - ID of the job. Use this ID while submitting share generated from this job.
//     prevhash[1] - Hash of previous block.
//     coinb1[2] - Initial part of coinbase transaction.
//     coinb2[3] - Final part of coinbase transaction.
//     merkle_branch[4] - List of hashes, will be used for calculation of merkle root. This is not a list of all transactions, it only contains prepared hashes of steps of merkle tree algorithm. Please read some materials for understanding how merkle trees calculation works. Unfortunately this example don't have any step hashes included, my bad!
//     version[5] - Bitcoin block version.
//     nbits[6] - Encoded current network difficulty
//     ntime[7] - Current ntime/
//     clean_jobs[8] - When true, server indicates that submitting shares from previous jobs don't have a sense and such shares will be rejected. When this flag is set, miner should also drop all previous jobs, so job_ids can be eventually rotated.

//     Example:
// {
//    "method":"mining.notify",
//    "params":[
//       "a3ea",
//       "4818f9a4fa448534986f75c8b77384ef7dd49aba000008410000000000000000",
//       "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff300345a20c142f6375626963685f706f6f6c2f6375626963682f12eaa353ece6582035f91b",
//       "ffffffff05be6bee070000000017a9146a2c207193af5b48250b93d0156923dcc0fde249877cd7dc0f0000000017a9146a2c207193af5b48250b93d0156923dcc0fde249879da1e50b0000000017a9146a2c207193af5b48250b93d0156923dcc0fde24987e135f7030000000017a9146a2c207193af5b48250b93d0156923dcc0fde249870000000000000000266a24aa21a9eda5c4f6d3a60662bf6d834b76b06b9cc10cc0d4224f57d0434909a24adb94563400000000",
//       [
//          "938c68b28a729d553eafa73af6f86f5b7f85c1771eb8b53e85cccb5f46aaa2b9",
//          "e20cb7a8c12ceb339580332bdcbc8d41171d8acd47eb6f443639dca0be74a541",
//          "ce9748b3e800329dd47949ac392bb415ceae5bd6bcb99963a1600c5a63835de5",
//          "eeb3a92ecd0005785472bbedeec03a7efbd597e31bb53156caea88a728958ec5",
//          "e72e04fa50abe1d27d8858e873874368ec23665486866c818938ca9b537cc1dc",
//          "a7d7f528f7ac6a7384e6a9acb19d5bf6a5f3fd2aa458f82c182452ba8f3db658",
//          "ab28930427dc7404d59f7a700c45fab34a89dea38723e212dc9a4e8b054ee5a4",
//          "3bc79fe447d2af0d1b855f8abe6036fa9ebe2e9962255a7838d4a5117ff234e4",
//          "390ead13914596b26bdc9a1a773e96b2c38a0f3143edc24749c4aa8a8aa0c29f",
//          "b3888243dafc8ed007f645f07be6d5e6e82fdb69187379a692eb05ffc9895d9f",
//          "00d44022f7b2c3d4992ad2b8081dd689dd5503e459f4351e07e5e9619cf98621",
//          "b69d8f729c4ed79037a90a3b59285464f1f15c861a84a74357471282005f60a1"
//       ],
//       "20000010",
//       "1704005a",
//       "65b7b39b",
//       false
//    ],
//    "id":null
// }
// Docs: https://ru.braiins.com/stratum-v1/docs
// if you have additional questions or write me
//     */

//     sendJob(job, fee = false, cleanJobs = false) {
//         /*
//         check if need to send clean previous jobs, it need when change between pools
//         */

//         if (cleanJobs)
//             job[8] = true

//         const obj = {
//             id: null,
//             method: 'mining.notify',
//             params: job
//         }
//         /*
//         Check is second/fee pool active, if active create array which contain jobs_id like 1-1f22
//         job_id usually is hex string which start and increases like 1,2,3,4 ... 9,f,1a,1b,1c
//         */

//         if (this.fee.active && fee) {
//             const uniqid = this.id + '-' + job[0]
//             this.fee.jobs.push({
//                 id: job[0],
//                 uniqid: uniqid
//             })
//             /*
//             Shift first job_id if while second/fee active and send many mining.notify
//             */

//             if (this.fee.jobs.length > 4)
//                 this.fee.jobs.shift()
//             job[0] = uniqid
//             /*
//             Send mining.notify to miner
//             */

//             this.send(obj)
//         }
//         /*
//         Send mining.notify to miner from main pool
//         */

//         if (!this.fee.active && !fee)
//             this.send(obj)
//     }

//     /*
//     sendDifficulty- send to miner mining.set_difficulty
//     params is threshold value below which pool will reject share(mining.submit)

//     Example:
// {
//    "id":null,
//    "method":"mining.set_difficulty",
//    "params":[
//       16100
//    ]
// }
//     */

//     sendDifficulty(difficulty, fee = false) {
//         if ((this.fee.active && fee) || (!this.fee.active && !fee)) {
//             this.send({
//                 id: null,
//                 method: 'mining.set_difficulty',
//                 params: [difficulty]
//             })
//         }
//     }

//     /*
//     sendVersionMask- send version mask to miner from main or second/fee pool

//     Example:
// {
//    "error":null,
//    "id":1,
//    "result":{
//       "version-rolling":true,
//       "version-rolling.mask":"18000000"
//    }
// }
// Or:
// {
//    "id":null,
//    "method":"mining.set_version_mask",
//    "params":[
//       "1fffe000"
//    ],
//    "error":null
// }
// Docs:
// https://en.bitcoin.it/wiki/BIP_0310
//     */

//     sendVersionMask(mask, fee = false) {
//         if ((this.fee.active && fee) || (!this.fee.active && !fee)) {
//             this.send({
//                 id: null,
//                 method: 'mining.set_version_mask',
//                 params: [mask],
//                 error: null
//             })
//         }
//     }

//     /*
// TODO ADD EXAMPLE
// Docs:
// https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.set_extranonce
// https://bitcoinwiki.org/wiki/stratum-mining-protocol
// https://github.com/DanielKrawisz/go-Stratum/blob/master/README.md maybe you can use this code? i dont check it
//     */

//     sendExtranonce(extranonce1, extranonce2_size, fee = false) {
//         if ((this.fee.active && fee) || (!this.fee.active && !fee)) {
//             this.send({
//                 id: null,
//                 method: 'mining.set_extranonce',
//                 params: [
//                     extranonce1,
//                     extranonce2_size
//                 ]
//             })
//         }
//     }

//     /*
//     Check if socket alive and send obj to miner(message)
//     */

//     send(obj) {
//         if (!this.socket || !this.socket.writable)
//             return
//         this.socket.write(JSON.stringify(obj) + '\n')
//     }

//     /*
//     setEvent- read incoming data from miner
//     */

//     setEvents() {
//         this.dataBuffer = ''

//         this.socket.on('data', data => {
//             this.dataBuffer += data
//             if (Buffer.byteLength(this.dataBuffer, 'utf8') > 102400) {
//                 logger.warn(this.logPrefix, 'Excessive packet size');
//                 this.socket.destroy()
//                 return
//             }
//             if (!this.dataBuffer.includes('\n'))
//                 return
//             const messages = this.dataBuffer.split('\n')
//             this.dataBuffer = this.dataBuffer.slice(-1) === '\n' ? '' : messages.pop()
//             for (let i = 0; i < messages.length; i++) {
//                 if (messages[i].trim() === '') continue
//                 /*
//                 Check if data valid and can be parsed to json
//                 */
//                 let obj
//                 try {
//                     obj = JSON.parse(messages[i])
//                 } catch (err) {
//                     if (messages[i].includes('method'))
//                         logger.warn(this.logPrefix, 'Malformed message from miner:', messages[i])
//                     this.socket.destroy()
//                     break
//                 }
//                 this.handleMessage(obj)
//             }
//         })

//         this.socket.on('error', err => {
//             if (err.code !== 'ECONNRESET')
//                 logger.warn(this.logPrefix, 'Miner connection error:', err)
//         })

//         this.socket.on('close', () => {
//             logger.info(this.logPrefix, 'Miner connection closed')
//             this.socket = null
//             this.dataBuffer = ''
//             this.fee.job = []
//             this.ctrl.deleteMiner(this.port, this.id)
//             clearInterval(this.fee.interval)
//             clearTimeout(this.fee.timeout)
//             if (this.pool.socket)
//                 this.pool.socket.destroy()
//             if (this.fee.pool.socket)
//                 this.fee.pool.socket.destroy()
//         })
//     }

//     /*
//     Handle message from miner
//     */

//     handleMessage(obj) {
//         // logger.info(this.logPrefix, `Method ${obj.method}`)

//         switch (obj.method) {
//             /*
//             handle mining.authorize
//             Example:
// {
//    "id":26314,
//    "method":"mining.authorize",
//    "params":[
//       "tests9k.72th",
//       "123"
//    ]
// }
//             */
//             case 'mining.authorize':
//                 /*
//                 check if username\account name is empty,if empty set custom,anyway send to main pool
//                 TODO add to config
//                 */

//                 if (obj.params[0] === '') {
//                     obj.params[0] = 'test'
//                     logger.info(this.logPrefix, `Null worker: ${obj.params[0]}`)
//                     break
//                 }
//                 logger.info(this.logPrefix, `Worker: ${obj.params[0]}`)
//                 break
//             /*
//             handle mining.subsribe
//             there are many variations,usually send one or two params
//             1.must name and version of mining software in the given format or empty string
//             2.must be session id received during a previous session or null if a new session should be initiated by the server
//             3.must be the host the client is trying to connect to or null
//             4.must be the port on the given host the client is trying to connect to or null
//             Example:
// {
//    "id":26313,
//    "method":"mining.subscribe",
//    "params":[
//       "whatsminer/v1.0",
//       "27dce2daa197369b"
//    ]
// }
// Docs:
// https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.subscribe
// ONLY FOR INFO,THIS IS NOT FOR SHA256 STRATUM :
// https://github.com/aeternity/protocol/blob/master/STRATUM.md#mining-subscribe
//             */

//             case 'mining.subscribe':
//                 this.fee.pool.status.push('mining.subscribe')
//                 logger.debug(this.logPrefix, `Subscribe : ${obj.params[0]}`)
//                 /*
//                 save useragent for future(setFee function)
//                  */
//                 this.pool.useragent = obj.params[0]

//                 this.fee.pool.send(obj)
//                 break
//             /* handle, parse mining.configure and send to main pool and second/fee pool, add to array status of miner connection,it need for setFee
// Example:
// {
//    "id":60584,
//    "method":"mining.configure",
//    "params":[
//       [
//          "version-rolling"
//       ],
//       {
//          "version-rolling.mask":"1fffe000"
//       }
//    ]
// }
// DOCS: check in file DOCS.md there are many params
// TODO ADD PARAMS
// https://en.bitcoin.it/wiki/BIP_0310#Request_.22mining.configure.22
//              */

//             case 'mining.configure':
//                 try {
//                     var params_mask = obj.params[1]
//                     var version = JSON.parse(JSON.stringify(params_mask))

//                     this.fee.pool.send(obj)
//                     this.pool.version_rolling_mask_miner = version['version-rolling.mask']
//                     this.fee.pool.version_rolling_mask_miner = version['version-rolling.mask']
//                     this.pool.status.push('mining.configure')
//                     this.fee.pool.status.push('mining.configure')
//                     if (version['version-rolling.min-bit-count'] !== undefined) {
//                         this.pool.version_rolling_min_bit_count = version['version-rolling.min-bit-count']
//                         this.fee.pool.version_rolling_min_bit_count = version['version-rolling.min-bit-count']

//                     }
//                 } catch (err) {
//                     logger.warn(this.logPrefix, 'Malformed mining.configure from miner:', obj)
//                 }

//                 // logger.info(this.logPrefix, `pool.version_rolling_mask_miner : ${this.pool.version_rolling_mask_miner}`)
//                 // logger.info(this.logPrefix, `fee.pool.version_rolling_mask_miner: ${this.fee.pool.version_rolling_mask_miner}`)
//                 // logger.info(this.logPrefix, `pool.version_rolling_min_bit_count: ${this.pool.version_rolling_min_bit_count}`)
//                 // logger.info(this.logPrefix, `fee.pool.version_rolling_min_bit_count: ${this.fee.pool.version_rolling_min_bit_count}`)
//                 break

//             /*
//             handle mining.submit
//             Check if job_id not in second/fee job_id array,if it exist,replace job_id and send to second\fee pool and send to miner result true,if not send to main pool
//              */

//             case 'mining.submit':
//                 let fee = this.fee.jobs.find(job => job.uniqid === obj.params[1])
//                 if (fee) {
//                     obj.params[0] = this.ctrl.params.fee.worker
//                     obj.params[1] = fee.id
//                     this.fee.pool.send(obj)
//                     this.send({id: obj.id, result: true, error: null})
//                     return
//                 }

//         }
//         /*
//         Send  messages to main pool
//          */

//         this.pool.send(obj)
//     }
// }
