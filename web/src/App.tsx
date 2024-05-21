import './App.css'
import configureCadence from '@/common/candance';

import {useEffect, useState} from 'react'
import Editor, {Monaco} from '@monaco-editor/react';
import {ProjectFile, ProjectService} from "@/common/project.service.ts";

function App() {
    const LANGUAGE_CADENCE = 'cadence';
    const [code, setCode] = useState('');
    const [args, setArgs] = useState('');
    const [projectFiles, setProjectFiles] = useState<ProjectFile[]>();
    const [executionResult, setExecutionResult] = useState<unknown>();
    const service = new ProjectService({baseUrl: "http://localhost:8080"});
    const urlParams = new URLSearchParams(window.location.search);
    const projectUrl = urlParams.get('projectUrl');

    async function onExecute() {
        const isScript = code.includes("pub fun main");
        const isTransaction = code.includes("transaction");

        // A very dump heuristic to determine if Cadence code is a transaction or script
        if (isScript) {
            setExecutionResult(await service.executeScript({
                source: code,
                arguments: args
            }))
        } else if (isTransaction) {
            setExecutionResult(await service.executeTransaction({
                source: code,
                arguments: args
            }))
        }
    }

    useEffect(() => {
        (async function () {
            if (projectUrl) {
                await service.openProject(projectUrl)
                setProjectFiles(await service.listProjectFiles())
            }
        })()
    }, [projectUrl]);

    const beforeEditorMount = (monaco: Monaco) => {
        configureCadence(monaco);
    }

    if (!projectUrl) {
        return (
            <div>
                Set `projectUrl` query parameter
            </div>
        )
    }

    return (
        <div className='flex flex-row'>
            <div className="flex flex-col gap-y-[10px] p-2">
                {projectFiles
                    ?.filter(file => !file.isDirectory)
                    .map(file => {
                        const fileName = file.path.split("/").reverse()[0]
                        return (
                            <div key={file.path} onClick={() => setCode(file.content)}
                                 className="max-w-[200px] truncate text-left">
                                {fileName}
                            </div>
                        )
                    })}
            </div>

            <div className="flex flex-col w-full">
                <Editor
                    theme='vs-dark'
                    language={LANGUAGE_CADENCE}
                    value={code}
                    onChange={code => setCode(code ?? "")}
                    className="h-[60vh] pt-2 w-full"
                    options={{automaticLayout: true}}
                    beforeMount={beforeEditorMount}
                />

                <div className="h-[40vh] flex flex-row">
                    <pre>
                        {JSON.stringify(executionResult, null, 4)}
                    </pre>
                    <div>
                        <label>
                            Arguments
                            <textarea rows={10} value={args} onChange={e => setArgs(e.target.value)}></textarea>
                        </label>
                        <button onClick={onExecute}>Execute</button>
                    </div>
                </div>
            </div>
        </div>
    )
}

export default App;
