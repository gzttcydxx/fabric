'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        for (let i=0; i<this.roundArguments.assets; i++) {
            const did = `did:example:${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Creating DID: ${did}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'CreateIdentity',
                invokerIdentity: 'user1',
                contractArguments: [did],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }

    async submitTransaction() {
        const randomId = Math.floor(Math.random()*this.roundArguments.assets);
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'ReadIdentity',
            invokerIdentity: 'user1',
            contractArguments: [`did:example:${this.workerIndex}_${randomId}`],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    async cleanupWorkloadModule() {
        for (let i=0; i<this.roundArguments.assets; i++) {
            const did = `did:example:${this.workerIndex}_${i}`;
            console.log(`Worker ${this.workerIndex}: Deleting DID: ${did}`);
            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'DeleteIdentity',
                invokerIdentity: 'user1',
                contractArguments: [did],
                readOnly: false
            };

            await this.sutAdapter.sendRequests(request);
        }
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
