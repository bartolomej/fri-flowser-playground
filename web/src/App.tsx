import './App.css'
import configureCadence from '@/common/candance';

import { useEffect, useState } from 'react'
import Editor, { Monaco } from '@monaco-editor/react';
import { BlockchainState, ProjectFile, ProjectLog, ProjectService } from "@/common/project.service.ts";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"


function App() {
    const LANGUAGE_CADENCE = 'cadence';
    const [openFile, setOpenFile] = useState<ProjectFile>()
    const [args, setArgs] = useState('');
    const [blockchainState, setBlockchainState] = useState<BlockchainState>();
    const [projectLogs, setProjectLogs] = useState<ProjectLog[]>();
    const [projectFiles, setProjectFiles] = useState<ProjectFile[]>();
    const [executionResult, setExecutionResult] = useState<unknown>();
    const service = new ProjectService({ baseUrl: "http://localhost:8080" });
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

    useEffect(() => {
        if (projectUrl) {
            const interval = setInterval(async () => {
                setBlockchainState(await service.getProjectBlockchainState())
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
        <div className='flex flex-row w-[100vw]'>
            <div className="flex flex-col gap-y-[10px] w-[20%] p-2">
                {projectFiles
                    ?.filter(file => !file.isDirectory)
                    .map(file => {
                        const fileName = file.path.split("/").reverse()[0]
                        return (
                            <div key={file.path} onClick={() => setOpenFile(file)}
                                className="max-w-[200px] truncate text-left hover:cursor-pointer hover:opacity-20">
                                {fileName}
                            </div>
                        )
                    })}
            </div>

            <div className="flex flex-col w-[80%]">
                {openFile ? (
                    <Editor
                        theme='vs-dark'
                        language={LANGUAGE_CADENCE}
                        value={openFile?.content ?? ""}
                        onChange={code => setOpenFile({ ...openFile, content: code ?? "" })}
                        className="h-[60vh] pt-2 w-full"
                        options={{ automaticLayout: true }}
                        beforeMount={beforeEditorMount}
                    />
                ) : (
                    <div>No files open</div>
                )}

                <div className="h-[38vh] flex flex-row w-[100%]">
                    <Tabs defaultValue="json" className="w-[100%]">
                        <TabsList className='bg-[#242424]'>
                            <TabsTrigger className='border border-1 border-white' value="json">JSON</TabsTrigger>
                            <TabsTrigger className='border border-1 border-white' value="log">Log</TabsTrigger>
                            <TabsTrigger className='border border-1 border-white' value="execute">Execute</TabsTrigger>
                        </TabsList>
                        <TabsContent className='w-[100%]' value="json">
                            <pre className="overflow-scroll h-[38vh]">
                                {JSON.stringify(blockchainState, null, 4)}
                            </pre>
                        </TabsContent>
                        <TabsContent value="log" className='w-[100%] h-[38vh]'>
                            <pre className='max-h-[38vh] h-[38vh] overflow-y-auto'>
                                {projectLogs
                                    ?.filter(log => log.level !== "debug")
                                    ?.sort((a, b) => b.time.getTime() - a.time.getTime())
                                    .map((log, i) => (
                                        <div key={i}>[{log.level}][{getFormattedTime(log)}] {log.msg}</div>
                                    ))}
                            </pre>
                        </TabsContent>
                        <TabsContent value="execute" className='w-[100%]'>
                            <div className='h-[38vh]'>
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
                        </TabsContent>
                    </Tabs>
                </div>
            </div>
        </div>
    )
}

function getFormattedTime(log: ProjectLog): string {
    return log.time.toISOString().split("T")[1].split(".")[0]
}

export default App;
