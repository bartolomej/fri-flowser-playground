export type ProjectFile = {
    path: string;
    isDirectory: boolean;
    content: string;
}

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
    constructor(private readonly config: Config) {}

    async listProjectFiles(): Promise<ProjectFile[]> {
        return fetch(`${this.config.baseUrl}/projects/files`).then(res => res.json());
    }

    async openProject(projectUrl: string): Promise<void> {
        await fetch(`${this.config.baseUrl}/projects`, {
            method: "POST",
            body: JSON.stringify({ projectUrl })
        });
    }

    async executeScript(request: ExecuteScriptRequest): Promise<unknown> {
        return fetch(`${this.config.baseUrl}/scripts`, {
            method: "POST",
            body: JSON.stringify(request)
        }).then(res => res.json());
    }

    async executeTransaction(request: ExecuteTransactionRequest): Promise<unknown> {
        return fetch(`${this.config.baseUrl}/transactions`, {
            method: "POST",
            body: JSON.stringify(request)
        }).then(res => res.json());
    }

}
