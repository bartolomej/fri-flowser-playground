export type ProjectFile = {
    path: string;
    isDirectory: boolean;
    content: string;
}

export type ProjectLog = Record<string, unknown> & {
    level: string;
    time: Date;
    msg: string;
}

// TODO: Define type if needed later
export type BlockchainState = object;

type ExecuteScriptRequest = {
    source: string;
    // Encoded using: https://cadence-lang.org/docs/json-cadence-spec
    arguments: any;
    // File location (full path)
    location: string;
}

type ExecuteTransactionRequest = {
    source: string;
    // Encoded using: https://cadence-lang.org/docs/json-cadence-spec
    arguments: any;
    // File location (full path)
    location: string;
}

type Config = {
    baseUrl: string;
}

export class ProjectService {
    constructor(private readonly config: Config) {
    }

    async getProjectBlockchainState(): Promise<BlockchainState> {
        return fetch(`${this.config.baseUrl}/projects/blockchain-state`).then(res => res.json());
    }

    async listProjectLogs(): Promise<ProjectLog[]> {
        return fetch(`${this.config.baseUrl}/projects/logs`).then(res => res.json()).then(logs => logs.map((log: string) => {
            const parsedLog = JSON.parse(log);

            return {
                ...parsedLog,
                time: new Date(parsedLog.time)
            }
        }));
    }

    async listProjectFiles(): Promise<ProjectFile[]> {
        return fetch(`${this.config.baseUrl}/projects/files`).then(res => res.json());
    }

    async openProject(projectUrl: string): Promise<void> {
        await fetch(`${this.config.baseUrl}/projects`, {
            method: "POST",
            body: JSON.stringify({projectUrl})
        });
    }

    async executeScript(request: ExecuteScriptRequest): Promise<unknown> {
        return fetch(`${this.config.baseUrl}/projects/scripts`, {
            method: "POST",
            body: JSON.stringify(request)
        }).then(res => res.json());
    }

    async executeTransaction(request: ExecuteTransactionRequest): Promise<unknown> {
        return fetch(`${this.config.baseUrl}/projects/transactions`, {
            method: "POST",
            body: JSON.stringify(request)
        }).then(res => res.json());
    }

}
