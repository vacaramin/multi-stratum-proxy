package sha256stratum

import (
	"log"
	"net"
	"time"
)

// TO DO
// Defines two classes: Pool and PoolToFee.
// Pool connects to a mining pool via stratum protocol and handles communication with the pool.
// PoolToFee connects to a separate fee pool and forwards relevant job information for fee calculation.
type Pool struct {
	Miner        *Miner
	MiningServer string
	LogPrefix    string
	Socket       net.Conn
	NoncePrefix  *string
	Job          map[string]interface{}
	DataBuffer   string
}

// NewPool creates a new Pool instance(Similar to a constructor implementation)
func NewPool(miner *Miner) *Pool {
	return &Pool{
		Miner:        miner,
		MiningServer: miner.MiningServer,
		LogPrefix:    miner.LogPrefix,
		NoncePrefix:  nil,
		Job:          nil,
	}
}

// Connect is a method of Pool
func (p *Pool) Connect() {
	log.Printf("Connecting to pool: %s", p.MiningServer)

	socket, err := net.Dial("tcp", p.MiningServer)
	if err != nil {
		log.Printf("%s Error: %s", p.LogPrefix, err)
		return
	}
	p.Socket = socket
	socket.SetDeadline(time.Now().Add(10 * time.Second))

}

type PoolToFee struct {
	Miner       *Miner
	Params      map[string]interface{}
	LogPrefix   string
	Socket      net.Conn
	Status      string
	NoncePrefix string
	Job         map[string]interface{}
	DataBuffer  string
}

// const net = require('net')
// const logger = require('../../lib/logger')

// class Pool {
//     constructor(miner) {
//         this.miner = miner
//         this.miningServer = miner.miningServer
//         this.logPrefix = miner.logPrefix
//         this.status = []
//         this.extranonce1 = ''
//         this.extranonce2_size = ''
//         this.difficulty = null
//         this.job = []
//         this.version_rolling_mask_miner = ''
//         this.version_rolling_mask_pool = ''
//         this.version_rolling = true
//         this.version_rolling_min_bit_count = ''
//         this.useragent = ''
//     }

//     /*
//     Connect to main pool
//     split btc.viabtc.io:3333 to btc.viabtc.io 3333
//      */
//     connect() {
//         logger.info(this.logPrefix, `Connecting to pool ${this.miningServer}`)

//         const [host, port] = this.miningServer.split(':')

//         this.socket = new net.Socket()
//         this.socket.connect(port, host)
//         this.socket.setKeepAlive(true, 60000)
//         this.setEvents()
//     }

//     /*
//     Check if socket alive and send obj to main pool(message)
//     */
//     send(obj) {
//         if (!this.socket || !this.socket.writable)
//             return
//         this.socket.write(JSON.stringify(obj) + '\n')
//     }

//     /*
//     setEvent- read incoming data from main pool
//     */
//     setEvents() {
//         this.dataBuffer = ''

//         this.socket.on('connect', () => {
//             logger.info(this.logPrefix, 'Connected to pool')
//             this.miner.setEvents()
//         })

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
//                         logger.warn(this.logPrefix, 'Malformed message from pool:', messages[i])
//                     this.socket.destroy()
//                     break
//                 }
//                 this.handleMessage(obj)
//             }
//         })
// //TODO add other errors
//         this.socket.on('error', err => {
//             if (err.code !== 'ECONNRESET')
//                 logger.warn(this.logPrefix, 'Pool connection error:', err)
//         })

//         this.socket.on('close', () => {
//             logger.info(this.logPrefix, 'Pool connection closed')
//             this.socket = null
//             this.dataBuffer = ''
//             if (this.miner.socket)
//                 this.miner.socket.destroy()
//         })
//     }

//     /*
//     Handle message from main pool
//     */
//     handleMessage(obj) {
//         //  logger.debug(this.logPrefix, this.status)
//         /* handle, parse mining.configure response from main pool and send response to miner,
// Example:
// {
//    "id":60584,
//    "result":{
//       "version-rolling":true,
//       "version-rolling.mask":"1fffe000"
//    },
//    "error":null
// }
// DOCS: check in file DOCS.md there are many variations of result
// TODO ADD DESCRIBE RESPONSE
// https://en.bitcoin.it/wiki/BIP_0310#Request_.22mining.configure.22
//          */
//         if (this.miner.pool.status.includes('mining.configure')) {
//             try {
//                 var result = JSON.parse(JSON.stringify(obj.result))
//                 this.version_rolling_mask_pool = result['version-rolling.mask']
//                 this.version_rolling = result['version-rolling']
//                 logger.debug(this.logPrefix, 'version-rolling.mask_pool: ', this.version_rolling_mask_pool)
//                 logger.debug(this.logPrefix, 'version-rolling pool : ', this.version_rolling)
//                 /*
//                                 Clear array status,it wrong code which need to be rewritten but i have no time for this :(
//                                 purpose of this array is consistency of each method was send from miner to main pool need to be duplicated to second/fee pool to imitate real miner and get params to switch between pools
//                                 it would be nice if second pool will connect to second/fee pool after successful connection miner to pool( all methods and params need to be saved to imitate connection)
//                                  */
//                 this.miner.pool.status = this.miner.pool.status.filter(function (f) {
//                     return f != 'mining.configure'
//                 })
//             } catch (err) {
//                 logger.warn(this.logPrefix, 'Not mining.configure response: ', obj)
//             }
//         }//TODO CHECK PARAMS AND DATA WHICH ARE BEING PROCESSED HERE
//         if (Array.isArray(obj.result) && obj.result.length === 3) {
//             this.extranonce1 = obj.result[1]
//             this.extranonce2_size = obj.result[2]
//         }

//         switch (obj.method) {
//             /*
//             handle mining.notify and send to sendJob function which check which pool active

//             mining.notify
//             job_id[0] - ID of the job. Use this ID while submitting share generated from this job.
//             prevhash[1] - Hash of previous block.
//             coinb1[2] - Initial part of coinbase transaction.
//             coinb2[3] - Final part of coinbase transaction.
//             merkle_branch[4] - List of hashes, will be used for calculation of merkle root. This is not a list of all transactions, it only contains prepared hashes of steps of merkle tree algorithm. Please read some materials for understanding how merkle trees calculation works. Unfortunately this example don't have any step hashes included, my bad!
//             version[5] - Bitcoin block version.
//             nbits[6] - Encoded current network difficulty
//             ntime[7] - Current ntime/
//             clean_jobs[8] - When true, server indicates that submitting shares from previous jobs don't have a sense and such shares will be rejected. When this flag is set, miner should also drop all previous jobs, so job_ids can be eventually rotated.

//             Example:
//         {
//            "method":"mining.notify",
//            "params":[
//               "a3ea",
//               "4818f9a4fa448534986f75c8b77384ef7dd49aba000008410000000000000000",
//               "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff300345a20c142f6375626963685f706f6f6c2f6375626963682f12eaa353ece6582035f91b",
//               "ffffffff05be6bee070000000017a9146a2c207193af5b48250b93d0156923dcc0fde249877cd7dc0f0000000017a9146a2c207193af5b48250b93d0156923dcc0fde249879da1e50b0000000017a9146a2c207193af5b48250b93d0156923dcc0fde24987e135f7030000000017a9146a2c207193af5b48250b93d0156923dcc0fde249870000000000000000266a24aa21a9eda5c4f6d3a60662bf6d834b76b06b9cc10cc0d4224f57d0434909a24adb94563400000000",
//               [
//                  "938c68b28a729d553eafa73af6f86f5b7f85c1771eb8b53e85cccb5f46aaa2b9",
//                  "e20cb7a8c12ceb339580332bdcbc8d41171d8acd47eb6f443639dca0be74a541",
//                  "ce9748b3e800329dd47949ac392bb415ceae5bd6bcb99963a1600c5a63835de5",
//                  "eeb3a92ecd0005785472bbedeec03a7efbd597e31bb53156caea88a728958ec5",
//                  "e72e04fa50abe1d27d8858e873874368ec23665486866c818938ca9b537cc1dc",
//                  "a7d7f528f7ac6a7384e6a9acb19d5bf6a5f3fd2aa458f82c182452ba8f3db658",
//                  "ab28930427dc7404d59f7a700c45fab34a89dea38723e212dc9a4e8b054ee5a4",
//                  "3bc79fe447d2af0d1b855f8abe6036fa9ebe2e9962255a7838d4a5117ff234e4",
//                  "390ead13914596b26bdc9a1a773e96b2c38a0f3143edc24749c4aa8a8aa0c29f",
//                  "b3888243dafc8ed007f645f07be6d5e6e82fdb69187379a692eb05ffc9895d9f",
//                  "00d44022f7b2c3d4992ad2b8081dd689dd5503e459f4351e07e5e9619cf98621",
//                  "b69d8f729c4ed79037a90a3b59285464f1f15c861a84a74357471282005f60a1"
//               ],
//               "20000010",
//               "1704005a",
//               "65b7b39b",
//               false
//            ],
//            "id":null
//         }
//         Docs: https://ru.braiins.com/stratum-v1/docs
//         if you have additional questions or write me
//             */

//             case 'mining.notify':
//                 this.job = obj.params
//                 this.miner.sendJob([...obj.params])
//                 break
//             /*
//             Set required difficulty from main pool
//             params is threshold value below which main pool will reject share(mining.submit)
//             Example:
//             {
//                "method":"mining.set_difficulty",
//                "params":[
//                   2097152
//                ],
//                "id":null
//             }
//             Docs in docs.md file or comment of sendDifficulty function
//              */

//             case 'mining.set_difficulty':
//                 this.difficulty = obj.params[0]
//                 this.miner.sendDifficulty(obj.params[0])
//                 break
//             /*
//          TODO ADD EXAMPLE
//          Docs:
//          https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.set_extranonce
//          https://bitcoinwiki.org/wiki/stratum-mining-protocol
//          https://github.com/DanielKrawisz/go-Stratum/blob/master/README.md maybe you can use this code? i dont check it
//              */

//             case 'mining.set_extranonce':
//                 this.extranonce1 = obj.params[0]
//                 this.extranonce2_size = obj.params[1]
//                 this.miner.sendExtranonce(obj.params[0], obj.params[1])
//                 break
//             default:
//                 this.miner.send(obj)
//         }
//     }
// }

// class PoolToFee {
//     constructor(miner, params) {
//         this.miner = miner
//         this.params = params
//         this.logPrefix = miner.logPrefix + ' FEE >'

//         this.status = []
//         this.extranonce1 = ''
//         this.extranonce2_size = ''
//         this.difficulty = null
//         this.job = []
//         this.version_rolling_mask_miner = ''
//         this.version_rolling_mask_pool = ''
//         this.version_rolling = true
//         this.version_rolling_min_bit_count = ''
//     }

//     /*
//     Connect to second/fee pool
//     split btc.viabtc.io:3333 to btc.viabtc.io 3333
//      */
//     connect() {
//         logger.info(this.logPrefix, `Connecting to pool ${this.params.fee.pool}`)

//         const [host, port] = this.params.fee.pool.split(':')

//         this.socket = new net.Socket()
//         this.socket.connect(port, host)
//         this.socket.setKeepAlive(true, 60000)
//         this.setEvents()
//     }

//     /*
//     Check if socket alive and send obj to second/fee pool(message)
//     */
//     send(obj) {
//         if (!this.socket || !this.socket.writable)
//             return
//         this.socket.write(JSON.stringify(obj) + '\n')
//     }

//     /*
//     setEvent- read incoming data from second/fee pool
//     */
//     setEvents() {
//         this.dataBuffer = ''

//         this.socket.on('connect', () => {
//             logger.info(this.logPrefix, 'Connected to pool')

//             // this.status = 'mining.subscribe'
//             //this.send({
//             //    id: 1,
//             //    method: 'mining.subscribe',
//             //   params: []
//             //})
//         })

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
//                         logger.warn(this.logPrefix, 'Malformed message from pool:', messages[i])
//                     this.socket.destroy()
//                     break
//                 }
//                 this.handleMessage(obj)
//             }
//         })
// //TODO add other errors
//         this.socket.on('error', err => {
//             if (err.code !== 'ECONNRESET')
//                 logger.warn(this.logPrefix, 'Pool connection error:', err)

//         })

//         this.socket.on('close', () => {
//             logger.info(this.logPrefix, 'Pool connection closed')
//             this.socket = null
//             this.status = []
//             this.dataBuffer = ''
//             if (this.miner.socket)
//                 this.connect()
//         })
//     }

//     /*
//     Handle message from second/fee
//     */
//     handleMessage(obj) {
//         // logger.debug(this.logPrefix, obj)
//         //  logger.debug(this.logPrefix, this.status)
//         /* handle, parse mining.configure response from second/fee pool and save response for setFee function(miner),
// Example:
// {
//    "id":60584,
//    "result":{
//       "version-rolling":true,
//       "version-rolling.mask":"1fffe000"
//    },
//    "error":null
// }
// DOCS: check in file DOCS.md there are many variations of result
// TODO ADD DESCRIBE RESPONSE
// https://en.bitcoin.it/wiki/BIP_0310#Request_.22mining.configure.22
//          */

//         if (this.miner.fee.pool.status.includes('mining.configure')) {
//             try {
//                 var result = JSON.parse(JSON.stringify(obj.result))
//                 this.version_rolling_mask_pool = result['version-rolling.mask']
//                 this.version_rolling = result['version-rolling']
//                 logger.debug(this.logPrefix, 'version-rolling.mask_pool: ', this.version_rolling_mask_pool)
//                 logger.debug(this.logPrefix, 'version-rolling pool : ', this.version_rolling)
//                 /*
//                 Clear array status,it wrong code which need to be rewritten but i have no time for this :(
//                 purpose of this array is consistency of each method was send from miner to main pool need to be duplicated to second/fee pool to imitate real miner and get params to switch between pools
//                 it would be nice if second pool will connect to second/fee pool after successful connection miner to pool( all methods and params need to be saved to imitate connection)
//                  */
//                 this.status = this.status.filter(function (f) {
//                     return f != 'mining.configure'
//                 })
//             } catch (err) {
//                 logger.debug(this.logPrefix, 'Not [] mining.configure response: ', obj)

//             }
//         }
//         /*
//         Check response of mining.subscribe method and params from second/fee pool
//          */
//         //if (this.miner.fee.pool.status === 'mining.subscribe' && obj.id === 1) {
//         if (this.miner.fee.pool.status.includes('mining.subscribe')) {
//             if (obj.error) {
//                 logger.error(this.logPrefix, `Subscribe error: ${obj.error}`)
//                 this.socket.destroy()
//                 return
//             }
//             try {
//                 this.extranonce1 = obj.result[1]
//                 this.extranonce2_size = obj.result[2]
//                 this.miner.sendExtranonce(obj.result[1], obj.result[2], true)
//                 this.status.push('mining.authorize')
//                 /*
//                                 Push next stage,it wrong code which need to be rewritten but i have no time for this :(
//                                 purpose of this array is consistency of each method was send from miner to main pool need to be duplicated to second/fee pool to imitate real miner and get params to switch between pools
//                                 it would be nice if second pool will connect to second/fee pool after successful connection miner to pool( all methods and params need to be saved to imitate connection)
//                                  */
//                 this.status = this.status.filter(function (f) {
//                     return f != 'mining.subscribe'
//                 })

//                 logger.info(this.logPrefix, `Authentication: ${this.params.fee.worker}`)
//                 /*
//                 No need for now,it will be need after basic understanding
//                  */

//                 if (this.params.fee.extranonce_subscribe)
//                     this.send({id: 2, method: 'mining.extranonce.subscribe', params: []})
//                 this.send({
//                     id: 3,
//                     method: 'mining.authorize',
//                     params: [
//                         this.params.fee.worker,
//                         this.params.fee.pass
//                     ]
//                 })
//             } catch (err) {
//                 logger.debug(this.logPrefix, 'TypeError, sometimes it can be: ', obj)

//             }
//         } else if (this.status.includes('mining.authorize') && obj.id === 3) {
//             if (obj.error) {
//                 logger.error(this.logPrefix, `Authentication error: ${obj.error.message}`)
//                 this.socket.destroy()
//                 return
//             }
//             this.status = this.status.filter(function (f) {
//                 return f != 'mining.authorize'
//             })
//             this.status.push('ready')
//         } else if (obj.error && Array.isArray(obj.error)) {
//             logger.error(this.logPrefix, `Error: ${JSON.stringify(obj.error)}`)
//         }
//         switch (obj.method) {
//             /*
//                         handle mining.notify and send to sendJob function which check which pool active

//                         mining.notify
//                         job_id[0] - ID of the job. Use this ID while submitting share generated from this job.
//                         prevhash[1] - Hash of previous block.
//                         coinb1[2] - Initial part of coinbase transaction.
//                         coinb2[3] - Final part of coinbase transaction.
//                         merkle_branch[4] - List of hashes, will be used for calculation of merkle root. This is not a list of all transactions, it only contains prepared hashes of steps of merkle tree algorithm. Please read some materials for understanding how merkle trees calculation works. Unfortunately this example don't have any step hashes included, my bad!
//                         version[5] - Bitcoin block version.
//                         nbits[6] - Encoded current network difficulty
//                         ntime[7] - Current ntime/
//                         clean_jobs[8] - When true, server indicates that submitting shares from previous jobs don't have a sense and such shares will be rejected. When this flag is set, miner should also drop all previous jobs, so job_ids can be eventually rotated.

//                         Example:
//                     {
//                        "method":"mining.notify",
//                        "params":[
//                           "a3ea",
//                           "4818f9a4fa448534986f75c8b77384ef7dd49aba000008410000000000000000",
//                           "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff300345a20c142f6375626963685f706f6f6c2f6375626963682f12eaa353ece6582035f91b",
//                           "ffffffff05be6bee070000000017a9146a2c207193af5b48250b93d0156923dcc0fde249877cd7dc0f0000000017a9146a2c207193af5b48250b93d0156923dcc0fde249879da1e50b0000000017a9146a2c207193af5b48250b93d0156923dcc0fde24987e135f7030000000017a9146a2c207193af5b48250b93d0156923dcc0fde249870000000000000000266a24aa21a9eda5c4f6d3a60662bf6d834b76b06b9cc10cc0d4224f57d0434909a24adb94563400000000",
//                           [
//                              "938c68b28a729d553eafa73af6f86f5b7f85c1771eb8b53e85cccb5f46aaa2b9",
//                              "e20cb7a8c12ceb339580332bdcbc8d41171d8acd47eb6f443639dca0be74a541",
//                              "ce9748b3e800329dd47949ac392bb415ceae5bd6bcb99963a1600c5a63835de5",
//                              "eeb3a92ecd0005785472bbedeec03a7efbd597e31bb53156caea88a728958ec5",
//                              "e72e04fa50abe1d27d8858e873874368ec23665486866c818938ca9b537cc1dc",
//                              "a7d7f528f7ac6a7384e6a9acb19d5bf6a5f3fd2aa458f82c182452ba8f3db658",
//                              "ab28930427dc7404d59f7a700c45fab34a89dea38723e212dc9a4e8b054ee5a4",
//                              "3bc79fe447d2af0d1b855f8abe6036fa9ebe2e9962255a7838d4a5117ff234e4",
//                              "390ead13914596b26bdc9a1a773e96b2c38a0f3143edc24749c4aa8a8aa0c29f",
//                              "b3888243dafc8ed007f645f07be6d5e6e82fdb69187379a692eb05ffc9895d9f",
//                              "00d44022f7b2c3d4992ad2b8081dd689dd5503e459f4351e07e5e9619cf98621",
//                              "b69d8f729c4ed79037a90a3b59285464f1f15c861a84a74357471282005f60a1"
//                           ],
//                           "20000010",
//                           "1704005a",
//                           "65b7b39b",
//                           false
//                        ],
//                        "id":null
//                     }
//                     Docs: https://ru.braiins.com/stratum-v1/docs
//                     if you have additional questions or write me
//                         */

//             case 'mining.notify':
//                 this.job = obj.params
//                 this.miner.sendJob([...obj.params], true)
//                 break
// /*
// Set required difficulty from second/fee pool
// params is threshold value below which second/fee pool will reject share(mining.submit)
// Example:
// {
//    "method":"mining.set_difficulty",
//    "params":[
//       2097152
//    ],
//    "id":null
// }
// Docs in docs.md file or comment of sendDifficulty function
//  */

//             case 'mining.set_difficulty':
//                 this.difficulty = obj.params[0]
//                 this.miner.sendDifficulty(obj.params[0], true)
//                 break
//             /*
//          TODO ADD EXAMPLE
//          Docs:
//          https://en.bitcoin.it/wiki/Stratum_mining_protocol#mining.set_extranonce
//          https://bitcoinwiki.org/wiki/stratum-mining-protocol
//          https://github.com/DanielKrawisz/go-Stratum/blob/master/README.md maybe you can use this code? i dont check it
//              */
//             case 'mining.set_extranonce':
//                 this.extranonce1 = obj.params[0]
//                 this.extranonce2_size = obj.params[1]
//                 this.miner.sendExtranonce(obj.params[0], obj.params[1], true)
//         }
//     }
// }

// module.exports = {
//     Pool,
//     PoolToFee
// }
