import './App.css'
import configureCadence from '@/common/candance';

import { useEffect, useState } from 'react'
import Editor, { Monaco } from '@monaco-editor/react';
import { BlockchainState, ProjectFile, ProjectLog, ProjectService } from "@/common/project.service.ts";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/tabs.tsx"
import {JsonView} from "@/components/JsonView.tsx";


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

    const sidebarWidth = 300;

    return (
        <div className='flex flex-row'>
            <div className="flex flex-col gap-y-[10px] flex-1 p-2" style={{width: sidebarWidth}}>
                <b>PROJECT FILES</b>
                {projectFiles
                    ?.filter(file => !file.isDirectory)
                    .map(file => {
                        const fileName = file.path.split("/").reverse()[0]
                        return (
                            <div key={file.path} onClick={() => setOpenFile(file)}
                                className="truncate text-left hover:cursor-pointer transition-all hover:translate-x-[5px]">
                                {fileName}
                            </div>
                        )
                    })}
            </div>

            <div className="flex flex-col w-full flex-1" style={{width: `calc(100vh - ${sidebarWidth}px)`}}>
                <div className="h-[60vh]">
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
                </div>

                <div className="h-[40vh] max-h-[40vh] flex flex-row w-[100%]">
                    <Tabs defaultValue="json" className="w-[100%]">
                        <TabsList className='bg-[#242424]'>
                            <TabsTrigger className='border border-1 border-white' value="state">State</TabsTrigger>
                            <TabsTrigger className='border border-1 border-white' value="logs">Logs</TabsTrigger>
                            <TabsTrigger className='border border-1 border-white' value="execute">Execute</TabsTrigger>
                        </TabsList>

                        <TabsContent className='w-[100%]' value="state">
                            {blockchainState ? <JsonView style={{height: "100%", borderRadius: 10, padding: 10}} name="blockchain" src={blockchainState} /> : "Loading..."}
                        </TabsContent>

                        <TabsContent value="logs" className="h-full">
                            <pre className='overflow-y-auto h-full'>
                                {projectLogs
                                    ?.filter(log => log.level !== "debug")
                                    ?.sort((a, b) => b.time.getTime() - a.time.getTime())
                                    .map((log, i) => (
                                        <div key={i}>[{log.level}][{getFormattedTime(log)}] {log.msg}</div>
                                    ))}
                            </pre>
                        </TabsContent>

                        <TabsContent value="execute" className='w-[100%]'>
                            <div className="flex gap-x-[10px]">
                                <div className="flex-1">
                                    Arguments:
                                    <Editor
                                        theme='vs-dark'
                                        language={"JavaScript"}
                                        value={args}
                                        onChange={code => setArgs(code ?? "")}
                                        options={{automaticLayout: true}}
                                        height={200}
                                    />
                                    <button onClick={onExecute}>Execute</button>
                                </div>
                                <div className="flex-1">
                                    {executionResult ? (
                                        <pre>
                                            <JsonView style={{height: "100%", borderRadius: 10, padding: 10}} name="result" src={executionResult} />
                                        </pre>
                                    ) : (
                                        <div className="flex justify-center items-center h-full">Execute to see results</div>
                                    )}
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
