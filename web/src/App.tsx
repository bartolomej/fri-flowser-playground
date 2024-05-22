import './App.css'
import configureCadence from '@/common/candance';

import {useEffect, useState} from 'react'
import Editor, {Monaco} from '@monaco-editor/react';
import {ProjectFile, ProjectLog, ProjectService} from "@/common/project.service.ts";

function App() {
    const LANGUAGE_CADENCE = 'cadence';
    const [openFile, setOpenFile] = useState<ProjectFile>()
    const [args, setArgs] = useState('');
    const [projectLogs, setProjectLogs] = useState<ProjectLog[]>();
    const [projectFiles, setProjectFiles] = useState<ProjectFile[]>();
    const [executionResult, setExecutionResult] = useState<unknown>();
    const service = new ProjectService({baseUrl: "http://localhost:8080"});
    const urlParams = new URLSearchParams(window.location.search);
    const projectUrl = urlParams.get('projectUrl');

    async function onExecute() {
        if (!openFile) {
            return;
        }

        const isScript = openFile.content.includes("pub fun main");
        const isTransaction = openFile.content.includes("transaction");

        // A very dump heuristic to determine if Cadence code is a transaction or script
        if (isScript) {
            setExecutionResult(await service.executeScript({
                source: openFile.content,
                arguments: args,
                location: openFile.path
            }))
        } else if (isTransaction) {
            setExecutionResult(await service.executeTransaction({
                source: openFile.content,
                arguments: args,
                location: openFile.path
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

    useEffect(() => {
        if (projectUrl) {
            const interval = setInterval(async () => {
                setProjectLogs(await service.listProjectLogs())
            }, 1000);

            return () => clearInterval(interval)
        }
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
                            <div key={file.path} onClick={() => setOpenFile(file)}
                                 className="max-w-[200px] truncate text-left">
                                {fileName}
                            </div>
                        )
                    })}
            </div>

            <div className="flex flex-col w-full">
                {openFile ? (
                    <Editor
                        theme='vs-dark'
                        language={LANGUAGE_CADENCE}
                        value={openFile?.content ?? ""}
                        onChange={code => setOpenFile({...openFile, content: code ?? ""})}
                        className="h-[60vh] pt-2 w-full"
                        options={{automaticLayout: true}}
                        beforeMount={beforeEditorMount}
                    />
                ) : (
                    <div>No files open</div>
                )}

                <div className="h-[40vh] flex flex-row">
                    <pre>
                        {projectLogs
                            ?.filter(log => log.level !== "debug")
                            ?.sort((a, b) => b.time.getTime() - a.time.getTime())
                            .map((log, i) => (
                                <div key={i}>[{log.level}][{getFormattedTime(log)}] {log.msg}</div>
                            ))}
                    </pre>
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

function getFormattedTime(log: ProjectLog): string {
    return log.time.toISOString().split("T")[1].split(".")[0]
}

export default App;
